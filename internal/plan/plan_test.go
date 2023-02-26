// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package plan_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/c4-project/c4t/internal/mutation"

	"github.com/c4-project/c4t/internal/plan/stage"

	"github.com/c4-project/c4t/internal/subject/corpus"

	"github.com/c4-project/c4t/internal/helper/testhelp"

	"github.com/stretchr/testify/assert"

	"github.com/c4-project/c4t/internal/model/service/compiler"

	"github.com/c4-project/c4t/internal/id"

	"github.com/c4-project/c4t/internal/plan"
)

// ExamplePlan_CompilerIDs is a runnable example for Plan.CompilerIDs.
func ExamplePlan_CompilerIDs() {
	p := plan.Plan{Compilers: compiler.InstanceMap{
		id.FromString("gcc.ppc"):   {Compiler: compiler.Compiler{Arch: id.ArchPPC}},
		id.FromString("clang.ppc"): {Compiler: compiler.Compiler{Arch: id.ArchPPC}},
		id.FromString("gcc"):       {Compiler: compiler.Compiler{Arch: id.ArchArm}},
		id.FromString("clang"):     {Compiler: compiler.Compiler{Arch: id.ArchArm}},
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
		"one compiler": {plan.Plan{Compilers: compiler.InstanceMap{
			id.FromString("gcc"): {Compiler: compiler.Compiler{Arch: id.ArchX8664}},
		}}, []id.ID{id.ArchX8664}},
		"same arch": {plan.Plan{Compilers: compiler.InstanceMap{
			id.FromString("gcc"):   {Compiler: compiler.Compiler{Arch: id.ArchArm}},
			id.FromString("clang"): {Compiler: compiler.Compiler{Arch: id.ArchArm}},
		}}, []id.ID{id.ArchArm}},
		"two arches": {plan.Plan{Compilers: compiler.InstanceMap{
			id.FromString("gcc-ppc"):   {Compiler: compiler.Compiler{Arch: id.ArchPPC}},
			id.FromString("clang-ppc"): {Compiler: compiler.Compiler{Arch: id.ArchPPC}},
			id.FromString("gcc"):       {Compiler: compiler.Compiler{Arch: id.ArchArm}},
			id.FromString("clang"):     {Compiler: compiler.Compiler{Arch: id.ArchArm}},
		}}, []id.ID{id.ArchArm, id.ArchPPC}},
	}

	for name, c := range cases {
		name := name
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

type badStage struct {
	err error
}

func (b badStage) Stage() stage.Stage {
	return stage.Last
}
func (b badStage) Close() error {
	return nil
}
func (b badStage) Run(context.Context, *plan.Plan) (*plan.Plan, error) {
	return nil, b.err
}

// TestPlan_RunStage_error tests that giving Plan.RunStage a body that errors propagates that error correctly.
func TestPlan_RunStage_error(t *testing.T) {
	t.Parallel()

	p := plan.Mock()
	want := errors.New("oops")

	_, got := p.RunStage(context.Background(), badStage{err: want})

	testhelp.ExpectErrorIs(t, got, want, "running stage with error")
}

// ExamplePlan_SetMutant is a runnable example for Plan.SetMutant.
func ExamplePlan_SetMutant() {
	p := plan.Mock()
	fmt.Printf("plan mutant: %s (is mutation test: %v)\n", p.Mutant(), p.IsMutationTest())
	p.SetMutant(mutation.Mutant{Name: mutation.Name{Operator: "XYZ", Variant: 1}, Index: 42})
	fmt.Printf("plan mutant: %s (is mutation test: %v)\n", p.Mutant(), p.IsMutationTest())
	p.Mutation = &mutation.Config{Enabled: true}
	fmt.Printf("plan mutant: %s (is mutation test: %v)\n", p.Mutant(), p.IsMutationTest())
	p.SetMutant(mutation.Mutant{Name: mutation.Name{Operator: "XYZ", Variant: 1}, Index: 42})
	fmt.Printf("plan mutant: %s (is mutation test: %v)\n", p.Mutant(), p.IsMutationTest())

	// Output:
	// plan mutant: 0 (is mutation test: false)
	// plan mutant: 0 (is mutation test: false)
	// plan mutant: 0 (is mutation test: true)
	// plan mutant: XYZ1:42 (is mutation test: true)
}
