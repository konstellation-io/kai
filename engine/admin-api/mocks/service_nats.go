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

// CreateGlobalKeyValueStore mocks base method.
func (m *MockNatsManagerService) CreateGlobalKeyValueStore(ctx context.Context, product string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateGlobalKeyValueStore", ctx, product)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateGlobalKeyValueStore indicates an expected call of CreateGlobalKeyValueStore.
func (mr *MockNatsManagerServiceMockRecorder) CreateGlobalKeyValueStore(ctx, product interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateGlobalKeyValueStore", reflect.TypeOf((*MockNatsManagerService)(nil).CreateGlobalKeyValueStore), ctx, product)
}

// CreateObjectStores mocks base method.
func (m *MockNatsManagerService) CreateObjectStores(ctx context.Context, product string, version *entity.Version) (*entity.VersionObjectStores, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateObjectStores", ctx, product, version)
	ret0, _ := ret[0].(*entity.VersionObjectStores)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateObjectStores indicates an expected call of CreateObjectStores.
func (mr *MockNatsManagerServiceMockRecorder) CreateObjectStores(ctx, product, version interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateObjectStores", reflect.TypeOf((*MockNatsManagerService)(nil).CreateObjectStores), ctx, product, version)
}

// CreateStreams mocks base method.
func (m *MockNatsManagerService) CreateStreams(ctx context.Context, product string, version *entity.Version) (*entity.VersionStreams, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateStreams", ctx, product, version)
	ret0, _ := ret[0].(*entity.VersionStreams)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateStreams indicates an expected call of CreateStreams.
func (mr *MockNatsManagerServiceMockRecorder) CreateStreams(ctx, product, version interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateStreams", reflect.TypeOf((*MockNatsManagerService)(nil).CreateStreams), ctx, product, version)
}

// CreateVersionKeyValueStores mocks base method.
func (m *MockNatsManagerService) CreateVersionKeyValueStores(ctx context.Context, product string, version *entity.Version) (*entity.KeyValueStores, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateVersionKeyValueStores", ctx, product, version)
	ret0, _ := ret[0].(*entity.KeyValueStores)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateVersionKeyValueStores indicates an expected call of CreateVersionKeyValueStores.
func (mr *MockNatsManagerServiceMockRecorder) CreateVersionKeyValueStores(ctx, product, version interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateVersionKeyValueStores", reflect.TypeOf((*MockNatsManagerService)(nil).CreateVersionKeyValueStores), ctx, product, version)
}

// DeleteObjectStores mocks base method.
func (m *MockNatsManagerService) DeleteObjectStores(ctx context.Context, product, versionTag string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteObjectStores", ctx, product, versionTag)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteObjectStores indicates an expected call of DeleteObjectStores.
func (mr *MockNatsManagerServiceMockRecorder) DeleteObjectStores(ctx, product, versionTag interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteObjectStores", reflect.TypeOf((*MockNatsManagerService)(nil).DeleteObjectStores), ctx, product, versionTag)
}

// DeleteStreams mocks base method.
func (m *MockNatsManagerService) DeleteStreams(ctx context.Context, product, versionTag string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteStreams", ctx, product, versionTag)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteStreams indicates an expected call of DeleteStreams.
func (mr *MockNatsManagerServiceMockRecorder) DeleteStreams(ctx, product, versionTag interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteStreams", reflect.TypeOf((*MockNatsManagerService)(nil).DeleteStreams), ctx, product, versionTag)
}

// UpdateKeyValueConfiguration mocks base method.
func (m *MockNatsManagerService) UpdateKeyValueConfiguration(ctx context.Context, configurations []entity.KeyValueConfiguration) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateKeyValueConfiguration", ctx, configurations)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateKeyValueConfiguration indicates an expected call of UpdateKeyValueConfiguration.
func (mr *MockNatsManagerServiceMockRecorder) UpdateKeyValueConfiguration(ctx, configurations interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateKeyValueConfiguration", reflect.TypeOf((*MockNatsManagerService)(nil).UpdateKeyValueConfiguration), ctx, configurations)
}
