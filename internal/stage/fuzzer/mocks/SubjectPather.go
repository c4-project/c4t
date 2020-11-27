// Code generated by mockery v2.3.0. DO NOT EDIT.

package mocks

import (
	fuzzer "github.com/MattWindsor91/c4t/internal/stage/fuzzer"
	mock "github.com/stretchr/testify/mock"
)

// SubjectPather is an autogenerated mock type for the SubjectPather type
type SubjectPather struct {
	mock.Mock
}

// Prepare provides a mock function with given fields:
func (_m *SubjectPather) Prepare() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SubjectLitmus provides a mock function with given fields: sc
func (_m *SubjectPather) SubjectLitmus(sc fuzzer.SubjectCycle) string {
	ret := _m.Called(sc)

	var r0 string
	if rf, ok := ret.Get(0).(func(fuzzer.SubjectCycle) string); ok {
		r0 = rf(sc)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// SubjectTrace provides a mock function with given fields: sc
func (_m *SubjectPather) SubjectTrace(sc fuzzer.SubjectCycle) string {
	ret := _m.Called(sc)

	var r0 string
	if rf, ok := ret.Get(0).(func(fuzzer.SubjectCycle) string); ok {
		r0 = rf(sc)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}
