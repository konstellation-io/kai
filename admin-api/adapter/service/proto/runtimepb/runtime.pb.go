// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/runtimepb/runtime.proto

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

type Runtime struct {
	Name                 string             `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Namespace            string             `protobuf:"bytes,2,opt,name=namespace,proto3" json:"namespace,omitempty"`
	Minio                *Runtime_MinioConf `protobuf:"bytes,3,opt,name=minio,proto3" json:"minio,omitempty"`
	Mongo                *Runtime_MongoConf `protobuf:"bytes,4,opt,name=mongo,proto3" json:"mongo,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *Runtime) Reset()         { *m = Runtime{} }
func (m *Runtime) String() string { return proto.CompactTextString(m) }
func (*Runtime) ProtoMessage()    {}
func (*Runtime) Descriptor() ([]byte, []int) {
	return fileDescriptor_cc98cfd24676e4d3, []int{0}
}

func (m *Runtime) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Runtime.Unmarshal(m, b)
}
func (m *Runtime) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Runtime.Marshal(b, m, deterministic)
}
func (m *Runtime) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Runtime.Merge(m, src)
}
func (m *Runtime) XXX_Size() int {
	return xxx_messageInfo_Runtime.Size(m)
}
func (m *Runtime) XXX_DiscardUnknown() {
	xxx_messageInfo_Runtime.DiscardUnknown(m)
}

var xxx_messageInfo_Runtime proto.InternalMessageInfo

func (m *Runtime) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Runtime) GetNamespace() string {
	if m != nil {
		return m.Namespace
	}
	return ""
}

func (m *Runtime) GetMinio() *Runtime_MinioConf {
	if m != nil {
		return m.Minio
	}
	return nil
}

func (m *Runtime) GetMongo() *Runtime_MongoConf {
	if m != nil {
		return m.Mongo
	}
	return nil
}

type Runtime_MongoConf struct {
	Username             string   `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	Password             string   `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
	SharedKey            string   `protobuf:"bytes,3,opt,name=sharedKey,proto3" json:"sharedKey,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Runtime_MongoConf) Reset()         { *m = Runtime_MongoConf{} }
func (m *Runtime_MongoConf) String() string { return proto.CompactTextString(m) }
func (*Runtime_MongoConf) ProtoMessage()    {}
func (*Runtime_MongoConf) Descriptor() ([]byte, []int) {
	return fileDescriptor_cc98cfd24676e4d3, []int{0, 0}
}

func (m *Runtime_MongoConf) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Runtime_MongoConf.Unmarshal(m, b)
}
func (m *Runtime_MongoConf) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Runtime_MongoConf.Marshal(b, m, deterministic)
}
func (m *Runtime_MongoConf) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Runtime_MongoConf.Merge(m, src)
}
func (m *Runtime_MongoConf) XXX_Size() int {
	return xxx_messageInfo_Runtime_MongoConf.Size(m)
}
func (m *Runtime_MongoConf) XXX_DiscardUnknown() {
	xxx_messageInfo_Runtime_MongoConf.DiscardUnknown(m)
}

var xxx_messageInfo_Runtime_MongoConf proto.InternalMessageInfo

func (m *Runtime_MongoConf) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *Runtime_MongoConf) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

func (m *Runtime_MongoConf) GetSharedKey() string {
	if m != nil {
		return m.SharedKey
	}
	return ""
}

type Runtime_MinioConf struct {
	AccessKey            string   `protobuf:"bytes,1,opt,name=accessKey,proto3" json:"accessKey,omitempty"`
	SecretKey            string   `protobuf:"bytes,2,opt,name=secretKey,proto3" json:"secretKey,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Runtime_MinioConf) Reset()         { *m = Runtime_MinioConf{} }
func (m *Runtime_MinioConf) String() string { return proto.CompactTextString(m) }
func (*Runtime_MinioConf) ProtoMessage()    {}
func (*Runtime_MinioConf) Descriptor() ([]byte, []int) {
	return fileDescriptor_cc98cfd24676e4d3, []int{0, 1}
}

func (m *Runtime_MinioConf) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Runtime_MinioConf.Unmarshal(m, b)
}
func (m *Runtime_MinioConf) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Runtime_MinioConf.Marshal(b, m, deterministic)
}
func (m *Runtime_MinioConf) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Runtime_MinioConf.Merge(m, src)
}
func (m *Runtime_MinioConf) XXX_Size() int {
	return xxx_messageInfo_Runtime_MinioConf.Size(m)
}
func (m *Runtime_MinioConf) XXX_DiscardUnknown() {
	xxx_messageInfo_Runtime_MinioConf.DiscardUnknown(m)
}

