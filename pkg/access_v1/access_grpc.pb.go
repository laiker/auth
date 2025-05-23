// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: access.proto

package access_v1

import (
	context "context"
	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// AccessV1Client is the client API for AccessV1 service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AccessV1Client interface {
	HasAccess(ctx context.Context, in *CheckRequest, opts ...grpc.CallOption) (*empty.Empty, error)
}

type accessV1Client struct {
	cc grpc.ClientConnInterface
}

func NewAccessV1Client(cc grpc.ClientConnInterface) AccessV1Client {
	return &accessV1Client{cc}
}

func (c *accessV1Client) HasAccess(ctx context.Context, in *CheckRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/access_v1.AccessV1/HasAccess", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AccessV1Server is the server API for AccessV1 service.
// All implementations must embed UnimplementedAccessV1Server
// for forward compatibility
type AccessV1Server interface {
	HasAccess(context.Context, *CheckRequest) (*empty.Empty, error)
	mustEmbedUnimplementedAccessV1Server()
}

// UnimplementedAccessV1Server must be embedded to have forward compatible implementations.
type UnimplementedAccessV1Server struct {
}

func (UnimplementedAccessV1Server) HasAccess(context.Context, *CheckRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HasAccess not implemented")
}
func (UnimplementedAccessV1Server) mustEmbedUnimplementedAccessV1Server() {}

// UnsafeAccessV1Server may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AccessV1Server will
// result in compilation errors.
type UnsafeAccessV1Server interface {
	mustEmbedUnimplementedAccessV1Server()
}

func RegisterAccessV1Server(s grpc.ServiceRegistrar, srv AccessV1Server) {
	s.RegisterService(&AccessV1_ServiceDesc, srv)
}

func _AccessV1_HasAccess_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccessV1Server).HasAccess(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/access_v1.AccessV1/HasAccess",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccessV1Server).HasAccess(ctx, req.(*CheckRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AccessV1_ServiceDesc is the grpc.ServiceDesc for AccessV1 service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AccessV1_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "access_v1.AccessV1",
	HandlerType: (*AccessV1Server)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "HasAccess",
			Handler:    _AccessV1_HasAccess_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "access.proto",
}
