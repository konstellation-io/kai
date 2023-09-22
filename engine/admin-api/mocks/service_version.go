// Code generated by MockGen. DO NOT EDIT.
// Source: version.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	entity "github.com/konstellation-io/kai/engine/admin-api/domain/entity"
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
func (m *MockVersionService) Publish(ctx context.Context, productID, versionTag string) (map[string]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Publish", ctx, productID, versionTag)
	ret0, _ := ret[0].(map[string]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Publish indicates an expected call of Publish.
func (mr *MockVersionServiceMockRecorder) Publish(ctx, productID, versionTag interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Publish", reflect.TypeOf((*MockVersionService)(nil).Publish), ctx, productID, versionTag)
}

// RegisterProcess mocks base method.
func (m *MockVersionService) RegisterProcess(ctx context.Context, processID, processImage string, file []byte) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterProcess", ctx, processID, processImage, file)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RegisterProcess indicates an expected call of RegisterProcess.
func (mr *MockVersionServiceMockRecorder) RegisterProcess(ctx, processID, processImage, file interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterProcess", reflect.TypeOf((*MockVersionService)(nil).RegisterProcess), ctx, processID, processImage, file)
}

// Start mocks base method.
func (m *MockVersionService) Start(ctx context.Context, productID string, version *entity.Version, versionConfig *entity.VersionConfig) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Start", ctx, productID, version, versionConfig)
	ret0, _ := ret[0].(error)
	return ret0
}

// Start indicates an expected call of Start.
func (mr *MockVersionServiceMockRecorder) Start(ctx, productID, version, versionConfig interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockVersionService)(nil).Start), ctx, productID, version, versionConfig)
}

// Stop mocks base method.
func (m *MockVersionService) Stop(ctx context.Context, productID string, version *entity.Version) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stop", ctx, productID, version)
	ret0, _ := ret[0].(error)
	return ret0
}

// Stop indicates an expected call of Stop.
func (mr *MockVersionServiceMockRecorder) Stop(ctx, productID, version interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockVersionService)(nil).Stop), ctx, productID, version)
}

// Unpublish mocks base method.
func (m *MockVersionService) Unpublish(ctx context.Context, productID string, version *entity.Version) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Unpublish", ctx, productID, version)
	ret0, _ := ret[0].(error)
	return ret0
}

// Unpublish indicates an expected call of Unpublish.
func (mr *MockVersionServiceMockRecorder) Unpublish(ctx, productID, version interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unpublish", reflect.TypeOf((*MockVersionService)(nil).Unpublish), ctx, productID, version)
}

// WatchProcessStatus mocks base method.
func (m *MockVersionService) WatchProcessStatus(ctx context.Context, productID, versionTag string) (<-chan *entity.Process, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WatchProcessStatus", ctx, productID, versionTag)
	ret0, _ := ret[0].(<-chan *entity.Process)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WatchProcessStatus indicates an expected call of WatchProcessStatus.
func (mr *MockVersionServiceMockRecorder) WatchProcessStatus(ctx, productID, versionTag interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WatchProcessStatus", reflect.TypeOf((*MockVersionService)(nil).WatchProcessStatus), ctx, productID, versionTag)
}
