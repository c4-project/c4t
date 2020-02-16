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

// SubjectPathset is the set of file paths associated with a fuzzer output.
type SubjectPathset struct {
	// FileLitmus is the path to this subject's fuzzed Litmus file.
	FileLitmus string
	// FileTrace is the path to this subject's fuzzer trace file.
	FileTrace string
}
