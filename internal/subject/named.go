// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package subject

// Named wraps a Subject with its name.
type Named struct {
	// Name is the name of the subject.
	Name string

	// Normalise embeds the subject itself.
	Subject
}

// AddName copies this subject into a new Named with the given name.
func (s *Subject) AddName(name string) *Named {
	return &Named{Name: name, Subject: *s}
}
