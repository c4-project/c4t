// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package fuzzer

import "errors"

var (
	// ErrConfigNil occurs when a fuzzer gets constructed using a nil config.
	ErrConfigNil = errors.New("config nil")

	// ErrDriverNil occurs when the fuzzer tries to use the nil pointer as its single-fuzz driver.
	ErrDriverNil = errors.New("driver nil")

	// ErrDriverNil occurs when the fuzzer tries to use the nil pointer as its stat dumper.
	ErrStatDumperNil = errors.New("stat dumper nil")
)
