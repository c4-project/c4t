// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package optlevel_test

import (
	"fmt"
	"testing"

	"github.com/c4-project/c4t/internal/model/service/compiler/optlevel"
)

// ExampleSelection_Override is a testable example for Override.
func ExampleSelection_Override() {
	sel := optlevel.Selection{
		Enabled:  []string{"g", "s"},
		Disabled: []string{"fast"},
	}
	for o := range sel.Override(map[string]struct{}{
		"1":    {},
		"2":    {},
		"3":    {},
		"fast": {},
	}) {
		fmt.Println(o)
	}

	// Unordered output:
	// 1
	// 2
	// 3
	// g
	// s
}

// TestSelection_Override tests Override with a variety of cases.
func TestSelection_Override(t *testing.T) {
	t.Parallel()
}
