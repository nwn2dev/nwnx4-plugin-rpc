// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.22.2
// source: ex.proto

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

// ExServiceClient is the client API for ExService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ExServiceClient interface {
	ExBuildGeneric(ctx context.Context, in *ExBuildGenericRequest, opts ...grpc.CallOption) (*ExBuildGenericResponse, error)
}

type exServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewExServiceClient(cc grpc.ClientConnInterface) ExServiceClient {
	return &exServiceClient{cc}
}

func (c *exServiceClient) ExBuildGeneric(ctx context.Context, in *ExBuildGenericRequest, opts ...grpc.CallOption) (*ExBuildGenericResponse, error) {
	out := new(ExBuildGenericResponse)
	err := c.cc.Invoke(ctx, "/NWNX4.RPC.ExService/ExBuildGeneric", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ExServiceServer is the server API for ExService service.
// All implementations must embed UnimplementedExServiceServer
// for forward compatibility
type ExServiceServer interface {
	ExBuildGeneric(context.Context, *ExBuildGenericRequest) (*ExBuildGenericResponse, error)
	mustEmbedUnimplementedExServiceServer()
}

// UnimplementedExServiceServer must be embedded to have forward compatible implementations.
type UnimplementedExServiceServer struct {
}

func (UnimplementedExServiceServer) ExBuildGeneric(context.Context, *ExBuildGenericRequest) (*ExBuildGenericResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExBuildGeneric not implemented")
}
func (UnimplementedExServiceServer) mustEmbedUnimplementedExServiceServer() {}

// UnsafeExServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ExServiceServer will
// result in compilation errors.
type UnsafeExServiceServer interface {
	mustEmbedUnimplementedExServiceServer()
}

func RegisterExServiceServer(s grpc.ServiceRegistrar, srv ExServiceServer) {
	s.RegisterService(&ExService_ServiceDesc, srv)
}

func _ExService_ExBuildGeneric_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExBuildGenericRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExServiceServer).ExBuildGeneric(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/NWNX4.RPC.ExService/ExBuildGeneric",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExServiceServer).ExBuildGeneric(ctx, req.(*ExBuildGenericRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ExService_ServiceDesc is the grpc.ServiceDesc for ExService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ExService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "NWNX4.RPC.ExService",
	HandlerType: (*ExServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ExBuildGeneric",
			Handler:    _ExService_ExBuildGeneric_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "ex.proto",
}