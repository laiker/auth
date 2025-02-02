package main

import (
	"flag"
	"log"

	"github.com/laiker/auth/internal/app"
	"github.com/laiker/auth/internal/config"
	"golang.org/x/net/context"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

func main() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, config.ConfigPathKey, configPath)

	a, err := app.NewApp(ctx)

	if err != nil {
		log.Fatal(err)
	}

	err = a.Run()

	if err != nil {
		log.Fatal(err)
	}
}
