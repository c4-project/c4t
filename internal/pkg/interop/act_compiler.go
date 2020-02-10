package interop

import (
	"bytes"
	"os"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"
)

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
