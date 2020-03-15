// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package obs concerns 'observations': the end result of running a test harness on a particular machine.
package obs

// Obs represents an observation in ACT's JSON-based format.
type Obs struct {
	// Flags contains any flags that are active on Obs.
	Flags Flag `json:"flags,omitempty" toml:"flags,omitzero"`

	// CounterExamples lists all states that passed validation.
	CounterExamples []State `json:"counter_examples,omitempty" toml:"counter_examples,omitempty"`

	// Witnesses lists all states that passed validation.
	Witnesses []State `json:"witnesses,omitempty" toml:"witnesses,omitempty"`

	// States lists all observed states.
	States []State `json:"states" toml:"states,omitempty"`
}

// Sat gets whether the observation satisfies its validation.
func (o *Obs) Sat() bool {
	return o.Flags.Has(Sat)
}

// Unsat gets whether the observation does not satisfy its validation.
func (o *Obs) Unsat() bool {
	return o.Flags.Has(Unsat)
}
