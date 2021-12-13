// Code generated by MockGen. DO NOT EDIT.
// Source: api_token.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	entity "github.com/konstellation-io/kre/engine/admin-api/domain/entity"
)

// MockAPITokenRepo is a mock of APITokenRepo interface.
type MockAPITokenRepo struct {
	ctrl     *gomock.Controller
	recorder *MockAPITokenRepoMockRecorder
}

// MockAPITokenRepoMockRecorder is the mock recorder for MockAPITokenRepo.
type MockAPITokenRepoMockRecorder struct {
	mock *MockAPITokenRepo
}

// NewMockAPITokenRepo creates a new mock instance.
func NewMockAPITokenRepo(ctrl *gomock.Controller) *MockAPITokenRepo {
	mock := &MockAPITokenRepo{ctrl: ctrl}
	mock.recorder = &MockAPITokenRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAPITokenRepo) EXPECT() *MockAPITokenRepoMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockAPITokenRepo) Create(ctx context.Context, apiToken entity.APIToken, code string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, apiToken, code)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockAPITokenRepoMockRecorder) Create(ctx, apiToken, code interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockAPITokenRepo)(nil).Create), ctx, apiToken, code)
}

// DeleteById mocks base method.
func (m *MockAPITokenRepo) DeleteById(ctx context.Context, token string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteById", ctx, token)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteById indicates an expected call of DeleteById.
func (mr *MockAPITokenRepoMockRecorder) DeleteById(ctx, token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteById", reflect.TypeOf((*MockAPITokenRepo)(nil).DeleteById), ctx, token)
}

// DeleteByUserIDs mocks base method.
func (m *MockAPITokenRepo) DeleteByUserIDs(ctx context.Context, userIDs []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteByUserIDs", ctx, userIDs)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteByUserIDs indicates an expected call of DeleteByUserIDs.
func (mr *MockAPITokenRepoMockRecorder) DeleteByUserIDs(ctx, userIDs interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteByUserIDs", reflect.TypeOf((*MockAPITokenRepo)(nil).DeleteByUserIDs), ctx, userIDs)
}

// GenerateCode mocks base method.
func (m *MockAPITokenRepo) GenerateCode(userID string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateCode", userID)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateCode indicates an expected call of GenerateCode.
func (mr *MockAPITokenRepoMockRecorder) GenerateCode(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateCode", reflect.TypeOf((*MockAPITokenRepo)(nil).GenerateCode), userID)
}

// GetByID mocks base method.
func (m *MockAPITokenRepo) GetByID(ctx context.Context, id string) (*entity.APIToken, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(*entity.APIToken)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockAPITokenRepoMockRecorder) GetByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockAPITokenRepo)(nil).GetByID), ctx, id)
}

// GetByToken mocks base method.
func (m *MockAPITokenRepo) GetByToken(ctx context.Context, token string) (*entity.APIToken, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByToken", ctx, token)
	ret0, _ := ret[0].(*entity.APIToken)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByToken indicates an expected call of GetByToken.
func (mr *MockAPITokenRepoMockRecorder) GetByToken(ctx, token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByToken", reflect.TypeOf((*MockAPITokenRepo)(nil).GetByToken), ctx, token)
}

// GetByUserID mocks base method.
func (m *MockAPITokenRepo) GetByUserID(ctx context.Context, userID string) ([]*entity.APIToken, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByUserID", ctx, userID)
	ret0, _ := ret[0].([]*entity.APIToken)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByUserID indicates an expected call of GetByUserID.
func (mr *MockAPITokenRepoMockRecorder) GetByUserID(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByUserID", reflect.TypeOf((*MockAPITokenRepo)(nil).GetByUserID), ctx, userID)
}

// UpdateLastActivity mocks base method.
func (m *MockAPITokenRepo) UpdateLastActivity(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateLastActivity", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateLastActivity indicates an expected call of UpdateLastActivity.
func (mr *MockAPITokenRepoMockRecorder) UpdateLastActivity(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateLastActivity", reflect.TypeOf((*MockAPITokenRepo)(nil).UpdateLastActivity), ctx, id)
}
