// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package fuzzer

// Configuration lets the tester pass information to the fuzzer.
type Configuration struct {
	// Params contains a low-level key-value map of parameters to pass to the fuzzer.
	Params map[string]string `toml:"params,omitempty"`
}
