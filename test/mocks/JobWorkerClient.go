// Code generated by mockery v2.3.0. DO NOT EDIT.

package mocks

import (
	context "context"

	grpc "google.golang.org/grpc"

	mock "github.com/stretchr/testify/mock"

	proto "github.com/andrejtokarcik/jobworker/proto"
)

// JobWorkerClient is an autogenerated mock type for the JobWorkerClient type
type JobWorkerClient struct {
	mock.Mock
}

// GetJob provides a mock function with given fields: ctx, in, opts
func (_m *JobWorkerClient) GetJob(ctx context.Context, in *proto.GetJobRequest, opts ...grpc.CallOption) (*proto.GetJobResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *proto.GetJobResponse
	if rf, ok := ret.Get(0).(func(context.Context, *proto.GetJobRequest, ...grpc.CallOption) *proto.GetJobResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*proto.GetJobResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *proto.GetJobRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// StartJob provides a mock function with given fields: ctx, in, opts
func (_m *JobWorkerClient) StartJob(ctx context.Context, in *proto.StartJobRequest, opts ...grpc.CallOption) (*proto.StartJobResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *proto.StartJobResponse
	if rf, ok := ret.Get(0).(func(context.Context, *proto.StartJobRequest, ...grpc.CallOption) *proto.StartJobResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*proto.StartJobResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *proto.StartJobRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// StopJob provides a mock function with given fields: ctx, in, opts
func (_m *JobWorkerClient) StopJob(ctx context.Context, in *proto.StopJobRequest, opts ...grpc.CallOption) (*proto.StopJobResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *proto.StopJobResponse
	if rf, ok := ret.Get(0).(func(context.Context, *proto.StopJobRequest, ...grpc.CallOption) *proto.StopJobResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*proto.StopJobResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *proto.StopJobRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
