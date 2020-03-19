// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package pathset_test

import (
	"fmt"
	"path/filepath"

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
