// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package iohelp

import (
	"io"
	"io/ioutil"
	"log"
)

// EnsureWriter passes through w if non-nil, or supplies ioutil.Discard otherwise.
func EnsureWriter(w io.Writer) io.Writer {
	if w == nil {
		return ioutil.Discard
	}
	return w
}

// EnsureLog passes through l if non-nil, or constructs a dummy Logger otherwise.
func EnsureLog(l *log.Logger) *log.Logger {
	if l == nil {
		return log.New(ioutil.Discard, "", 0)
	}
	return l
}
