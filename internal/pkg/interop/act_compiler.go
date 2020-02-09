package interop

import (
	"bytes"
	"os"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"
)

// CompilerLister is the interface of things that can query compiler information.
type CompilerLister interface {
	// ListCompilers asks the compiler inspector to list all available compilers given the filter f.
	ListCompilers(f model.CompilerFilter) ([]*model.Compiler, error)
}

// BinActCompiler is the name of the ACT compiler services binary.
const BinActCompiler = "act-compiler"

func (a ActRunner) ListCompilers(f model.CompilerFilter) ([]*model.Compiler, error) {
	argv := append([]string{"list"}, f.ToArgv()...)
	sargs := StandardArgs{Verbose: false}

	var obuf bytes.Buffer

	if err := a.Run(BinActCompiler, nil, &obuf, os.Stderr, sargs, argv...); err != nil {
		return nil, err
	}

	return model.ParseCompilerList(&obuf)
}
