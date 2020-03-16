// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package forward describes the JSON-based protocol used to 'forward' messages
// and errors from a machine-runner to a director, potentially over SSH.
package forward

import "github.com/MattWindsor91/act-tester/internal/pkg/corpus/builder"

// Forward describes a 'forwarded' message or error.
type Forward struct {
	// Error carries an error.
	Error error `json:"error,omitempty"`

	// BuildStart carries information about the beginning of a corpus build.
	BuildStart *builder.Manifest `json:"build_start,omitempty"`

	// BuildUpdate carries an update from a corpus build.
	BuildUpdate *builder.Request `json:"build_update,omitempty"`

	// BuildEnd being true signifies that this is the end of a corpus build.
	BuildEnd bool `json:"build_end,omitempty"`
}
