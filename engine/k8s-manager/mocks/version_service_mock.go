// Code generated by mockery v2.32.0. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/konstellation-io/kai/engine/k8s-manager/internal/domain"
	mock "github.com/stretchr/testify/mock"

	usecase "github.com/konstellation-io/kai/engine/k8s-manager/internal/application/usecase"
)

// VersionServiceMock is an autogenerated mock type for the VersionService type
type VersionServiceMock struct {
	mock.Mock
}

type VersionServiceMock_Expecter struct {
	mock *mock.Mock
}

func (_m *VersionServiceMock) EXPECT() *VersionServiceMock_Expecter {
	return &VersionServiceMock_Expecter{mock: &_m.Mock}
}

// PublishVersion provides a mock function with given fields: ctx, product, version
func (_m *VersionServiceMock) PublishVersion(ctx context.Context, product string, version string) (map[string]string, error) {
	ret := _m.Called(ctx, product, version)

	var r0 map[string]string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (map[string]string, error)); ok {
		return rf(ctx, product, version)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) map[string]string); ok {
		r0 = rf(ctx, product, version)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]string)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, product, version)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// VersionServiceMock_PublishVersion_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'PublishVersion'
type VersionServiceMock_PublishVersion_Call struct {
	*mock.Call
}

// PublishVersion is a helper method to define mock.On call
//   - ctx context.Context
//   - product string
//   - version string
func (_e *VersionServiceMock_Expecter) PublishVersion(ctx interface{}, product interface{}, version interface{}) *VersionServiceMock_PublishVersion_Call {
	return &VersionServiceMock_PublishVersion_Call{Call: _e.mock.On("PublishVersion", ctx, product, version)}
}

func (_c *VersionServiceMock_PublishVersion_Call) Run(run func(ctx context.Context, product string, version string)) *VersionServiceMock_PublishVersion_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *VersionServiceMock_PublishVersion_Call) Return(_a0 map[string]string, _a1 error) *VersionServiceMock_PublishVersion_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *VersionServiceMock_PublishVersion_Call) RunAndReturn(run func(context.Context, string, string) (map[string]string, error)) *VersionServiceMock_PublishVersion_Call {
	_c.Call.Return(run)
	return _c
}

// StartVersion provides a mock function with given fields: ctx, version
func (_m *VersionServiceMock) StartVersion(ctx context.Context, version domain.Version) error {
	ret := _m.Called(ctx, version)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, domain.Version) error); ok {
		r0 = rf(ctx, version)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// VersionServiceMock_StartVersion_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'StartVersion'
type VersionServiceMock_StartVersion_Call struct {
	*mock.Call
}

// StartVersion is a helper method to define mock.On call
//   - ctx context.Context
//   - version domain.Version
func (_e *VersionServiceMock_Expecter) StartVersion(ctx interface{}, version interface{}) *VersionServiceMock_StartVersion_Call {
	return &VersionServiceMock_StartVersion_Call{Call: _e.mock.On("StartVersion", ctx, version)}
}

func (_c *VersionServiceMock_StartVersion_Call) Run(run func(ctx context.Context, version domain.Version)) *VersionServiceMock_StartVersion_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(domain.Version))
	})
	return _c
}

func (_c *VersionServiceMock_StartVersion_Call) Return(_a0 error) *VersionServiceMock_StartVersion_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *VersionServiceMock_StartVersion_Call) RunAndReturn(run func(context.Context, domain.Version) error) *VersionServiceMock_StartVersion_Call {
	_c.Call.Return(run)
	return _c
}

// StopVersion provides a mock function with given fields: ctx, params
func (_m *VersionServiceMock) StopVersion(ctx context.Context, params usecase.StopParams) error {
	ret := _m.Called(ctx, params)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, usecase.StopParams) error); ok {
		r0 = rf(ctx, params)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// VersionServiceMock_StopVersion_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'StopVersion'
type VersionServiceMock_StopVersion_Call struct {
	*mock.Call
}

// StopVersion is a helper method to define mock.On call
//   - ctx context.Context
//   - params usecase.StopParams
func (_e *VersionServiceMock_Expecter) StopVersion(ctx interface{}, params interface{}) *VersionServiceMock_StopVersion_Call {
	return &VersionServiceMock_StopVersion_Call{Call: _e.mock.On("StopVersion", ctx, params)}
}

func (_c *VersionServiceMock_StopVersion_Call) Run(run func(ctx context.Context, params usecase.StopParams)) *VersionServiceMock_StopVersion_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(usecase.StopParams))
	})
	return _c
}

func (_c *VersionServiceMock_StopVersion_Call) Return(_a0 error) *VersionServiceMock_StopVersion_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *VersionServiceMock_StopVersion_Call) RunAndReturn(run func(context.Context, usecase.StopParams) error) *VersionServiceMock_StopVersion_Call {
	_c.Call.Return(run)
	return _c
}

// UnpublishVersion provides a mock function with given fields: ctx, product, version
func (_m *VersionServiceMock) UnpublishVersion(ctx context.Context, product string, version string) error {
	ret := _m.Called(ctx, product, version)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, product, version)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// VersionServiceMock_UnpublishVersion_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UnpublishVersion'
type VersionServiceMock_UnpublishVersion_Call struct {
	*mock.Call
}

// UnpublishVersion is a helper method to define mock.On call
//   - ctx context.Context
//   - product string
//   - version string
func (_e *VersionServiceMock_Expecter) UnpublishVersion(ctx interface{}, product interface{}, version interface{}) *VersionServiceMock_UnpublishVersion_Call {
	return &VersionServiceMock_UnpublishVersion_Call{Call: _e.mock.On("UnpublishVersion", ctx, product, version)}
}

func (_c *VersionServiceMock_UnpublishVersion_Call) Run(run func(ctx context.Context, product string, version string)) *VersionServiceMock_UnpublishVersion_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *VersionServiceMock_UnpublishVersion_Call) Return(_a0 error) *VersionServiceMock_UnpublishVersion_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *VersionServiceMock_UnpublishVersion_Call) RunAndReturn(run func(context.Context, string, string) error) *VersionServiceMock_UnpublishVersion_Call {
	_c.Call.Return(run)
	return _c
}

// NewVersionServiceMock creates a new instance of VersionServiceMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewVersionServiceMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *VersionServiceMock {
	mock := &VersionServiceMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
