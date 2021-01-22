// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package compiler_test

import (
	"fmt"
	"testing"

	"github.com/c4-project/c4t/internal/model/service/compiler/optlevel"

	"github.com/c4-project/c4t/internal/helper/testhelp"

	"github.com/c4-project/c4t/internal/model/id"
	"github.com/c4-project/c4t/internal/model/service/compiler"
)

// ExampleInstance_AddName is a runnable example for Instance.AddName.
func ExampleInstance_AddName() {
	c := compiler.Instance{
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

// ExampleInstance_AddNameString is a runnable example for Instance.AddName.
func ExampleInstance_AddNameString() {
	c := compiler.Instance{
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

// TestInstance_AddNameString_error exercises Instance.AddNameString's error path.
func TestInstance_AddNameString_error(t *testing.T) {
	t.Parallel()

	c := compiler.Instance{
		Compiler: compiler.Compiler{
			Style: id.CStyleGCC,
			Arch:  id.ArchX8664,
		},
	}
	_, err := c.AddNameString("foo..bar")
	testhelp.ExpectErrorIs(t, err, id.ErrTagEmpty, "AddNameString with empty-tag id")
}

// ExampleNamed_FullID is a runnable example for Named.FullID.
func ExampleNamed_FullID() {
	c := compiler.Instance{SelectedMOpt: "arch=skylake", SelectedOpt: &optlevel.Named{
		Name:  "3",
		Level: optlevel.Level{},
	}}
	n, _ := c.AddNameString("gcc.4")
	i, _ := n.FullID()
	fmt.Println(i)

	// Output:
	// gcc.4.o3.march=skylake
}
