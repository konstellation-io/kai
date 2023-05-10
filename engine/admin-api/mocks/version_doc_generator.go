// Code generated by MockGen. DO NOT EDIT.
// Source: doc_generator.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockDocGenerator is a mock of DocGenerator interface.
type MockDocGenerator struct {
	ctrl     *gomock.Controller
	recorder *MockDocGeneratorMockRecorder
}

// MockDocGeneratorMockRecorder is the mock recorder for MockDocGenerator.
type MockDocGeneratorMockRecorder struct {
	mock *MockDocGenerator
}

// NewMockDocGenerator creates a new mock instance.
func NewMockDocGenerator(ctrl *gomock.Controller) *MockDocGenerator {
	mock := &MockDocGenerator{ctrl: ctrl}
	mock.recorder = &MockDocGeneratorMockRecorder{mock}

	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDocGenerator) EXPECT() *MockDocGeneratorMockRecorder {
	return m.recorder
}

// Generate mocks base method.
func (m *MockDocGenerator) Generate(versionName, docFolder string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Generate", versionName, docFolder)
	ret0, _ := ret[0].(error)

	return ret0
}

// Generate indicates an expected call of Generate.
func (mr *MockDocGeneratorMockRecorder) Generate(versionName, docFolder interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Generate", reflect.TypeOf((*MockDocGenerator)(nil).Generate), versionName, docFolder)
}
