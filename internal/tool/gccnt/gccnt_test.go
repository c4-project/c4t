// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package gccnt

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/1set/gut/ystring"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGccnt_DryRun tests gccn't by dry-running on various input configurations.
func TestGccnt_DryRun(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		in  Gccnt
		out string
	}{
		"passthrough": {
			in:  Gccnt{Bin: "gcc", In: []string{"hello.c"}, Out: "a.out"},
			out: "invocation: gcc -o a.out -O hello.c",
		},
		"passthrough-opts": {
			in: Gccnt{
				Bin:         "gcc",
				In:          []string{"hello.c"},
				Out:         "a.out",
				DivergeOpts: []string{"2", "3"},
				ErrorOpts:   []string{"1"},
			},
			out: `The following optimisation levels will trigger an error: 1
			The following optimisation levels will trigger divergence: 2 3
			invocation: gcc -o a.out -O hello.c`,
		},
		"diverge": {
			in: Gccnt{
				Bin:         "gcc",
				In:          []string{"hello.c"},
				Out:         "a.out",
				OptLevel:    "3",
				DivergeOpts: []string{"2", "3"},
				ErrorOpts:   []string{"1"},
			},
			out: `The following optimisation levels will trigger an error: 1
			The following optimisation levels will trigger divergence: 2 3
            gccn't would diverge here`,
		},
		"error": {
			in: Gccnt{
				Bin:         "gcc",
				In:          []string{"hello.c"},
				Out:         "a.out",
				OptLevel:    "1",
				DivergeOpts: []string{"2", "3"},
				ErrorOpts:   []string{"1"},
			},
			out: `The following optimisation levels will trigger an error: 1
			The following optimisation levels will trigger divergence: 2 3
            gccn't would error here`,
		},
	}

	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			err := c.in.DryRun(context.Background(), &buf)
			require.NoError(t, err)

			assert.Equal(t, massageString(c.out), massageString(buf.String()), "dry run output differs")
		})
	}
}

func massageString(s string) string {
	return strings.TrimSpace(ystring.Shrink(s, " "))
}
