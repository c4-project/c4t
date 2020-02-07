package interop

import (
	"bytes"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"
)

// CompilerInspectable is the interface of things that can query compiler information.
type CompilerInspectable interface {
	// List asks the compiler inspector to list all available compilers given the filter f.
	List(f model.CompilerFilter) ([]model.Compiler, error)
}

func (a *ActRunner) List(f model.CompilerFilter) ([]*model.Compiler, error) {
	argv := append([]string{"list"}, f.ToArgv()...)
	sargs := StandardArgs{Verbose: false}

	var obuf bytes.Buffer
	var ebuf bytes.Buffer

	if err := a.Run(BinActCompiler, nil, &obuf, &ebuf, sargs, argv...); err != nil {
		return nil, err
	}

	return model.ParseCompilerList(&obuf)
}
