// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package fuzzer

// Config configures the fuzzer.
type Config struct {
	// Disabled, if set true, disables the fuzzer stage in the main tester.
	Disabled bool `toml:"disabled,omitempty"`

	// FuzzesPerSubject specifies the default number of times the fuzzer will be invoked per subject.
	// If zero, the default number is used.
	FuzzesPerSubject int `toml:"fuzzes_per_subject,omitempty"`

	// Params contains a low-level key-value map of parameters to pass to the fuzzer.
	Params map[string]string `toml:"params,omitempty"`
}
