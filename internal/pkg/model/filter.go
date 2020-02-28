// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package model

// CompilerFilter specifies filtering predicates used to find compilers.
type CompilerFilter struct {
	// CompPred is the compiler predicate.
	CompPred string

	// MachPred is the machine predicate.
	MachPred string
}

// ToArgv converts c to an argument vector fragment.
func (c CompilerFilter) ToArgv() []string {
	var argv []string
	if c.CompPred != "" {
		argv = append(argv, "-filter-compilers", c.CompPred)
	}
	if c.MachPred != "" {
		argv = append(argv, "-filter-machines", c.MachPred)
	}
	return argv
}
