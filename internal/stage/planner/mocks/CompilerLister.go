// Code generated by mockery v2.1.0. DO NOT EDIT.

package mocks

import (
	context "context"

	compiler "github.com/MattWindsor91/act-tester/internal/model/service/compiler"

	id "github.com/MattWindsor91/act-tester/internal/model/id"

	mock "github.com/stretchr/testify/mock"
)

// CompilerLister is an autogenerated mock type for the CompilerLister type
type CompilerLister struct {
	mock.Mock
}

// ListCompilers provides a mock function with given fields: ctx, mid
func (_m *CompilerLister) ListCompilers(ctx context.Context, mid id.ID) (map[string]compiler.Compiler, error) {
	ret := _m.Called(ctx, mid)

	var r0 map[string]compiler.Compiler
	if rf, ok := ret.Get(0).(func(context.Context, id.ID) map[string]compiler.Compiler); ok {
		r0 = rf(ctx, mid)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]compiler.Compiler)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, id.ID) error); ok {
		r1 = rf(ctx, mid)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
