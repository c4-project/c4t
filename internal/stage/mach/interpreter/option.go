// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package interpreter

import (
	"io"

	"github.com/MattWindsor91/c4t/internal/helper/iohelp"
)

// Option is the type of options to the interpreter.
type Option func(*Interpreter)

// Options bundles the options os into one option.
func Options(os ...Option) Option {
	return func(i *Interpreter) {
		for _, o := range os {
			o(i)
		}
	}
}

// LogTo logs compiler error output to w.
func LogTo(w io.Writer) Option {
	return func(i *Interpreter) { i.logw = iohelp.EnsureWriter(w) }
}

// SetMaxObjs sets the maximum number of object files the interpreter can create.
func SetMaxObjs(cap uint64) Option {
	return func(i *Interpreter) { i.maxobjs = cap }
}
