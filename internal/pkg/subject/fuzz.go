// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package subject

// FuzzFileset is the set of file paths associated with a fuzzer output.
type FuzzFileset struct {
	// Litmus is the path to this subject's fuzzed Litmus file.
	Litmus string `toml:"litmus,omitempty"`
	// Trace is the path to this subject's fuzzer trace file.
	Trace string `toml:"trace,omitempty"`
}
