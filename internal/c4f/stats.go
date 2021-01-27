// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package c4f

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/c4-project/c4t/internal/model/litmus"

	"github.com/c4-project/c4t/internal/model/id"
)

// ErrStatsetParse occurs when there is a parse error reading a statset.
var ErrStatsetParse = errors.New("statistic parse error")

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

// ParseStats parses a statistics set from r into statistics set s.
// Each statistic should be in the form "name value\n".
func ParseStats(s *litmus.Statset, r io.Reader) error {
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		if err := parseLine(s, sc.Text()); err != nil {
			return err
		}
	}
	return sc.Err()
}

func parseLine(s *litmus.Statset, l string) error {
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
	return setByID(s, nid, val)
}

func errBadSubcategory(cat string, rest id.ID) error {
	return fmt.Errorf("%w: unexpected subcategory for %q: %s", ErrStatsetParse, cat, rest)
}

func errEmptyId(cat string) error {
	return fmt.Errorf("%w: empty %q sub-id", ErrStatsetParse, cat)
}

func setByID(s *litmus.Statset, nid id.ID, val int) error {
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
		return setAtomicType(s, nrest, val)
	case catMemOrders:
		return setMemOrder(s, nrest, val)
	case catLiterals:
		return setLiteral(s, nrest, val)
	}
	return nil
}

func setLiteral(s *litmus.Statset, lid id.ID, val int) error {
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

func setAtomicType(s *litmus.Statset, tid id.ID, val int) error {
	return setAtomicRelated(tid, val, catAtomics, s.AtomicExpressions.AddType, s.AtomicStatements.AddType)
}

func setMemOrder(s *litmus.Statset, mid id.ID, val int) error {
	return setAtomicRelated(mid, val, catMemOrders, s.AtomicExpressions.AddMemOrder, s.AtomicStatements.AddMemOrder)
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
