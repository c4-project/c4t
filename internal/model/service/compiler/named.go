// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package compiler

import (
	"strings"

	"github.com/c4-project/c4t/internal/model/id"
)

// Named wraps an Instance with its ID.
type Named struct {
	// ID is the ID of the compiler.
	ID id.ID `toml:"id" json:"id"`

	Instance
}

// AddName names this Instance with ID name, lifting it to a Named.
func (c Instance) AddName(name id.ID) *Named {
	return &Named{ID: name, Instance: c}
}

// AddNameString tries to resolve name into an ID then name this Instance with it.
func (c Instance) AddNameString(name string) (*Named, error) {
	nid, err := id.TryFromString(name)
	if err != nil {
		return nil, err
	}
	return c.AddName(nid), err
}

// FullID gets a fully qualified identifier for this configuration, consisting of the compiler name, followed by
// 'oOpt' where 'Opt' is its selected optimisation name, and 'mMopt' where 'Mopt' is its selected machine profile.
//
// Where Opt or Mopt contain '.', these become '_'.  This behaviour may change.
func (n Named) FullID() (id.ID, error) {
	// In case of things like '-march=armv8.1-a'.
	repl := strings.NewReplacer(".", "_")
	o := repl.Replace(n.SelectedOptName())
	m := repl.Replace(n.SelectedMOpt)

	// We don't append in the config time, which means that this ID doesn't fully capture the compiler specification;
	// that said, maybe the config time being a part of the specification is a rare enough case that we needn't worry.
	return id.New(append(n.ID.Tags(), "o"+o, "m"+m)...)
}
