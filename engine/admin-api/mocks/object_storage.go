// Code generated by mockery v2.32.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockObjectStorage is an autogenerated mock type for the ObjectStorage type
type MockObjectStorage struct {
	mock.Mock
}

type MockObjectStorage_Expecter struct {
	mock *mock.Mock
}

func (_m *MockObjectStorage) EXPECT() *MockObjectStorage_Expecter {
	return &MockObjectStorage_Expecter{mock: &_m.Mock}
}

// CreateBucket provides a mock function with given fields: ctx, bucket
func (_m *MockObjectStorage) CreateBucket(ctx context.Context, bucket string) error {
	ret := _m.Called(ctx, bucket)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, bucket)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockObjectStorage_CreateBucket_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateBucket'
type MockObjectStorage_CreateBucket_Call struct {
	*mock.Call
}

// CreateBucket is a helper method to define mock.On call
//   - ctx context.Context
//   - bucket string
func (_e *MockObjectStorage_Expecter) CreateBucket(ctx interface{}, bucket interface{}) *MockObjectStorage_CreateBucket_Call {
	return &MockObjectStorage_CreateBucket_Call{Call: _e.mock.On("CreateBucket", ctx, bucket)}
}

func (_c *MockObjectStorage_CreateBucket_Call) Run(run func(ctx context.Context, bucket string)) *MockObjectStorage_CreateBucket_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockObjectStorage_CreateBucket_Call) Return(_a0 error) *MockObjectStorage_CreateBucket_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockObjectStorage_CreateBucket_Call) RunAndReturn(run func(context.Context, string) error) *MockObjectStorage_CreateBucket_Call {
	_c.Call.Return(run)
	return _c
}

// CreateBucketPolicy provides a mock function with given fields: ctx, bucket
func (_m *MockObjectStorage) CreateBucketPolicy(ctx context.Context, bucket string) (string, error) {
	ret := _m.Called(ctx, bucket)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (string, error)); ok {
		return rf(ctx, bucket)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(ctx, bucket)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, bucket)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockObjectStorage_CreateBucketPolicy_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateBucketPolicy'
type MockObjectStorage_CreateBucketPolicy_Call struct {
	*mock.Call
}

// CreateBucketPolicy is a helper method to define mock.On call
//   - ctx context.Context
//   - bucket string
func (_e *MockObjectStorage_Expecter) CreateBucketPolicy(ctx interface{}, bucket interface{}) *MockObjectStorage_CreateBucketPolicy_Call {
	return &MockObjectStorage_CreateBucketPolicy_Call{Call: _e.mock.On("CreateBucketPolicy", ctx, bucket)}
}

func (_c *MockObjectStorage_CreateBucketPolicy_Call) Run(run func(ctx context.Context, bucket string)) *MockObjectStorage_CreateBucketPolicy_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockObjectStorage_CreateBucketPolicy_Call) Return(_a0 string, _a1 error) *MockObjectStorage_CreateBucketPolicy_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockObjectStorage_CreateBucketPolicy_Call) RunAndReturn(run func(context.Context, string) (string, error)) *MockObjectStorage_CreateBucketPolicy_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteBucket provides a mock function with given fields: ctx, bucket
func (_m *MockObjectStorage) DeleteBucket(ctx context.Context, bucket string) error {
	ret := _m.Called(ctx, bucket)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, bucket)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockObjectStorage_DeleteBucket_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteBucket'
type MockObjectStorage_DeleteBucket_Call struct {
	*mock.Call
}

// DeleteBucket is a helper method to define mock.On call
//   - ctx context.Context
//   - bucket string
func (_e *MockObjectStorage_Expecter) DeleteBucket(ctx interface{}, bucket interface{}) *MockObjectStorage_DeleteBucket_Call {
	return &MockObjectStorage_DeleteBucket_Call{Call: _e.mock.On("DeleteBucket", ctx, bucket)}
}

func (_c *MockObjectStorage_DeleteBucket_Call) Run(run func(ctx context.Context, bucket string)) *MockObjectStorage_DeleteBucket_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockObjectStorage_DeleteBucket_Call) Return(_a0 error) *MockObjectStorage_DeleteBucket_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockObjectStorage_DeleteBucket_Call) RunAndReturn(run func(context.Context, string) error) *MockObjectStorage_DeleteBucket_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteBucketPolicy provides a mock function with given fields: ctx, policyName
func (_m *MockObjectStorage) DeleteBucketPolicy(ctx context.Context, policyName string) error {
	ret := _m.Called(ctx, policyName)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, policyName)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockObjectStorage_DeleteBucketPolicy_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteBucketPolicy'
