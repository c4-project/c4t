// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package fuzzer

import "errors"

var (
	// ErrDriverNil occurs when the fuzzer tries to use the nil pointer as its single-fuzz driver.
	ErrDriverNil = errors.New("driver nil")
)
