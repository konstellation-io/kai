// Code generated by MockGen. DO NOT EDIT.
// Source: access_control.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	token "github.com/konstellation-io/kre/engine/admin-api/delivery/http/token"
	auth "github.com/konstellation-io/kre/engine/admin-api/domain/usecase/auth"
)

// MockAccessControl is a mock of AccessControl interface.
type MockAccessControl struct {
	ctrl     *gomock.Controller
	recorder *MockAccessControlMockRecorder
}

// MockAccessControlMockRecorder is the mock recorder for MockAccessControl.
type MockAccessControlMockRecorder struct {
	mock *MockAccessControl
}

// NewMockAccessControl creates a new mock instance.
func NewMockAccessControl(ctrl *gomock.Controller) *MockAccessControl {
	mock := &MockAccessControl{ctrl: ctrl}
	mock.recorder = &MockAccessControlMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAccessControl) EXPECT() *MockAccessControlMockRecorder {
	return m.recorder
}

// CheckPermission mocks base method.
func (m *MockAccessControl) CheckPermission(user *token.UserRoles, product string, action auth.AccessControlAction) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckPermission", user, product, action)
	ret0, _ := ret[0].(error)
	return ret0
}

// CheckPermission indicates an expected call of CheckPermission.
func (mr *MockAccessControlMockRecorder) CheckPermission(user, product, action interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckPermission", reflect.TypeOf((*MockAccessControl)(nil).CheckPermission), user, product, action)
}
