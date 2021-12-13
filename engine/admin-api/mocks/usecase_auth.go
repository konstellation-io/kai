// Code generated by MockGen. DO NOT EDIT.
// Source: auth.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	entity "github.com/konstellation-io/kre/engine/admin-api/domain/entity"
)

// MockAuthInteracter is a mock of AuthInteracter interface.
type MockAuthInteracter struct {
	ctrl     *gomock.Controller
	recorder *MockAuthInteracterMockRecorder
}

// MockAuthInteracterMockRecorder is the mock recorder for MockAuthInteracter.
type MockAuthInteracterMockRecorder struct {
	mock *MockAuthInteracter
}

// NewMockAuthInteracter creates a new mock instance.
func NewMockAuthInteracter(ctrl *gomock.Controller) *MockAuthInteracter {
	mock := &MockAuthInteracter{ctrl: ctrl}
	mock.recorder = &MockAuthInteracterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthInteracter) EXPECT() *MockAuthInteracterMockRecorder {
	return m.recorder
}

// CheckSessionIsActive mocks base method.
func (m *MockAuthInteracter) CheckSessionIsActive(token string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckSessionIsActive", token)
	ret0, _ := ret[0].(error)
	return ret0
}

// CheckSessionIsActive indicates an expected call of CheckSessionIsActive.
func (mr *MockAuthInteracterMockRecorder) CheckSessionIsActive(token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckSessionIsActive", reflect.TypeOf((*MockAuthInteracter)(nil).CheckSessionIsActive), token)
}

// CountUserSessions mocks base method.
func (m *MockAuthInteracter) CountUserSessions(ctx context.Context, userID string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CountUserSessions", ctx, userID)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CountUserSessions indicates an expected call of CountUserSessions.
func (mr *MockAuthInteracterMockRecorder) CountUserSessions(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CountUserSessions", reflect.TypeOf((*MockAuthInteracter)(nil).CountUserSessions), ctx, userID)
}

// CreateSession mocks base method.
func (m *MockAuthInteracter) CreateSession(session entity.Session) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSession", session)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateSession indicates an expected call of CreateSession.
func (mr *MockAuthInteracterMockRecorder) CreateSession(session interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSession", reflect.TypeOf((*MockAuthInteracter)(nil).CreateSession), session)
}

// Logout mocks base method.
func (m *MockAuthInteracter) Logout(userID, token string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Logout", userID, token)
	ret0, _ := ret[0].(error)
	return ret0
}

// Logout indicates an expected call of Logout.
func (mr *MockAuthInteracterMockRecorder) Logout(userID, token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Logout", reflect.TypeOf((*MockAuthInteracter)(nil).Logout), userID, token)
}

// RevokeUserSessions mocks base method.
func (m *MockAuthInteracter) RevokeUserSessions(userIDs []string, loggedUser, comment string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RevokeUserSessions", userIDs, loggedUser, comment)
	ret0, _ := ret[0].(error)
	return ret0
}

// RevokeUserSessions indicates an expected call of RevokeUserSessions.
func (mr *MockAuthInteracterMockRecorder) RevokeUserSessions(userIDs, loggedUser, comment interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RevokeUserSessions", reflect.TypeOf((*MockAuthInteracter)(nil).RevokeUserSessions), userIDs, loggedUser, comment)
}

// SignIn mocks base method.
func (m *MockAuthInteracter) SignIn(ctx context.Context, email string, verificationCodeDurationInMinutes int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignIn", ctx, email, verificationCodeDurationInMinutes)
	ret0, _ := ret[0].(error)
	return ret0
}

// SignIn indicates an expected call of SignIn.
func (mr *MockAuthInteracterMockRecorder) SignIn(ctx, email, verificationCodeDurationInMinutes interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignIn", reflect.TypeOf((*MockAuthInteracter)(nil).SignIn), ctx, email, verificationCodeDurationInMinutes)
}

// UpdateLastActivity mocks base method.
func (m *MockAuthInteracter) UpdateLastActivity(loggedUserID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateLastActivity", loggedUserID)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateLastActivity indicates an expected call of UpdateLastActivity.
func (mr *MockAuthInteracterMockRecorder) UpdateLastActivity(loggedUserID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateLastActivity", reflect.TypeOf((*MockAuthInteracter)(nil).UpdateLastActivity), loggedUserID)
}

// VerifyAPIToken mocks base method.
func (m *MockAuthInteracter) VerifyAPIToken(ctx context.Context, apiToken string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifyAPIToken", ctx, apiToken)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// VerifyAPIToken indicates an expected call of VerifyAPIToken.
func (mr *MockAuthInteracterMockRecorder) VerifyAPIToken(ctx, apiToken interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyAPIToken", reflect.TypeOf((*MockAuthInteracter)(nil).VerifyAPIToken), ctx, apiToken)
}

// VerifyCode mocks base method.
func (m *MockAuthInteracter) VerifyCode(code string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifyCode", code)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// VerifyCode indicates an expected call of VerifyCode.
func (mr *MockAuthInteracterMockRecorder) VerifyCode(code interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyCode", reflect.TypeOf((*MockAuthInteracter)(nil).VerifyCode), code)
}
