// Code generated by MockGen. DO NOT EDIT.
// Source: nats_manager.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	entity "github.com/konstellation-io/kai/engine/nats-manager/internal/entity"
)

// MockNatsManager is a mock of NatsManager interface.
type MockNatsManager struct {
	ctrl     *gomock.Controller
	recorder *MockNatsManagerMockRecorder
}

// MockNatsManagerMockRecorder is the mock recorder for MockNatsManager.
type MockNatsManagerMockRecorder struct {
	mock *MockNatsManager
}

// NewMockNatsManager creates a new mock instance.
func NewMockNatsManager(ctrl *gomock.Controller) *MockNatsManager {
	mock := &MockNatsManager{ctrl: ctrl}
	mock.recorder = &MockNatsManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNatsManager) EXPECT() *MockNatsManagerMockRecorder {
	return m.recorder
}

// CreateVersionKeyValueStores mocks base method.
func (m *MockNatsManager) CreateVersionKeyValueStores(productID, versionTag string, workflows []entity.Workflow) (*entity.VersionKeyValueStores, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateVersionKeyValueStores", productID, versionTag, workflows)
	ret0, _ := ret[0].(*entity.VersionKeyValueStores)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateKeyValueStores indicates an expected call of CreateKeyValueStores.
func (mr *MockNatsManagerMockRecorder) CreateKeyValueStores(productID, versionTag, workflows interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateVersionKeyValueStores", reflect.TypeOf((*MockNatsManager)(nil).CreateVersionKeyValueStores), productID, versionTag, workflows)
}

// CreateObjectStores mocks base method.
func (m *MockNatsManager) CreateObjectStores(productID, versionTag string, workflows []entity.Workflow) (entity.WorkflowsObjectStoresConfig, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateObjectStores", productID, versionTag, workflows)
	ret0, _ := ret[0].(entity.WorkflowsObjectStoresConfig)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateObjectStores indicates an expected call of CreateObjectStores.
func (mr *MockNatsManagerMockRecorder) CreateObjectStores(productID, versionTag, workflows interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateObjectStores", reflect.TypeOf((*MockNatsManager)(nil).CreateObjectStores), productID, versionTag, workflows)
}

// CreateStreams mocks base method.
func (m *MockNatsManager) CreateStreams(productID, versionTag string, workflows []entity.Workflow) (entity.WorkflowsStreamsConfig, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateStreams", productID, versionTag, workflows)
	ret0, _ := ret[0].(entity.WorkflowsStreamsConfig)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateStreams indicates an expected call of CreateStreams.
func (mr *MockNatsManagerMockRecorder) CreateStreams(productID, versionTag, workflows interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateStreams", reflect.TypeOf((*MockNatsManager)(nil).CreateStreams), productID, versionTag, workflows)
}

// DeleteObjectStores mocks base method.
func (m *MockNatsManager) DeleteObjectStores(productID, versionTag string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteObjectStores", productID, versionTag)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteObjectStores indicates an expected call of DeleteObjectStores.
func (mr *MockNatsManagerMockRecorder) DeleteObjectStores(productID, versionTag interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteObjectStores", reflect.TypeOf((*MockNatsManager)(nil).DeleteObjectStores), productID, versionTag)
}

// DeleteStreams mocks base method.
func (m *MockNatsManager) DeleteStreams(productID, versionTag string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteStreams", productID, versionTag)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteStreams indicates an expected call of DeleteStreams.
func (mr *MockNatsManagerMockRecorder) DeleteStreams(productID, versionTag interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteStreams", reflect.TypeOf((*MockNatsManager)(nil).DeleteStreams), productID, versionTag)
}
