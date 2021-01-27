// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package interpreter

import "github.com/c4-project/c4t/internal/model/service/compiler"

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

// CompileWith sets the interpreter's compiler to driver d and config c.
// This is required to interpret recipes that involve compilation.
func CompileWith(d Driver, c *compiler.Instance) Option {
	return func(i *Interpreter) { i.driver = d; i.compiler = c }
}

// SetMaxObjs sets the maximum number of object files the interpreter can create.
func SetMaxObjs(cap uint64) Option {
	return func(i *Interpreter) { i.maxobjs = cap }
}
