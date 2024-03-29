// Code generated by MockGen. DO NOT EDIT.
// Source: logs.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	entity "github.com/konstellation-io/kai/engine/admin-api/domain/entity"
)

// MockLogsService is a mock of LogsService interface.
type MockLogsService struct {
	ctrl     *gomock.Controller
	recorder *MockLogsServiceMockRecorder
}

// MockLogsServiceMockRecorder is the mock recorder for MockLogsService.
type MockLogsServiceMockRecorder struct {
	mock *MockLogsService
}

// NewMockLogsService creates a new mock instance.
func NewMockLogsService(ctrl *gomock.Controller) *MockLogsService {
	mock := &MockLogsService{ctrl: ctrl}
	mock.recorder = &MockLogsServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLogsService) EXPECT() *MockLogsServiceMockRecorder {
	return m.recorder
}

// GetLogs mocks base method.
func (m *MockLogsService) GetLogs(logFilters entity.LogFilters) ([]*entity.Log, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLogs", logFilters)
	ret0, _ := ret[0].([]*entity.Log)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLogs indicates an expected call of GetLogs.
func (mr *MockLogsServiceMockRecorder) GetLogs(logFilters interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLogs", reflect.TypeOf((*MockLogsService)(nil).GetLogs), logFilters)
}
