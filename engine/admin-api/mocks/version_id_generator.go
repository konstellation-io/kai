// Code generated by MockGen. DO NOT EDIT.
// Source: id_generator.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockIDGenerator is a mock of IDGenerator interface.
type MockIDGenerator struct {
	ctrl     *gomock.Controller
	recorder *MockIDGeneratorMockRecorder
}

// MockIDGeneratorMockRecorder is the mock recorder for MockIDGenerator.
type MockIDGeneratorMockRecorder struct {
	mock *MockIDGenerator
}

// NewMockIDGenerator creates a new mock instance.
func NewMockIDGenerator(ctrl *gomock.Controller) *MockIDGenerator {
	mock := &MockIDGenerator{ctrl: ctrl}
	mock.recorder = &MockIDGeneratorMockRecorder{mock}

	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIDGenerator) EXPECT() *MockIDGeneratorMockRecorder {
	return m.recorder
}

// NewID mocks base method.
func (m *MockIDGenerator) NewID() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewID")
	ret0, _ := ret[0].(string)

	return ret0
}

// NewID indicates an expected call of NewID.
func (mr *MockIDGeneratorMockRecorder) NewID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewID", reflect.TypeOf((*MockIDGenerator)(nil).NewID))
}
