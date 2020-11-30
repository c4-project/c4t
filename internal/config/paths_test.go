// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package config_test

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/MattWindsor91/c4t/internal/config"
)

// ExamplePathset_OutPath is a runnable example for Pathset.OutPath.
func ExamplePathset_OutPath() {
	// OutDir can contain ~ for home expansions, but we can't easily show that in an example.
	ps := config.Pathset{OutDir: filepath.Join("foo", "bar", "baz")}
	if r, err := ps.OutPath("test.yaml"); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(filepath.ToSlash(r))
	}

	// Output:
	// foo/bar/baz/test.yaml
}

// ExamplePathset_FallbackToInputs is a runnable example for Pathset.FallbackToInputs
func ExamplePathset_FallbackToInputs() {
	// Inputs can contain ~ for home expansions, but we can't easily show that in an example.
	ps := config.Pathset{Inputs: []string{"foo", "bar", "baz"}}
	is1, _ := ps.FallbackToInputs(nil)
	fmt.Println(strings.Join(is1, ", "))
	is2, _ := ps.FallbackToInputs([]string{})
	fmt.Println(strings.Join(is2, ", "))
	is3, _ := ps.FallbackToInputs([]string{"foobaz", "barbaz"})
	fmt.Println(strings.Join(is3, ", "))

	// Output:
	// foo, bar, baz
	// foo, bar, baz
	// foobaz, barbaz
}
