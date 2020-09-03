// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package obs_test

import (
	"fmt"
	"testing"

	"github.com/MattWindsor91/act-tester/internal/subject/status"

	"github.com/MattWindsor91/act-tester/internal/subject/obs"

	"github.com/MattWindsor91/act-tester/internal/helper/testhelp"
)

// ExampleObs_Sat is a testable example for Obs.Sat.
func ExampleObs_Sat() {
	fmt.Println("empty:", (&obs.Obs{}).Sat())
	fmt.Println("unsat:", (&obs.Obs{Flags: obs.Unsat}).Sat())
	fmt.Println("e-sat:", (&obs.Obs{Flags: obs.Sat | obs.Exist}).Sat())

	// Output:
	// empty: false
	// unsat: false
	// e-sat: true
}

// ExampleObs_Unsat is a testable example for Obs.Unsat.
func ExampleObs_Unsat() {
	fmt.Println("empty:", (&obs.Obs{}).Unsat())
	fmt.Println("unsat:", (&obs.Obs{Flags: obs.Unsat}).Unsat())
	fmt.Println("e-sat:", (&obs.Obs{Flags: obs.Sat | obs.Exist}).Unsat())

	// Output:
	// empty: false
	// unsat: true
	// e-sat: false
}

// Used to avoid compiler-optimising-out of the benchmark below.
var stat status.Status

// BenchmarkObs_Status benchmarks Obs.Status.
func BenchmarkObs_Status(b *testing.B) {
	cases := map[string]*obs.Obs{
		"empty":   {},
		"undef":   {Flags: obs.Undef},
		"sat":     {Flags: obs.Sat},
		"unsat":   {Flags: obs.Unsat},
		"e-sat":   {Flags: obs.Sat | obs.Exist},
		"e-unsat": {Flags: obs.Unsat | obs.Exist},
	}

	for name, c := range cases {
		c := c
		b.Run(name, func(b *testing.B) {
			s := stat
			for i := 0; i < b.N; i++ {
				s = c.Status()
			}
			stat = s
		})

	}

}

// ExampleObs_Status is a testable example for Obs.Status.
func ExampleObs_Status() {
	fmt.Println("empty:  ", (&obs.Obs{}).Status())
	fmt.Println("undef:  ", (&obs.Obs{Flags: obs.Undef}).Status())
	fmt.Println("sat:    ", (&obs.Obs{Flags: obs.Sat}).Status())
	fmt.Println("unsat:  ", (&obs.Obs{Flags: obs.Unsat}).Status())
	fmt.Println("e-sat:  ", (&obs.Obs{Flags: obs.Sat | obs.Exist}).Status())
	fmt.Println("e-unsat:", (&obs.Obs{Flags: obs.Unsat | obs.Exist}).Status())

	// output:
	// empty:   Flagged
	// undef:   Flagged
	// sat:     Ok
	// unsat:   Flagged
	// e-sat:   Flagged
	// e-unsat: Ok
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
