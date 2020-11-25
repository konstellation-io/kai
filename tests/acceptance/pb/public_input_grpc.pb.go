// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// EntrypointClient is the client API for Entrypoint service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EntrypointClient interface {
	WorkflowTest(ctx context.Context, in *Message, opts ...grpc.CallOption) (*Message, error)
}

type entrypointClient struct {
	cc grpc.ClientConnInterface
}

func NewEntrypointClient(cc grpc.ClientConnInterface) EntrypointClient {
	return &entrypointClient{cc}
}

func (c *entrypointClient) WorkflowTest(ctx context.Context, in *Message, opts ...grpc.CallOption) (*Message, error) {
	out := new(Message)
	err := c.cc.Invoke(ctx, "/entrypoint.Entrypoint/WorkflowTest", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EntrypointServer is the server API for Entrypoint service.
// All implementations must embed UnimplementedEntrypointServer
// for forward compatibility
type EntrypointServer interface {
	WorkflowTest(context.Context, *Message) (*Message, error)
	mustEmbedUnimplementedEntrypointServer()
}

// UnimplementedEntrypointServer must be embedded to have forward compatible implementations.
type UnimplementedEntrypointServer struct {
}

func (UnimplementedEntrypointServer) WorkflowTest(context.Context, *Message) (*Message, error) {
	return nil, status.Errorf(codes.Unimplemented, "method WorkflowTest not implemented")
}
func (UnimplementedEntrypointServer) mustEmbedUnimplementedEntrypointServer() {}

// UnsafeEntrypointServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EntrypointServer will
// result in compilation errors.
type UnsafeEntrypointServer interface {
	mustEmbedUnimplementedEntrypointServer()
}

func RegisterEntrypointServer(s grpc.ServiceRegistrar, srv EntrypointServer) {
	s.RegisterService(&_Entrypoint_serviceDesc, srv)
}

func _Entrypoint_WorkflowTest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Message)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EntrypointServer).WorkflowTest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/entrypoint.Entrypoint/WorkflowTest",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EntrypointServer).WorkflowTest(ctx, req.(*Message))
	}
	return interceptor(ctx, in, info, handler)
}

var _Entrypoint_serviceDesc = grpc.ServiceDesc{
	ServiceName: "entrypoint.Entrypoint",
	HandlerType: (*EntrypointServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "WorkflowTest",
			Handler:    _Entrypoint_WorkflowTest_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "public_input.proto",
}
