// Code generated by MockGen. DO NOT EDIT.
// Source: user_activity.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	entity "github.com/konstellation-io/kai/engine/admin-api/domain/entity"
)

// MockUserActivityInteracter is a mock of UserActivityInteracter interface.
type MockUserActivityInteracter struct {
	ctrl     *gomock.Controller
	recorder *MockUserActivityInteracterMockRecorder
}

// MockUserActivityInteracterMockRecorder is the mock recorder for MockUserActivityInteracter.
type MockUserActivityInteracterMockRecorder struct {
	mock *MockUserActivityInteracter
}

// NewMockUserActivityInteracter creates a new mock instance.
func NewMockUserActivityInteracter(ctrl *gomock.Controller) *MockUserActivityInteracter {
	mock := &MockUserActivityInteracter{ctrl: ctrl}
	mock.recorder = &MockUserActivityInteracterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserActivityInteracter) EXPECT() *MockUserActivityInteracterMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockUserActivityInteracter) Get(ctx context.Context, loggedUserID string, userEmail *string, types []entity.UserActivityType, versionIDs []string, fromDate, toDate, lastID *string) ([]*entity.UserActivity, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, loggedUserID, userEmail, types, versionIDs, fromDate, toDate, lastID)
	ret0, _ := ret[0].([]*entity.UserActivity)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockUserActivityInteracterMockRecorder) Get(ctx, loggedUserID, userEmail, types, versionIDs, fromDate, toDate, lastID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockUserActivityInteracter)(nil).Get), ctx, loggedUserID, userEmail, types, versionIDs, fromDate, toDate, lastID)
}

// RegisterCreateAction mocks base method.
func (m *MockUserActivityInteracter) RegisterCreateAction(userID, runtimeID string, version *entity.Version) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterCreateAction", userID, runtimeID, version)
	ret0, _ := ret[0].(error)
	return ret0
}

// RegisterCreateAction indicates an expected call of RegisterCreateAction.
func (mr *MockUserActivityInteracterMockRecorder) RegisterCreateAction(userID, runtimeID, version interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterCreateAction", reflect.TypeOf((*MockUserActivityInteracter)(nil).RegisterCreateAction), userID, runtimeID, version)
}

// RegisterCreateRuntime mocks base method.
func (m *MockUserActivityInteracter) RegisterCreateRuntime(userID string, runtime *entity.Runtime) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterCreateRuntime", userID, runtime)
	ret0, _ := ret[0].(error)
	return ret0
}

// RegisterCreateRuntime indicates an expected call of RegisterCreateRuntime.
func (mr *MockUserActivityInteracterMockRecorder) RegisterCreateRuntime(userID, runtime interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterCreateRuntime", reflect.TypeOf((*MockUserActivityInteracter)(nil).RegisterCreateRuntime), userID, runtime)
}

// RegisterPublishAction mocks base method.
func (m *MockUserActivityInteracter) RegisterPublishAction(userID, runtimeID string, version, prev *entity.Version, comment string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterPublishAction", userID, runtimeID, version, prev, comment)
	ret0, _ := ret[0].(error)
	return ret0
}

// RegisterPublishAction indicates an expected call of RegisterPublishAction.
func (mr *MockUserActivityInteracterMockRecorder) RegisterPublishAction(userID, runtimeID, version, prev, comment interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterPublishAction", reflect.TypeOf((*MockUserActivityInteracter)(nil).RegisterPublishAction), userID, runtimeID, version, prev, comment)
}

// RegisterStartAction mocks base method.
func (m *MockUserActivityInteracter) RegisterStartAction(userID, runtimeID string, version *entity.Version, comment string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterStartAction", userID, runtimeID, version, comment)
	ret0, _ := ret[0].(error)
	return ret0
}

// RegisterStartAction indicates an expected call of RegisterStartAction.
func (mr *MockUserActivityInteracterMockRecorder) RegisterStartAction(userID, runtimeID, version, comment interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterStartAction", reflect.TypeOf((*MockUserActivityInteracter)(nil).RegisterStartAction), userID, runtimeID, version, comment)
}

// RegisterStopAction mocks base method.
func (m *MockUserActivityInteracter) RegisterStopAction(userID, runtimeID string, version *entity.Version, comment string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterStopAction", userID, runtimeID, version, comment)
	ret0, _ := ret[0].(error)
	return ret0
}

// RegisterStopAction indicates an expected call of RegisterStopAction.
func (mr *MockUserActivityInteracterMockRecorder) RegisterStopAction(userID, runtimeID, version, comment interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterStopAction", reflect.TypeOf((*MockUserActivityInteracter)(nil).RegisterStopAction), userID, runtimeID, version, comment)
}

// RegisterUnpublishAction mocks base method.
func (m *MockUserActivityInteracter) RegisterUnpublishAction(userID, runtimeID string, version *entity.Version, comment string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterUnpublishAction", userID, runtimeID, version, comment)
	ret0, _ := ret[0].(error)
	return ret0
}

// RegisterUnpublishAction indicates an expected call of RegisterUnpublishAction.
func (mr *MockUserActivityInteracterMockRecorder) RegisterUnpublishAction(userID, runtimeID, version, comment interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterUnpublishAction", reflect.TypeOf((*MockUserActivityInteracter)(nil).RegisterUnpublishAction), userID, runtimeID, version, comment)
}

// RegisterUpdateProductGrants mocks base method.
func (m *MockUserActivityInteracter) RegisterUpdateProductGrants(userID, targetUserID, product string, productGrants []string, comment string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterUpdateProductGrants", userID, targetUserID, product, productGrants, comment)
	ret0, _ := ret[0].(error)
	return ret0
}

// RegisterUpdateProductGrants indicates an expected call of RegisterUpdateProductGrants.
func (mr *MockUserActivityInteracterMockRecorder) RegisterUpdateProductGrants(userID, targetUserID, product, productGrants, comment interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterUpdateProductGrants", reflect.TypeOf((*MockUserActivityInteracter)(nil).RegisterUpdateProductGrants), userID, targetUserID, product, productGrants, comment)
}
