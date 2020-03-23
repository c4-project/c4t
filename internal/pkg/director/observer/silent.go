// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package observer

import (
	"time"

	"github.com/MattWindsor91/act-tester/internal/pkg/model/corpus/collate"

	"github.com/MattWindsor91/act-tester/internal/pkg/model/corpus/builder"
)

// Silent wraps the builder silent-observer to add the additional Instance functions.
type Silent struct{ builder.SilentObserver }

// OnIteration does nothing.
func (o Silent) OnIteration(uint64, time.Time) {}

// OnCollation does nothing.
func (o Silent) OnCollation(*collate.Collation) {}

// OnCopyStart does nothing.
func (o Silent) OnCopyStart(int) {}

// OnCopy does nothing.
func (o Silent) OnCopy(string, string) {}

// OnCopyFinish does nothing.
func (o Silent) OnCopyFinish() {}
