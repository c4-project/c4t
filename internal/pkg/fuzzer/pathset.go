package fuzzer

import "path"

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

func NewPathset(root string) *Pathset {
	return &Pathset{
		DirRoot:   root,
		DirLitmus: path.Join(root, segLitmus),
		DirTrace:  path.Join(root, segTrace),
	}
}
