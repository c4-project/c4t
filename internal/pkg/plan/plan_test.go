package plan_test

import (
	"context"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/MattWindsor91/act-tester/internal/pkg/corpus"

	"github.com/MattWindsor91/act-tester/internal/pkg/testhelp"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"
	"github.com/MattWindsor91/act-tester/internal/pkg/plan"
)

var testPlan = plan.Plan{
	Creation: time.Time{},
	Seed:     0,
	Machines: testMachines,
	Corpus:   testCorpus,
}

var testMachines = map[string]plan.MachinePlan{
	"localhost": {
		Machine: model.Machine{Cores: 4},
		Backend: model.Backend{
			ID:          model.IDFromString("litmus7"),
			IDQualified: true,
			MachineID:   nil,
			Style:       model.IDFromString("litmus"),
		},
		Compilers: map[string]model.Compiler{
			"gcc": {
				Style: model.IDFromString("gcc"),
				Arch:  model.ArchX8664,
			},
		},
	},
	"holly": {
		Machine: model.Machine{Cores: 10},
		Backend: model.Backend{
			ID:          model.IDFromString("litmus"),
			IDQualified: true,
			MachineID:   nil,
			Style:       model.IDFromString("litmus"),
		},
		Compilers: map[string]model.Compiler{
			"clang": {
				Style: model.IDFromString("clang"),
				Arch:  model.ArchArm,
			},
		},
	},
}

var testCorpus = corpus.New("a.litmus", "b.litmus", "c.litmus")

// TestPlan_ParMachines tests ParMachines by spinning up a basic computation on the test machine set.
func TestPlan_ParMachines(t *testing.T) {
	var sm sync.Map

	// This should store each machine's number of cores to sm.
	err := testPlan.ParMachines(context.Background(),
		func(ctx context.Context, id model.ID, machinePlan plan.MachinePlan) error {
			sm.Store(id.String(), machinePlan.Cores)
			return nil
		},
		func(ctx context.Context) error {
			sm.Store("aux", 5)
			return nil
		},
	)
	if err != nil {
		t.Fatalf("unexpected err in ParMachines: %v", err)
	}
	for ids, p := range testMachines {
		want := p.Cores
		gint, ok := sm.Load(ids)
		if !ok {
			t.Errorf("ParMachines didn't run for machine %s?", ids)
		}
		got, gok := gint.(int)
		if !gok {
			t.Errorf("ParMachines didn't store int for machine %s, got=%v", ids, gint)
		}
		if got != want {
			t.Errorf("ParMachines stored %d machine %s; want %d", got, ids, want)
		}
	}
}

// TestPlan_Machine_Present tests the Machine method for a present machine.
func TestPlan_Machine_Present(t *testing.T) {
	// TODO(@MattWindsor91): test ID
	_, got, err := testPlan.Machine(model.IDFromString("localhost"))
	if err != nil {
		t.Fatalf("error retrieving known-present machine: %v", err)
	}
	want := testMachines["localhost"]
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Machine() for known-present machine got=%v; want=%v", got, want)
	}
}

// TestPlan_Machine_Present tests the Machine method for an absent machine.
func TestPlan_Machine_Absent(t *testing.T) {
	_, _, err := testPlan.Machine(model.IDFromString("mu"))
	testhelp.ExpectErrorIs(t, err, plan.ErrNoMachine, "Machine() for known-absent machine")
}
