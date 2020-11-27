// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package service

import "errors"

var (
	// ErrNil is a standard error for when a service record is nil, but shouldn't be.
	ErrNil = errors.New("service record missing")
)
