// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package fuzzer

import (
	"github.com/MattWindsor91/act-tester/internal/model/subject"
	"github.com/stretchr/testify/mock"
)

// MockPathset mocks the SubjectPather interface.
type MockPathset struct {
	mock.Mock
}

func (m *MockPathset) Prepare() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockPathset) SubjectPaths(sc SubjectCycle) subject.FuzzFileset {
	args := m.Called(sc)
	return args.Get(0).(subject.FuzzFileset)
}
