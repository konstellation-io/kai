// Code generated by MockGen. DO NOT EDIT.
// Source: metric.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
	entity "github.com/konstellation-io/kai/engine/admin-api/domain/entity"
)

// MockMetricRepo is a mock of MetricRepo interface.
type MockMetricRepo struct {
	ctrl     *gomock.Controller
	recorder *MockMetricRepoMockRecorder
}

// MockMetricRepoMockRecorder is the mock recorder for MockMetricRepo.
type MockMetricRepoMockRecorder struct {
	mock *MockMetricRepo
}

// NewMockMetricRepo creates a new mock instance.
func NewMockMetricRepo(ctrl *gomock.Controller) *MockMetricRepo {
	mock := &MockMetricRepo{ctrl: ctrl}
	mock.recorder = &MockMetricRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMetricRepo) EXPECT() *MockMetricRepoMockRecorder {
	return m.recorder
}

// CreateIndexes mocks base method.
func (m *MockMetricRepo) CreateIndexes(ctx context.Context, runtimeID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateIndexes", ctx, runtimeID)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateIndexes indicates an expected call of CreateIndexes.
func (mr *MockMetricRepoMockRecorder) CreateIndexes(ctx, runtimeID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateIndexes", reflect.TypeOf((*MockMetricRepo)(nil).CreateIndexes), ctx, runtimeID)
}

// GetMetrics mocks base method.
func (m *MockMetricRepo) GetMetrics(ctx context.Context, startDate, endDate time.Time, runtimeID, versionTag string) ([]entity.ClassificationMetric, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMetrics", ctx, startDate, endDate, runtimeID, versionTag)
	ret0, _ := ret[0].([]entity.ClassificationMetric)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMetrics indicates an expected call of GetMetrics.
func (mr *MockMetricRepoMockRecorder) GetMetrics(ctx, startDate, endDate, runtimeID, versionTag interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMetrics", reflect.TypeOf((*MockMetricRepo)(nil).GetMetrics), ctx, startDate, endDate, runtimeID, versionTag)
}
