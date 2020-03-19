// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package director

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
	Machine(id id.ID) MachineObserver
}

// MachineObserver is an interface for types that observe a director machine loop.
type MachineObserver interface {
	// OnIteration lets the observer know that the machine loop has started anew.
	// iter is, modulo eventual overflow, the current iteration number;
	// time is the time at which the iteration started.
	OnIteration(iter uint64, time time.Time)

	builder.Observer
}

// SilentObserver wraps the builder silent-observer to add the additional MachineObserver functions.
type SilentObserver struct{ builder.SilentObserver }

// OnIteration does nothing.
func (o SilentObserver) OnIteration(uint64, time.Time) {
}
