// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package fuzzer_test

import (
	"testing"

	"github.com/MattWindsor91/act-tester/internal/pkg/fuzzer"
	"github.com/MattWindsor91/act-tester/internal/pkg/plan"
	"github.com/MattWindsor91/act-tester/internal/pkg/testhelp"
)

// makeConfig makes a 'valid' fuzzer config.
func makeConfig() fuzzer.Config {
	return fuzzer.Config{
		Driver: fuzzer.NopFuzzer{},
		Paths:  &fuzzer.MockPathset{},
		Quantities: fuzzer.QuantitySet{
			CorpusSize:    0,
			SubjectCycles: 10,
		},
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
