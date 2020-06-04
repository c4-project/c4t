// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package litmus

import (
	"context"

	"github.com/MattWindsor91/act-tester/internal/model/id"
)

// StatDumper is the interface of things that can dump statistics for a litmus test.
type StatDumper interface {
	// DumpStats populates s with statistics gleaned from the Litmus file at path.
	DumpStats(ctx context.Context, s *Statset, path string) error
}

//go:generate mockery -name StatDumper

// AtomicStatset contains a set of statistics about atomics (expressions or statements).
type AtomicStatset struct {
	// Types gives the types of atomic, categorised by type.
	Types map[string]int `toml:"types,omitzero,omitempty" json:"types,omitempty"`
	// MemOrders gives the types of memory order, categorised by type.
	MemOrders map[string]int `toml:"mem_orders,omitzero,omitempty" json:"mem_orders,omitempty"`
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
	Threads int `toml:"threads,omitzero" json:"threads,omitempty"`

	// Returns is the number of return statements.
	Returns int `toml:"returns,omitzero" json:"returns,omitempty"`

	// LiteralBools is the number of Boolean literals (true, false, etc).
	LiteralBools int `toml:"literal_bools,omitzero" json:"literal_bools,omitempty"`

	// AtomicExpressions gives information about atomic statements.
	AtomicExpressions AtomicStatset `toml:"atomic_expressions,omitempty,omitzero" json:"atomic_expressions,omitempty"`

	// AtomicStatements gives information about atomic statements.
	AtomicStatements AtomicStatset `toml:"atomic_statements,omitempty,omitzero" json:"atomic_statements,omitempty"`
}
