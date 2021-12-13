// Code generated by MockGen. DO NOT EDIT.
// Source: runtime.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	entity "github.com/konstellation-io/kre/engine/admin-api/domain/entity"
)

// MockRuntimeRepo is a mock of RuntimeRepo interface.
type MockRuntimeRepo struct {
	ctrl     *gomock.Controller
	recorder *MockRuntimeRepoMockRecorder
}

// MockRuntimeRepoMockRecorder is the mock recorder for MockRuntimeRepo.
type MockRuntimeRepoMockRecorder struct {
	mock *MockRuntimeRepo
}

// NewMockRuntimeRepo creates a new mock instance.
func NewMockRuntimeRepo(ctrl *gomock.Controller) *MockRuntimeRepo {
	mock := &MockRuntimeRepo{ctrl: ctrl}
	mock.recorder = &MockRuntimeRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRuntimeRepo) EXPECT() *MockRuntimeRepoMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockRuntimeRepo) Create(ctx context.Context, runtime *entity.Runtime) (*entity.Runtime, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, runtime)
	ret0, _ := ret[0].(*entity.Runtime)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockRuntimeRepoMockRecorder) Create(ctx, runtime interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockRuntimeRepo)(nil).Create), ctx, runtime)
}

// Get mocks base method.
func (m *MockRuntimeRepo) Get(ctx context.Context) (*entity.Runtime, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx)
	ret0, _ := ret[0].(*entity.Runtime)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockRuntimeRepoMockRecorder) Get(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockRuntimeRepo)(nil).Get), ctx)
}
