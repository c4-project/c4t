// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package obs concerns 'observations': the end result of running a test harness on a particular machine.
package obs

import "github.com/MattWindsor91/act-tester/internal/model/status"

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

// AddState adds the state s in accordance with the tag t.
func (o *Obs) AddState(t Tag, s State) {
	o.States = append(o.States, s)
	switch t {
	case TagWitness:
		o.Witnesses = append(o.Witnesses, s)
	case TagCounter:
		o.CounterExamples = append(o.CounterExamples, s)
	}
}

// Status determines the status of an observation o.
func (o *Obs) Status() status.Status {
	// TODO(@MattWindsor91): allow interestingness criteria
	if o.Unsat() {
		return status.Flagged
	}
	return status.Ok
}
