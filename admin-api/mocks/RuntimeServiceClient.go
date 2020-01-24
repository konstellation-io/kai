// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import context "context"
import grpc "google.golang.org/grpc"
import mock "github.com/stretchr/testify/mock"
import runtimepb "gitlab.com/konstellation/konstellation-ce/kre/admin-api/runtimepb"

// RuntimeServiceClient is an autogenerated mock type for the RuntimeServiceClient type
type RuntimeServiceClient struct {
	mock.Mock
}

// ActivateVersion provides a mock function with given fields: ctx, in, opts
func (_m *RuntimeServiceClient) ActivateVersion(ctx context.Context, in *runtimepb.ActivateVersionRequest, opts ...grpc.CallOption) (*runtimepb.ActivateVersionResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *runtimepb.ActivateVersionResponse
	if rf, ok := ret.Get(0).(func(context.Context, *runtimepb.ActivateVersionRequest, ...grpc.CallOption) *runtimepb.ActivateVersionResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*runtimepb.ActivateVersionResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *runtimepb.ActivateVersionRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeactivateVersion provides a mock function with given fields: ctx, in, opts
func (_m *RuntimeServiceClient) DeactivateVersion(ctx context.Context, in *runtimepb.DeactivateVersionRequest, opts ...grpc.CallOption) (*runtimepb.DeactivateVersionResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *runtimepb.DeactivateVersionResponse
	if rf, ok := ret.Get(0).(func(context.Context, *runtimepb.DeactivateVersionRequest, ...grpc.CallOption) *runtimepb.DeactivateVersionResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*runtimepb.DeactivateVersionResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *runtimepb.DeactivateVersionRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeployVersion provides a mock function with given fields: ctx, in, opts
func (_m *RuntimeServiceClient) DeployVersion(ctx context.Context, in *runtimepb.DeployVersionRequest, opts ...grpc.CallOption) (*runtimepb.DeployVersionResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *runtimepb.DeployVersionResponse
	if rf, ok := ret.Get(0).(func(context.Context, *runtimepb.DeployVersionRequest, ...grpc.CallOption) *runtimepb.DeployVersionResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*runtimepb.DeployVersionResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *runtimepb.DeployVersionRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// StopVersion provides a mock function with given fields: ctx, in, opts
func (_m *RuntimeServiceClient) StopVersion(ctx context.Context, in *runtimepb.StopVersionRequest, opts ...grpc.CallOption) (*runtimepb.StopVersionResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *runtimepb.StopVersionResponse
	if rf, ok := ret.Get(0).(func(context.Context, *runtimepb.StopVersionRequest, ...grpc.CallOption) *runtimepb.StopVersionResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*runtimepb.StopVersionResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *runtimepb.StopVersionRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateVersionConfig provides a mock function with given fields: ctx, in, opts
func (_m *RuntimeServiceClient) UpdateVersionConfig(ctx context.Context, in *runtimepb.UpdateVersionConfigRequest, opts ...grpc.CallOption) (*runtimepb.UpdateVersionConfigResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *runtimepb.UpdateVersionConfigResponse
	if rf, ok := ret.Get(0).(func(context.Context, *runtimepb.UpdateVersionConfigRequest, ...grpc.CallOption) *runtimepb.UpdateVersionConfigResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*runtimepb.UpdateVersionConfigResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *runtimepb.UpdateVersionConfigRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// WatchNodeLogs provides a mock function with given fields: ctx, in, opts
func (_m *RuntimeServiceClient) WatchNodeLogs(ctx context.Context, in *runtimepb.WatchNodeLogsRequest, opts ...grpc.CallOption) (runtimepb.RuntimeService_WatchNodeLogsClient, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 runtimepb.RuntimeService_WatchNodeLogsClient
	if rf, ok := ret.Get(0).(func(context.Context, *runtimepb.WatchNodeLogsRequest, ...grpc.CallOption) runtimepb.RuntimeService_WatchNodeLogsClient); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(runtimepb.RuntimeService_WatchNodeLogsClient)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *runtimepb.WatchNodeLogsRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
