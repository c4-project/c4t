// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package subject

import "time"

// Fuzz is the set of file paths, and other metadata, associated with a fuzzer output.
type Fuzz struct {
	// Duration is the length of time it took to fuzz this file.
	Duration time.Duration `toml:"duration,omitzero"`

	// Files is the set of files produced by this fuzzing.
	Files FuzzFileset `toml:"files"`
}

// FuzzFileset is the set of files associated with a fuzzer output.
type FuzzFileset struct {
	// Litmus is the path to this subject's fuzzed Litmus file.
	Litmus string `toml:"litmus,omitempty"`
	// Trace is the path to this subject's fuzzer trace file.
	Trace string `toml:"trace,omitempty"`
}
