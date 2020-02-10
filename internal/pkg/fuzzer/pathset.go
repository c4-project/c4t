package fuzzer

import (
	"fmt"
	"os"
	"path"

	"github.com/sirupsen/logrus"
)

const (
	// segLitmus is the directory element added to the root directory to form the litmus directory.
	segLitmus = "litmus"

	// segTrace is the directory element added to the root directory to form the trace directory.
	segTrace = "trace"
)

// Pathset contains the pre-computed paths used by a run of the fuzzer.
type Pathset struct {
	// DirRoot is the root directory of this fuzzer's path set.
	DirRoot string

	// DirLitmus is the directory to which litmus tests will be written.
	DirLitmus string

	// DirTrace is the directory to which traces will be written.
	DirTrace string
}

// NewPathset constructs a new pathset from the directory root.
func NewPathset(root string) *Pathset {
	return &Pathset{
		DirRoot:   root,
		DirLitmus: path.Join(root, segLitmus),
		DirTrace:  path.Join(root, segTrace),
	}
}

// Mkdirs tries to make each directory in pathset.
func (p Pathset) Mkdirs() error {
	for _, dir := range []string{p.DirRoot, p.DirLitmus, p.DirTrace} {
		logrus.Debugf("mkdir %s\n", dir)
		if err := os.Mkdir(dir, 0744); err != nil {
			return err
		}
	}
	return nil
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
