// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package iohelp

import (
	"io"
	"io/ioutil"
)

// NopWriteCloser is like a NopCloser, but implements WriteCloser rather than ReadCloser.
type NopWriteCloser struct {
	io.Writer
}

// Close does nothing.
func (n NopWriteCloser) Close() error {
	return nil
}

// DiscardCloser is like ioutil.Discard, but implements WriteCloser.
func DiscardCloser() io.WriteCloser {
	return NopWriteCloser{ioutil.Discard}
}
