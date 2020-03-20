// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package pathset_test

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/MattWindsor91/act-tester/internal/pkg/director/pathset"
)

// ExampleNewSaved is a runnable example for NewSaved.
func ExampleNewSaved() {
	p := pathset.NewSaved("saved")

	fmt.Println("timeouts:        ", filepath.ToSlash(p.DirTimeouts))
	fmt.Println("run failures:    ", filepath.ToSlash(p.DirRunFailures))
	fmt.Println("flagged:         ", filepath.ToSlash(p.DirFlagged))
	fmt.Println("compile failures:", filepath.ToSlash(p.DirCompileFailures))

	// Output:
	// timeouts:         saved/timeout
	// run failures:     saved/run_fail
	// flagged:          saved/flagged
	// compile failures: saved/compile_fail
}

// ExampleSaved_CompileFailureTarFile is a runnable example for CompileFailureTarFile.
func ExampleSaved_CompileFailureTarFile() {
	p := pathset.NewSaved("saved")
	t := time.Date(2015, time.October, 21, 7, 28, 0, 0, time.FixedZone("UTC-8", -8*60*60))
	fmt.Println(filepath.ToSlash(p.CompileFailureTarFile("foo", t)))

	// Output:
	// saved/compile_fail/2015/10/21/072800/foo.tar.gz
}

// ExampleSaved_FlaggedTarFile is a runnable example for FlaggedTarFile.
func ExampleSaved_FlaggedTarFile() {
	p := pathset.NewSaved("saved")
	t := time.Date(2015, time.October, 21, 7, 28, 0, 0, time.FixedZone("UTC-8", -8*60*60))
	fmt.Println(filepath.ToSlash(p.FlaggedTarFile("foo", t)))

	// Output:
	// saved/flagged/2015/10/21/072800/foo.tar.gz
}

// ExampleSaved_RunFailureTarFile is a runnable example for RunFailureTarFile.
func ExampleSaved_RunFailureTarFile() {
	p := pathset.NewSaved("saved")
	t := time.Date(2015, time.October, 21, 7, 28, 0, 0, time.FixedZone("UTC-8", -8*60*60))
	fmt.Println(filepath.ToSlash(p.RunFailureTarFile("foo", t)))

	// Output:
	// saved/run_fail/2015/10/21/072800/foo.tar.gz
}

// ExampleSaved_TimeoutTarFile is a runnable example for TimeoutTarFile.
func ExampleSaved_TimeoutTarFile() {
	p := pathset.NewSaved("saved")
	t := time.Date(2015, time.October, 21, 7, 28, 0, 0, time.FixedZone("UTC-8", -8*60*60))
	fmt.Println(filepath.ToSlash(p.TimeoutTarFile("foo", t)))

	// Output:
	// saved/timeout/2015/10/21/072800/foo.tar.gz
}
