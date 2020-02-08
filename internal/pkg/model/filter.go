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
		argv = append(argv, "-compiler-predicate", c.CompPred)
	}
	if c.MachPred != "" {
		argv = append(argv, "-machine-predicate", c.MachPred)
	}
	return argv
}
