// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package obs

// Tag classifies a state line.
type Tag int

const (
	// TagUnknown represents a state that is not known to be either a witness or a counter-example.
	TagUnknown Tag = iota
	// TagWitness represents a state that validates a condition.
	TagWitness
	// TagCounter represents a state that goes against a condition.
	TagCounter
)
