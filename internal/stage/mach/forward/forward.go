// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

// Package forward describes the JSON-based protocol used to 'forward' messages  and errors from a machine node to its
// invoker, potentially over SSH.
package forward

import (
	"github.com/MattWindsor91/c4t/internal/stage/mach/observer"
	"github.com/MattWindsor91/c4t/internal/subject/corpus/builder"
)

// Forward describes a 'forwarded' message or error.
type Forward struct {
	// Error carries an error's Error string.
	Error string `json:"error,omitempty"`

	// Build carries information about a corpus build.
	Build *builder.Message `json:"build,omitempty"`

	// Action carries information about a machine-node action.
	Action *observer.Message `json:"action.omitempty"`
}
