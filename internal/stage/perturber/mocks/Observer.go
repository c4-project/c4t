// Code generated by mockery v2.3.0. DO NOT EDIT.

package mocks

import (
	compiler "github.com/c4-project/c4t/internal/model/service/compiler"
	builder "github.com/c4-project/c4t/internal/subject/corpus/builder"

	mock "github.com/stretchr/testify/mock"

	perturber "github.com/c4-project/c4t/internal/stage/perturber"
)

// Observer is an autogenerated mock type for the Observer type
type Observer struct {
	mock.Mock
}

// OnBuild provides a mock function with given fields: _a0
func (_m *Observer) OnBuild(_a0 builder.Message) {
	_m.Called(_a0)
}

// OnCompilerConfig provides a mock function with given fields: _a0
func (_m *Observer) OnCompilerConfig(_a0 compiler.Message) {
	_m.Called(_a0)
}

// OnPerturb provides a mock function with given fields: m
func (_m *Observer) OnPerturb(m perturber.Message) {
	_m.Called(m)
}
