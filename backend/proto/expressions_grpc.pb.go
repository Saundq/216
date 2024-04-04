// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v5.26.1
// source: proto/expressions.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// ExpressionClient is the client API for Expression service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ExpressionClient interface {
	Do(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error)
}

type expressionClient struct {
	cc grpc.ClientConnInterface
}

func NewExpressionClient(cc grpc.ClientConnInterface) ExpressionClient {
	return &expressionClient{cc}
}

func (c *expressionClient) Do(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/expressions.Expression/Do", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ExpressionServer is the server API for Expression service.
// All implementations must embed UnimplementedExpressionServer
// for forward compatibility
type ExpressionServer interface {
	Do(context.Context, *Request) (*Response, error)
	mustEmbedUnimplementedExpressionServer()
}

// UnimplementedExpressionServer must be embedded to have forward compatible implementations.
type UnimplementedExpressionServer struct {
}

func (UnimplementedExpressionServer) Do(context.Context, *Request) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Do not implemented")
}
func (UnimplementedExpressionServer) mustEmbedUnimplementedExpressionServer() {}

// UnsafeExpressionServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ExpressionServer will
// result in compilation errors.
type UnsafeExpressionServer interface {
	mustEmbedUnimplementedExpressionServer()
}

func RegisterExpressionServer(s grpc.ServiceRegistrar, srv ExpressionServer) {
	s.RegisterService(&Expression_ServiceDesc, srv)
}

func _Expression_Do_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExpressionServer).Do(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/expressions.Expression/Do",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExpressionServer).Do(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

// Expression_ServiceDesc is the grpc.ServiceDesc for Expression service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Expression_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "expressions.Expression",
	HandlerType: (*ExpressionServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Do",
			Handler:    _Expression_Do_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/expressions.proto",
}