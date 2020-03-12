// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package director_test

import (
	"fmt"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"

	"github.com/MattWindsor91/act-tester/internal/pkg/director"
)

// ExampleNewPathset is a runnable example for NewPathset.
func ExampleNewPathset() {
	p := director.NewPathset("/home/enid/tests")

	fmt.Println(p.DirScratch)
	fmt.Println(p.DirSaved)

	// Output:
	// /home/enid/tests/scratch
	// /home/enid/tests/saved
}

// ExamplePathset_MachineScratch is a runnable example for MachineScratch.
func ExamplePathset_MachineScratch() {
	p := director.Pathset{DirSaved: "saved", DirScratch: "scratch"}
	mid := model.IDFromString("foo.bar.baz")
	mp := p.MachineScratch(mid)

	fmt.Println(mp.DirFuzz)
	fmt.Println(mp.DirLift)
	fmt.Println(mp.DirPlan)

	// Output:
	// scratch/foo/bar/baz/fuzz
	// scratch/foo/bar/baz/lift
	// scratch/foo/bar/baz/plan
}

// ExampleMachinePathset_PlanForStage is a runnable example for PlanForStage.
func ExampleMachinePathset_PlanForStage() {
	mp := director.MachinePathset{
		DirPlan: "foo/plan",
		DirFuzz: "foo/fuzz",
		DirLift: "foo/lift",
	}
	fmt.Print(mp.PlanForStage("fuzz"))
	// Output:
	// foo/plan/plan.fuzz.toml
}
