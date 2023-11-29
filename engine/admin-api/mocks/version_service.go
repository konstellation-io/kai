// Code generated by MockGen. DO NOT EDIT.
// Source: ../proto/versionpb/version_grpc.pb.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	versionpb "github.com/konstellation-io/kai/engine/admin-api/adapter/service/proto/versionpb"
	grpc "google.golang.org/grpc"
	metadata "google.golang.org/grpc/metadata"
)

// MockVersionServiceClient is a mock of VersionServiceClient interface.
type MockVersionServiceClient struct {
	ctrl     *gomock.Controller
	recorder *MockVersionServiceClientMockRecorder
}

// MockVersionServiceClientMockRecorder is the mock recorder for MockVersionServiceClient.
type MockVersionServiceClientMockRecorder struct {
	mock *MockVersionServiceClient
}

// NewMockVersionServiceClient creates a new mock instance.
func NewMockVersionServiceClient(ctrl *gomock.Controller) *MockVersionServiceClient {
	mock := &MockVersionServiceClient{ctrl: ctrl}
	mock.recorder = &MockVersionServiceClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockVersionServiceClient) EXPECT() *MockVersionServiceClientMockRecorder {
	return m.recorder
}

// Publish mocks base method.
func (m *MockVersionServiceClient) Publish(ctx context.Context, in *versionpb.PublishRequest, opts ...grpc.CallOption) (*versionpb.PublishResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Publish", varargs...)
	ret0, _ := ret[0].(*versionpb.PublishResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Publish indicates an expected call of Publish.
func (mr *MockVersionServiceClientMockRecorder) Publish(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Publish", reflect.TypeOf((*MockVersionServiceClient)(nil).Publish), varargs...)
}

// RegisterProcess mocks base method.
func (m *MockVersionServiceClient) RegisterProcess(ctx context.Context, in *versionpb.RegisterProcessRequest, opts ...grpc.CallOption) (*versionpb.RegisterProcessResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "RegisterProcess", varargs...)
	ret0, _ := ret[0].(*versionpb.RegisterProcessResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RegisterProcess indicates an expected call of RegisterProcess.
func (mr *MockVersionServiceClientMockRecorder) RegisterProcess(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterProcess", reflect.TypeOf((*MockVersionServiceClient)(nil).RegisterProcess), varargs...)
}

// Start mocks base method.
func (m *MockVersionServiceClient) Start(ctx context.Context, in *versionpb.StartRequest, opts ...grpc.CallOption) (*versionpb.Response, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Start", varargs...)
	ret0, _ := ret[0].(*versionpb.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Start indicates an expected call of Start.
func (mr *MockVersionServiceClientMockRecorder) Start(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockVersionServiceClient)(nil).Start), varargs...)
}

// Stop mocks base method.
func (m *MockVersionServiceClient) Stop(ctx context.Context, in *versionpb.StopRequest, opts ...grpc.CallOption) (*versionpb.Response, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Stop", varargs...)
	ret0, _ := ret[0].(*versionpb.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Stop indicates an expected call of Stop.
func (mr *MockVersionServiceClientMockRecorder) Stop(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockVersionServiceClient)(nil).Stop), varargs...)
}

// Unpublish mocks base method.
func (m *MockVersionServiceClient) Unpublish(ctx context.Context, in *versionpb.UnpublishRequest, opts ...grpc.CallOption) (*versionpb.Response, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Unpublish", varargs...)
	ret0, _ := ret[0].(*versionpb.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Unpublish indicates an expected call of Unpublish.
func (mr *MockVersionServiceClientMockRecorder) Unpublish(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unpublish", reflect.TypeOf((*MockVersionServiceClient)(nil).Unpublish), varargs...)
}

// WatchProcessStatus mocks base method.
func (m *MockVersionServiceClient) WatchProcessStatus(ctx context.Context, in *versionpb.ProcessStatusRequest, opts ...grpc.CallOption) (versionpb.VersionService_WatchProcessStatusClient, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "WatchProcessStatus", varargs...)
	ret0, _ := ret[0].(versionpb.VersionService_WatchProcessStatusClient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WatchProcessStatus indicates an expected call of WatchProcessStatus.
func (mr *MockVersionServiceClientMockRecorder) WatchProcessStatus(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WatchProcessStatus", reflect.TypeOf((*MockVersionServiceClient)(nil).WatchProcessStatus), varargs...)
}

// MockVersionService_WatchProcessStatusClient is a mock of VersionService_WatchProcessStatusClient interface.
type MockVersionService_WatchProcessStatusClient struct {
	ctrl     *gomock.Controller
	recorder *MockVersionService_WatchProcessStatusClientMockRecorder
}

// MockVersionService_WatchProcessStatusClientMockRecorder is the mock recorder for MockVersionService_WatchProcessStatusClient.
type MockVersionService_WatchProcessStatusClientMockRecorder struct {
	mock *MockVersionService_WatchProcessStatusClient
}

// NewMockVersionService_WatchProcessStatusClient creates a new mock instance.
func NewMockVersionService_WatchProcessStatusClient(ctrl *gomock.Controller) *MockVersionService_WatchProcessStatusClient {
	mock := &MockVersionService_WatchProcessStatusClient{ctrl: ctrl}
	mock.recorder = &MockVersionService_WatchProcessStatusClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockVersionService_WatchProcessStatusClient) EXPECT() *MockVersionService_WatchProcessStatusClientMockRecorder {
	return m.recorder
}

// CloseSend mocks base method.
func (m *MockVersionService_WatchProcessStatusClient) CloseSend() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloseSend")
	ret0, _ := ret[0].(error)
	return ret0
}

// CloseSend indicates an expected call of CloseSend.
func (mr *MockVersionService_WatchProcessStatusClientMockRecorder) CloseSend() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseSend", reflect.TypeOf((*MockVersionService_WatchProcessStatusClient)(nil).CloseSend))
}

// Context mocks base method.
func (m *MockVersionService_WatchProcessStatusClient) Context() context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Context indicates an expected call of Context.
func (mr *MockVersionService_WatchProcessStatusClientMockRecorder) Context() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockVersionService_WatchProcessStatusClient)(nil).Context))
}

