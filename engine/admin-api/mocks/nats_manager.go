// Code generated by MockGen. DO NOT EDIT.
// Source: ../proto/natspb/nats_grpc.pb.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	natspb "github.com/konstellation-io/kai/engine/admin-api/adapter/service/proto/natspb"
	grpc "google.golang.org/grpc"
)

// MockNatsManagerServiceClient is a mock of NatsManagerServiceClient interface.
type MockNatsManagerServiceClient struct {
	ctrl     *gomock.Controller
	recorder *MockNatsManagerServiceClientMockRecorder
}

// MockNatsManagerServiceClientMockRecorder is the mock recorder for MockNatsManagerServiceClient.
type MockNatsManagerServiceClientMockRecorder struct {
	mock *MockNatsManagerServiceClient
}

// NewMockNatsManagerServiceClient creates a new mock instance.
func NewMockNatsManagerServiceClient(ctrl *gomock.Controller) *MockNatsManagerServiceClient {
	mock := &MockNatsManagerServiceClient{ctrl: ctrl}
	mock.recorder = &MockNatsManagerServiceClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNatsManagerServiceClient) EXPECT() *MockNatsManagerServiceClientMockRecorder {
	return m.recorder
}

// CreateKeyValueStores mocks base method.
func (m *MockNatsManagerServiceClient) CreateKeyValueStores(ctx context.Context, in *natspb.CreateKeyValueStoresRequest, opts ...grpc.CallOption) (*natspb.CreateKeyValueStoreResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CreateVersionKeyValueStores", varargs...)
	ret0, _ := ret[0].(*natspb.CreateKeyValueStoreResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateKeyValueStores indicates an expected call of CreateKeyValueStores.
func (mr *MockNatsManagerServiceClientMockRecorder) CreateKeyValueStores(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateVersionKeyValueStores", reflect.TypeOf((*MockNatsManagerServiceClient)(nil).CreateKeyValueStores), varargs...)
}

// CreateObjectStores mocks base method.
func (m *MockNatsManagerServiceClient) CreateObjectStores(ctx context.Context, in *natspb.CreateObjectStoresRequest, opts ...grpc.CallOption) (*natspb.CreateObjectStoresResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CreateObjectStores", varargs...)
	ret0, _ := ret[0].(*natspb.CreateObjectStoresResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateObjectStores indicates an expected call of CreateObjectStores.
func (mr *MockNatsManagerServiceClientMockRecorder) CreateObjectStores(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateObjectStores", reflect.TypeOf((*MockNatsManagerServiceClient)(nil).CreateObjectStores), varargs...)
}

// CreateStreams mocks base method.
func (m *MockNatsManagerServiceClient) CreateStreams(ctx context.Context, in *natspb.CreateStreamsRequest, opts ...grpc.CallOption) (*natspb.CreateStreamsResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CreateStreams", varargs...)
	ret0, _ := ret[0].(*natspb.CreateStreamsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateStreams indicates an expected call of CreateStreams.
func (mr *MockNatsManagerServiceClientMockRecorder) CreateStreams(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateStreams", reflect.TypeOf((*MockNatsManagerServiceClient)(nil).CreateStreams), varargs...)
}

// DeleteObjectStores mocks base method.
func (m *MockNatsManagerServiceClient) DeleteObjectStores(ctx context.Context, in *natspb.DeleteObjectStoresRequest, opts ...grpc.CallOption) (*natspb.DeleteResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeleteObjectStores", varargs...)
	ret0, _ := ret[0].(*natspb.DeleteResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteObjectStores indicates an expected call of DeleteObjectStores.
func (mr *MockNatsManagerServiceClientMockRecorder) DeleteObjectStores(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteObjectStores", reflect.TypeOf((*MockNatsManagerServiceClient)(nil).DeleteObjectStores), varargs...)
}

// DeleteStreams mocks base method.
func (m *MockNatsManagerServiceClient) DeleteStreams(ctx context.Context, in *natspb.DeleteStreamsRequest, opts ...grpc.CallOption) (*natspb.DeleteResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeleteStreams", varargs...)
	ret0, _ := ret[0].(*natspb.DeleteResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteStreams indicates an expected call of DeleteStreams.
func (mr *MockNatsManagerServiceClientMockRecorder) DeleteStreams(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteStreams", reflect.TypeOf((*MockNatsManagerServiceClient)(nil).DeleteStreams), varargs...)
}

// MockNatsManagerServiceServer is a mock of NatsManagerServiceServer interface.
type MockNatsManagerServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockNatsManagerServiceServerMockRecorder
}

// MockNatsManagerServiceServerMockRecorder is the mock recorder for MockNatsManagerServiceServer.
type MockNatsManagerServiceServerMockRecorder struct {
	mock *MockNatsManagerServiceServer
}

// NewMockNatsManagerServiceServer creates a new mock instance.
func NewMockNatsManagerServiceServer(ctrl *gomock.Controller) *MockNatsManagerServiceServer {
	mock := &MockNatsManagerServiceServer{ctrl: ctrl}
	mock.recorder = &MockNatsManagerServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNatsManagerServiceServer) EXPECT() *MockNatsManagerServiceServerMockRecorder {
	return m.recorder
}

// CreateKeyValueStores mocks base method.
func (m *MockNatsManagerServiceServer) CreateKeyValueStores(arg0 context.Context, arg1 *natspb.CreateKeyValueStoresRequest) (*natspb.CreateKeyValueStoreResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateVersionKeyValueStores", arg0, arg1)
	ret0, _ := ret[0].(*natspb.CreateKeyValueStoreResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateKeyValueStores indicates an expected call of CreateKeyValueStores.
func (mr *MockNatsManagerServiceServerMockRecorder) CreateKeyValueStores(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateVersionKeyValueStores", reflect.TypeOf((*MockNatsManagerServiceServer)(nil).CreateKeyValueStores), arg0, arg1)
}

// CreateObjectStores mocks base method.
func (m *MockNatsManagerServiceServer) CreateObjectStores(arg0 context.Context, arg1 *natspb.CreateObjectStoresRequest) (*natspb.CreateObjectStoresResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateObjectStores", arg0, arg1)
	ret0, _ := ret[0].(*natspb.CreateObjectStoresResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateObjectStores indicates an expected call of CreateObjectStores.
func (mr *MockNatsManagerServiceServerMockRecorder) CreateObjectStores(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateObjectStores", reflect.TypeOf((*MockNatsManagerServiceServer)(nil).CreateObjectStores), arg0, arg1)
}

// CreateStreams mocks base method.
func (m *MockNatsManagerServiceServer) CreateStreams(arg0 context.Context, arg1 *natspb.CreateStreamsRequest) (*natspb.CreateStreamsResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateStreams", arg0, arg1)
	ret0, _ := ret[0].(*natspb.CreateStreamsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateStreams indicates an expected call of CreateStreams.
func (mr *MockNatsManagerServiceServerMockRecorder) CreateStreams(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateStreams", reflect.TypeOf((*MockNatsManagerServiceServer)(nil).CreateStreams), arg0, arg1)
}

// DeleteObjectStores mocks base method.
func (m *MockNatsManagerServiceServer) DeleteObjectStores(arg0 context.Context, arg1 *natspb.DeleteObjectStoresRequest) (*natspb.DeleteResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteObjectStores", arg0, arg1)
	ret0, _ := ret[0].(*natspb.DeleteResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteObjectStores indicates an expected call of DeleteObjectStores.
func (mr *MockNatsManagerServiceServerMockRecorder) DeleteObjectStores(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteObjectStores", reflect.TypeOf((*MockNatsManagerServiceServer)(nil).DeleteObjectStores), arg0, arg1)
}

// DeleteStreams mocks base method.
func (m *MockNatsManagerServiceServer) DeleteStreams(arg0 context.Context, arg1 *natspb.DeleteStreamsRequest) (*natspb.DeleteResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteStreams", arg0, arg1)
	ret0, _ := ret[0].(*natspb.DeleteResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteStreams indicates an expected call of DeleteStreams.
func (mr *MockNatsManagerServiceServerMockRecorder) DeleteStreams(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteStreams", reflect.TypeOf((*MockNatsManagerServiceServer)(nil).DeleteStreams), arg0, arg1)
}

// mustEmbedUnimplementedNatsManagerServiceServer mocks base method.
func (m *MockNatsManagerServiceServer) mustEmbedUnimplementedNatsManagerServiceServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedNatsManagerServiceServer")
}

// mustEmbedUnimplementedNatsManagerServiceServer indicates an expected call of mustEmbedUnimplementedNatsManagerServiceServer.
func (mr *MockNatsManagerServiceServerMockRecorder) mustEmbedUnimplementedNatsManagerServiceServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedNatsManagerServiceServer", reflect.TypeOf((*MockNatsManagerServiceServer)(nil).mustEmbedUnimplementedNatsManagerServiceServer))
}

// MockUnsafeNatsManagerServiceServer is a mock of UnsafeNatsManagerServiceServer interface.
type MockUnsafeNatsManagerServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockUnsafeNatsManagerServiceServerMockRecorder
}

// MockUnsafeNatsManagerServiceServerMockRecorder is the mock recorder for MockUnsafeNatsManagerServiceServer.
type MockUnsafeNatsManagerServiceServerMockRecorder struct {
	mock *MockUnsafeNatsManagerServiceServer
}

// NewMockUnsafeNatsManagerServiceServer creates a new mock instance.
func NewMockUnsafeNatsManagerServiceServer(ctrl *gomock.Controller) *MockUnsafeNatsManagerServiceServer {
	mock := &MockUnsafeNatsManagerServiceServer{ctrl: ctrl}
	mock.recorder = &MockUnsafeNatsManagerServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUnsafeNatsManagerServiceServer) EXPECT() *MockUnsafeNatsManagerServiceServerMockRecorder {
	return m.recorder
}

// mustEmbedUnimplementedNatsManagerServiceServer mocks base method.
func (m *MockUnsafeNatsManagerServiceServer) mustEmbedUnimplementedNatsManagerServiceServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedNatsManagerServiceServer")
}

// mustEmbedUnimplementedNatsManagerServiceServer indicates an expected call of mustEmbedUnimplementedNatsManagerServiceServer.
func (mr *MockUnsafeNatsManagerServiceServerMockRecorder) mustEmbedUnimplementedNatsManagerServiceServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedNatsManagerServiceServer", reflect.TypeOf((*MockUnsafeNatsManagerServiceServer)(nil).mustEmbedUnimplementedNatsManagerServiceServer))
}
