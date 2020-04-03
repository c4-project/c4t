// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package subject

import "time"

// Result is the base structure for things that represent the result of an external process.
type Result struct {
	// Time is the time at which the process commenced.
	Time time.Time `toml:"time,omitzero"`

	// Duration is the rough duration of the process.
	Duration time.Duration `toml:"duration,omitzero"`

	// Status is the status of the process.
	Status Status `toml:"status"`
}