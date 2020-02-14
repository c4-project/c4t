package fuzzer

import (
	"fmt"
	"path"
)

const (
	// segLitmus is the directory element added to the root directory to form the litmus directory.
	segLitmus = "litmus"

	// segTrace is the directory element added to the root directory to form the trace directory.
	segTrace = "trace"
)

// Pathset contains the pre-computed paths used by a run of the fuzzer.
type Pathset struct {
	// DirLitmus is the directory to which litmus tests will be written.
	DirLitmus string

	// DirTrace is the directory to which traces will be written.
	DirTrace string
}

// NewPathset constructs a new pathset from the directory root.
func NewPathset(root string) *Pathset {
	return &Pathset{
		DirLitmus: path.Join(root, segLitmus),
		DirTrace:  path.Join(root, segTrace),
	}
}

// Dirs gets a list of all directories in the pathset.
func (p Pathset) Dirs() []string {
	return []string{p.DirTrace, p.DirLitmus}
}

// OnSubject gets the litmus and trace file paths for the subject with the given name and fuzzing cycle.
func (p Pathset) OnSubject(name string, cycle int) (litmus, trace string) {
	base := CycledName(name, cycle)
	return path.Join(p.DirLitmus, base+".litmus"), path.Join(p.DirTrace, base+".trace")
}

// CycledName gets the new name of subject name given the current cycle.
func CycledName(name string, cycle int) string {
	return fmt.Sprintf("%s_%d", name, cycle)
}
