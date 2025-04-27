package app

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/laiker/auth/internal/closer"
	"github.com/laiker/auth/internal/config"
	"github.com/laiker/auth/internal/interceptor"
	"github.com/laiker/auth/pkg/access_v1"
	"github.com/laiker/auth/pkg/auth_v1"
	"github.com/laiker/auth/pkg/user_v1"
	_ "github.com/laiker/auth/statik"
	"github.com/pkg/errors"
	"github.com/rakyll/statik/fs"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

type App struct {
	serviceProvider *ServiceProvider
	grpcServer      *grpc.Server
	httpServer      *http.Server
	swaggerServer   *http.Server
	logger          *slog.Logger
}

func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		defer wg.Done()

		err := a.runGRPCServer()
		if err != nil {
			log.Fatalf("failed to run GRPC server: %v", err)
		}
	}()

	go func() {
		defer wg.Done()

		err := a.runHTTPServer()
		if err != nil {
			log.Fatalf("failed to run HTTP server: %v", err)
		}
	}()

	go func() {
		defer wg.Done()

		err := a.runSwaggerServer()
		if err != nil {
			log.Fatalf("failed to run Swagger server: %v", err)
		}
	}()

	wg.Wait()

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initLogger,
		a.initGRPCServer,
		a.initHTTPServer,
		a.initSwaggerServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(ctx context.Context) error {
	err := config.Load(ctx.Value(config.ConfigPathKey).(string))
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initLogger(_ context.Context) error {
	a.logger = a.serviceProvider.Logger()
	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {

	crds, err := credentials.NewServerTLSFromFile("service.pem", "service.key")

	if err != nil {
		return errors.New(fmt.Sprintf("Failed to generate credentials %v", err))
	}

	a.grpcServer = grpc.NewServer(
		grpc.Creds(crds),
		grpc.UnaryInterceptor(interceptor.ValidateInterceptor()),
	)

	reflection.Register(a.grpcServer)

	user_v1.RegisterUserV1Server(a.grpcServer, a.serviceProvider.UserApi(ctx))
	auth_v1.RegisterAuthV1Server(a.grpcServer, a.serviceProvider.AuthApi(ctx))
	access_v1.RegisterAccessV1Server(a.grpcServer, a.serviceProvider.AccessApi(ctx))

	return nil
}

func (a *App) initHTTPServer(ctx context.Context) error {
	mux := runtime.NewServeMux()

	crds, err := credentials.NewClientTLSFromFile("service.pem", "localhost")

	if err != nil {
		return errors.New(fmt.Sprintf("Failed to generate credentials %v", err))
	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(crds),
	}

	err = user_v1.RegisterUserV1HandlerFromEndpoint(ctx, mux, a.serviceProvider.GRPCConfig().Address(), opts)
	if err != nil {
		return err
	}

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		//AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		//AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Authorization"},
		//AllowCredentials: true,
	})

	a.httpServer = &http.Server{
		Addr:              a.serviceProvider.HTTPConfig().Address(),
		Handler:           corsMiddleware.Handler(mux),
		ReadHeaderTimeout: time.Duration(10) * time.Second,
	}

	return nil
}

func (a *App) runGRPCServer() error {
	a.logger.Info(fmt.Sprintf("GRPC server is running on %s", a.serviceProvider.GRPCConfig().Address()))

	list, err := net.Listen("tcp", a.serviceProvider.GRPCConfig().Address())
	if err != nil {
		return err
	}

	err = a.grpcServer.Serve(list)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) runHTTPServer() error {
	a.logger.Info(fmt.Sprintf("HTTP server is running on %s", a.serviceProvider.HTTPConfig().Address()))

	err := a.httpServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initSwaggerServer(_ context.Context) error {
	statikFs, err := fs.New()
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.StripPrefix("/", http.FileServer(statikFs)))
	mux.HandleFunc("/api.swagger.json", serveSwaggerFile("/api.swagger.json"))

	a.swaggerServer = &http.Server{
		Addr:              a.serviceProvider.SwaggerConfig().Address(),
		Handler:           mux,
		ReadHeaderTimeout: time.Duration(10) * time.Second,
	}

	return nil
}

func (a *App) runSwaggerServer() error {
	a.logger.Info(fmt.Sprintf("Swagger server is running on %s", a.serviceProvider.SwaggerConfig().Address()))

	err := a.swaggerServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func serveSwaggerFile(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Serving swagger file: %s\n", path)

		statikFs, err := fs.New()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Open swagger file: %s", path)

		file, err := statikFs.Open(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer func(file http.File) {
			err := file.Close()
			if err != nil {
				log.Printf("Error closing swagger file: %v", err)
			}
		}(file)

		log.Printf("Read swagger file: %s", path)

		content, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Write swagger file: %s", path)

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Served swagger file: %s", path)
	}
}
