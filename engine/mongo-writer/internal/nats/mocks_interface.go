// Code generated by MockGen. DO NOT EDIT.
// Source: interface.go

// Package nats is a generated GoMock package.
package nats

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	nats "github.com/nats-io/nats.go"
)

// MockManager is a mock of Manager interface.
type MockManager struct {
	ctrl     *gomock.Controller
	recorder *MockManagerMockRecorder
}

// MockManagerMockRecorder is the mock recorder for MockManager.
type MockManagerMockRecorder struct {
	mock *MockManager
}

// NewMockManager creates a new mock instance.
func NewMockManager(ctrl *gomock.Controller) *MockManager {
	mock := &MockManager{ctrl: ctrl}
	mock.recorder = &MockManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockManager) EXPECT() *MockManagerMockRecorder {
	return m.recorder
}

// Connect mocks base method.
func (m *MockManager) Connect() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Connect")
	ret0, _ := ret[0].(error)
	return ret0
}

// Connect indicates an expected call of Connect.
func (mr *MockManagerMockRecorder) Connect() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Connect", reflect.TypeOf((*MockManager)(nil).Connect))
}

// Disconnect mocks base method.
func (m *MockManager) Disconnect() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Disconnect")
}

// Disconnect indicates an expected call of Disconnect.
func (mr *MockManagerMockRecorder) Disconnect() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Disconnect", reflect.TypeOf((*MockManager)(nil).Disconnect))
}

// IncreaseTotalMsgs mocks base method.
func (m *MockManager) IncreaseTotalMsgs(amount int64) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "IncreaseTotalMsgs", amount)
}

// IncreaseTotalMsgs indicates an expected call of IncreaseTotalMsgs.
func (mr *MockManagerMockRecorder) IncreaseTotalMsgs(amount interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IncreaseTotalMsgs", reflect.TypeOf((*MockManager)(nil).IncreaseTotalMsgs), amount)
}

// SubscribeToChannel mocks base method.
func (m *MockManager) SubscribeToChannel(channel string) chan *nats.Msg {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubscribeToChannel", channel)
	ret0, _ := ret[0].(chan *nats.Msg)
	return ret0
}

// SubscribeToChannel indicates an expected call of SubscribeToChannel.
func (mr *MockManagerMockRecorder) SubscribeToChannel(channel interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubscribeToChannel", reflect.TypeOf((*MockManager)(nil).SubscribeToChannel), channel)
}

// TotalMsgs mocks base method.
func (m *MockManager) TotalMsgs() int64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TotalMsgs")
	ret0, _ := ret[0].(int64)
	return ret0
}

// TotalMsgs indicates an expected call of TotalMsgs.
func (mr *MockManagerMockRecorder) TotalMsgs() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TotalMsgs", reflect.TypeOf((*MockManager)(nil).TotalMsgs))
}
