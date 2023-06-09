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

// MockVersionRepo is a mock of VersionRepo interface.
type MockVersionRepo struct {
	ctrl     *gomock.Controller
	recorder *MockVersionRepoMockRecorder
}

// MockVersionRepoMockRecorder is the mock recorder for MockVersionRepo.
type MockVersionRepoMockRecorder struct {
	mock *MockVersionRepo
}

// NewMockVersionRepo creates a new mock instance.
func NewMockVersionRepo(ctrl *gomock.Controller) *MockVersionRepo {
	mock := &MockVersionRepo{ctrl: ctrl}
	mock.recorder = &MockVersionRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockVersionRepo) EXPECT() *MockVersionRepoMockRecorder {
	return m.recorder
}

// ClearPublishedVersion mocks base method.
func (m *MockVersionRepo) ClearPublishedVersion(ctx context.Context, productID string) (*entity.Version, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ClearPublishedVersion", ctx, productID)
	ret0, _ := ret[0].(*entity.Version)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ClearPublishedVersion indicates an expected call of ClearPublishedVersion.
func (mr *MockVersionRepoMockRecorder) ClearPublishedVersion(ctx, productID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ClearPublishedVersion", reflect.TypeOf((*MockVersionRepo)(nil).ClearPublishedVersion), ctx, productID)
}

// Create mocks base method.
func (m *MockVersionRepo) Create(userID, productID string, version *entity.Version) (*entity.Version, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", userID, productID, version)
	ret0, _ := ret[0].(*entity.Version)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockVersionRepoMockRecorder) Create(userID, productID, version interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockVersionRepo)(nil).Create), userID, productID, version)
}

// CreateIndexes mocks base method.
func (m *MockVersionRepo) CreateIndexes(ctx context.Context, productID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateIndexes", ctx, productID)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateIndexes indicates an expected call of CreateIndexes.
func (mr *MockVersionRepoMockRecorder) CreateIndexes(ctx, productID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateIndexes", reflect.TypeOf((*MockVersionRepo)(nil).CreateIndexes), ctx, productID)
}

// GetByID mocks base method.
func (m *MockVersionRepo) GetByID(productID, versionID string) (*entity.Version, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", productID, versionID)
	ret0, _ := ret[0].(*entity.Version)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockVersionRepoMockRecorder) GetByID(productID, versionID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockVersionRepo)(nil).GetByID), productID, versionID)
}

// GetByName mocks base method.
func (m *MockVersionRepo) GetByName(ctx context.Context, productID, name string) (*entity.Version, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByName", ctx, productID, name)
	ret0, _ := ret[0].(*entity.Version)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByName indicates an expected call of GetByName.
func (mr *MockVersionRepoMockRecorder) GetByName(ctx, productID, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByName", reflect.TypeOf((*MockVersionRepo)(nil).GetByName), ctx, productID, name)
}

// GetByProduct mocks base method.
func (m *MockVersionRepo) GetByProduct(ctx context.Context, productID string) ([]*entity.Version, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByProduct", ctx, productID)
	ret0, _ := ret[0].([]*entity.Version)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByProduct indicates an expected call of GetByProduct.
func (mr *MockVersionRepoMockRecorder) GetByProduct(ctx, productID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByProduct", reflect.TypeOf((*MockVersionRepo)(nil).GetByProduct), ctx, productID)
}

// SetErrors mocks base method.
func (m *MockVersionRepo) SetErrors(ctx context.Context, productID string, version *entity.Version, errorMessages []string) (*entity.Version, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetErrors", ctx, productID, version, errorMessages)
	ret0, _ := ret[0].(*entity.Version)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SetErrors indicates an expected call of SetErrors.
func (mr *MockVersionRepoMockRecorder) SetErrors(ctx, productID, version, errorMessages interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetErrors", reflect.TypeOf((*MockVersionRepo)(nil).SetErrors), ctx, productID, version, errorMessages)
}

// SetStatus mocks base method.
func (m *MockVersionRepo) SetStatus(ctx context.Context, productID, versionID string, status entity.VersionStatus) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetStatus", ctx, productID, versionID, status)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetStatus indicates an expected call of SetStatus.
func (mr *MockVersionRepoMockRecorder) SetStatus(ctx, productID, versionID, status interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetStatus", reflect.TypeOf((*MockVersionRepo)(nil).SetStatus), ctx, productID, versionID, status)
}

// Update mocks base method.
func (m *MockVersionRepo) Update(productID string, version *entity.Version) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", productID, version)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockVersionRepoMockRecorder) Update(productID, version interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockVersionRepo)(nil).Update), productID, version)
}

// UploadKRTFile mocks base method.
func (m *MockVersionRepo) UploadKRTFile(productID string, version *entity.Version, file string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadKRTFile", productID, version, file)
	ret0, _ := ret[0].(error)
	return ret0
}

// UploadKRTFile indicates an expected call of UploadKRTFile.
func (mr *MockVersionRepoMockRecorder) UploadKRTFile(productID, version, file interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadKRTFile", reflect.TypeOf((*MockVersionRepo)(nil).UploadKRTFile), productID, version, file)
}
