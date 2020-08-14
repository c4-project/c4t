// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package quantity_test

import (
	"log"
	"os"
	"time"

	"github.com/MattWindsor91/act-tester/internal/quantity"
)

// ExampleRootSet_Log is a runnable example for RootSet.Log.
func ExampleRootSet_Log() {
	qs := quantity.RootSet{
		MachineSet: quantity.MachineSet{
			Fuzz: quantity.FuzzSet{
				CorpusSize:    10,
				SubjectCycles: 5,
				NWorkers:      4,
			},
			Mach: quantity.MachNodeSet{
				Compiler: quantity.BatchSet{
					Timeout:  quantity.Timeout(3 * time.Minute),
					NWorkers: 6,
				},
				Runner: quantity.BatchSet{
					Timeout:  quantity.Timeout(2 * time.Minute),
					NWorkers: 7,
				},
			},
			Perturb: quantity.PerturbSet{
				CorpusSize: 80,
			},
		},
		Plan: quantity.PlanSet{
			NWorkers: 9,
		},
	}

	qs.Log(log.New(os.Stdout, "", 0))

	// Output:
	// [Plan]
	// running across 9 workers
	// [Perturb]
	// target corpus size: 80 subjects
	// [Fuzz]
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
}
