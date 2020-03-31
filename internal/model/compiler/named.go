// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package compiler

import (
	"github.com/MattWindsor91/act-tester/internal/model/id"
)

// Named wraps a Compiler with its ID.
type Named struct {
	// ID is the ID of the compiler.
	ID id.ID `toml:"id"`

	Compiler
}

// AddName names this Compiler with ID name, lifting it to a Named.
func (c *Compiler) AddName(name id.ID) *Named {
	return &Named{ID: name, Compiler: *c}
}

// AddNameString tries to resolve name into an ID then name this Compiler with it.
func (c *Compiler) AddNameString(name string) (*Named, error) {
	nid, err := id.TryFromString(name)
	if err != nil {
		return nil, err
	}
	return c.AddName(nid), err
}
