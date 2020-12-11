// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package obs_test

import (
	"fmt"
	"os"

	"github.com/MattWindsor91/c4t/internal/subject/obs"
)

// ExamplePretty_dnf shows the result of pretty-printing an observation to a DNF postcondition.
func ExamplePretty_dnf() {
	o := obs.Obs{
		Flags: obs.Sat,
		CounterExamples: []obs.State{
			{"x": "1", "y": "0"},
			{"x": "0", "y": "0"},
		},
		Witnesses: []obs.State{
			{"x": "0", "y": "1"},
			{"x": "1", "y": "1"},
		},
		States: []obs.State{
			{"x": "1", "y": "0"},
			{"x": "0", "y": "0"},
			{"x": "0", "y": "1"},
			{"x": "1", "y": "1"},
		},
	}

	if err := obs.Pretty(os.Stdout, o, obs.PrettyMode{Dnf: true}); err != nil {
		fmt.Println("error:", err)
	}

	// Output:
	// forall (
	//      (x == 1 /\ y == 0)
	//   \/ (x == 0 /\ y == 0)
	//   \/ (x == 0 /\ y == 1)
	//   \/ (x == 1 /\ y == 1)
	// )
}

// ExamplePretty_dnf_empty shows the result of pretty-printing an empty observation to a DNF postcondition.
func ExamplePretty_dnf_empty() {
	if err := obs.Pretty(os.Stdout, (obs.Obs{}), obs.PrettyMode{Dnf: true}); err != nil {
		fmt.Println("error:", err)
	}

	// Output:
	// forall (
	//   true
	// )
}
