package interceptor

import (
	"context"
	"fmt"

	"github.com/bufbuild/protovalidate-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type validator interface {
	ProtoReflect() protoreflect.Message
}

func ValidateInterceptor() grpc.UnaryServerInterceptor {

	vd, err := protovalidate.New()
	if err != nil {
		panic(fmt.Sprintf("failed to initialize validator: %v", err))
	}

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if msg, ok := req.(validator); ok {
			if err := vd.Validate(msg); err != nil {
				return nil, status.Errorf(codes.InvalidArgument, "validation failed: %v", err)
			}
		}

		return handler(ctx, req)
	}
}
