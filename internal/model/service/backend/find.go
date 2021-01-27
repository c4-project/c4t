// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package backend

import (
	"errors"
	"fmt"
	"strings"

	"github.com/c4-project/c4t/internal/model/id"
)

// ErrNoMatch occurs when we can't find a backend that matches the given criteria.
var ErrNoMatch = errors.New("no matching backend found")

// Finder is the interface of things that can find backends for machines.
type Finder interface {
	// FindBackend asks for a backend that matches the given criteria.
	FindBackend(c Criteria) (*NamedSpec, error)
}

// Criteria contains the criteria for which a backend should be found.
//
// The criteria is a conjunction of each present criterion in the struct.
type Criteria struct {
	// IDGlob is a glob pattern for the identifier of the backend.
	IDGlob id.ID
	// StyleGlob is a glob pattern for the style of the backend.
	StyleGlob id.ID

	// TODO(@MattWindsor91): it'd be nice to have target/source resolution, but doing this without having a dependency
	// cycle seems difficult.
}

// String outputs a string representation of a set of criteria.
func (c Criteria) String() string {
	parts := make([]string, 0, 2)
	for _, g := range []struct {
		name string
		glob id.ID
	}{
		{name: "id", glob: c.IDGlob},
		{name: "style", glob: c.StyleGlob},
	} {
		if g.glob.IsEmpty() {
			continue
		}
		parts = append(parts, fmt.Sprintf("%s=%s", g.name, g.glob))
	}
	if len(parts) == 0 {
		return "any"
	}
	return strings.Join(parts, ", ")
}

// Matches tries to see if s matches this criteria.
func (c Criteria) Matches(s NamedSpec) (bool, error) {
	for _, g := range []struct{ id, glob id.ID }{
		{id: s.ID, glob: c.IDGlob},
		{id: s.Style, glob: c.StyleGlob},
	} {
		if g.glob.IsEmpty() {
			continue
		}
		if match, err := g.id.Matches(g.glob); !match || err != nil {
			return false, err
		}
	}
	return true, nil
}

// Find tries to find a matching backend in a list of specifications specs.
func (c Criteria) Find(specs []NamedSpec) (*NamedSpec, error) {
	for _, s := range specs {
		match, err := c.Matches(s)
		if err != nil {
			return nil, err
		}
		if match {
			return &s, nil
		}
	}
	return nil, fmt.Errorf("%w: %s", ErrNoMatch, c)
}
