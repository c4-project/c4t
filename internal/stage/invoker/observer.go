// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package invoker

import (
	"github.com/c4-project/c4t/internal/copier"
	"github.com/c4-project/c4t/internal/subject/corpus/builder"
)

// Observer is the union of the various interfaces of observers used by invoker.
type Observer interface {
	copier.Observer
	builder.Observer
}
