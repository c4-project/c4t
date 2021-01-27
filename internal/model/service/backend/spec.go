// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package backend

import (
	"github.com/c4-project/c4t/internal/model/id"
	"github.com/c4-project/c4t/internal/model/service"
)

// Spec tells the tester how to run a backend.
type Spec struct {
	// Style is the declared style of the backend.
	Style id.ID `toml:"style" json:"style"`

	// Run contains information on how to run the backend; if given, this overrides any default RunInfo for the backend.
	Run *service.RunInfo `toml:"run,omitempty" json:"run,omitempty"`
}

// NamedSpec wraps a Spec with its ID.
type NamedSpec struct {
	// ID is the ID of the backend.
	ID id.ID `toml:"id"`

	Spec
}
