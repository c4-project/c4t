// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package config_test

import (
	"log"
	"os"
	"time"

	"github.com/MattWindsor91/act-tester/internal/config"
	"github.com/MattWindsor91/act-tester/internal/stage/fuzzer"
	"github.com/MattWindsor91/act-tester/internal/stage/mach/quantity"
	"github.com/MattWindsor91/act-tester/internal/stage/mach/timeout"
	"github.com/MattWindsor91/act-tester/internal/stage/perturber"
	"github.com/MattWindsor91/act-tester/internal/stage/planner"
)

// ExampleQuantitySet_Log is a runnable example for QuantitySet.Log.
func ExampleQuantitySet_Log() {
	qs := config.QuantitySet{
		Fuzz: fuzzer.QuantitySet{
			CorpusSize:    10,
			SubjectCycles: 5,
			NWorkers:      4,
		},
		Mach: quantity.Set{
			Compiler: quantity.SingleSet{
				Timeout:  timeout.Timeout(3 * time.Minute),
				NWorkers: 6,
			},
			Runner: quantity.SingleSet{
				Timeout:  timeout.Timeout(2 * time.Minute),
				NWorkers: 7,
			},
		},
		Plan: planner.QuantitySet{
			NWorkers: 9,
		},
		Perturb: perturber.QuantitySet{
			CorpusSize: 80,
		},
	}

	qs.Log(log.New(os.Stdout, "", 0))

	// Output:
	// [Fuzzer]
	// running across 4 workers
	// fuzzing each subject 5 times
	// target corpus size: 10 subjects
	// [Mach]
	// [Compiler]
	// running across 6 workers
	// timeout at 3m0s
	// [Runner]
	// running across 7 workers
	// timeout at 2m0s
	// [Plan]
	// running across 9 workers
	// [Perturb]
	// target corpus size: 80 subjects
}
