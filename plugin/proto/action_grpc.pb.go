// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.22.2
// source: action.proto

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

// ActionServiceClient is the client API for ActionService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ActionServiceClient interface {
	CallAction(ctx context.Context, in *CallActionRequest, opts ...grpc.CallOption) (*CallActionResponse, error)
}

type actionServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewActionServiceClient(cc grpc.ClientConnInterface) ActionServiceClient {
	return &actionServiceClient{cc}
}

func (c *actionServiceClient) CallAction(ctx context.Context, in *CallActionRequest, opts ...grpc.CallOption) (*CallActionResponse, error) {
	out := new(CallActionResponse)
	err := c.cc.Invoke(ctx, "/NWNX4.RPC.ActionService/CallAction", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ActionServiceServer is the server API for ActionService service.
// All implementations must embed UnimplementedActionServiceServer
// for forward compatibility
type ActionServiceServer interface {
	CallAction(context.Context, *CallActionRequest) (*CallActionResponse, error)
	mustEmbedUnimplementedActionServiceServer()
}

// UnimplementedActionServiceServer must be embedded to have forward compatible implementations.
type UnimplementedActionServiceServer struct {
}

func (UnimplementedActionServiceServer) CallAction(context.Context, *CallActionRequest) (*CallActionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CallAction not implemented")
}
func (UnimplementedActionServiceServer) mustEmbedUnimplementedActionServiceServer() {}

// UnsafeActionServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ActionServiceServer will
// result in compilation errors.
type UnsafeActionServiceServer interface {
	mustEmbedUnimplementedActionServiceServer()
}

func RegisterActionServiceServer(s grpc.ServiceRegistrar, srv ActionServiceServer) {
	s.RegisterService(&ActionService_ServiceDesc, srv)
}

func _ActionService_CallAction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CallActionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ActionServiceServer).CallAction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/NWNX4.RPC.ActionService/CallAction",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ActionServiceServer).CallAction(ctx, req.(*CallActionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ActionService_ServiceDesc is the grpc.ServiceDesc for ActionService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ActionService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "NWNX4.RPC.ActionService",
	HandlerType: (*ActionServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CallAction",
			Handler:    _ActionService_CallAction_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "action.proto",
}