package plan

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"
)

// ExampleMachinePlan_CompilerIDs is a runnable example for CompilerIDs.
func ExampleMachinePlan_CompilerIDs() {
	plan := MachinePlan{Compilers: map[string]model.Compiler{
		"gcc.ppc":   {Arch: model.ArchPPC},
		"clang.ppc": {Arch: model.ArchPPC},
		"gcc":       {Arch: model.ArchArm},
		"clang":     {Arch: model.ArchArm},
	}}
	for _, c := range plan.CompilerIDs() {
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
	tests := []struct {
		name string
		plan MachinePlan
		want []model.ID
	}{
		{"no arches", MachinePlan{}, []model.ID{}},
		{"one compiler", MachinePlan{Compilers: map[string]model.Compiler{
			"gcc": {Arch: model.ArchX8664},
		}}, []model.ID{model.ArchX8664}},
		{"same arch", MachinePlan{Compilers: map[string]model.Compiler{
			"gcc":   {Arch: model.ArchArm},
			"clang": {Arch: model.ArchArm},
		}}, []model.ID{model.ArchArm}},
		{"two arches", MachinePlan{Compilers: map[string]model.Compiler{
			"gcc-ppc":   {Arch: model.ArchPPC},
			"clang-ppc": {Arch: model.ArchPPC},
			"gcc":       {Arch: model.ArchArm},
			"clang":     {Arch: model.ArchArm},
		}}, []model.ID{model.ArchArm, model.ArchPPC}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.plan.Arches()
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("%s: Arches=%v; want %v", test.name, got, test.want)
			}
		})
	}
}
