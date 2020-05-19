// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package save_test

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/MattWindsor91/act-tester/internal/controller/analyse/save"

	"github.com/MattWindsor91/act-tester/internal/model/status"
)

// ExampleNewPathset is a runnable example for NewPathset.
func ExampleNewPathset() {
	p := save.NewPathset("saved")

	for s := status.FirstBad; s < status.Num; s++ {
		fmt.Printf("%s: %s\n", s, p.Dirs[s])
	}

	// Unordered output:
	// flagged: saved/flagged
	// compile/fail: saved/compile_fail
	// compile/timeout: saved/compile_timeout
	// run/fail: saved/run_fail
	// run/timeout: saved/run_timeout
}

// ExamplePathset_PlanFile is a runnable example for PlanFile.
func ExamplePathset_PlanFile() {
	p := save.NewPathset("saved")
	t := time.Date(2015, time.October, 21, 7, 28, 0, 0, time.FixedZone("UTC-8", -8*60*60))
	stf, _ := p.PlanFile(status.CompileFail, t)
	fmt.Println(filepath.ToSlash(stf))

	// Output:
	// saved/compile_fail/2015/10/21/07_28_00/plan.json
}

// ExamplePathset_SubjectTarFile is a runnable example for SubjectTarFile.
func ExamplePathset_SubjectTarFile() {
	p := save.NewPathset("saved")
	t := time.Date(2015, time.October, 21, 7, 28, 0, 0, time.FixedZone("UTC-8", -8*60*60))
	stf, _ := p.SubjectTarFile("foo", status.CompileFail, t)
	fmt.Println(filepath.ToSlash(stf))

	// Output:
	// saved/compile_fail/2015/10/21/07_28_00/foo.tar.gz
}
