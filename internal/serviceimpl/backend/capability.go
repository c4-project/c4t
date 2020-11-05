// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package backend

// Capability is the enumeration of things that a backend claims to be able to do.
type Capability uint8

const (
	// This backend can lift test-cases into recipes.
	CanLift Capability = 1 << iota
	// This backend can be run without a compile/
	CanRunStandalone
	// This backend's recipes can produce executable harnesses.
	CanProduceExecutables
)
