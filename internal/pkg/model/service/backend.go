// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package service

import "github.com/MattWindsor91/act-tester/internal/pkg/model/id"

// Backend tells the tester how to run a backend.
type Backend struct {
	// Style is the declared style of the backend.
	Style id.ID `toml:"style"`

	// Run contains information on how to run the compiler.
	Run *RunInfo `toml:"run,omitempty"`
}

// NamedBackend wraps a Backend with its ID.
type NamedBackend struct {
	// ID is the ID of the backend.
	ID id.ID `toml:"id"`

	Backend
}
