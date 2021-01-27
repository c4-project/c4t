// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package config_test

import (
	"fmt"

	"github.com/c4-project/c4t/internal/quantity"

	"github.com/c4-project/c4t/internal/config"
)

// ExampleConfig_DisableFuzz is a runnable example for Config.DisableFuzz.
func ExampleConfig_DisableFuzz() {
	var c config.Config
	c.DisableFuzz()
	fmt.Println(c.Fuzz.Disabled)
	c.Fuzz.Disabled = false
	fmt.Println(c.Fuzz.Disabled)
	c.DisableFuzz()
	fmt.Println(c.Fuzz.Disabled)

	// Output:
	// true
	// false
	// true
}

func ExampleConfig_OverrideQuantities() {
	c := config.Config{
		Quantities: quantity.RootSet{
			Plan: quantity.PlanSet{
				NWorkers: 10,
			},
			MachineSet: quantity.MachineSet{
				Fuzz: quantity.FuzzSet{
					CorpusSize:    20,
					SubjectCycles: 4,
					NWorkers:      5,
				},
			},
		},
	}
	c.OverrideQuantities(quantity.RootSet{
		MachineSet: quantity.MachineSet{
			Fuzz: quantity.FuzzSet{
				SubjectCycles: 8,
			},
		},
	})

	fmt.Println(c.Quantities.Plan.NWorkers)
	fmt.Println(c.Quantities.MachineSet.Fuzz.CorpusSize)
	fmt.Println(c.Quantities.MachineSet.Fuzz.SubjectCycles)

	// Output:
	// 10
	// 20
	// 8
}
