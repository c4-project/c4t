// Code generated by mockery v2.3.0. DO NOT EDIT.

package mocks

import (
	context "context"

	"github.com/MattWindsor91/c4t/internal/model/service/fuzzer"

	mock "github.com/stretchr/testify/mock"
)

// SingleFuzzer is an autogenerated mock type for the SingleFuzzer type
type SingleFuzzer struct {
	mock.Mock
}

// Fuzz provides a mock function with given fields: _a0, _a1
func (_m *SingleFuzzer) Fuzz(_a0 context.Context, _a1 fuzzer.Job) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, fuzzer.Job) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
