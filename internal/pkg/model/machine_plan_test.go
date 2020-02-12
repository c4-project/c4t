package model

import (
	"reflect"
	"testing"
)

// TestMachinePlan_Arches tests the Arches method on MachinePlan.
func TestMachinePlan_Arches(t *testing.T) {
	tests := []struct {
		name string
		plan MachinePlan
		want []ID
	}{
		{"no arches", MachinePlan{}, []ID{}},
		{"one compiler", MachinePlan{Compilers: []Compiler{
			Compiler{Service: Service{ID: IDFromString("gcc")}, Arch: ArchX8664},
		}}, []ID{ArchX8664}},
		{"same arch", MachinePlan{Compilers: []Compiler{
			Compiler{Service: Service{ID: IDFromString("gcc")}, Arch: ArchArm},
			Compiler{Service: Service{ID: IDFromString("clang")}, Arch: ArchArm},
		}}, []ID{ArchArm}},
		{"two arches", MachinePlan{Compilers: []Compiler{
			Compiler{Service: Service{ID: IDFromString("gcc-ppc")}, Arch: ArchPPC},
			Compiler{Service: Service{ID: IDFromString("clang-ppc")}, Arch: ArchPPC},
			Compiler{Service: Service{ID: IDFromString("gcc")}, Arch: ArchArm},
			Compiler{Service: Service{ID: IDFromString("clang")}, Arch: ArchArm},
		}}, []ID{ArchArm, ArchPPC}},
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