// Header mocks base method.
func (m *MockVersionService_WatchProcessStatusClient) Header() (metadata.MD, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Header")
	ret0, _ := ret[0].(metadata.MD)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Header indicates an expected call of Header.
func (mr *MockVersionService_WatchProcessStatusClientMockRecorder) Header() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Header", reflect.TypeOf((*MockVersionService_WatchProcessStatusClient)(nil).Header))
}

// Recv mocks base method.
func (m *MockVersionService_WatchProcessStatusClient) Recv() (*versionpb.ProcessStatusResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Recv")
	ret0, _ := ret[0].(*versionpb.ProcessStatusResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Recv indicates an expected call of Recv.
func (mr *MockVersionService_WatchProcessStatusClientMockRecorder) Recv() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Recv", reflect.TypeOf((*MockVersionService_WatchProcessStatusClient)(nil).Recv))
}

// RecvMsg mocks base method.
func (m_2 *MockVersionService_WatchProcessStatusClient) RecvMsg(m any) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "RecvMsg", m)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecvMsg indicates an expected call of RecvMsg.
func (mr *MockVersionService_WatchProcessStatusClientMockRecorder) RecvMsg(m interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecvMsg", reflect.TypeOf((*MockVersionService_WatchProcessStatusClient)(nil).RecvMsg), m)
}

// SendMsg mocks base method.
func (m_2 *MockVersionService_WatchProcessStatusClient) SendMsg(m any) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "SendMsg", m)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMsg indicates an expected call of SendMsg.
func (mr *MockVersionService_WatchProcessStatusClientMockRecorder) SendMsg(m interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMsg", reflect.TypeOf((*MockVersionService_WatchProcessStatusClient)(nil).SendMsg), m)
}

// Trailer mocks base method.
func (m *MockVersionService_WatchProcessStatusClient) Trailer() metadata.MD {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Trailer")
	ret0, _ := ret[0].(metadata.MD)
	return ret0
}

// Trailer indicates an expected call of Trailer.
func (mr *MockVersionService_WatchProcessStatusClientMockRecorder) Trailer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Trailer", reflect.TypeOf((*MockVersionService_WatchProcessStatusClient)(nil).Trailer))
}

// MockVersionServiceServer is a mock of VersionServiceServer interface.
type MockVersionServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockVersionServiceServerMockRecorder
}

// MockVersionServiceServerMockRecorder is the mock recorder for MockVersionServiceServer.
type MockVersionServiceServerMockRecorder struct {
	mock *MockVersionServiceServer
}

// NewMockVersionServiceServer creates a new mock instance.
func NewMockVersionServiceServer(ctrl *gomock.Controller) *MockVersionServiceServer {
	mock := &MockVersionServiceServer{ctrl: ctrl}
	mock.recorder = &MockVersionServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockVersionServiceServer) EXPECT() *MockVersionServiceServerMockRecorder {
	return m.recorder
}

