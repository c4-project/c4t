package interop

import (
	"bytes"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"
)

// BinActCompiler is the name of the ACT compiler services binary.
const BinActCompiler = "act-compiler"

// ListCompilers queries ACT for a list of compilers satisfying f.
func (a ActRunner) ListCompilers(f model.CompilerFilter) (map[string]map[string]model.Compiler, error) {
	sargs := StandardArgs{Verbose: false}

	var obuf bytes.Buffer

	cmd := a.Command(BinActCompiler, "list", sargs, f.ToArgv()...)
	cmd.Stdout = &obuf

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	return model.ParseCompilerList(&obuf)
}
