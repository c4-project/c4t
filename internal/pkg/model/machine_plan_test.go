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
		want []Id
	}{
		{"no arches", MachinePlan{}, []Id{}},
		{"one compiler", MachinePlan{Compilers: []Compiler{
			Compiler{Service: Service{Id: IdFromString("gcc")}, Arch: ArchX8664},
		}}, []Id{ArchX8664}},
		{"same arch", MachinePlan{Compilers: []Compiler{
			Compiler{Service: Service{Id: IdFromString("gcc")}, Arch: ArchArm},
			Compiler{Service: Service{Id: IdFromString("clang")}, Arch: ArchArm},
		}}, []Id{ArchArm}},
		{"two arches", MachinePlan{Compilers: []Compiler{
			Compiler{Service: Service{Id: IdFromString("gcc-ppc")}, Arch: ArchPPC},
			Compiler{Service: Service{Id: IdFromString("clang-ppc")}, Arch: ArchPPC},
			Compiler{Service: Service{Id: IdFromString("gcc")}, Arch: ArchArm},
			Compiler{Service: Service{Id: IdFromString("clang")}, Arch: ArchArm},
		}}, []Id{ArchArm, ArchPPC}},
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