// Publish mocks base method.
func (m *MockVersionServiceServer) Publish(arg0 context.Context, arg1 *versionpb.PublishRequest) (*versionpb.PublishResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Publish", arg0, arg1)
	ret0, _ := ret[0].(*versionpb.PublishResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Publish indicates an expected call of Publish.
func (mr *MockVersionServiceServerMockRecorder) Publish(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Publish", reflect.TypeOf((*MockVersionServiceServer)(nil).Publish), arg0, arg1)
}

// RegisterProcess mocks base method.
func (m *MockVersionServiceServer) RegisterProcess(arg0 context.Context, arg1 *versionpb.RegisterProcessRequest) (*versionpb.RegisterProcessResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterProcess", arg0, arg1)
	ret0, _ := ret[0].(*versionpb.RegisterProcessResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RegisterProcess indicates an expected call of RegisterProcess.
func (mr *MockVersionServiceServerMockRecorder) RegisterProcess(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterProcess", reflect.TypeOf((*MockVersionServiceServer)(nil).RegisterProcess), arg0, arg1)
}

// Start mocks base method.
func (m *MockVersionServiceServer) Start(arg0 context.Context, arg1 *versionpb.StartRequest) (*versionpb.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Start", arg0, arg1)
	ret0, _ := ret[0].(*versionpb.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Start indicates an expected call of Start.
func (mr *MockVersionServiceServerMockRecorder) Start(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockVersionServiceServer)(nil).Start), arg0, arg1)
}

// Stop mocks base method.
func (m *MockVersionServiceServer) Stop(arg0 context.Context, arg1 *versionpb.StopRequest) (*versionpb.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stop", arg0, arg1)
	ret0, _ := ret[0].(*versionpb.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Stop indicates an expected call of Stop.
func (mr *MockVersionServiceServerMockRecorder) Stop(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockVersionServiceServer)(nil).Stop), arg0, arg1)
}

// Unpublish mocks base method.
func (m *MockVersionServiceServer) Unpublish(arg0 context.Context, arg1 *versionpb.UnpublishRequest) (*versionpb.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Unpublish", arg0, arg1)
	ret0, _ := ret[0].(*versionpb.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Unpublish indicates an expected call of Unpublish.
func (mr *MockVersionServiceServerMockRecorder) Unpublish(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unpublish", reflect.TypeOf((*MockVersionServiceServer)(nil).Unpublish), arg0, arg1)
}

// WatchProcessStatus mocks base method.
func (m *MockVersionServiceServer) WatchProcessStatus(arg0 *versionpb.ProcessStatusRequest, arg1 versionpb.VersionService_WatchProcessStatusServer) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WatchProcessStatus", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// WatchProcessStatus indicates an expected call of WatchProcessStatus.
func (mr *MockVersionServiceServerMockRecorder) WatchProcessStatus(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WatchProcessStatus", reflect.TypeOf((*MockVersionServiceServer)(nil).WatchProcessStatus), arg0, arg1)
}

// mustEmbedUnimplementedVersionServiceServer mocks base method.
func (m *MockVersionServiceServer) mustEmbedUnimplementedVersionServiceServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedVersionServiceServer")
}

// mustEmbedUnimplementedVersionServiceServer indicates an expected call of mustEmbedUnimplementedVersionServiceServer.
func (mr *MockVersionServiceServerMockRecorder) mustEmbedUnimplementedVersionServiceServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedVersionServiceServer", reflect.TypeOf((*MockVersionServiceServer)(nil).mustEmbedUnimplementedVersionServiceServer))
}

// MockUnsafeVersionServiceServer is a mock of UnsafeVersionServiceServer interface.
type MockUnsafeVersionServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockUnsafeVersionServiceServerMockRecorder
}

// MockUnsafeVersionServiceServerMockRecorder is the mock recorder for MockUnsafeVersionServiceServer.
type MockUnsafeVersionServiceServerMockRecorder struct {
	mock *MockUnsafeVersionServiceServer
}

// NewMockUnsafeVersionServiceServer creates a new mock instance.
func NewMockUnsafeVersionServiceServer(ctrl *gomock.Controller) *MockUnsafeVersionServiceServer {
	mock := &MockUnsafeVersionServiceServer{ctrl: ctrl}
	mock.recorder = &MockUnsafeVersionServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUnsafeVersionServiceServer) EXPECT() *MockUnsafeVersionServiceServerMockRecorder {
	return m.recorder
}

// mustEmbedUnimplementedVersionServiceServer mocks base method.
func (m *MockUnsafeVersionServiceServer) mustEmbedUnimplementedVersionServiceServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedVersionServiceServer")
}

