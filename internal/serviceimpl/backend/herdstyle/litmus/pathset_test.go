// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package litmus_test

import (
	"fmt"
	"testing"

	litmus2 "github.com/MattWindsor91/c4t/internal/serviceimpl/backend/herdstyle/litmus"

	"github.com/MattWindsor91/c4t/internal/helper/testhelp"
)

// ExamplePathset_Args is a runnable example for Args.
func ExamplePathset_Args() {
	for _, arg := range (&litmus2.Pathset{
		FileIn: "/home/foo/bar/baz.litmus",
		DirOut: "/tmp/scratch/",
	}).Args() {
		fmt.Println(arg)
	}

	// Output:
	// -o
	// /tmp/scratch/
	// /home/foo/bar/baz.litmus
}

// ExamplePathset_MainCFile is a runnable example for MainCFile.
func ExamplePathset_MainCFile() {
	fmt.Println((&litmus2.Pathset{
		FileIn: "/home/foo/bar/baz.litmus",
		DirOut: "/tmp/scratch/",
	}).MainCFile())

	// Output:
	// /tmp/scratch/baz.c
}

// TestPathset_Check makes sure Check catches various pathset errors.
func TestPathset_Check(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		in  litmus2.Pathset
		err error
	}{
		"no-file": {
			in:  litmus2.Pathset{DirOut: "/tmp/scratch"},
			err: litmus2.ErrNoFileIn,
		},
		"no-dir": {
			in:  litmus2.Pathset{FileIn: "/home/foo/bar/baz.litmus"},
			err: litmus2.ErrNoDirOut,
		},
		"ok": {
			in: litmus2.Pathset{
				FileIn: "/home/foo/bar/baz.litmus",
				DirOut: "/tmp/scratch/",
			},
		},
	}

	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := c.in.Check()
			testhelp.ExpectErrorIs(t, got, c.err, "checking pathset")
		})
	}
}