var xxx_messageInfo_Runtime_MinioConf proto.InternalMessageInfo

func (m *Runtime_MinioConf) GetAccessKey() string {
	if m != nil {
		return m.AccessKey
	}
	return ""
}

func (m *Runtime_MinioConf) GetSecretKey() string {
	if m != nil {
		return m.SecretKey
	}
	return ""
}

type Request struct {
	Runtime              *Runtime `protobuf:"bytes,1,opt,name=runtime,proto3" json:"runtime,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Request) Reset()         { *m = Request{} }
func (m *Request) String() string { return proto.CompactTextString(m) }
func (*Request) ProtoMessage()    {}
func (*Request) Descriptor() ([]byte, []int) {
	return fileDescriptor_cc98cfd24676e4d3, []int{1}
}

func (m *Request) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Request.Unmarshal(m, b)
}
func (m *Request) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Request.Marshal(b, m, deterministic)
}
func (m *Request) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Request.Merge(m, src)
}
func (m *Request) XXX_Size() int {
	return xxx_messageInfo_Request.Size(m)
}
func (m *Request) XXX_DiscardUnknown() {
	xxx_messageInfo_Request.DiscardUnknown(m)
}

var xxx_messageInfo_Request proto.InternalMessageInfo

func (m *Request) GetRuntime() *Runtime {
	if m != nil {
		return m.Runtime
	}
	return nil
}

type Response struct {
	Message              string   `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Response) Reset()         { *m = Response{} }
func (m *Response) String() string { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()    {}
func (*Response) Descriptor() ([]byte, []int) {
	return fileDescriptor_cc98cfd24676e4d3, []int{2}
}

func (m *Response) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Response.Unmarshal(m, b)
}
func (m *Response) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Response.Marshal(b, m, deterministic)
}
func (m *Response) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Response.Merge(m, src)
}
func (m *Response) XXX_Size() int {
	return xxx_messageInfo_Response.Size(m)
}
func (m *Response) XXX_DiscardUnknown() {
	xxx_messageInfo_Response.DiscardUnknown(m)
}

var xxx_messageInfo_Response proto.InternalMessageInfo

func (m *Response) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

