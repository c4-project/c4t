// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package fuzzer

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// ErrNotSubjectCycleName occurs when we try to ParseSubjectCycle on a string that isn't of the right format.
var ErrNotSubjectCycleName = errors.New("not a valid subject-cycle name")

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

// ParseSubjectCycle tries to back-form a SubjectCycle from s.
func ParseSubjectCycle(s string) (SubjectCycle, error) {
	var (
		sc  SubjectCycle
		err error
	)

	end := strings.LastIndexByte(s, '_')
	if end < 0 {
		return sc, fmt.Errorf("%w: no underscore in %s", ErrNotSubjectCycleName, s)
	}

	sc.Name = s[:end]
	sc.Cycle, err = strconv.Atoi(s[end+1:])
	return sc, err
}
