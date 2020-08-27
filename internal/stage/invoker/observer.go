// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package invoker

import (
	"github.com/MattWindsor91/act-tester/internal/copier"
	"github.com/MattWindsor91/act-tester/internal/subject/corpus/builder"
)

// Observer is the union of the various interfaces of observers used by invoker.
type Observer interface {
	copier.Observer
	builder.Observer
}
