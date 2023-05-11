// Code generated by MockGen. DO NOT EDIT.
// Source: nodelog.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	entity "github.com/konstellation-io/kai/engine/admin-api/domain/entity"
)

// MockNodeLogRepository is a mock of NodeLogRepository interface.
type MockNodeLogRepository struct {
	ctrl     *gomock.Controller
	recorder *MockNodeLogRepositoryMockRecorder
}

// MockNodeLogRepositoryMockRecorder is the mock recorder for MockNodeLogRepository.
type MockNodeLogRepositoryMockRecorder struct {
	mock *MockNodeLogRepository
}

// NewMockNodeLogRepository creates a new mock instance.
func NewMockNodeLogRepository(ctrl *gomock.Controller) *MockNodeLogRepository {
	mock := &MockNodeLogRepository{ctrl: ctrl}
	mock.recorder = &MockNodeLogRepositoryMockRecorder{mock}

	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNodeLogRepository) EXPECT() *MockNodeLogRepositoryMockRecorder {
	return m.recorder
}

// CreateIndexes mocks base method.
func (m *MockNodeLogRepository) CreateIndexes(ctx context.Context, runtimeID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateIndexes", ctx, runtimeID)
	ret0, _ := ret[0].(error)

	return ret0
}

// CreateIndexes indicates an expected call of CreateIndexes.
func (mr *MockNodeLogRepositoryMockRecorder) CreateIndexes(ctx, runtimeID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateIndexes", reflect.TypeOf((*MockNodeLogRepository)(nil).CreateIndexes), ctx, runtimeID)
}

// PaginatedSearch mocks base method.
func (m *MockNodeLogRepository) PaginatedSearch(ctx context.Context, runtimeID string, searchOpts *entity.SearchLogsOptions) (*entity.SearchLogsResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PaginatedSearch", ctx, runtimeID, searchOpts)
	ret0, _ := ret[0].(*entity.SearchLogsResult)
	ret1, _ := ret[1].(error)

	return ret0, ret1
}

// PaginatedSearch indicates an expected call of PaginatedSearch.
func (mr *MockNodeLogRepositoryMockRecorder) PaginatedSearch(ctx, runtimeID, searchOpts interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PaginatedSearch", reflect.TypeOf((*MockNodeLogRepository)(nil).PaginatedSearch), ctx, runtimeID, searchOpts)
}

// WatchNodeLogs mocks base method.
func (m *MockNodeLogRepository) WatchNodeLogs(ctx context.Context, runtimeID, versionName string, filters entity.LogFilters) (<-chan *entity.NodeLog, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WatchNodeLogs", ctx, runtimeID, versionName, filters)
	ret0, _ := ret[0].(<-chan *entity.NodeLog)
	ret1, _ := ret[1].(error)

	return ret0, ret1
}

// WatchNodeLogs indicates an expected call of WatchNodeLogs.
func (mr *MockNodeLogRepositoryMockRecorder) WatchNodeLogs(ctx, runtimeID, versionName, filters interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WatchNodeLogs", reflect.TypeOf((*MockNodeLogRepository)(nil).WatchNodeLogs), ctx, runtimeID, versionName, filters)
}
