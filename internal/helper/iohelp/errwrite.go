// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package iohelp

// ErrWriter provides an io.Writer that fails with the wrapped error.
type ErrWriter struct {
	Err error
}

func (e ErrWriter) Write(_ []byte) (n int, err error) {
	return 0, e.Err
}
