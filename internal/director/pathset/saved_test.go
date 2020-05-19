// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package pathset_test

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/MattWindsor91/act-tester/internal/model/status"

	"github.com/MattWindsor91/act-tester/internal/director/pathset"
)

// ExampleNewSaved is a runnable example for NewSaved.
func ExampleNewSaved() {
	p := pathset.NewSaved("saved")

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

// ExampleSaved_PlanFile is a runnable example for PlanFile.
func ExampleSaved_PlanFile() {
	p := pathset.NewSaved("saved")
	t := time.Date(2015, time.October, 21, 7, 28, 0, 0, time.FixedZone("UTC-8", -8*60*60))
	stf, _ := p.PlanFile(status.CompileFail, t)
	fmt.Println(filepath.ToSlash(stf))

	// Output:
	// saved/compile_fail/2015/10/21/072800/plan.json
}

// ExampleSaved_SubjectTarFile is a runnable example for SubjectTarFile.
func ExampleSaved_SubjectTarFile() {
	p := pathset.NewSaved("saved")
	t := time.Date(2015, time.October, 21, 7, 28, 0, 0, time.FixedZone("UTC-8", -8*60*60))
	stf, _ := p.SubjectTarFile("foo", status.CompileFail, t)
	fmt.Println(filepath.ToSlash(stf))

	// Output:
	// saved/compile_fail/2015/10/21/072800/foo.tar.gz
}
