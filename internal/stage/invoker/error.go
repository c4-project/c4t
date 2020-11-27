// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package invoker

import "errors"

var (
	// ErrDirEmpty occurs when the local directory filepath is empty.
	ErrDirEmpty = errors.New("local dir is empty string")
)
