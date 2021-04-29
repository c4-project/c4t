// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package backend

import "strings"

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

// Satisfies is true if this capability contains all capabilities in c2.
func (c Capability) Satisfies(c2 Capability) bool {
	return c&c2 == c2
}

// String gets a stringified form of the capability
func (c Capability) String() string {
	parts := make([]string, 0, 4)
	for _, s := range []struct {
		cap Capability
		str string
	}{
		{cap: CanRunStandalone, str: "run-standalone"},
		{cap: CanProduceObj, str: "produce-obj"},
		{cap: CanProduceExe, str: "produce-exe"},
		{cap: CanLiftLitmus, str: "lift-litmus"},
	} {
		if c.Satisfies(s.cap) {
			parts = append(parts, s.str)
		}
	}
	return strings.Join(parts, "+")
}
