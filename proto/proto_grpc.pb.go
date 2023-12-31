// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.24.3
// source: proto/proto.proto

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

const (
	RAService_RA_FullMethodName = "/proto.RAService/RA"
)

// RAServiceClient is the client API for RAService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RAServiceClient interface {
	RA(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Reply, error)
}

type rAServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewRAServiceClient(cc grpc.ClientConnInterface) RAServiceClient {
	return &rAServiceClient{cc}
}

func (c *rAServiceClient) RA(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Reply, error) {
	out := new(Reply)
	err := c.cc.Invoke(ctx, RAService_RA_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RAServiceServer is the server API for RAService service.
// All implementations must embed UnimplementedRAServiceServer
// for forward compatibility
type RAServiceServer interface {
	RA(context.Context, *Request) (*Reply, error)
	mustEmbedUnimplementedRAServiceServer()
}

// UnimplementedRAServiceServer must be embedded to have forward compatible implementations.
type UnimplementedRAServiceServer struct {
}

func (UnimplementedRAServiceServer) RA(context.Context, *Request) (*Reply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RA not implemented")
}
func (UnimplementedRAServiceServer) mustEmbedUnimplementedRAServiceServer() {}

// UnsafeRAServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RAServiceServer will
// result in compilation errors.
type UnsafeRAServiceServer interface {
	mustEmbedUnimplementedRAServiceServer()
}

func RegisterRAServiceServer(s grpc.ServiceRegistrar, srv RAServiceServer) {
	s.RegisterService(&RAService_ServiceDesc, srv)
}

func _RAService_RA_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RAServiceServer).RA(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RAService_RA_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RAServiceServer).RA(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

// RAService_ServiceDesc is the grpc.ServiceDesc for RAService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RAService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.RAService",
	HandlerType: (*RAServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RA",
			Handler:    _RAService_RA_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/proto.proto",
}
