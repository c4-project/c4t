// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package compiler

import (
	"github.com/MattWindsor91/act-tester/internal/model/id"
	"github.com/MattWindsor91/act-tester/internal/model/subject"
	"github.com/stretchr/testify/mock"
)

// MockSubjectPather mocks SubjectPather.
type MockSubjectPather struct {
	mock.Mock
}

// Prepare mocks the eponymous method.
func (m *MockSubjectPather) Prepare(compilers []id.ID) error {
	args := m.Called(compilers)
	return args.Error(0)
}

// SubjectPaths mocks the eponymous method.
func (m *MockSubjectPather) SubjectPaths(sc SubjectCompile) subject.CompileFileset {
	args := m.Called(sc)
	return args.Get(0).(subject.CompileFileset)
}
