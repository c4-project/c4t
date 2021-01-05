// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	context "context"

	backend "github.com/c4-project/c4t/internal/model/service/backend"

	mock "github.com/stretchr/testify/mock"

	service "github.com/c4-project/c4t/internal/model/service"
)

// Resolver is an autogenerated mock type for the Resolver type
type Resolver struct {
	mock.Mock
}

// Probe provides a mock function with given fields: ctx, sr
func (_m *Resolver) Probe(ctx context.Context, sr service.Runner) ([]backend.NamedSpec, error) {
	ret := _m.Called(ctx, sr)

	var r0 []backend.NamedSpec
	if rf, ok := ret.Get(0).(func(context.Context, service.Runner) []backend.NamedSpec); ok {
		r0 = rf(ctx, sr)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]backend.NamedSpec)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, service.Runner) error); ok {
		r1 = rf(ctx, sr)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Resolve provides a mock function with given fields: s
func (_m *Resolver) Resolve(s backend.Spec) (backend.Backend, error) {
	ret := _m.Called(s)

	var r0 backend.Backend
	if rf, ok := ret.Get(0).(func(backend.Spec) backend.Backend); ok {
		r0 = rf(s)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(backend.Backend)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(backend.Spec) error); ok {
		r1 = rf(s)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}