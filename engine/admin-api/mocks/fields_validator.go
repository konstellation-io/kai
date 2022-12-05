// Code generated by MockGen. DO NOT EDIT.
// Source: fields_validator.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockFieldsValidator is a mock of FieldsValidator interface.
type MockFieldsValidator struct {
	ctrl     *gomock.Controller
	recorder *MockFieldsValidatorMockRecorder
}

// MockFieldsValidatorMockRecorder is the mock recorder for MockFieldsValidator.
type MockFieldsValidatorMockRecorder struct {
	mock *MockFieldsValidator
}

// NewMockFieldsValidator creates a new mock instance.
func NewMockFieldsValidator(ctrl *gomock.Controller) *MockFieldsValidator {
	mock := &MockFieldsValidator{ctrl: ctrl}
	mock.recorder = &MockFieldsValidatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFieldsValidator) EXPECT() *MockFieldsValidatorMockRecorder {
	return m.recorder
}

// Run mocks base method.
func (m *MockFieldsValidator) Run(yaml interface{}) []error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Run", yaml)
	ret0, _ := ret[0].([]error)
	return ret0
}

// Run indicates an expected call of Run.
func (mr *MockFieldsValidatorMockRecorder) Run(yaml interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*MockFieldsValidator)(nil).Run), yaml)
}