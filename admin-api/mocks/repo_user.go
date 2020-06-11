// Code generated by MockGen. DO NOT EDIT.
// Source: user.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	entity "gitlab.com/konstellation/kre/admin-api/domain/entity"
	reflect "reflect"
)

// MockUserRepo is a mock of UserRepo interface
type MockUserRepo struct {
	ctrl     *gomock.Controller
	recorder *MockUserRepoMockRecorder
}

// MockUserRepoMockRecorder is the mock recorder for MockUserRepo
type MockUserRepoMockRecorder struct {
	mock *MockUserRepo
}

// NewMockUserRepo creates a new mock instance
func NewMockUserRepo(ctrl *gomock.Controller) *MockUserRepo {
	mock := &MockUserRepo{ctrl: ctrl}
	mock.recorder = &MockUserRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockUserRepo) EXPECT() *MockUserRepoMockRecorder {
	return m.recorder
}

// GetByEmail mocks base method
func (m *MockUserRepo) GetByEmail(email string) (*entity.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByEmail", email)
	ret0, _ := ret[0].(*entity.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByEmail indicates an expected call of GetByEmail
func (mr *MockUserRepoMockRecorder) GetByEmail(email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByEmail", reflect.TypeOf((*MockUserRepo)(nil).GetByEmail), email)
}

// Create mocks base method
func (m *MockUserRepo) Create(ctx context.Context, email string, accessLevel entity.AccessLevel) (*entity.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, email, accessLevel)
	ret0, _ := ret[0].(*entity.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create
func (mr *MockUserRepoMockRecorder) Create(ctx, email, accessLevel interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockUserRepo)(nil).Create), ctx, email, accessLevel)
}

// GetByID mocks base method
func (m *MockUserRepo) GetByID(userID string) (*entity.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", userID)
	ret0, _ := ret[0].(*entity.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID
func (mr *MockUserRepoMockRecorder) GetByID(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockUserRepo)(nil).GetByID), userID)
}

// GetByIDs mocks base method
func (m *MockUserRepo) GetByIDs(keys []string) ([]*entity.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByIDs", keys)
	ret0, _ := ret[0].([]*entity.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByIDs indicates an expected call of GetByIDs
func (mr *MockUserRepoMockRecorder) GetByIDs(keys interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByIDs", reflect.TypeOf((*MockUserRepo)(nil).GetByIDs), keys)
}

// GetAll mocks base method
func (m *MockUserRepo) GetAll(ctx context.Context, returnDeleted bool) ([]*entity.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", ctx, returnDeleted)
	ret0, _ := ret[0].([]*entity.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll
func (mr *MockUserRepoMockRecorder) GetAll(ctx, returnDeleted interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockUserRepo)(nil).GetAll), ctx, returnDeleted)
}

// UpdateAccessLevel mocks base method
func (m *MockUserRepo) UpdateAccessLevel(ctx context.Context, userIDs []string, accessLevel entity.AccessLevel) ([]*entity.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateAccessLevel", ctx, userIDs, accessLevel)
	ret0, _ := ret[0].([]*entity.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateAccessLevel indicates an expected call of UpdateAccessLevel
func (mr *MockUserRepoMockRecorder) UpdateAccessLevel(ctx, userIDs, accessLevel interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateAccessLevel", reflect.TypeOf((*MockUserRepo)(nil).UpdateAccessLevel), ctx, userIDs, accessLevel)
}

// MarkAsDeleted mocks base method
func (m *MockUserRepo) MarkAsDeleted(ctx context.Context, userIDs []string) ([]*entity.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MarkAsDeleted", ctx, userIDs)
	ret0, _ := ret[0].([]*entity.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MarkAsDeleted indicates an expected call of MarkAsDeleted
func (mr *MockUserRepoMockRecorder) MarkAsDeleted(ctx, userIDs interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MarkAsDeleted", reflect.TypeOf((*MockUserRepo)(nil).MarkAsDeleted), ctx, userIDs)
}

// UpdateLastAccess mocks base method
func (m *MockUserRepo) UpdateLastAccess(userID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateLastAccess", userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateLastAccess indicates an expected call of UpdateLastAccess
func (mr *MockUserRepoMockRecorder) UpdateLastAccess(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateLastAccess", reflect.TypeOf((*MockUserRepo)(nil).UpdateLastAccess), userID)
}
