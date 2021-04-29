// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package backend_test

import (
	"fmt"

	"github.com/c4-project/c4t/internal/model/service/backend"
)

// ExampleCapability_Satisfies is a runnable example for Capability.Satisfies.
func ExampleCapability_Satisfies() {
	c := backend.Capability(backend.CanLiftLitmus | backend.CanRunStandalone)
	fmt.Println(c.Satisfies(0))
	fmt.Println(c.Satisfies(backend.CanLiftLitmus))
	fmt.Println(c.Satisfies(backend.CanRunStandalone))
	fmt.Println(c.Satisfies(backend.CanProduceObj))
	fmt.Println(c.Satisfies(backend.CanLiftLitmus | backend.CanProduceObj))

	// output:
	// true
	// true
	// true
	// false
	// false
}

// ExampleCapability_String is a runnable example for Capability.String.
func ExampleCapability_String() {
	fmt.Println(backend.Capability(0))
	fmt.Println(backend.Capability(backend.CanLiftLitmus))
	fmt.Println(backend.Capability(backend.CanRunStandalone))
	fmt.Println(backend.Capability(backend.CanProduceObj))
	fmt.Println(backend.Capability(backend.CanLiftLitmus | backend.CanProduceObj))

	// output:
	//
	// lift-litmus
	// run-standalone
	// produce-obj
	// produce-obj+lift-litmus
}
