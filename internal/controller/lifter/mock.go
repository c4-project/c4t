// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package lifter

import (
	"context"
	"io"
	"path"

	"github.com/MattWindsor91/act-tester/internal/model/job"

	"github.com/MattWindsor91/act-tester/internal/model/id"
)

// MockHarnessMaker mocks HarnessMaker.
type MockHarnessMaker struct {
	// SeenSpecs collects the HarnessSpecs that the harness maker has seen.
	SeenSpecs []job.Harness

	// Err is the error to return on calls to MakeHarness.
	Err error
}

// MakeHarness mocks MakeHarness.
func (m *MockHarnessMaker) MakeHarness(_ context.Context, spec job.Harness, _ io.Writer) (outFiles []string, err error) {
	m.SeenSpecs = append(m.SeenSpecs, spec)
	return []string{path.Join(spec.OutDir, "out.c")}, m.Err
}

// MockPather mocks Pather.
type MockPather struct {
	// Arches captures the last-prepared set of architecture IDs.
	Arches []id.ID

	// Subjects captures the last-prepared set of subject names.
	Subjects []string
}

// Prepare pretends to prepare a MockPather.
func (m *MockPather) Prepare(arches []id.ID, subjects []string) error {
	m.Arches = arches
	m.Subjects = subjects
	return nil
}

// Path pretends to resolve a path.
func (m *MockPather) Path(_ id.ID, _ string) (string, error) {
	return "foo", nil
}
