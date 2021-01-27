// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package builder

// Manifest describes the layout of a corpus build task.
type Manifest struct {
	// Name describes the build task that is starting.
	Name string `json:"name"`

	// NReqs is the number of requests in the build task.
	// The corpus builder will wait for this number of requests to arrive.
	NReqs int `json:"nreqs"`
}
