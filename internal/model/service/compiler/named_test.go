// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package compiler_test

import (
	"fmt"

	"github.com/c4-project/c4t/internal/model/service/compiler/optlevel"

	"github.com/c4-project/c4t/internal/id"
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

// ExampleNamed_FullID is a runnable example for Named.FullID.
func ExampleNamed_FullID() {
	c := compiler.Instance{SelectedMOpt: "arch=skylake", SelectedOpt: &optlevel.Named{
		Name:  "3",
		Level: optlevel.Level{},
	}}
	i, _ := c.AddName(id.FromString("gcc.4")).FullID()
	fmt.Println(i)

	// Output:
	// gcc.4.o3.march=skylake
}

// ExampleNamed_FullID_dotarch is a runnable example for Named.FullID where the mopt contains dots.
func ExampleNamed_FullID_dotarch() {
	c := compiler.Instance{SelectedMOpt: "arch=armv8.1-a"}
	i, _ := c.AddName(id.FromString("gcc.8")).FullID()
	fmt.Println(i)

	// Output:
	// gcc.8.o.march=armv8_1-a
}
