// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package act

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/MattWindsor91/act-tester/internal/model/id"
	"io"
	"strconv"
	"strings"
)

// ErrStatsetParse occurs when there is a parse error reading a statset.
var ErrStatsetParse = errors.New("statistic parse error")

// StatDumper is the interface of things that can dump statistics for a litmus test.
type StatDumper interface {
	// DumpStats populates s with statistics gleaned from the Litmus file at path.
	DumpStats(ctx context.Context, s *Statset, path string) error
}

// AtomicStatset contains a set of statistics about atomics (expressions or statements).
type AtomicStatset struct {
	// Types gives the types of atomic, categorised by type.
	Types map[string]int
	// MemOrders gives the types of memory order, categorised by type.
	MemOrders map[string]int
}

// AddType adds k to the type with ID id.
func (s *AtomicStatset) AddType(id id.ID, k int) {
	if s.Types == nil {
		s.Types = make(map[string]int)
	}
	s.Types[id.String()] += k
}

// AddMemOrder adds k to the memory order with ID id.
func (s *AtomicStatset) AddMemOrder(id id.ID, k int) {
	if s.MemOrders == nil {
		s.MemOrders = make(map[string]int)
	}
	s.MemOrders[id.String()] += k
}

// Statset contains a set of statistics acquired from `act-c dump-stats`.
type Statset struct {
	// Threads is the number of threads.
	Threads int

	// Returns is the number of return statements.
	Returns int

	// LiteralBools is the number of Boolean literals (true, false, etc).
	LiteralBools int

	// AtomicExpressions gives information about atomic statements.
	AtomicExpressions AtomicStatset

	// AtomicStatements gives information about atomic statements.
	AtomicStatements AtomicStatset
}

const (
	catAtomics           = "atomics"
	catAtomicsExpression = "expression"
	catAtomicsStatement  = "statement"
	catLiterals          = "literals"
	catLiteralsBool      = "bool"
	catMemOrders         = "mem-orders"
	catThreads           = "threads"
	catReturns           = "returns"
)

// Parse parses a statistics set from r into this statistics set.
// Each statistic should be in the form "name value\n".
func (s *Statset) Parse(r io.Reader) error {
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		if err := s.parseLine(sc.Text()); err != nil {
			return err
		}
	}
	return sc.Err()
}

func (s *Statset) parseLine(l string) error {
	fs := strings.Fields(l)
	if len(fs) != 2 {
		return fmt.Errorf("%w: %d fields in %q; want 2", ErrStatsetParse, len(fs), fs)
	}

	name := fs[0]
	val, err := strconv.Atoi(fs[1])
	if err != nil {
		return err
	}
	nid, err := id.TryFromString(name)
	if err != nil {
		return err
	}
	return s.setByID(nid, val)
}

func errBadSubcategory(cat string, rest id.ID) error {
	return fmt.Errorf("%w: unexpected subcategory for %q: %s", ErrStatsetParse, cat, rest)
}

func errEmptyId(cat string) error {
	return fmt.Errorf("%w: empty %q sub-id", ErrStatsetParse, cat)
}

func (s *Statset) setByID(nid id.ID, val int) error {
	ncat, nrest, ok := nid.Uncons()
	if !ok {
		return fmt.Errorf("%w: empty id", ErrStatsetParse)
	}
	switch ncat {
	case catThreads:
		if !nrest.IsEmpty() {
			return errBadSubcategory(catThreads, nrest)
		}
		s.Threads = val
	case catReturns:
		if !nrest.IsEmpty() {
			return errBadSubcategory(catReturns, nrest)
		}
		s.Returns = val
	case catAtomics:
		return s.setAtomicType(nrest, val)
	case catMemOrders:
		return s.setMemOrder(nrest, val)
	case catLiterals:
		return s.setLiteral(nrest, val)
	}
	return nil
}

func (s *Statset) setLiteral(lid id.ID, val int) error {
	lcat, lrest, ok := lid.Uncons()
	if !ok {
		return errEmptyId(catLiterals)
	}
	switch lcat {
	case catLiteralsBool:
		if !lrest.IsEmpty() {
			return fmt.Errorf("%w: unexpected subcategory for '%s.%s': %s", ErrStatsetParse, catLiterals, catLiteralsBool, lrest)
		}
		s.LiteralBools = val
	}
	return nil
}

func setAtomicRelated(aid id.ID, val int, pcat string, exp, stm func(id.ID, int)) error {
	cat, rest, ok := aid.Uncons()
	if !ok {
		return errEmptyId(pcat)
	}
	switch cat {
	case catAtomicsExpression:
		exp(rest, val)
	case catAtomicsStatement:
		stm(rest, val)
	}
	return nil
}

func (s *Statset) setAtomicType(tid id.ID, val int) error {
	return setAtomicRelated(tid, val, catAtomics, s.AtomicExpressions.AddType, s.AtomicStatements.AddType)
}

func (s *Statset) setMemOrder(mid id.ID, val int) error {
	return setAtomicRelated(mid, val, catMemOrders, s.AtomicExpressions.AddMemOrder, s.AtomicStatements.AddMemOrder)
}
