// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package fuzzer

import (
	"path"

	"github.com/MattWindsor91/act-tester/internal/pkg/model/subject"
)

// MockPathset mocks the SubjectPather interface.
type MockPathset struct {
	HasPrepared   bool
	SubjectCycles []SubjectCycle
}

func (m *MockPathset) Prepare() error {
	m.HasPrepared = true
	return nil
}

func (m *MockPathset) SubjectPaths(sc SubjectCycle) subject.FuzzFileset {
	m.SubjectCycles = append(m.SubjectCycles, sc)
	return subject.FuzzFileset{
		Litmus: path.Join("litmus", sc.String()),
		Trace:  path.Join("trace", sc.String()),
	}
}
