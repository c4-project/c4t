// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package subject

// Named wraps a Normalise with its name.
type Named struct {
	// Name is the name of the subject.
	Name string

	// Normalise embeds the subject itself.
	Subject
}
