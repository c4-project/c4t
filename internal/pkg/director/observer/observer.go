// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package observer defines interfaces and basic implementations of the director's 'observer' pattern.
package observer

import (
	"context"
	"time"

	"github.com/MattWindsor91/act-tester/internal/pkg/corpus/builder"
	"github.com/MattWindsor91/act-tester/internal/pkg/model/id"
)

// Observer is an interface for types that implement multi-machine test progress observation.
type Observer interface {
	// Run runs the observer in a blocking manner using context ctx.
	// It will use cancel to cancel ctx if needed.
	Run(ctx context.Context, cancel func()) error

	// Instance gets a sub-observer for the machine with ID id.
	// It can fail if no such observer is available.
	Instance(id id.ID) (Instance, error)
}

// Instance is an interface for types that observe a director loop.
type Instance interface {
	// OnIteration lets the observer know that the machine loop has started anew.
	// iter is, modulo eventual overflow, the current iteration number;
	// time is the time at which the iteration started.
	OnIteration(iter uint64, time time.Time)

	// Instance observers can observe corpus building operations.
	builder.Observer
}
