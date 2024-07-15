// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v5.27.2
// source: proto/agent.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// CalculatorServiceClient is the client API for CalculatorService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CalculatorServiceClient interface {
	TaskGet(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*TaskGetResponse, error)
	TaskPost(ctx context.Context, in *TaskPostRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type calculatorServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCalculatorServiceClient(cc grpc.ClientConnInterface) CalculatorServiceClient {
	return &calculatorServiceClient{cc}
}

func (c *calculatorServiceClient) TaskGet(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*TaskGetResponse, error) {
	out := new(TaskGetResponse)
	err := c.cc.Invoke(ctx, "/calculator.CalculatorService/TaskGet", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *calculatorServiceClient) TaskPost(ctx context.Context, in *TaskPostRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/calculator.CalculatorService/TaskPost", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CalculatorServiceServer is the server API for CalculatorService service.
// All implementations must embed UnimplementedCalculatorServiceServer
// for forward compatibility
type CalculatorServiceServer interface {
	TaskGet(context.Context, *emptypb.Empty) (*TaskGetResponse, error)
	TaskPost(context.Context, *TaskPostRequest) (*emptypb.Empty, error)
	mustEmbedUnimplementedCalculatorServiceServer()
}

// UnimplementedCalculatorServiceServer must be embedded to have forward compatible implementations.
type UnimplementedCalculatorServiceServer struct {
}

func (UnimplementedCalculatorServiceServer) TaskGet(context.Context, *emptypb.Empty) (*TaskGetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TaskGet not implemented")
}
func (UnimplementedCalculatorServiceServer) TaskPost(context.Context, *TaskPostRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TaskPost not implemented")
}
func (UnimplementedCalculatorServiceServer) mustEmbedUnimplementedCalculatorServiceServer() {}

// UnsafeCalculatorServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CalculatorServiceServer will
// result in compilation errors.
type UnsafeCalculatorServiceServer interface {
	mustEmbedUnimplementedCalculatorServiceServer()
}

func RegisterCalculatorServiceServer(s grpc.ServiceRegistrar, srv CalculatorServiceServer) {
	s.RegisterService(&CalculatorService_ServiceDesc, srv)
}

func _CalculatorService_TaskGet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalculatorServiceServer).TaskGet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/calculator.CalculatorService/TaskGet",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalculatorServiceServer).TaskGet(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _CalculatorService_TaskPost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TaskPostRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CalculatorServiceServer).TaskPost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/calculator.CalculatorService/TaskPost",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CalculatorServiceServer).TaskPost(ctx, req.(*TaskPostRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// CalculatorService_ServiceDesc is the grpc.ServiceDesc for CalculatorService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CalculatorService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "calculator.CalculatorService",
	HandlerType: (*CalculatorServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "TaskGet",
			Handler:    _CalculatorService_TaskGet_Handler,
		},
		{
			MethodName: "TaskPost",
			Handler:    _CalculatorService_TaskPost_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/agent.proto",
}
