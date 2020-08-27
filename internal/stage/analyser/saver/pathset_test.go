// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package saver_test

import (
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/MattWindsor91/act-tester/internal/helper/testhelp"

	"github.com/MattWindsor91/act-tester/internal/stage/analyser/saver"

	"github.com/MattWindsor91/act-tester/internal/subject/status"
)

// ExampleNewPathset is a runnable example for NewPathset.
func ExampleNewPathset() {
	p := saver.NewPathset("saved")

	for s := status.FirstBad; s <= status.Last; s++ {
		fmt.Printf("%s: %s\n", s, p.Dirs[s])
	}

	// Unordered output:
	// Flagged: saved/flagged
	// CompileFail: saved/compile_fail
	// CompileTimeout: saved/compile_timeout
	// RunFail: saved/run_fail
	// RunTimeout: saved/run_timeout
}

// ExamplePathset_SubjectRun is a runnable example for SubjectRun.
func ExamplePathset_SubjectRun() {
	p := saver.NewPathset("saved")
	t := time.Date(2015, time.October, 21, 7, 28, 0, 0, time.FixedZone("UTC-8", -8*60*60))
	stf, _ := p.SubjectRun(status.CompileFail, t)
	fmt.Println("root:", filepath.ToSlash(stf.DirRoot))
	fmt.Println("plan:", filepath.ToSlash(stf.FilePlan))

	// Output:
	// root: saved/compile_fail/2015/10/21/07_28_00
	// plan: saved/compile_fail/2015/10/21/07_28_00/plan.json.gz
}

// ExampleRunPathset_SubjectTarFile is a runnable example for SubjectTarFile.
func ExampleRunPathset_SubjectTarFile() {
	p := saver.NewPathset("saved")
	t := time.Date(2015, time.October, 21, 7, 28, 0, 0, time.FixedZone("UTC-8", -8*60*60))
	rp, _ := p.SubjectRun(status.CompileFail, t)
	fmt.Println(filepath.ToSlash(rp.SubjectTarFile("foo")))

	// Output:
	// saved/compile_fail/2015/10/21/07_28_00/foo.tar.gz
}

// TestPathset_SubjectRun_errors tests several error cases for SubjectRun.
func TestPathset_SubjectRun_errors(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		in  status.Status
		err error
	}{
		"ok":      {in: status.Ok, err: status.ErrBad},
		"unknown": {in: status.Unknown, err: status.ErrBad},
		"oob":     {in: status.Last + 1, err: status.ErrBad},
	}

	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			p := saver.NewPathset("saved")
			d := time.Date(2015, time.October, 21, 7, 28, 0, 0, time.FixedZone("UTC-8", -8*60*60))
			_, err := p.SubjectRun(c.in, d)
			testhelp.ExpectErrorIs(t, err, c.err, "SubjectRun")
		})
	}
}
