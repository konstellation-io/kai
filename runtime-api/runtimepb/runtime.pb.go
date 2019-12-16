// Code generated by protoc-gen-go. DO NOT EDIT.
// source: runtimepb/runtime.proto

package runtimepb

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type RuntimeVersion struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RuntimeVersion) Reset()         { *m = RuntimeVersion{} }
func (m *RuntimeVersion) String() string { return proto.CompactTextString(m) }
func (*RuntimeVersion) ProtoMessage()    {}
func (*RuntimeVersion) Descriptor() ([]byte, []int) {
	return fileDescriptor_d0e5095094a8d27f, []int{0}
}

func (m *RuntimeVersion) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RuntimeVersion.Unmarshal(m, b)
}
func (m *RuntimeVersion) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RuntimeVersion.Marshal(b, m, deterministic)
}
func (m *RuntimeVersion) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RuntimeVersion.Merge(m, src)
}
func (m *RuntimeVersion) XXX_Size() int {
	return xxx_messageInfo_RuntimeVersion.Size(m)
}
func (m *RuntimeVersion) XXX_DiscardUnknown() {
	xxx_messageInfo_RuntimeVersion.DiscardUnknown(m)
}

var xxx_messageInfo_RuntimeVersion proto.InternalMessageInfo

func (m *RuntimeVersion) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type CreateRuntimeVersionRequest struct {
	RuntimeVersion       *RuntimeVersion `protobuf:"bytes,1,opt,name=runtimeVersion,proto3" json:"runtimeVersion,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *CreateRuntimeVersionRequest) Reset()         { *m = CreateRuntimeVersionRequest{} }
func (m *CreateRuntimeVersionRequest) String() string { return proto.CompactTextString(m) }
func (*CreateRuntimeVersionRequest) ProtoMessage()    {}
func (*CreateRuntimeVersionRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_d0e5095094a8d27f, []int{1}
}

func (m *CreateRuntimeVersionRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateRuntimeVersionRequest.Unmarshal(m, b)
}
func (m *CreateRuntimeVersionRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateRuntimeVersionRequest.Marshal(b, m, deterministic)
}
func (m *CreateRuntimeVersionRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateRuntimeVersionRequest.Merge(m, src)
}
func (m *CreateRuntimeVersionRequest) XXX_Size() int {
	return xxx_messageInfo_CreateRuntimeVersionRequest.Size(m)
}
func (m *CreateRuntimeVersionRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateRuntimeVersionRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CreateRuntimeVersionRequest proto.InternalMessageInfo

func (m *CreateRuntimeVersionRequest) GetRuntimeVersion() *RuntimeVersion {
	if m != nil {
		return m.RuntimeVersion
	}
	return nil
}

type CreateRuntimeVersionResponse struct {
	Success              bool     `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Message              string   `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreateRuntimeVersionResponse) Reset()         { *m = CreateRuntimeVersionResponse{} }
func (m *CreateRuntimeVersionResponse) String() string { return proto.CompactTextString(m) }
func (*CreateRuntimeVersionResponse) ProtoMessage()    {}
func (*CreateRuntimeVersionResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_d0e5095094a8d27f, []int{2}
}

func (m *CreateRuntimeVersionResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateRuntimeVersionResponse.Unmarshal(m, b)
}
func (m *CreateRuntimeVersionResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateRuntimeVersionResponse.Marshal(b, m, deterministic)
}
func (m *CreateRuntimeVersionResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateRuntimeVersionResponse.Merge(m, src)
}
func (m *CreateRuntimeVersionResponse) XXX_Size() int {
	return xxx_messageInfo_CreateRuntimeVersionResponse.Size(m)
}
func (m *CreateRuntimeVersionResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateRuntimeVersionResponse.DiscardUnknown(m)
}

var xxx_messageInfo_CreateRuntimeVersionResponse proto.InternalMessageInfo

func (m *CreateRuntimeVersionResponse) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

func (m *CreateRuntimeVersionResponse) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

type CheckRuntimeVersionIsCreatedRequest struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CheckRuntimeVersionIsCreatedRequest) Reset()         { *m = CheckRuntimeVersionIsCreatedRequest{} }
func (m *CheckRuntimeVersionIsCreatedRequest) String() string { return proto.CompactTextString(m) }
func (*CheckRuntimeVersionIsCreatedRequest) ProtoMessage()    {}
func (*CheckRuntimeVersionIsCreatedRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_d0e5095094a8d27f, []int{3}
}

func (m *CheckRuntimeVersionIsCreatedRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CheckRuntimeVersionIsCreatedRequest.Unmarshal(m, b)
}
func (m *CheckRuntimeVersionIsCreatedRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CheckRuntimeVersionIsCreatedRequest.Marshal(b, m, deterministic)
}
func (m *CheckRuntimeVersionIsCreatedRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CheckRuntimeVersionIsCreatedRequest.Merge(m, src)
}
func (m *CheckRuntimeVersionIsCreatedRequest) XXX_Size() int {
	return xxx_messageInfo_CheckRuntimeVersionIsCreatedRequest.Size(m)
}
func (m *CheckRuntimeVersionIsCreatedRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CheckRuntimeVersionIsCreatedRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CheckRuntimeVersionIsCreatedRequest proto.InternalMessageInfo

func (m *CheckRuntimeVersionIsCreatedRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type CheckRuntimeVersionIsCreatedResponse struct {
	Success              bool     `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Message              string   `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CheckRuntimeVersionIsCreatedResponse) Reset()         { *m = CheckRuntimeVersionIsCreatedResponse{} }
func (m *CheckRuntimeVersionIsCreatedResponse) String() string { return proto.CompactTextString(m) }
func (*CheckRuntimeVersionIsCreatedResponse) ProtoMessage()    {}
func (*CheckRuntimeVersionIsCreatedResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_d0e5095094a8d27f, []int{4}
}

func (m *CheckRuntimeVersionIsCreatedResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CheckRuntimeVersionIsCreatedResponse.Unmarshal(m, b)
}
func (m *CheckRuntimeVersionIsCreatedResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CheckRuntimeVersionIsCreatedResponse.Marshal(b, m, deterministic)
}
func (m *CheckRuntimeVersionIsCreatedResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CheckRuntimeVersionIsCreatedResponse.Merge(m, src)
}
func (m *CheckRuntimeVersionIsCreatedResponse) XXX_Size() int {
	return xxx_messageInfo_CheckRuntimeVersionIsCreatedResponse.Size(m)
}
func (m *CheckRuntimeVersionIsCreatedResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_CheckRuntimeVersionIsCreatedResponse.DiscardUnknown(m)
}

var xxx_messageInfo_CheckRuntimeVersionIsCreatedResponse proto.InternalMessageInfo

func (m *CheckRuntimeVersionIsCreatedResponse) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

func (m *CheckRuntimeVersionIsCreatedResponse) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func init() {
	proto.RegisterType((*RuntimeVersion)(nil), "runtime.RuntimeVersion")
	proto.RegisterType((*CreateRuntimeVersionRequest)(nil), "runtime.CreateRuntimeVersionRequest")
	proto.RegisterType((*CreateRuntimeVersionResponse)(nil), "runtime.CreateRuntimeVersionResponse")
	proto.RegisterType((*CheckRuntimeVersionIsCreatedRequest)(nil), "runtime.CheckRuntimeVersionIsCreatedRequest")
	proto.RegisterType((*CheckRuntimeVersionIsCreatedResponse)(nil), "runtime.CheckRuntimeVersionIsCreatedResponse")
}

func init() { proto.RegisterFile("runtimepb/runtime.proto", fileDescriptor_d0e5095094a8d27f) }

var fileDescriptor_d0e5095094a8d27f = []byte{
	// 258 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x2f, 0x2a, 0xcd, 0x2b,
	0xc9, 0xcc, 0x4d, 0x2d, 0x48, 0xd2, 0x87, 0xb2, 0xf4, 0x0a, 0x8a, 0xf2, 0x4b, 0xf2, 0x85, 0xd8,
	0xa1, 0x5c, 0x25, 0x15, 0x2e, 0xbe, 0x20, 0x08, 0x33, 0x2c, 0xb5, 0xa8, 0x38, 0x33, 0x3f, 0x4f,
	0x48, 0x88, 0x8b, 0x25, 0x2f, 0x31, 0x37, 0x55, 0x82, 0x51, 0x81, 0x51, 0x83, 0x33, 0x08, 0xcc,
	0x56, 0x8a, 0xe3, 0x92, 0x76, 0x2e, 0x4a, 0x4d, 0x2c, 0x49, 0x45, 0x55, 0x1b, 0x94, 0x5a, 0x58,
	0x9a, 0x5a, 0x5c, 0x22, 0x64, 0xcf, 0xc5, 0x57, 0x84, 0x22, 0x01, 0xd6, 0xcc, 0x6d, 0x24, 0xae,
	0x07, 0xb3, 0x15, 0x4d, 0x1f, 0x9a, 0x72, 0xa5, 0x20, 0x2e, 0x19, 0xec, 0xe6, 0x17, 0x17, 0xe4,
	0xe7, 0x15, 0xa7, 0x0a, 0x49, 0x70, 0xb1, 0x17, 0x97, 0x26, 0x27, 0xa7, 0x16, 0x17, 0x83, 0x4d,
	0xe6, 0x08, 0x82, 0x71, 0x41, 0x32, 0xb9, 0xa9, 0xc5, 0xc5, 0x89, 0xe9, 0xa9, 0x12, 0x4c, 0x60,
	0x07, 0xc3, 0xb8, 0x4a, 0x96, 0x5c, 0xca, 0xce, 0x19, 0xa9, 0xc9, 0xd9, 0xa8, 0x46, 0x7a, 0x16,
	0x43, 0x2c, 0x4a, 0x81, 0xb9, 0x1d, 0x9b, 0x77, 0xa3, 0xb8, 0x54, 0xf0, 0x6b, 0x25, 0xdf, 0x59,
	0x46, 0x7f, 0x19, 0xb9, 0x44, 0x51, 0xcd, 0x0d, 0x4e, 0x2d, 0x2a, 0xcb, 0x4c, 0x4e, 0x15, 0x4a,
	0xe5, 0x12, 0xc1, 0x16, 0x08, 0x42, 0x2a, 0xf0, 0x50, 0xc4, 0x13, 0x07, 0x52, 0xaa, 0x04, 0x54,
	0x41, 0x9c, 0xac, 0xc4, 0x20, 0x54, 0xcb, 0x25, 0x83, 0xcf, 0x73, 0x42, 0x3a, 0x08, 0x83, 0x08,
	0x07, 0x9f, 0x94, 0x2e, 0x91, 0xaa, 0x61, 0xd6, 0x3b, 0x71, 0x47, 0x71, 0xc2, 0x13, 0x65, 0x12,
	0x1b, 0x38, 0x35, 0x1a, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff, 0xc7, 0xcc, 0xc6, 0xc1, 0xa8, 0x02,
	0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// RuntimeVersionServiceClient is the client API for RuntimeVersionService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type RuntimeVersionServiceClient interface {
	CreateRuntimeVersion(ctx context.Context, in *CreateRuntimeVersionRequest, opts ...grpc.CallOption) (*CreateRuntimeVersionResponse, error)
	CheckRuntimeVersionIsCreated(ctx context.Context, in *CheckRuntimeVersionIsCreatedRequest, opts ...grpc.CallOption) (*CheckRuntimeVersionIsCreatedResponse, error)
}

type runtimeVersionServiceClient struct {
	cc *grpc.ClientConn
}

func NewRuntimeVersionServiceClient(cc *grpc.ClientConn) RuntimeVersionServiceClient {
	return &runtimeVersionServiceClient{cc}
}

func (c *runtimeVersionServiceClient) CreateRuntimeVersion(ctx context.Context, in *CreateRuntimeVersionRequest, opts ...grpc.CallOption) (*CreateRuntimeVersionResponse, error) {
	out := new(CreateRuntimeVersionResponse)
	err := c.cc.Invoke(ctx, "/runtime.RuntimeVersionService/CreateRuntimeVersion", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *runtimeVersionServiceClient) CheckRuntimeVersionIsCreated(ctx context.Context, in *CheckRuntimeVersionIsCreatedRequest, opts ...grpc.CallOption) (*CheckRuntimeVersionIsCreatedResponse, error) {
	out := new(CheckRuntimeVersionIsCreatedResponse)
	err := c.cc.Invoke(ctx, "/runtime.RuntimeVersionService/CheckRuntimeVersionIsCreated", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RuntimeVersionServiceServer is the server API for RuntimeVersionService service.
type RuntimeVersionServiceServer interface {
	CreateRuntimeVersion(context.Context, *CreateRuntimeVersionRequest) (*CreateRuntimeVersionResponse, error)
	CheckRuntimeVersionIsCreated(context.Context, *CheckRuntimeVersionIsCreatedRequest) (*CheckRuntimeVersionIsCreatedResponse, error)
}

// UnimplementedRuntimeVersionServiceServer can be embedded to have forward compatible implementations.
type UnimplementedRuntimeVersionServiceServer struct {
}

func (*UnimplementedRuntimeVersionServiceServer) CreateRuntimeVersion(ctx context.Context, req *CreateRuntimeVersionRequest) (*CreateRuntimeVersionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateRuntimeVersion not implemented")
}
func (*UnimplementedRuntimeVersionServiceServer) CheckRuntimeVersionIsCreated(ctx context.Context, req *CheckRuntimeVersionIsCreatedRequest) (*CheckRuntimeVersionIsCreatedResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckRuntimeVersionIsCreated not implemented")
}

func RegisterRuntimeVersionServiceServer(s *grpc.Server, srv RuntimeVersionServiceServer) {
	s.RegisterService(&_RuntimeVersionService_serviceDesc, srv)
}

func _RuntimeVersionService_CreateRuntimeVersion_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateRuntimeVersionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RuntimeVersionServiceServer).CreateRuntimeVersion(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/runtime.RuntimeVersionService/CreateRuntimeVersion",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RuntimeVersionServiceServer).CreateRuntimeVersion(ctx, req.(*CreateRuntimeVersionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RuntimeVersionService_CheckRuntimeVersionIsCreated_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckRuntimeVersionIsCreatedRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RuntimeVersionServiceServer).CheckRuntimeVersionIsCreated(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/runtime.RuntimeVersionService/CheckRuntimeVersionIsCreated",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RuntimeVersionServiceServer).CheckRuntimeVersionIsCreated(ctx, req.(*CheckRuntimeVersionIsCreatedRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _RuntimeVersionService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "runtime.RuntimeVersionService",
	HandlerType: (*RuntimeVersionServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateRuntimeVersion",
			Handler:    _RuntimeVersionService_CreateRuntimeVersion_Handler,
		},
		{
			MethodName: "CheckRuntimeVersionIsCreated",
			Handler:    _RuntimeVersionService_CheckRuntimeVersionIsCreated_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "runtimepb/runtime.proto",
}
