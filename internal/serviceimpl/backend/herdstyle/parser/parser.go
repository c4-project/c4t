// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

// Package parser contains logic for parsing Herd-style observations.
package parser

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/c4-project/c4t/internal/subject/obs"
)

// Parse parses an observation from r into o using i.
func Parse(i Impl, r io.Reader, o *obs.Obs) error {
	p := parser{impl: i, o: o}
	return p.parse(r)
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

// parse parses r into this parser.
func (p *parser) parse(r io.Reader) error {
	if p.impl == nil {
		return ErrNoImpl
	}

	if err := p.parseLines(r); err != nil {
		return err
	}

	return p.checkFinalState()
}

// parseLine processes the lines of a Herdtools observation from reader r.
func (p *parser) parseLines(r io.Reader) error {
	s := bufio.NewScanner(r)
	lineno := 1
	for s.Scan() {
		if err := p.parseLine(s.Text()); err != nil {
			return fmt.Errorf("line %d (%q): %w", lineno, s.Text(), err)
		}
		lineno++
	}
	return s.Err()
}

// parseLine processes a single line of a Herdtools observation.
func (p *parser) parseLine(line string) error {
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
	if nf == 0 {
		return nil
	}
	if fields[0] != "Test" {
		return p.processPreTestImplHooks(fields)
	}
	if nf != 3 {
		return fmt.Errorf("%w: expected three fields, got %d", ErrBadTestType, nf)
	}
	return p.processTestType(fields)
}

func (p *parser) processTestType(fields []string) error {
	var err error
	if p.tt, err = parseTestType(fields[2]); err != nil {
		return err
	}
	p.o.Flags |= p.tt.Flags()

	return p.afterPreTest()
}

// processPreTestImplHooks handles passing a pre-test line to the implementation to scan for flags.
func (p *parser) processPreTestImplHooks(fields []string) error {
	f, err := p.impl.ParsePreTestLine(fields)
	if err != nil {
		return err
	}
	p.o.Flags |= f
	return nil
}

func (p *parser) processPreamble(fields []string) error {
	nstates, ok, err := p.impl.ParseStateCount(fields)
	if err != nil {
		return err
	}
	if !ok {
		// Skip this line.
		return nil
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
	// Herd and Litmus themselves always follow the final state line with a summary;
	// the leniency here, as often is the case, is mainly for rmem and other herd-a-likes.
	nf := len(fields)
	// Some summary lines might be 'Flag (description)'.
	if nf == 0 {
		return nil
	}
	f, ok := parseFlag(fields[0])
	if !ok {
		// We want to catch possible mismatches between state count and state lines.
		return p.errorIfStateLine(fields)
	}
	// Making sure not to override any partiality flags already parsed.
	p.o.Flags |= f
	return p.afterSummary()
}

func (p *parser) errorIfStateLine(fields []string) error {
	if _, err := p.impl.ParseStateLine(p.tt, fields); err != nil {
		// Intentional
		return nil
	}
	return fmt.Errorf("%w: possible extraneous state line", ErrBadSummary)
}

// parseFlag parses the summary flag f as an observation flag.
func parseFlag(f string) (flag obs.Flag, ok bool) {
	ok = true
	switch f {
	case "Ok":
		flag = obs.Sat
	case "No":
		flag = obs.Unsat
	case "Undef":
		flag = obs.Undef
	default:
		ok = false
	}
	return
}
