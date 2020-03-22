// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package observer

import (
	"time"

	"github.com/MattWindsor91/act-tester/internal/pkg/corpus/builder"
)

// SilentObserver wraps the builder silent-observer to add the additional Instance functions.
type SilentObserver struct{ builder.SilentObserver }

// OnIteration does nothing.
func (o SilentObserver) OnIteration(uint64, time.Time) {
}
