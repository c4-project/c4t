// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package obs_test

import (
	"fmt"
	"os"

	"github.com/c4-project/c4t/internal/subject/obs"
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

// ExamplePretty_dnf_empty shows the result of pretty-printing an interesting existential.
func ExamplePretty_interesting_exists() {
	o := obs.Obs{
		Flags: obs.Sat | obs.Exist,
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

	if err := obs.Pretty(os.Stdout, o, obs.PrettyMode{Interesting: true}); err != nil {
		fmt.Println("error:", err)
	}

	// Output:
	// postcondition witnessed by:
	//   x = 0, y = 1
	//   x = 1, y = 1
}
