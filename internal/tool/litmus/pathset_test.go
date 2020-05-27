// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package litmus_test

import (
	"fmt"
	"testing"

	"github.com/MattWindsor91/act-tester/internal/helper/testhelp"

	"github.com/MattWindsor91/act-tester/internal/tool/litmus"
)

// ExamplePathset_Args is a runnable example for Args.
func ExamplePathset_Args() {
	for _, arg := range (&litmus.Pathset{
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
	fmt.Println((&litmus.Pathset{
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
		in  litmus.Pathset
		err error
	}{
		"no-file": {
			in:  litmus.Pathset{DirOut: "/tmp/scratch"},
			err: litmus.ErrNoFileIn,
		},
		"no-dir": {
			in:  litmus.Pathset{FileIn: "/home/foo/bar/baz.litmus"},
			err: litmus.ErrNoDirOut,
		},
		"ok": {
			in: litmus.Pathset{
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
