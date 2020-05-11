// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package model

import (
	"context"

	"github.com/MattWindsor91/act-tester/internal/model/id"
)

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
