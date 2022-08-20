// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.4
// source: nwscript/nwnx.proto

package nwscript

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

// NWNXServiceClient is the client API for NWNXService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type NWNXServiceClient interface {
	NWNXGetInt(ctx context.Context, in *NWNXGetIntRequest, opts ...grpc.CallOption) (*Int, error)
	NWNXSetInt(ctx context.Context, in *NWNXSetIntRequest, opts ...grpc.CallOption) (*Void, error)
	NWNXGetFloat(ctx context.Context, in *NWNXGetFloatRequest, opts ...grpc.CallOption) (*Float, error)
	NWNXSetFloat(ctx context.Context, in *NWNXSetFloatRequest, opts ...grpc.CallOption) (*Void, error)
	NWNXGetString(ctx context.Context, in *NWNXGetStringRequest, opts ...grpc.CallOption) (*String, error)
	NWNXSetString(ctx context.Context, in *NWNXSetStringRequest, opts ...grpc.CallOption) (*Void, error)
}

type nWNXServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewNWNXServiceClient(cc grpc.ClientConnInterface) NWNXServiceClient {
	return &nWNXServiceClient{cc}
}

func (c *nWNXServiceClient) NWNXGetInt(ctx context.Context, in *NWNXGetIntRequest, opts ...grpc.CallOption) (*Int, error) {
	out := new(Int)
	err := c.cc.Invoke(ctx, "/NWNX4.RPC.Proto.NWScript.NWNXService/NWNXGetInt", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nWNXServiceClient) NWNXSetInt(ctx context.Context, in *NWNXSetIntRequest, opts ...grpc.CallOption) (*Void, error) {
	out := new(Void)
	err := c.cc.Invoke(ctx, "/NWNX4.RPC.Proto.NWScript.NWNXService/NWNXSetInt", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nWNXServiceClient) NWNXGetFloat(ctx context.Context, in *NWNXGetFloatRequest, opts ...grpc.CallOption) (*Float, error) {
	out := new(Float)
	err := c.cc.Invoke(ctx, "/NWNX4.RPC.Proto.NWScript.NWNXService/NWNXGetFloat", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nWNXServiceClient) NWNXSetFloat(ctx context.Context, in *NWNXSetFloatRequest, opts ...grpc.CallOption) (*Void, error) {
	out := new(Void)
	err := c.cc.Invoke(ctx, "/NWNX4.RPC.Proto.NWScript.NWNXService/NWNXSetFloat", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nWNXServiceClient) NWNXGetString(ctx context.Context, in *NWNXGetStringRequest, opts ...grpc.CallOption) (*String, error) {
	out := new(String)
	err := c.cc.Invoke(ctx, "/NWNX4.RPC.Proto.NWScript.NWNXService/NWNXGetString", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nWNXServiceClient) NWNXSetString(ctx context.Context, in *NWNXSetStringRequest, opts ...grpc.CallOption) (*Void, error) {
	out := new(Void)
	err := c.cc.Invoke(ctx, "/NWNX4.RPC.Proto.NWScript.NWNXService/NWNXSetString", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NWNXServiceServer is the server API for NWNXService service.
// All implementations must embed UnimplementedNWNXServiceServer
// for forward compatibility
type NWNXServiceServer interface {
	NWNXGetInt(context.Context, *NWNXGetIntRequest) (*Int, error)
	NWNXSetInt(context.Context, *NWNXSetIntRequest) (*Void, error)
	NWNXGetFloat(context.Context, *NWNXGetFloatRequest) (*Float, error)
	NWNXSetFloat(context.Context, *NWNXSetFloatRequest) (*Void, error)
	NWNXGetString(context.Context, *NWNXGetStringRequest) (*String, error)
	NWNXSetString(context.Context, *NWNXSetStringRequest) (*Void, error)
	mustEmbedUnimplementedNWNXServiceServer()
}

// UnimplementedNWNXServiceServer must be embedded to have forward compatible implementations.
type UnimplementedNWNXServiceServer struct {
}

func (UnimplementedNWNXServiceServer) NWNXGetInt(context.Context, *NWNXGetIntRequest) (*Int, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NWNXGetInt not implemented")
}
func (UnimplementedNWNXServiceServer) NWNXSetInt(context.Context, *NWNXSetIntRequest) (*Void, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NWNXSetInt not implemented")
}
func (UnimplementedNWNXServiceServer) NWNXGetFloat(context.Context, *NWNXGetFloatRequest) (*Float, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NWNXGetFloat not implemented")
}
func (UnimplementedNWNXServiceServer) NWNXSetFloat(context.Context, *NWNXSetFloatRequest) (*Void, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NWNXSetFloat not implemented")
}
func (UnimplementedNWNXServiceServer) NWNXGetString(context.Context, *NWNXGetStringRequest) (*String, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NWNXGetString not implemented")
}
func (UnimplementedNWNXServiceServer) NWNXSetString(context.Context, *NWNXSetStringRequest) (*Void, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NWNXSetString not implemented")
}
func (UnimplementedNWNXServiceServer) mustEmbedUnimplementedNWNXServiceServer() {}

// UnsafeNWNXServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to NWNXServiceServer will
// result in compilation errors.
type UnsafeNWNXServiceServer interface {
	mustEmbedUnimplementedNWNXServiceServer()
}

func RegisterNWNXServiceServer(s grpc.ServiceRegistrar, srv NWNXServiceServer) {
	s.RegisterService(&NWNXService_ServiceDesc, srv)
}

func _NWNXService_NWNXGetInt_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NWNXGetIntRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NWNXServiceServer).NWNXGetInt(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/NWNX4.RPC.Proto.NWScript.NWNXService/NWNXGetInt",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NWNXServiceServer).NWNXGetInt(ctx, req.(*NWNXGetIntRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NWNXService_NWNXSetInt_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NWNXSetIntRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NWNXServiceServer).NWNXSetInt(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/NWNX4.RPC.Proto.NWScript.NWNXService/NWNXSetInt",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NWNXServiceServer).NWNXSetInt(ctx, req.(*NWNXSetIntRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NWNXService_NWNXGetFloat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NWNXGetFloatRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NWNXServiceServer).NWNXGetFloat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/NWNX4.RPC.Proto.NWScript.NWNXService/NWNXGetFloat",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NWNXServiceServer).NWNXGetFloat(ctx, req.(*NWNXGetFloatRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NWNXService_NWNXSetFloat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NWNXSetFloatRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NWNXServiceServer).NWNXSetFloat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/NWNX4.RPC.Proto.NWScript.NWNXService/NWNXSetFloat",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NWNXServiceServer).NWNXSetFloat(ctx, req.(*NWNXSetFloatRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NWNXService_NWNXGetString_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NWNXGetStringRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NWNXServiceServer).NWNXGetString(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/NWNX4.RPC.Proto.NWScript.NWNXService/NWNXGetString",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NWNXServiceServer).NWNXGetString(ctx, req.(*NWNXGetStringRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NWNXService_NWNXSetString_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NWNXSetStringRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NWNXServiceServer).NWNXSetString(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/NWNX4.RPC.Proto.NWScript.NWNXService/NWNXSetString",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NWNXServiceServer).NWNXSetString(ctx, req.(*NWNXSetStringRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// NWNXService_ServiceDesc is the grpc.ServiceDesc for NWNXService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var NWNXService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "NWNX4.RPC.Proto.NWScript.NWNXService",
	HandlerType: (*NWNXServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "NWNXGetInt",
			Handler:    _NWNXService_NWNXGetInt_Handler,
		},
		{
			MethodName: "NWNXSetInt",
			Handler:    _NWNXService_NWNXSetInt_Handler,
		},
		{
			MethodName: "NWNXGetFloat",
			Handler:    _NWNXService_NWNXGetFloat_Handler,
		},
		{
			MethodName: "NWNXSetFloat",
			Handler:    _NWNXService_NWNXSetFloat_Handler,
		},
		{
			MethodName: "NWNXGetString",
			Handler:    _NWNXService_NWNXGetString_Handler,
		},
		{
			MethodName: "NWNXSetString",
			Handler:    _NWNXService_NWNXSetString_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "nwscript/nwnx.proto",
}
