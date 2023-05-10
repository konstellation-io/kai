// Code generated by MockGen. DO NOT EDIT.
// Source: version.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	entity "github.com/konstellation-io/kre/engine/admin-api/domain/entity"
)

// MockVersionService is a mock of VersionService interface.
type MockVersionService struct {
	ctrl     *gomock.Controller
	recorder *MockVersionServiceMockRecorder
}

// MockVersionServiceMockRecorder is the mock recorder for MockVersionService.
type MockVersionServiceMockRecorder struct {
	mock *MockVersionService
}

// NewMockVersionService creates a new mock instance.
func NewMockVersionService(ctrl *gomock.Controller) *MockVersionService {
	mock := &MockVersionService{ctrl: ctrl}
	mock.recorder = &MockVersionServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockVersionService) EXPECT() *MockVersionServiceMockRecorder {
	return m.recorder
}

// Publish mocks base method.
func (m *MockVersionService) Publish(runtimeID string, version *entity.Version) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Publish", runtimeID, version)
	ret0, _ := ret[0].(error)
	return ret0
}

// Publish indicates an expected call of Publish.
func (mr *MockVersionServiceMockRecorder) Publish(runtimeID, version interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Publish", reflect.TypeOf((*MockVersionService)(nil).Publish), runtimeID, version)
}

// Start mocks base method.
func (m *MockVersionService) Start(ctx context.Context, runtimeID string, version *entity.Version, versionConfig *entity.VersionConfig) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Start", ctx, runtimeID, version, versionConfig)
	ret0, _ := ret[0].(error)
	return ret0
}

// Start indicates an expected call of Start.
func (mr *MockVersionServiceMockRecorder) Start(ctx, runtimeID, version, versionConfig interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockVersionService)(nil).Start), ctx, runtimeID, version, versionConfig)
}

// Stop mocks base method.
func (m *MockVersionService) Stop(ctx context.Context, runtimeID string, version *entity.Version) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stop", ctx, runtimeID, version)
	ret0, _ := ret[0].(error)
	return ret0
}

// Stop indicates an expected call of Stop.
func (mr *MockVersionServiceMockRecorder) Stop(ctx, runtimeID, version interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockVersionService)(nil).Stop), ctx, runtimeID, version)
}

// Unpublish mocks base method.
func (m *MockVersionService) Unpublish(runtimeID string, version *entity.Version) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Unpublish", runtimeID, version)
	ret0, _ := ret[0].(error)
	return ret0
}

// Unpublish indicates an expected call of Unpublish.
func (mr *MockVersionServiceMockRecorder) Unpublish(runtimeID, version interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unpublish", reflect.TypeOf((*MockVersionService)(nil).Unpublish), runtimeID, version)
}

// UpdateConfig mocks base method.
func (m *MockVersionService) UpdateConfig(runtimeID string, version *entity.Version) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateConfig", runtimeID, version)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateConfig indicates an expected call of UpdateConfig.
func (mr *MockVersionServiceMockRecorder) UpdateConfig(runtimeID, version interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateConfig", reflect.TypeOf((*MockVersionService)(nil).UpdateConfig), runtimeID, version)
}

// WatchNodeStatus mocks base method.
func (m *MockVersionService) WatchNodeStatus(ctx context.Context, runtimeID, versionName string) (<-chan *entity.Node, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WatchNodeStatus", ctx, runtimeID, versionName)
	ret0, _ := ret[0].(<-chan *entity.Node)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WatchNodeStatus indicates an expected call of WatchNodeStatus.
func (mr *MockVersionServiceMockRecorder) WatchNodeStatus(ctx, runtimeID, versionName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WatchNodeStatus", reflect.TypeOf((*MockVersionService)(nil).WatchNodeStatus), ctx, runtimeID, versionName)
}
