// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package parser

import (
	"fmt"
	"strings"

	"github.com/MattWindsor91/act-tester/internal/subject/obs"
)

// StateLine is the struct that an implementation of the Herdtools parser must return when parsing a state line.
type StateLine struct {
	// NOccurs is the number of times this state has been observed.
	// If zero, there was no information about occurrences.
	NOccurs uint64
	// Tag is the observation tag (witness, counter-example, etc) of this observation.
	Tag obs.Tag
	// Rest is the remaining fields of the state line, which should be of the form 'x=y;'.
	Rest []string
}

// processStateLine performs the herdtools-common processing needed to add sl into an observation.
func (p *parser) processStateLine(sl *StateLine) error {
	s, err := parseState(sl.Rest)
	if err != nil {
		return err
	}
	// TODO(@MattWindsor91): number of occurrences?
	p.o.AddState(sl.Tag, s)
	return nil
}

// parseState parses a state from the mappings in fields.
func parseState(fields []string) (obs.State, error) {
	s := make(obs.State, len(fields))
	for _, f := range fields {
		k, v, err := parseStateMapping(f)
		if err != nil {
			return nil, err
		}
		s[k] = v
	}
	return s, nil
}

// parseStateMapping parses a state mapping in the form 'x=y;' at raw.
func parseStateMapping(raw string) (key string, val string, err error) {
	chop := strings.Split(strings.TrimSuffix(raw, ";"), "=")
	if len(chop) != 2 {
		return "", "", fmt.Errorf("%w: expected mapping of form 'x=y;', got %q", ErrBadStateLine, raw)
	}
	return strings.TrimSpace(chop[0]), strings.TrimSpace(chop[1]), nil
}
