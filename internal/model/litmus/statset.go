// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package litmus

import (
	"context"

	"github.com/c4-project/c4t/internal/id"
)

// StatDumper is the interface of things that can dump statistics for a litmus test.
type StatDumper interface {
	// DumpStats populates s with statistics gleaned from the Litmus file at filepath path.
	DumpStats(ctx context.Context, s *Statset, path string) error
}

//go:generate mockery --name=StatDumper

// AtomicStatset contains a set of statistics about atomics (expressions or statements).
type AtomicStatset struct {
	// Types gives the types of atomic, categorised by type.
	Types map[id.ID]int `toml:"types,omitzero,omitempty" json:"types,omitempty"`
	// MemOrders gives the types of memory order, categorised by type.
	MemOrders map[id.ID]int `toml:"mem_orders,omitzero,omitempty" json:"mem_orders,omitempty"`
}

// AddType adds k to the type with ID ty.
func (s *AtomicStatset) AddType(ty id.ID, k int) {
	if s.Types == nil {
		s.Types = make(map[id.ID]int)
	}
	s.Types[ty] += k
}

// AddMemOrder adds k to the memory order with ID mo.
func (s *AtomicStatset) AddMemOrder(mo id.ID, k int) {
	if s.MemOrders == nil {
		s.MemOrders = make(map[id.ID]int)
	}
	s.MemOrders[mo] += k
}

// Statset contains a set of statistics acquired from `c4f-c dump-stats`.
type Statset struct {
	// Threads is the number of threads.
	Threads int `json:"threads,omitempty"`

	// Returns is the number of return statements.
	Returns int `json:"returns,omitempty"`

	// LiteralBools is the number of Boolean literals (true, false, etc).
	LiteralBools int `json:"literal_bools,omitempty"`

	// AtomicExpressions gives information about atomic statements.
	AtomicExpressions AtomicStatset `json:"atomic_expressions,omitempty"`

	// AtomicStatements gives information about atomic statements.
	AtomicStatements AtomicStatset `json:"atomic_statements,omitempty"`
}
