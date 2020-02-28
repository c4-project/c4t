// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package plan_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"
	"github.com/MattWindsor91/act-tester/internal/pkg/plan"
)

// ExamplePlan_CompilerIDs is a runnable example for CompilerIDs.
func ExamplePlan_CompilerIDs() {
	p := plan.Plan{Compilers: map[string]model.Compiler{
		"gcc.ppc":   {Arch: model.ArchPPC},
		"clang.ppc": {Arch: model.ArchPPC},
		"gcc":       {Arch: model.ArchArm},
		"clang":     {Arch: model.ArchArm},
	}}
	for _, c := range p.CompilerIDs() {
		fmt.Println(c.String())
	}

	// Output:
	// clang
	// clang.ppc
	// gcc
	// gcc.ppc
}

// TestMachinePlan_Arches tests the Arches method on MachinePlan.
func TestMachinePlan_Arches(t *testing.T) {
	cases := map[string]struct {
		plan plan.Plan
		want []model.ID
	}{
		"no arches": {plan.Plan{}, []model.ID{}},
		"one compiler": {plan.Plan{Compilers: map[string]model.Compiler{
			"gcc": {Arch: model.ArchX8664},
		}}, []model.ID{model.ArchX8664}},
		"same arch": {plan.Plan{Compilers: map[string]model.Compiler{
			"gcc":   {Arch: model.ArchArm},
			"clang": {Arch: model.ArchArm},
		}}, []model.ID{model.ArchArm}},
		"two arches": {plan.Plan{Compilers: map[string]model.Compiler{
			"gcc-ppc":   {Arch: model.ArchPPC},
			"clang-ppc": {Arch: model.ArchPPC},
			"gcc":       {Arch: model.ArchArm},
			"clang":     {Arch: model.ArchArm},
		}}, []model.ID{model.ArchArm, model.ArchPPC}},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			got := c.plan.Arches()
			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("%s: Arches=%v; want %v", name, got, c.want)
			}
		})
	}
}