type RuntimeStatusResponse struct {
	Status               string   `protobuf:"bytes,1,opt,name=Status,proto3" json:"Status,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RuntimeStatusResponse) Reset()         { *m = RuntimeStatusResponse{} }
func (m *RuntimeStatusResponse) String() string { return proto.CompactTextString(m) }
func (*RuntimeStatusResponse) ProtoMessage()    {}
func (*RuntimeStatusResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_cc98cfd24676e4d3, []int{3}
}

func (m *RuntimeStatusResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RuntimeStatusResponse.Unmarshal(m, b)
}
func (m *RuntimeStatusResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RuntimeStatusResponse.Marshal(b, m, deterministic)
}
func (m *RuntimeStatusResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RuntimeStatusResponse.Merge(m, src)
}
func (m *RuntimeStatusResponse) XXX_Size() int {
	return xxx_messageInfo_RuntimeStatusResponse.Size(m)
}
func (m *RuntimeStatusResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_RuntimeStatusResponse.DiscardUnknown(m)
}

var xxx_messageInfo_RuntimeStatusResponse proto.InternalMessageInfo

func (m *RuntimeStatusResponse) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func init() {
	proto.RegisterType((*Runtime)(nil), "runtime.Runtime")
	proto.RegisterType((*Runtime_MongoConf)(nil), "runtime.Runtime.MongoConf")
	proto.RegisterType((*Runtime_MinioConf)(nil), "runtime.Runtime.MinioConf")
	proto.RegisterType((*Request)(nil), "runtime.Request")
	proto.RegisterType((*Response)(nil), "runtime.Response")
	proto.RegisterType((*RuntimeStatusResponse)(nil), "runtime.RuntimeStatusResponse")
}

func init() {
	proto.RegisterFile("proto/runtimepb/runtime.proto", fileDescriptor_cc98cfd24676e4d3)
}

var fileDescriptor_cc98cfd24676e4d3 = []byte{
	// 342 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x74, 0x52, 0x3d, 0x4e, 0xf3, 0x40,
	0x10, 0x4d, 0xf2, 0xe5, 0x4b, 0xe2, 0x89, 0x40, 0xb0, 0x12, 0xc8, 0xb2, 0x00, 0x21, 0x8b, 0x02,
	0x51, 0x24, 0x28, 0x88, 0x0b, 0x90, 0x82, 0x02, 0xd1, 0x98, 0x8e, 0x6e, 0xe3, 0x0c, 0xc1, 0x45,
	0x76, 0xcd, 0xee, 0x1a, 0xc4, 0x05, 0xb8, 0x1c, 0x97, 0x62, 0xff, 0xe3, 0x44, 0x50, 0x79, 0xe6,
	0xcd, 0xbc, 0x99, 0x37, 0xcf, 0x0b, 0xa7, 0xb5, 0xe0, 0x8a, 0x4f, 0x45, 0xc3, 0x54, 0xb5, 0xc6,
	0x7a, 0x11, 0xa2, 0x89, 0xc5, 0xc9, 0xd0, 0xa7, 0xf9, 0x77, 0x0f, 0x86, 0x85, 0x8b, 0x09, 0x81,
	0x3e, 0xa3, 0x6b, 0x4c, 0xbb, 0xe7, 0xdd, 0xcb, 0xa4, 0xb0, 0x31, 0x39, 0x81, 0xc4, 0x7c, 0x65,
	0x4d, 0x4b, 0x4c, 0x7b, 0xb6, 0xb0, 0x01, 0xc8, 0x35, 0xfc, 0x5f, 0x57, 0xac, 0xe2, 0xe9, 0x3f,
	0x5d, 0x19, 0xcf, 0xb2, 0x49, 0xd8, 0xe2, 0x47, 0x4e, 0x1e, 0x4d, 0x75, 0xce, 0xd9, 0x4b, 0xe1,
	0x1a, 0x2d, 0x83, 0xb3, 0x15, 0x4f, 0xfb, 0x7f, 0x31, 0x4c, 0xd5, 0x33, 0x4c, 0x98, 0x51, 0x48,
	0x22, 0x46, 0x32, 0x18, 0x35, 0x12, 0x45, 0x4b, 0x66, 0xcc, 0x4d, 0xad, 0xa6, 0x52, 0x7e, 0x70,
	0xb1, 0xf4, 0x4a, 0x63, 0x6e, 0xce, 0x90, 0xaf, 0x54, 0xe0, 0xf2, 0x01, 0x3f, 0xad, 0x58, 0x7d,
	0x46, 0x04, 0xb2, 0x7b, 0xbd, 0x22, 0x08, 0x35, 0xad, 0xb4, 0x2c, 0x51, 0x4a, 0xd3, 0xea, 0x76,
	0x6c, 0x00, 0x3b, 0x08, 0x4b, 0x81, 0xca, 0x54, 0xbd, 0x1f, 0x11, 0xc8, 0x6f, 0xb5, 0x99, 0xf8,
	0xd6, 0xa0, 0x54, 0xe4, 0x0a, 0x82, 0xc7, 0x76, 0xc8, 0x78, 0x76, 0xb0, 0x7b, 0x6a, 0x11, 0x7f,
	0xc2, 0x05, 0x8c, 0x0a, 0xed, 0x28, 0x67, 0x12, 0x49, 0x0a, 0x43, 0xed, 0xae, 0xa4, 0xab, 0x60,
	0x77, 0x48, 0xf3, 0x29, 0x1c, 0x79, 0xe6, 0x93, 0xa2, 0xaa, 0x91, 0x91, 0x72, 0x0c, 0x03, 0x87,
	0x78, 0xb9, 0x3e, 0x9b, 0x7d, 0x75, 0x61, 0x3f, 0x30, 0x50, 0xbc, 0x57, 0xfa, 0x87, 0x4d, 0x61,
	0x30, 0x17, 0x48, 0x15, 0x92, 0x96, 0x1c, 0xa7, 0x38, 0x3b, 0x6c, 0x21, 0x6e, 0x72, 0xde, 0x21,
	0x73, 0xd8, 0xdb, 0x5a, 0xfa, 0x0b, 0xef, 0x6c, 0xf7, 0xb0, 0x6d, 0x79, 0x79, 0xe7, 0x6e, 0xfc,
	0x9c, 0xc4, 0x87, 0xb8, 0x18, 0xd8, 0x17, 0x78, 0xf3, 0x13, 0x00, 0x00, 0xff, 0xff, 0xce, 0xe9,
	0xa0, 0x1e, 0xa2, 0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// RuntimeServiceClient is the client API for RuntimeService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type RuntimeServiceClient interface {
	Create(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error)
	RuntimeStatus(ctx context.Context, in *Request, opts ...grpc.CallOption) (*RuntimeStatusResponse, error)
}

type runtimeServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewRuntimeServiceClient(cc grpc.ClientConnInterface) RuntimeServiceClient {
	return &runtimeServiceClient{cc}
}

func (c *runtimeServiceClient) Create(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/runtime.RuntimeService/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *runtimeServiceClient) RuntimeStatus(ctx context.Context, in *Request, opts ...grpc.CallOption) (*RuntimeStatusResponse, error) {
	out := new(RuntimeStatusResponse)
	err := c.cc.Invoke(ctx, "/runtime.RuntimeService/RuntimeStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RuntimeServiceServer is the server API for RuntimeService service.
type RuntimeServiceServer interface {
	Create(context.Context, *Request) (*Response, error)
	RuntimeStatus(context.Context, *Request) (*RuntimeStatusResponse, error)
}

// UnimplementedRuntimeServiceServer can be embedded to have forward compatible implementations.
type UnimplementedRuntimeServiceServer struct {
}

func (*UnimplementedRuntimeServiceServer) Create(ctx context.Context, req *Request) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (*UnimplementedRuntimeServiceServer) RuntimeStatus(ctx context.Context, req *Request) (*RuntimeStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RuntimeStatus not implemented")
}

func RegisterRuntimeServiceServer(s *grpc.Server, srv RuntimeServiceServer) {
	s.RegisterService(&_RuntimeService_serviceDesc, srv)
}

func _RuntimeService_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RuntimeServiceServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/runtime.RuntimeService/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RuntimeServiceServer).Create(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _RuntimeService_RuntimeStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RuntimeServiceServer).RuntimeStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/runtime.RuntimeService/RuntimeStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RuntimeServiceServer).RuntimeStatus(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

var _RuntimeService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "runtime.RuntimeService",
	HandlerType: (*RuntimeServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _RuntimeService_Create_Handler,
		},
		{
			MethodName: "RuntimeStatus",
			Handler:    _RuntimeService_RuntimeStatus_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/runtimepb/runtime.proto",
}
