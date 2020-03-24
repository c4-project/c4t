// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package herdtools

import "fmt"

// parserState is the state of the parsing FSA.
type parserState int

const (
	// psEmpty states that we haven't read anything yet.
	psEmpty parserState = iota
	// psPreTest states that we haven't hit the actual test yet.
	psPreTest
	// psPreamble states that we're in the pre-state matter.
	psPreamble
	// psState states that we're in a state block.
	psState
	// psSummary states that we're reading the summary tag.
	psSummary
	// psPostamble states that we're in the post-summary matter.
	psPostamble
)

// afterBegin advances the parser state after determining the observation is non-empty.
func (p *parser) afterBegin() error {
	return p.transition(psEmpty, psPreTest)
}

// afterPreTest advances the parser state after finishing the pre-test matter.
func (p *parser) afterPreTest() error {
	return p.transition(psPreTest, psPreamble)
}

// afterPreamble advances the parser state after parsing the state count.
func (p *parser) afterPreamble(nstates uint64) error {
	err := p.checkState(psPreamble)
	if err != nil {
		return err
	}

	p.setStateCount(nstates)
	return nil
}

// setStateCount sets the state and state count of the parser according to nstates.
func (p *parser) setStateCount(nstates uint64) {
	if nstates == 0 {
		p.state = psSummary
	} else {
		p.state = psState
		p.nstates = nstates
	}
}

// afterStateLine advances the parser state after a state line.
func (p *parser) afterStateLine() error {
	if err := p.checkState(psState); err != nil {
		return nil
	}

	p.nstates--
	if p.nstates == 0 {
		p.state = psSummary
	}
	return nil
}

func (p *parser) afterSummary() error {
	return p.transition(psSummary, psPostamble)
}

// transition handles a simple state transition between from and to.
// It returns an error if the current state isn't from.
func (p *parser) transition(from, to parserState) error {
	err := p.checkState(from)
	p.state = to
	return err
}

// checkFinalState checks to see if the parser has ended in an appropriate state, and returns an error if not.
func (p *parser) checkFinalState() error {
	switch p.state {
	case psEmpty:
		return ErrInputEmpty
	case psState:
		return fmt.Errorf("%w: %d state(s) remain", ErrNotEnoughStates, p.nstates)
	case psPreTest:
		return ErrNoTest
	case psPreamble:
		return ErrNoStates
	case psSummary:
		return ErrNoSummary
	case psPostamble:
		return nil
	default:
		return fmt.Errorf("%w: %v", ErrBadState, p.state)
	}
}

// checkState returns with an error if the current automaton state isn't want.
func (p *parser) checkState(want parserState) error {
	if p.state != want {
		return fmt.Errorf("%w: got=%v, want=%v", ErrBadTransition, p.state, want)
	}
	return nil
}
