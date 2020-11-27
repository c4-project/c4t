// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package copier

import (
	"io"
	"os"
)

// Local implements Copier through os.
type Local struct{}

// Create calls os.Create on path.
func (l Local) Create(path string) (io.WriteCloser, error) {
	return os.Create(path)
}

// Open calls os.Open on path.
func (l Local) Open(path string) (io.ReadCloser, error) {
	return os.Open(path)
}

// MkdirAll calls os.MkdirAll on path, with vaguely sensible permissions.
func (l Local) MkdirAll(dir string) error {
	return os.MkdirAll(dir, 0744)
}
