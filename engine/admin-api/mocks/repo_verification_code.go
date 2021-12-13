// Code generated by MockGen. DO NOT EDIT.
// Source: verification_code.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
	entity "github.com/konstellation-io/kre/engine/admin-api/domain/entity"
)

// MockVerificationCodeRepo is a mock of VerificationCodeRepo interface.
type MockVerificationCodeRepo struct {
	ctrl     *gomock.Controller
	recorder *MockVerificationCodeRepoMockRecorder
}

// MockVerificationCodeRepoMockRecorder is the mock recorder for MockVerificationCodeRepo.
type MockVerificationCodeRepoMockRecorder struct {
	mock *MockVerificationCodeRepo
}

// NewMockVerificationCodeRepo creates a new mock instance.
func NewMockVerificationCodeRepo(ctrl *gomock.Controller) *MockVerificationCodeRepo {
	mock := &MockVerificationCodeRepo{ctrl: ctrl}
	mock.recorder = &MockVerificationCodeRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockVerificationCodeRepo) EXPECT() *MockVerificationCodeRepoMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockVerificationCodeRepo) Delete(code string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", code)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockVerificationCodeRepoMockRecorder) Delete(code interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockVerificationCodeRepo)(nil).Delete), code)
}

// Get mocks base method.
func (m *MockVerificationCodeRepo) Get(code string) (*entity.VerificationCode, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", code)
	ret0, _ := ret[0].(*entity.VerificationCode)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockVerificationCodeRepoMockRecorder) Get(code interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockVerificationCodeRepo)(nil).Get), code)
}

// Store mocks base method.
func (m *MockVerificationCodeRepo) Store(code, uid string, ttl time.Duration) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Store", code, uid, ttl)
	ret0, _ := ret[0].(error)
	return ret0
}

// Store indicates an expected call of Store.
func (mr *MockVerificationCodeRepoMockRecorder) Store(code, uid, ttl interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Store", reflect.TypeOf((*MockVerificationCodeRepo)(nil).Store), code, uid, ttl)
}
