// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package iohelp

import (
	"io"
)

// NopWriteCloser is like a NopCloser, but implements WriteCloser rather than ReadCloser.
type NopWriteCloser struct {
	io.Writer
}

// Close does nothing.
func (n NopWriteCloser) Close() error {
	return nil
}

// DiscardCloser is like io.Discard, but implements WriteCloser.
func DiscardCloser() io.WriteCloser {
	return NopWriteCloser{io.Discard}
}
