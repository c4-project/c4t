// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package gccnt_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/c4-project/c4t/internal/tool/gccnt"

	"github.com/1set/gut/ystring"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGccnt_DryRun tests gccn't by dry-running on various input configurations.
func TestGccnt_DryRun(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		in  gccnt.Gccnt
		out string
	}{
		"passthrough": {
			in:  gccnt.Gccnt{Bin: "gcc", In: []string{"hello.c"}, Out: "a.out"},
			out: "invocation: gcc -o a.out -O hello.c",
		},
		"passthrough-pthread": {
			in:  gccnt.Gccnt{Bin: "gcc", In: []string{"hello.c"}, Out: "a.out", Pthread: true},
			out: "invocation: gcc -o a.out -O -pthread hello.c",
		},
		"passthrough-oflag": {
			in:  gccnt.Gccnt{Bin: "gcc", In: []string{"hello.c"}, Out: "a.out", OptLevel: "3"},
			out: "invocation: gcc -o a.out -O3 hello.c",
		},
		"passthrough-march": {
			in:  gccnt.Gccnt{Bin: "gcc", In: []string{"hello.c"}, Out: "a.out", March: "nehalem"},
			out: "invocation: gcc -o a.out -O -march=nehalem hello.c",
		},
		"passthrough-mcpu": {
			in:  gccnt.Gccnt{Bin: "gcc", In: []string{"hello.c"}, Out: "a.out", Mcpu: "power9"},
			out: "invocation: gcc -o a.out -O -mcpu=power9 hello.c",
		},
		"passthrough-addopts": {
			in: gccnt.Gccnt{
				Bin: "gcc",
				In:  []string{"hello.c"},
				Out: "a.out",
				Conds: gccnt.ConditionSet{
					Diverge: gccnt.Condition{Opts: []string{"2", "3"}},
					Error:   gccnt.Condition{Opts: []string{"1"}},
				},
			},
			out: `The following optimisation levels will trigger divergence: 2 3
            The following optimisation levels will trigger an error: 1
			invocation: gcc -o a.out -O hello.c`,
		},
		"passthrough-mutant-miss-noperiods": {
			in: gccnt.Gccnt{
				Bin:    "gcc",
				In:     []string{"hello.c"},
				Out:    "a.out",
				Mutant: 2,
			},
			out: `MUTATION SELECTED: 2
			invocation: gcc -o a.out -O hello.c`,
		},
		"passthrough-mutant-miss": {
			in: gccnt.Gccnt{
				Bin:    "gcc",
				In:     []string{"hello.c"},
				Out:    "a.out",
				Mutant: 2,
				Conds: gccnt.ConditionSet{
					MutHitPeriod: 4,
					Diverge:      gccnt.Condition{MutPeriod: 3},
					Error:        gccnt.Condition{MutPeriod: 5},
				},
			},
			out: `MUTATION SELECTED: 2
			Mutation numbers that are multiples of 3 will trigger divergence
			Mutation numbers that are multiples of 5 will trigger an error
			invocation: gcc -o a.out -O hello.c`,
		},
		"passthrough-mutant-hit": {
			in: gccnt.Gccnt{
				Bin:    "gcc",
				In:     []string{"hello.c"},
				Out:    "a.out",
				Mutant: 4,
				Conds: gccnt.ConditionSet{
					MutHitPeriod: 4,
					Diverge:      gccnt.Condition{MutPeriod: 3},
					Error:        gccnt.Condition{MutPeriod: 5},
				},
			},
			out: `MUTATION SELECTED: 4
            MUTATION HIT: 4
			Mutation numbers that are multiples of 3 will trigger divergence
			Mutation numbers that are multiples of 5 will trigger an error
			invocation: gcc -o a.out -O hello.c`,
		},
		"diverge": {
			in: gccnt.Gccnt{
				Bin:      "gcc",
				In:       []string{"hello.c"},
				Out:      "a.out",
				OptLevel: "3",
				Conds: gccnt.ConditionSet{
					Diverge: gccnt.Condition{Opts: []string{"2", "3"}},
					Error:   gccnt.Condition{Opts: []string{"1"}},
				},
			},
			out: `The following optimisation levels will trigger divergence: 2 3
            The following optimisation levels will trigger an error: 1
            gccn't would diverge here`,
		},
		"diverge-mutant": {
			in: gccnt.Gccnt{
				Bin:    "gcc",
				In:     []string{"hello.c"},
				Out:    "a.out",
				Mutant: 3,
				Conds: gccnt.ConditionSet{
					MutHitPeriod: 4,
					Diverge:      gccnt.Condition{MutPeriod: 3},
					Error:        gccnt.Condition{MutPeriod: 5},
				},
			},
			out: `MUTATION SELECTED: 3
            MUTATION HIT: 3
			Mutation numbers that are multiples of 3 will trigger divergence
			Mutation numbers that are multiples of 5 will trigger an error
			gccn't would diverge here`,
		},
		"error": {
			in: gccnt.Gccnt{
				Bin:      "gcc",
				In:       []string{"hello.c"},
				Out:      "a.out",
				OptLevel: "1",
				Conds: gccnt.ConditionSet{
					Diverge: gccnt.Condition{Opts: []string{"2", "3"}},
					Error:   gccnt.Condition{Opts: []string{"1"}},
				},
			},
			out: `The following optimisation levels will trigger divergence: 2 3
            The following optimisation levels will trigger an error: 1
            gccn't would error here`,
		},
		"error-mutant": {
			in: gccnt.Gccnt{
				Bin:    "gcc",
				In:     []string{"hello.c"},
				Out:    "a.out",
				Mutant: 5,
				Conds: gccnt.ConditionSet{
					MutHitPeriod: 4,
					Diverge:      gccnt.Condition{MutPeriod: 3},
					Error:        gccnt.Condition{MutPeriod: 5},
				},
			},
			out: `MUTATION SELECTED: 5
            MUTATION HIT: 5
			Mutation numbers that are multiples of 3 will trigger divergence
			Mutation numbers that are multiples of 5 will trigger an error
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
