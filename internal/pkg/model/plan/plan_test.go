// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package plan_test

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"github.com/MattWindsor91/act-tester/internal/pkg/model/service"

	"github.com/MattWindsor91/act-tester/internal/pkg/model/id"

	"github.com/BurntSushi/toml"

	"github.com/MattWindsor91/act-tester/internal/pkg/model/plan"
)

// ExamplePlan_CompilerIDs is a runnable example for CompilerIDs.
func ExamplePlan_CompilerIDs() {
	p := plan.Plan{Compilers: map[string]service.Compiler{
		"gcc.ppc":   {Arch: id.ArchPPC},
		"clang.ppc": {Arch: id.ArchPPC},
		"gcc":       {Arch: id.ArchArm},
		"clang":     {Arch: id.ArchArm},
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

// TestMachinePlan_Arches tests the Arches method on MachinePlan.
func TestMachinePlan_Arches(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		plan plan.Plan
		want []id.ID
	}{
		"no arches": {plan.Plan{}, []id.ID{}},
		"one compiler": {plan.Plan{Compilers: map[string]service.Compiler{
			"gcc": {Arch: id.ArchX8664},
		}}, []id.ID{id.ArchX8664}},
		"same arch": {plan.Plan{Compilers: map[string]service.Compiler{
			"gcc":   {Arch: id.ArchArm},
			"clang": {Arch: id.ArchArm},
		}}, []id.ID{id.ArchArm}},
		"two arches": {plan.Plan{Compilers: map[string]service.Compiler{
			"gcc-ppc":   {Arch: id.ArchPPC},
			"clang-ppc": {Arch: id.ArchPPC},
			"gcc":       {Arch: id.ArchArm},
			"clang":     {Arch: id.ArchArm},
		}}, []id.ID{id.ArchArm, id.ArchPPC}},
	}

	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := c.plan.Arches()
			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("%s: Arches=%v; want %v", name, got, c.want)
			}
		})
	}
}

// TestPlan_Dump_roundTrip exercises Dump by doing a round-trip and checking if the reconstituted plan is similar.
func TestPlan_Dump_roundTrip(t *testing.T) {
	t.Parallel()

	p := plan.Mock()

	var b bytes.Buffer
	if err := p.Dump(&b); err != nil {
		t.Fatal("error dumping:", err)
	}

	var p2 plan.Plan
	if _, err := toml.DecodeReader(&b, &p2); err != nil {
		t.Fatal("error un-dumping:", err)
	}

	// TODO(@MattWindsor91): more comparisons?
	if !p.Header.Creation.Equal(p2.Header.Creation) {
		t.Errorf("date not equal after round-trip: send=%v, recv=%v", p.Header.Creation, p2.Header.Creation)
	}
}
