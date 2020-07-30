// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package compiler_test

import (
	"fmt"
	"testing"

	"github.com/MattWindsor91/act-tester/internal/helper/testhelp"

	"github.com/MattWindsor91/act-tester/internal/model/id"
	"github.com/MattWindsor91/act-tester/internal/model/service/compiler"
)

// ExampleConfiguration_AddName is a runnable example for AddName.
func ExampleConfiguration_AddName() {
	c := compiler.Configuration{
		Compiler: compiler.Compiler{
			Style: id.CStyleGCC,
			Arch:  id.ArchX8664,
		},
	}
	nc := c.AddName(id.FromString("gcc.native"))

	fmt.Println(nc.ID)
	fmt.Println(nc.Arch)
	fmt.Println(nc.Style)

	// Output:
	// gcc.native
	// x86.64
	// gcc
}

// ExampleConfiguration_AddNameString is a runnable example for AddName.
func ExampleConfiguration_AddNameString() {
	c := compiler.Configuration{
		Compiler: compiler.Compiler{
			Style: id.CStyleGCC,
			Arch:  id.ArchX8664,
		},
	}
	nc, _ := c.AddNameString("gcc.8")

	fmt.Println(nc.ID)
	fmt.Println(nc.Arch)
	fmt.Println(nc.Style)

	// Output:
	// gcc.8
	// x86.64
	// gcc
}

// TestConfiguration_AddNameString_error exercises AddNameString's error path.
func TestConfiguration_AddNameString_error(t *testing.T) {
	t.Parallel()

	c := compiler.Configuration{
		Compiler: compiler.Compiler{
			Style: id.CStyleGCC,
			Arch:  id.ArchX8664,
		},
	}
	_, err := c.AddNameString("foo..bar")
	testhelp.ExpectErrorIs(t, err, id.ErrTagEmpty, "AddNameString with empty-tag id")
}
