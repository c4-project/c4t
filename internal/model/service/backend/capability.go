// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package backend

// Capability is the enumeration of things that a backend claims to be able to do.
type Capability uint8

const (
	// CanRunStandalone states that the backend can produce ToStandalone targets.
	CanRunStandalone = 1 << iota
	// CanProduceObj states that the backend's recipes can ToObjRecipe targets.
	CanProduceObj
	// CanProduceExe states that the backend's recipes can produce ToExeRecipe targets.
	CanProduceExe
	// CanLiftLitmus states that the backend can consume LiftLitmus sources.
	CanLiftLitmus
)
