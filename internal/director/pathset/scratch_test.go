// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package pathset_test

import (
	"fmt"
	"path/filepath"

	"github.com/MattWindsor91/act-tester/internal/plan/stage"

	"github.com/MattWindsor91/act-tester/internal/director/pathset"
)

// ExampleNewScratch is a runnable example for NewScratch.
func ExampleNewScratch() {
	p := pathset.NewScratch("scratch")

	fmt.Println("run: ", filepath.ToSlash(p.DirRun))
	fmt.Println("lift:", filepath.ToSlash(p.DirLift))
	fmt.Println("fuzz:", filepath.ToSlash(p.DirFuzz))
	fmt.Println("plan:", filepath.ToSlash(p.DirPlan))

	// Output:
	// run:  scratch/run
	// lift: scratch/lift
	// fuzz: scratch/fuzz
	// plan: scratch/plan
}

// ExampleScratch_PlanForStage is a runnable example for PlanForStage.
func ExampleScratch_PlanForStage() {
	mp := pathset.Scratch{
		DirPlan: filepath.FromSlash("foo/plan"),
		DirFuzz: filepath.FromSlash("foo/fuzz"),
		DirLift: filepath.FromSlash("foo/lift"),
	}
	fmt.Print(filepath.ToSlash(mp.PlanForStage(stage.Fuzz)))
	// Output:
	// foo/plan/plan.fuzz.json
}
