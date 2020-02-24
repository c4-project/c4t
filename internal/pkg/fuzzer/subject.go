// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package fuzzer

import "fmt"

// SubjectCycle describes the unique name of a particular instance of the batch fuzzer.
type SubjectCycle struct {
	// Name is the name of the subject.
	Name string
	// Cycle is the current fuzz cycle (zero based).
	Cycle int
}

// String is a filepath-suitable string representation of a fuzz name.
func (f SubjectCycle) String() string {
	return fmt.Sprintf("%s_%d", f.Name, f.Cycle)
}
