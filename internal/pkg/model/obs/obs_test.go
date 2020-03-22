// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package obs_test

import (
	"fmt"
	"testing"

	"github.com/MattWindsor91/act-tester/internal/pkg/model/obs"

	"github.com/MattWindsor91/act-tester/internal/pkg/testhelp"
)

// ExampleObs_Sat is a testable example for Sat.
func ExampleObs_Sat() {
	o1 := obs.Obs{}
	fmt.Println("o1:", o1.Sat())

	o2 := obs.Obs{Flags: obs.Unsat}
	fmt.Println("o2:", o2.Sat())

	o3 := obs.Obs{Flags: obs.Sat}
	fmt.Println("o3:", o3.Sat())

	// Output:
	// o1: false
	// o2: false
	// o3: true
}

// ExampleObs_Unsat is a testable example for Unsat.
func ExampleObs_Unsat() {
	o1 := obs.Obs{}
	fmt.Println("o1:", o1.Unsat())

	o2 := obs.Obs{Flags: obs.Unsat}
	fmt.Println("o2:", o2.Unsat())

	o3 := obs.Obs{Flags: obs.Sat}
	fmt.Println("o3:", o3.Unsat())

	// Output:
	// o1: false
	// o2: true
	// o3: false
}

func TestObs_TOML_RoundTrip(t *testing.T) {
	t.Parallel()

	cases := map[string]obs.Obs{
		"empty":         {},
		"undef-nostate": {Flags: obs.Undef},
		"multiple-flags": {
			Flags: obs.Sat | obs.Undef,
			States: []obs.State{
				{"x": "27", "y": "53"},
				{"x": "27", "y": "42"},
			},
			Witnesses: []obs.State{
				{"x": "27", "y": "53"},
			},
		},
	}
	for name, want := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			testhelp.TestTomlRoundTrip(t, want, "Obs")
		})
	}
}
