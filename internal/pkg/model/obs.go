// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package model

// An observed state.
type ObsState map[string]string

// Obs represents an observation in ACT's JSON-based format.
type Obs struct {
	// Flags contains any flags that are active on Obs.
	Flags ObsFlag `json:"flags,omitempty" toml:"flags,omitzero"`

	// CounterExamples lists all states that passed validation.
	CounterExamples []ObsState `json:"counter_examples,omitempty" toml:"counter_examples,omitempty"`

	// Witnesses lists all states that passed validation.
	Witnesses []ObsState `json:"witnesses,omitempty" toml:"witnesses,omitempty"`

	// States lists all observed states.
	States []ObsState `json:"states" toml:"states,omitempty"`
}

// Sat gets whether the observation satisfies its validation.
func (o *Obs) Sat() bool {
	return o.Flags.Has(ObsSat)
}

// Unsat gets whether the observation does not satisfy its validation.
func (o *Obs) Unsat() bool {
	return o.Flags.Has(ObsUnsat)
}
