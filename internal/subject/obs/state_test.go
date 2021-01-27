// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package obs_test

import (
	"fmt"

	"github.com/c4-project/c4t/internal/subject/obs"
)

// ExampleState_Vars is a runnable example for State.Vars.
func ExampleState_Vars() {
	s := obs.State{
		"x": "1",
		"a": "2",
		"b": "3",
		"y": "4",
	}
	for _, v := range s.Vars() {
		fmt.Println(v)
	}

	// Output:
	// a
	// b
	// x
	// y
}
