// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package plan_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/MattWindsor91/act-tester/internal/plan/stage"

	"github.com/MattWindsor91/act-tester/internal/model/corpus"

	"github.com/MattWindsor91/act-tester/internal/helper/testhelp"

	"github.com/stretchr/testify/assert"

	"github.com/MattWindsor91/act-tester/internal/model/service/compiler"

	"github.com/MattWindsor91/act-tester/internal/model/id"

	"github.com/MattWindsor91/act-tester/internal/plan"
)

// ExamplePlan_CompilerIDs is a runnable example for Plan.CompilerIDs.
func ExamplePlan_CompilerIDs() {
	p := plan.Plan{Compilers: map[string]compiler.Configuration{
		"gcc.ppc":   {Compiler: compiler.Compiler{Arch: id.ArchPPC}},
		"clang.ppc": {Compiler: compiler.Compiler{Arch: id.ArchPPC}},
		"gcc":       {Compiler: compiler.Compiler{Arch: id.ArchArm}},
		"clang":     {Compiler: compiler.Compiler{Arch: id.ArchArm}},
	}}
	cids, _ := p.CompilerIDs()
	for _, c := range cids {
		fmt.Println(c)
	}

	// Output:
	// clang
	// clang.ppc
	// gcc
	// gcc.ppc
}

// TestPlan_Arches tests Plan.Arches.
func TestPlan_Arches(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		plan plan.Plan
		want []id.ID
	}{
		"no arches": {plan.Plan{}, []id.ID{}},
		"one compiler": {plan.Plan{Compilers: map[string]compiler.Configuration{
			"gcc": {Compiler: compiler.Compiler{Arch: id.ArchX8664}},
		}}, []id.ID{id.ArchX8664}},
		"same arch": {plan.Plan{Compilers: map[string]compiler.Configuration{
			"gcc":   {Compiler: compiler.Compiler{Arch: id.ArchArm}},
			"clang": {Compiler: compiler.Compiler{Arch: id.ArchArm}},
		}}, []id.ID{id.ArchArm}},
		"two arches": {plan.Plan{Compilers: map[string]compiler.Configuration{
			"gcc-ppc":   {Compiler: compiler.Compiler{Arch: id.ArchPPC}},
			"clang-ppc": {Compiler: compiler.Compiler{Arch: id.ArchPPC}},
			"gcc":       {Compiler: compiler.Compiler{Arch: id.ArchArm}},
			"clang":     {Compiler: compiler.Compiler{Arch: id.ArchArm}},
		}}, []id.ID{id.ArchArm, id.ArchPPC}},
	}

	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := c.plan.Arches()
			assert.Equalf(t, c.want, got, "%s: Arches=%v; want %v", name, got, c.want)
		})
	}
}

// TestPlan_Check tests Plan.Check.
func TestPlan_Check(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		in  plan.Plan
		err error
	}{
		"no version": {
			in:  plan.Plan{},
			err: plan.ErrVersionMismatch,
		},
		"version too low": {
			in: plan.Plan{Metadata: plan.Metadata{
				Version: plan.CurrentVer - 1,
			}},
			err: plan.ErrVersionMismatch,
		},
		"version too high": {
			in: plan.Plan{Metadata: plan.Metadata{
				Version: plan.CurrentVer + 1,
			}},
			err: plan.ErrVersionMismatch,
		},
		"no corpus": {
			in: plan.Plan{Metadata: plan.Metadata{
				Version: plan.CurrentVer,
			}},
			err: corpus.ErrNone,
		},
		"known good plan": {
			in: *plan.Mock(),
		},
	}

	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := c.in.Check()
			testhelp.ExpectErrorIs(t, err, c.err, "check plan")
		})
	}
}

// TestPlan_RunStage_error tests that giving Plan.RunStage a body that errors propagates that error correctly.
func TestPlan_RunStage_error(t *testing.T) {
	t.Parallel()

	p := plan.Mock()
	want := errors.New("oops")

	_, got := p.RunStage(context.Background(), stage.Last, func(context.Context, *plan.Plan) (*plan.Plan, error) {
		return nil, want
	})

	testhelp.ExpectErrorIs(t, got, want, "running stage with error")
}