type MockObjectStorage_DeleteBucketPolicy_Call struct {
	*mock.Call
}

// DeleteBucketPolicy is a helper method to define mock.On call
//   - ctx context.Context
//   - policyName string
func (_e *MockObjectStorage_Expecter) DeleteBucketPolicy(ctx interface{}, policyName interface{}) *MockObjectStorage_DeleteBucketPolicy_Call {
	return &MockObjectStorage_DeleteBucketPolicy_Call{Call: _e.mock.On("DeleteBucketPolicy", ctx, policyName)}
}

func (_c *MockObjectStorage_DeleteBucketPolicy_Call) Run(run func(ctx context.Context, policyName string)) *MockObjectStorage_DeleteBucketPolicy_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockObjectStorage_DeleteBucketPolicy_Call) Return(_a0 error) *MockObjectStorage_DeleteBucketPolicy_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockObjectStorage_DeleteBucketPolicy_Call) RunAndReturn(run func(context.Context, string) error) *MockObjectStorage_DeleteBucketPolicy_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteImageSources provides a mock function with given fields: ctx, product, image
func (_m *MockObjectStorage) DeleteImageSources(ctx context.Context, product string, image string) error {
	ret := _m.Called(ctx, product, image)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, product, image)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockObjectStorage_DeleteImageSources_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteImageSources'
type MockObjectStorage_DeleteImageSources_Call struct {
	*mock.Call
}

// DeleteImageSources is a helper method to define mock.On call
//   - ctx context.Context
//   - product string
//   - image string
func (_e *MockObjectStorage_Expecter) DeleteImageSources(ctx interface{}, product interface{}, image interface{}) *MockObjectStorage_DeleteImageSources_Call {
	return &MockObjectStorage_DeleteImageSources_Call{Call: _e.mock.On("DeleteImageSources", ctx, product, image)}
}

func (_c *MockObjectStorage_DeleteImageSources_Call) Run(run func(ctx context.Context, product string, image string)) *MockObjectStorage_DeleteImageSources_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *MockObjectStorage_DeleteImageSources_Call) Return(_a0 error) *MockObjectStorage_DeleteImageSources_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockObjectStorage_DeleteImageSources_Call) RunAndReturn(run func(context.Context, string, string) error) *MockObjectStorage_DeleteImageSources_Call {
	_c.Call.Return(run)
	return _c
}

// UploadImageSources provides a mock function with given fields: ctx, product, image, sources
func (_m *MockObjectStorage) UploadImageSources(ctx context.Context, product string, image string, sources []byte) error {
	ret := _m.Called(ctx, product, image, sources)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, []byte) error); ok {
		r0 = rf(ctx, product, image, sources)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockObjectStorage_UploadImageSources_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UploadImageSources'
type MockObjectStorage_UploadImageSources_Call struct {
	*mock.Call
}

// UploadImageSources is a helper method to define mock.On call
//   - ctx context.Context
//   - product string
//   - image string
//   - sources []byte
func (_e *MockObjectStorage_Expecter) UploadImageSources(ctx interface{}, product interface{}, image interface{}, sources interface{}) *MockObjectStorage_UploadImageSources_Call {
	return &MockObjectStorage_UploadImageSources_Call{Call: _e.mock.On("UploadImageSources", ctx, product, image, sources)}
}

func (_c *MockObjectStorage_UploadImageSources_Call) Run(run func(ctx context.Context, product string, image string, sources []byte)) *MockObjectStorage_UploadImageSources_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string), args[3].([]byte))
	})
	return _c
}

func (_c *MockObjectStorage_UploadImageSources_Call) Return(_a0 error) *MockObjectStorage_UploadImageSources_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockObjectStorage_UploadImageSources_Call) RunAndReturn(run func(context.Context, string, string, []byte) error) *MockObjectStorage_UploadImageSources_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockObjectStorage creates a new instance of MockObjectStorage. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockObjectStorage(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockObjectStorage {
	mock := &MockObjectStorage{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
