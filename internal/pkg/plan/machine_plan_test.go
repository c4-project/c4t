package plan

import (
	"reflect"
	"testing"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"
)

// TestMachinePlan_Arches tests the Arches method on MachinePlan.
func TestMachinePlan_Arches(t *testing.T) {
	tests := []struct {
		name string
		plan MachinePlan
		want []model.ID
	}{
		{"no arches", MachinePlan{}, []model.ID{}},
		{"one compiler", MachinePlan{Compilers: []model.Compiler{
			model.Compiler{Service: model.Service{ID: model.IDFromString("gcc")}, Arch: model.ArchX8664},
		}}, []model.ID{model.ArchX8664}},
		{"same arch", MachinePlan{Compilers: []model.Compiler{
			model.Compiler{Service: model.Service{ID: model.IDFromString("gcc")}, Arch: model.ArchArm},
			model.Compiler{Service: model.Service{ID: model.IDFromString("clang")}, Arch: model.ArchArm},
		}}, []model.ID{model.ArchArm}},
		{"two arches", MachinePlan{Compilers: []model.Compiler{
			model.Compiler{Service: model.Service{ID: model.IDFromString("gcc-ppc")}, Arch: model.ArchPPC},
			model.Compiler{Service: model.Service{ID: model.IDFromString("clang-ppc")}, Arch: model.ArchPPC},
			model.Compiler{Service: model.Service{ID: model.IDFromString("gcc")}, Arch: model.ArchArm},
			model.Compiler{Service: model.Service{ID: model.IDFromString("clang")}, Arch: model.ArchArm},
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
