// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package lifter

import (
	"context"
	"path"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"
)

// MockHarnessMaker mocks HarnessMaker.
type MockHarnessMaker struct {
	// SeenSpecs collects the HarnessSpecs that the harness maker has seen.
	SeenSpecs []model.HarnessSpec

	// Err is the error to return on calls to MakeHarness.
	Err error
}

// MakeHarness mocks MakeHarness.
func (m *MockHarnessMaker) MakeHarness(_ context.Context, spec model.HarnessSpec) (outFiles []string, err error) {
	m.SeenSpecs = append(m.SeenSpecs, spec)
	return []string{path.Join(spec.OutDir, "out.c")}, m.Err
}
