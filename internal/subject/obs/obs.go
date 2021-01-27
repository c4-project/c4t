// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

// Package obs concerns 'observations': the end result of running a test on a particular machine.
package obs

import "github.com/c4-project/c4t/internal/subject/status"

// Obs represents an observation in C4's JSON-based format.
type Obs struct {
	// Flags contains any flags that are active on Obs.
	Flags Flag `json:"flags,omitempty"`
	// States lists all states in this observation.
	States []State `json:"states,omitempty"`
}

// Status determines the status of an observation o.
//
// Currently, an observation is considered to be 'ok' if it is a satisfied universal or unsatisfied existential,
// and 'flagged' otherwise.
func (o Obs) Status() status.Status {
	if o.Flags.IsInteresting() {
		return status.Flagged
	}
	return status.Ok
}

// Witnesses gets the list of witnessing states in this observation.
func (o Obs) Witnesses() []State {
	return o.WithTag(TagWitness)
}

// CounterExamples gets the list of counter-example states in this observation.
func (o Obs) CounterExamples() []State {
	return o.WithTag(TagCounter)
}

// WithTag gets the list of states with tag t in this observation.
func (o Obs) WithTag(t Tag) []State {
	xs := make([]State, 0, len(o.States))
	for _, s := range o.States {
		if s.Tag == t {
			xs = append(xs, s)
		}
	}
	return xs
}

// State represents a single state in C4's JSON-based format.
type State struct {
	// Tag is the kind of state this is.
	Tag Tag `json:"tag,omitempty"`
	// Occurrences is the number of times this state was observed.
	// If this number is zero, there was no occurrence reporting for this state;
	// states which were observed zero times will not appear in the observation at all.
	Occurrences uint64 `json:"occurrences,omitempty"`
	// Values is the valuation for this state.
	Values Valuation `json:"values,omitempty"`
}
