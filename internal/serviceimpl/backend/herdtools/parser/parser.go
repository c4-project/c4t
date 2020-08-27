// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package parser contains logic for parsing Herd and Litmus
package parser

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/MattWindsor91/act-tester/internal/subject/obs"
)

// Parse parses an observation from r into o using i.
func Parse(i Impl, r io.Reader, o *obs.Obs) error {
	p := parser{impl: i, o: o}
	s := bufio.NewScanner(r)
	lineno := 1
	for s.Scan() {
		if err := p.processLine(s.Text()); err != nil {
			return fmt.Errorf("line %d: %w", lineno, err)
		}
		lineno++
	}
	if err := s.Err(); err != nil {
		return err
	}
	return p.checkFinalState()
}

// parser holds the state for a Herdtools parser.
type parser struct {
	// impl tells us how to perform the Herd/Litmus-specific parts of the parsing set-up.
	impl Impl

	// o is the observation we're creating.
	o *obs.Obs

	// tt is the test type, if any.
	tt TestType

	// state is the current state of the parsing FSA.
	state state

	// nstates is the number of states left to read if we're in psState.
	nstates uint64
}

// processLine processes a single line of a Herdtools observation.
func (p *parser) processLine(line string) error {
	if p.impl == nil {
		return ErrNoImpl
	}

	fields := strings.Fields(line)
	switch p.state {
	case psEmpty:
		return p.processEmpty(fields)
	case psPreTest:
		return p.processPreTest(fields)
	case psPreamble:
		return p.processPreamble(fields)
	case psState:
		return p.processState(fields)
	case psSummary:
		return p.processSummary(fields)
	case psPostamble:
		// TODO(@MattWindsor91): do something with this?
		return nil
	default:
		return fmt.Errorf("%w: %d", ErrBadState, p.state)
	}
}

func (p *parser) processEmpty(fields []string) error {
	if err := p.afterBegin(); err != nil {
		return err
	}
	return p.processPreTest(fields)
}

func (p *parser) processPreTest(fields []string) error {
	nf := len(fields)
	if nf == 0 || fields[0] != "Test" {
		return nil
	}

	if nf != 3 {
		return fmt.Errorf("%w: expected three fields, got %d", ErrBadTestType, nf)
	}

	var err error
	p.tt, err = parseTestType(fields[2])
	if err != nil {
		return err
	}

	return p.afterPreTest()
}

func (p *parser) processPreamble(fields []string) error {
	nstates, err := p.impl.ParseStateCount(fields)
	if err != nil {
		return err
	}
	return p.afterPreamble(nstates)
}

func (p *parser) processState(fields []string) error {
	sl, err := p.impl.ParseStateLine(p.tt, fields)
	if err != nil {
		return err
	}
	if err := p.processStateLine(sl); err != nil {
		return err
	}
	return p.afterStateLine()
}

func (p *parser) processSummary(fields []string) error {
	if nf := len(fields); nf != 1 {
		return fmt.Errorf("%w: expected one field, got %d", ErrBadSummary, nf)
	}
	var err error
	if p.o.Flags, err = parseFlag(fields[0]); err != nil {
		return err
	}
	return p.afterSummary()
}

// parseFlag parses the summary flag f as an observation flag.
func parseFlag(f string) (obs.Flag, error) {
	switch f {
	case "Yes": // seen in practice?
		fallthrough
	case "Ok":
		return obs.Sat, nil
	case "No":
		return obs.Unsat, nil
	case "Undef":
		return obs.Undef, nil
	default:
		return 0, fmt.Errorf("%w: bad flag %s", ErrBadSummary, f)
	}
}
