package compiler

import (
	"path"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"
)

const (
	segBins = "bins"
	segLogs = "logs"
)

// Pathset contains the various directories used by the test compiler.
type Pathset struct {
	// DirBins is the directory into which compiled binaries should go.
	DirBins string

	// DirLogs is the directory into which compiler logs should go.
	DirLogs string
}

// NewPathset constructs a new pathset from the directory root.
func NewPathset(root string) *Pathset {
	return &Pathset{
		DirBins: path.Join(root, segBins),
		DirLogs: path.Join(root, segLogs),
	}
}

// Dirs gets all of the directories mentioned by this pathset.
func (p *Pathset) Dirs() []string {
	return []string{p.DirBins, p.DirLogs}
}

// OnCompiler gets the binary and log file paths for subject as compiled by the compiler with CompilerID compiler.
func (p *Pathset) OnCompiler(compiler model.ID, subject string) (bin, log string) {
	csub := append(compiler.Tags(), subject)
	bpath := append([]string{p.DirBins}, csub...)
	lpath := append([]string{p.DirLogs}, csub...)
	return path.Join(bpath...), path.Join(lpath...)
}
