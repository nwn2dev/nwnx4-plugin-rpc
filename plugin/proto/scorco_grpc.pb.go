// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.22.2
// source: scorco.proto

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

// SCORCOServiceClient is the client API for SCORCOService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SCORCOServiceClient interface {
	SCORCOGetGFFSize(ctx context.Context, in *SCORCOGetGFFSizeRequest, opts ...grpc.CallOption) (*SCORCOGetGFFSizeResponse, error)
	SCORCOGetGFF(ctx context.Context, in *SCORCOGetGFFRequest, opts ...grpc.CallOption) (*SCORCOGetGFFResponse, error)
	SCORCOSetGFF(ctx context.Context, in *SCORCOSetGFFRequest, opts ...grpc.CallOption) (*SCORCOSetGFFResponse, error)
}

type sCORCOServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSCORCOServiceClient(cc grpc.ClientConnInterface) SCORCOServiceClient {
	return &sCORCOServiceClient{cc}
}

func (c *sCORCOServiceClient) SCORCOGetGFFSize(ctx context.Context, in *SCORCOGetGFFSizeRequest, opts ...grpc.CallOption) (*SCORCOGetGFFSizeResponse, error) {
	out := new(SCORCOGetGFFSizeResponse)
	err := c.cc.Invoke(ctx, "/NWNX4.RPC.SCORCOService/SCORCOGetGFFSize", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sCORCOServiceClient) SCORCOGetGFF(ctx context.Context, in *SCORCOGetGFFRequest, opts ...grpc.CallOption) (*SCORCOGetGFFResponse, error) {
	out := new(SCORCOGetGFFResponse)
	err := c.cc.Invoke(ctx, "/NWNX4.RPC.SCORCOService/SCORCOGetGFF", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sCORCOServiceClient) SCORCOSetGFF(ctx context.Context, in *SCORCOSetGFFRequest, opts ...grpc.CallOption) (*SCORCOSetGFFResponse, error) {
	out := new(SCORCOSetGFFResponse)
	err := c.cc.Invoke(ctx, "/NWNX4.RPC.SCORCOService/SCORCOSetGFF", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SCORCOServiceServer is the server API for SCORCOService service.
// All implementations must embed UnimplementedSCORCOServiceServer
// for forward compatibility
type SCORCOServiceServer interface {
	SCORCOGetGFFSize(context.Context, *SCORCOGetGFFSizeRequest) (*SCORCOGetGFFSizeResponse, error)
	SCORCOGetGFF(context.Context, *SCORCOGetGFFRequest) (*SCORCOGetGFFResponse, error)
	SCORCOSetGFF(context.Context, *SCORCOSetGFFRequest) (*SCORCOSetGFFResponse, error)
	mustEmbedUnimplementedSCORCOServiceServer()
}

// UnimplementedSCORCOServiceServer must be embedded to have forward compatible implementations.
type UnimplementedSCORCOServiceServer struct {
}

func (UnimplementedSCORCOServiceServer) SCORCOGetGFFSize(context.Context, *SCORCOGetGFFSizeRequest) (*SCORCOGetGFFSizeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SCORCOGetGFFSize not implemented")
}
func (UnimplementedSCORCOServiceServer) SCORCOGetGFF(context.Context, *SCORCOGetGFFRequest) (*SCORCOGetGFFResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SCORCOGetGFF not implemented")
}
func (UnimplementedSCORCOServiceServer) SCORCOSetGFF(context.Context, *SCORCOSetGFFRequest) (*SCORCOSetGFFResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SCORCOSetGFF not implemented")
}
func (UnimplementedSCORCOServiceServer) mustEmbedUnimplementedSCORCOServiceServer() {}

// UnsafeSCORCOServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SCORCOServiceServer will
// result in compilation errors.
type UnsafeSCORCOServiceServer interface {
	mustEmbedUnimplementedSCORCOServiceServer()
}

func RegisterSCORCOServiceServer(s grpc.ServiceRegistrar, srv SCORCOServiceServer) {
	s.RegisterService(&SCORCOService_ServiceDesc, srv)
}

func _SCORCOService_SCORCOGetGFFSize_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SCORCOGetGFFSizeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SCORCOServiceServer).SCORCOGetGFFSize(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/NWNX4.RPC.SCORCOService/SCORCOGetGFFSize",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SCORCOServiceServer).SCORCOGetGFFSize(ctx, req.(*SCORCOGetGFFSizeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SCORCOService_SCORCOGetGFF_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SCORCOGetGFFRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SCORCOServiceServer).SCORCOGetGFF(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/NWNX4.RPC.SCORCOService/SCORCOGetGFF",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SCORCOServiceServer).SCORCOGetGFF(ctx, req.(*SCORCOGetGFFRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SCORCOService_SCORCOSetGFF_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SCORCOSetGFFRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SCORCOServiceServer).SCORCOSetGFF(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/NWNX4.RPC.SCORCOService/SCORCOSetGFF",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SCORCOServiceServer).SCORCOSetGFF(ctx, req.(*SCORCOSetGFFRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// SCORCOService_ServiceDesc is the grpc.ServiceDesc for SCORCOService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SCORCOService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "NWNX4.RPC.SCORCOService",
	HandlerType: (*SCORCOServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SCORCOGetGFFSize",
			Handler:    _SCORCOService_SCORCOGetGFFSize_Handler,
		},
		{
			MethodName: "SCORCOGetGFF",
			Handler:    _SCORCOService_SCORCOGetGFF_Handler,
		},
		{
			MethodName: "SCORCOSetGFF",
			Handler:    _SCORCOService_SCORCOSetGFF_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "scorco.proto",
}
