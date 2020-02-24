// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package fuzzer_test

import (
	"path"
	"testing"

	"github.com/MattWindsor91/act-tester/internal/pkg/subject"

	"github.com/MattWindsor91/act-tester/internal/pkg/fuzzer"
	"github.com/MattWindsor91/act-tester/internal/pkg/plan"
	"github.com/MattWindsor91/act-tester/internal/pkg/testhelp"
)

// MockPathset mocks the SubjectPather interface.
type MockPathset struct {
	HasPrepared   bool
	SubjectCycles []fuzzer.SubjectCycle
}

func (m *MockPathset) Prepare() error {
	m.HasPrepared = true
	return nil
}

func (m *MockPathset) SubjectPaths(sc fuzzer.SubjectCycle) subject.FuzzFileset {
	m.SubjectCycles = append(m.SubjectCycles, sc)
	return subject.FuzzFileset{
		Litmus: path.Join("litmus", sc.String()),
		Trace:  path.Join("trace", sc.String()),
	}
}

// makeConfig makes a 'valid' fuzzer config.
func makeConfig() fuzzer.Config {
	return fuzzer.Config{
		Driver:        fuzzer.NopFuzzer{},
		Paths:         &MockPathset{},
		CorpusSize:    0,
		SubjectCycles: 10,
	}
}

// TestNewFuzzer_DriverNil makes sure fuzzer creation on a nil driver fails.
func TestNewFuzzer_DriverNil(t *testing.T) {
	c := makeConfig()
	c.Driver = nil
	_, err := fuzzer.New(&c, &plan.Plan{})
	testhelp.ExpectErrorIs(t, err, fuzzer.ErrDriverNil, "fuzzer.New on nil driver")
}

// TestNewFuzzer_PlanNil makes sure fuzzer creation on a nil plan fails.
func TestNewFuzzer_PlanNil(t *testing.T) {
	c := makeConfig()
	_, err := fuzzer.New(&c, nil)
	testhelp.ExpectErrorIs(t, err, plan.ErrNil, "fuzzer.New on nil plan")
}
