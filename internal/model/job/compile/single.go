// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package compile

// Single represents a single compile job.
type Single struct {
	Compile

	// Kind is the kind of file being produced by this job.
	Kind Kind
}

// Single produces a single-compile job with this job's compiler information and the given output kind.
func (j *Compile) Single(kind Kind) Single {
	return Single{
		Compile: *j,
		Kind:    kind,
	}
}
