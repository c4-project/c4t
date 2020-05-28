// Code generated by mockery v2.0.0-alpha.2. DO NOT EDIT.

package mocks

import (
	io "io"

	mock "github.com/stretchr/testify/mock"
)

// Copier is an autogenerated mock type for the Copier type
type Copier struct {
	mock.Mock
}

// Create provides a mock function with given fields: path
func (_m *Copier) Create(path string) (io.WriteCloser, error) {
	ret := _m.Called(path)

	var r0 io.WriteCloser
	if rf, ok := ret.Get(0).(func(string) io.WriteCloser); ok {
		r0 = rf(path)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(io.WriteCloser)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(path)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MkdirAll provides a mock function with given fields: dir
func (_m *Copier) MkdirAll(dir string) error {
	ret := _m.Called(dir)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(dir)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Open provides a mock function with given fields: path
func (_m *Copier) Open(path string) (io.ReadCloser, error) {
	ret := _m.Called(path)

	var r0 io.ReadCloser
	if rf, ok := ret.Get(0).(func(string) io.ReadCloser); ok {
		r0 = rf(path)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(io.ReadCloser)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(path)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}