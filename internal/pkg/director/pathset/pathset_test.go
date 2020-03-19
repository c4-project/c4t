// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package pathset_test

import (
	"fmt"
	"path/filepath"

	"github.com/MattWindsor91/act-tester/internal/pkg/director/pathset"

	"github.com/MattWindsor91/act-tester/internal/pkg/model/id"
)

// ExampleNew is a runnable example for New.
func ExampleNew() {
	p := pathset.New(filepath.FromSlash("tests"))

	fmt.Println(filepath.ToSlash(p.DirScratch))
	fmt.Println(filepath.ToSlash(p.DirSaved))

	// Output:
	// tests/scratch
	// tests/saved
}

// ExamplePathset_MachineSaved is a runnable example for MachineSaved.
func ExamplePathset_MachineSaved() {
	p := pathset.Pathset{DirSaved: "saved", DirScratch: "scratch"}
	mid := id.FromString("foo.bar.baz")
	mp := p.MachineSaved(mid)

	fmt.Println(filepath.ToSlash(mp.DirCompileFailures))
	fmt.Println(filepath.ToSlash(mp.DirFlagged))
	fmt.Println(filepath.ToSlash(mp.DirRunFailures))
	fmt.Println(filepath.ToSlash(mp.DirTimeouts))

	// Output:
	// saved/foo/bar/baz/compile_fail
	// saved/foo/bar/baz/flagged
	// saved/foo/bar/baz/run_fail
	// saved/foo/bar/baz/timeout
}

// ExamplePathset_MachineScratch is a runnable example for MachineScratch.
func ExamplePathset_MachineScratch() {
	p := pathset.Pathset{DirSaved: "saved", DirScratch: "scratch"}
	mid := id.FromString("foo.bar.baz")
	mp := p.MachineScratch(mid)

	fmt.Println(filepath.ToSlash(mp.DirFuzz))
	fmt.Println(filepath.ToSlash(mp.DirLift))
	fmt.Println(filepath.ToSlash(mp.DirPlan))
	fmt.Println(filepath.ToSlash(mp.DirRun))

	// Output:
	// scratch/foo/bar/baz/fuzz
	// scratch/foo/bar/baz/lift
	// scratch/foo/bar/baz/plan
	// scratch/foo/bar/baz/run
}
