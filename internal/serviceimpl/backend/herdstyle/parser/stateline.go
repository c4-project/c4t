// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package parser

import (
	"fmt"
	"strings"

	"github.com/c4-project/c4t/internal/subject/obs"
)

// StateLine is the struct that an implementation of the Herdtools parser must return when parsing a state line.
type StateLine struct {
	// State is the state parsed, less the valuation.
	State obs.State
	// Rest is the remaining fields of the state line, which should be of the form 'x=y;'.
	Rest []string
}

// processStateLine performs the herdtools-common processing needed to add sl into an observation.
func (p *parser) processStateLine(sl *StateLine) error {
	var err error
	if sl.State.Values, err = parseValuation(sl.Rest); err != nil {
		return err
	}
	// TODO(@MattWindsor91): number of occurrences?
	p.o.States = append(p.o.States, sl.State)
	return nil
}

// parseValuation parses a valuation from the mappings in fields.
func parseValuation(fields []string) (obs.Valuation, error) {
	s := make(obs.Valuation, len(fields))
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