// mustEmbedUnimplementedVersionServiceServer indicates an expected call of mustEmbedUnimplementedVersionServiceServer.
func (mr *MockUnsafeVersionServiceServerMockRecorder) mustEmbedUnimplementedVersionServiceServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedVersionServiceServer", reflect.TypeOf((*MockUnsafeVersionServiceServer)(nil).mustEmbedUnimplementedVersionServiceServer))
}

// MockVersionService_WatchProcessStatusServer is a mock of VersionService_WatchProcessStatusServer interface.
type MockVersionService_WatchProcessStatusServer struct {
	ctrl     *gomock.Controller
	recorder *MockVersionService_WatchProcessStatusServerMockRecorder
}

// MockVersionService_WatchProcessStatusServerMockRecorder is the mock recorder for MockVersionService_WatchProcessStatusServer.
type MockVersionService_WatchProcessStatusServerMockRecorder struct {
	mock *MockVersionService_WatchProcessStatusServer
}

// NewMockVersionService_WatchProcessStatusServer creates a new mock instance.
func NewMockVersionService_WatchProcessStatusServer(ctrl *gomock.Controller) *MockVersionService_WatchProcessStatusServer {
	mock := &MockVersionService_WatchProcessStatusServer{ctrl: ctrl}
	mock.recorder = &MockVersionService_WatchProcessStatusServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockVersionService_WatchProcessStatusServer) EXPECT() *MockVersionService_WatchProcessStatusServerMockRecorder {
	return m.recorder
}

// Context mocks base method.
func (m *MockVersionService_WatchProcessStatusServer) Context() context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Context indicates an expected call of Context.
func (mr *MockVersionService_WatchProcessStatusServerMockRecorder) Context() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockVersionService_WatchProcessStatusServer)(nil).Context))
}

// RecvMsg mocks base method.
func (m_2 *MockVersionService_WatchProcessStatusServer) RecvMsg(m any) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "RecvMsg", m)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecvMsg indicates an expected call of RecvMsg.
func (mr *MockVersionService_WatchProcessStatusServerMockRecorder) RecvMsg(m interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecvMsg", reflect.TypeOf((*MockVersionService_WatchProcessStatusServer)(nil).RecvMsg), m)
}

// Send mocks base method.
func (m *MockVersionService_WatchProcessStatusServer) Send(arg0 *versionpb.ProcessStatusResponse) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Send indicates an expected call of Send.
func (mr *MockVersionService_WatchProcessStatusServerMockRecorder) Send(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockVersionService_WatchProcessStatusServer)(nil).Send), arg0)
}

// SendHeader mocks base method.
func (m *MockVersionService_WatchProcessStatusServer) SendHeader(arg0 metadata.MD) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendHeader", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendHeader indicates an expected call of SendHeader.
func (mr *MockVersionService_WatchProcessStatusServerMockRecorder) SendHeader(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendHeader", reflect.TypeOf((*MockVersionService_WatchProcessStatusServer)(nil).SendHeader), arg0)
}

// SendMsg mocks base method.
func (m_2 *MockVersionService_WatchProcessStatusServer) SendMsg(m any) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "SendMsg", m)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMsg indicates an expected call of SendMsg.
func (mr *MockVersionService_WatchProcessStatusServerMockRecorder) SendMsg(m interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMsg", reflect.TypeOf((*MockVersionService_WatchProcessStatusServer)(nil).SendMsg), m)
}

// SetHeader mocks base method.
func (m *MockVersionService_WatchProcessStatusServer) SetHeader(arg0 metadata.MD) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetHeader", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetHeader indicates an expected call of SetHeader.
func (mr *MockVersionService_WatchProcessStatusServerMockRecorder) SetHeader(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetHeader", reflect.TypeOf((*MockVersionService_WatchProcessStatusServer)(nil).SetHeader), arg0)
}

// SetTrailer mocks base method.
func (m *MockVersionService_WatchProcessStatusServer) SetTrailer(arg0 metadata.MD) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetTrailer", arg0)
}

// SetTrailer indicates an expected call of SetTrailer.
func (mr *MockVersionService_WatchProcessStatusServerMockRecorder) SetTrailer(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetTrailer", reflect.TypeOf((*MockVersionService_WatchProcessStatusServer)(nil).SetTrailer), arg0)
}
