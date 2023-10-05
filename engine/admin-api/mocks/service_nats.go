// Code generated by MockGen. DO NOT EDIT.
// Source: nats.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	entity "github.com/konstellation-io/kai/engine/admin-api/domain/entity"
)

// MockNatsManagerService is a mock of NatsManagerService interface.
type MockNatsManagerService struct {
	ctrl     *gomock.Controller
	recorder *MockNatsManagerServiceMockRecorder
}

// MockNatsManagerServiceMockRecorder is the mock recorder for MockNatsManagerService.
type MockNatsManagerServiceMockRecorder struct {
	mock *MockNatsManagerService
}

// NewMockNatsManagerService creates a new mock instance.
func NewMockNatsManagerService(ctrl *gomock.Controller) *MockNatsManagerService {
	mock := &MockNatsManagerService{ctrl: ctrl}
	mock.recorder = &MockNatsManagerServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNatsManagerService) EXPECT() *MockNatsManagerServiceMockRecorder {
	return m.recorder
}

// CreateKeyValueStores mocks base method.
func (m *MockNatsManagerService) CreateKeyValueStores(ctx context.Context, runtimeID string, version *entity.Version) (*entity.KeyValueStoresConfig, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateVersionKeyValueStores", ctx, runtimeID, version)
	ret0, _ := ret[0].(*entity.KeyValueStoresConfig)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateKeyValueStores indicates an expected call of CreateKeyValueStores.
func (mr *MockNatsManagerServiceMockRecorder) CreateKeyValueStores(ctx, runtimeID, version interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateVersionKeyValueStores", reflect.TypeOf((*MockNatsManagerService)(nil).CreateKeyValueStores), ctx, runtimeID, version)
}

// CreateObjectStores mocks base method.
func (m *MockNatsManagerService) CreateObjectStores(ctx context.Context, runtimeID string, version *entity.Version) (*entity.VersionObjectStoresConfig, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateObjectStores", ctx, runtimeID, version)
	ret0, _ := ret[0].(*entity.VersionObjectStoresConfig)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateObjectStores indicates an expected call of CreateObjectStores.
func (mr *MockNatsManagerServiceMockRecorder) CreateObjectStores(ctx, runtimeID, version interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateObjectStores", reflect.TypeOf((*MockNatsManagerService)(nil).CreateObjectStores), ctx, runtimeID, version)
}

// CreateStreams mocks base method.
func (m *MockNatsManagerService) CreateStreams(ctx context.Context, runtimeID string, version *entity.Version) (*entity.VersionStreamsConfig, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateStreams", ctx, runtimeID, version)
	ret0, _ := ret[0].(*entity.VersionStreamsConfig)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateStreams indicates an expected call of CreateStreams.
func (mr *MockNatsManagerServiceMockRecorder) CreateStreams(ctx, runtimeID, version interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateStreams", reflect.TypeOf((*MockNatsManagerService)(nil).CreateStreams), ctx, runtimeID, version)
}

// DeleteObjectStores mocks base method.
func (m *MockNatsManagerService) DeleteObjectStores(ctx context.Context, runtimeID, versionTag string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteObjectStores", ctx, runtimeID, versionTag)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteObjectStores indicates an expected call of DeleteObjectStores.
func (mr *MockNatsManagerServiceMockRecorder) DeleteObjectStores(ctx, runtimeID, versionTag interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteObjectStores", reflect.TypeOf((*MockNatsManagerService)(nil).DeleteObjectStores), ctx, runtimeID, versionTag)
}

// DeleteStreams mocks base method.
func (m *MockNatsManagerService) DeleteStreams(ctx context.Context, runtimeID, versionTag string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteStreams", ctx, runtimeID, versionTag)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteStreams indicates an expected call of DeleteStreams.
func (mr *MockNatsManagerServiceMockRecorder) DeleteStreams(ctx, runtimeID, versionTag interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteStreams", reflect.TypeOf((*MockNatsManagerService)(nil).DeleteStreams), ctx, runtimeID, versionTag)
}
