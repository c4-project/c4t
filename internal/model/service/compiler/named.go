// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package compiler

import (
	"github.com/c4-project/c4t/internal/model/id"
)

// Named wraps a Configuration with its ID.
type Named struct {
	// ID is the ID of the compiler.
	ID id.ID `toml:"id" json:"id"`

	Configuration
}

// AddName names this Configuration with ID name, lifting it to a Named.
func (c Configuration) AddName(name id.ID) *Named {
	return &Named{ID: name, Configuration: c}
}

// AddNameString tries to resolve name into an ID then name this Configuration with it.
func (c Configuration) AddNameString(name string) (*Named, error) {
	nid, err := id.TryFromString(name)
	if err != nil {
		return nil, err
	}
	return c.AddName(nid), err
}

// FullID gets a fully qualified identifier for this configuration, consisting of the compiler name, followed by
// 'oOpt' where 'Opt' is its selected optimisation name, and 'mMopt' where 'Mopt' is its selected machine profile.
func (n Named) FullID() (id.ID, error) {
	return id.New(append(n.ID.Tags(), "o"+n.SelectedOptName(), "m"+n.SelectedMOpt)...)
}
