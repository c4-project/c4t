package interop

import (
	"bytes"
	"io"

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

func (a ActRunner) RunCompiler(c *model.NamedCompiler, infiles []string, outfile string, errw io.Writer) error {
	sargs := StandardArgs{Verbose: false}

	argv := runCompilerArgv(c.ID, infiles, outfile)
	cmd := a.Command(BinActCompiler, "run", sargs, argv...)
	cmd.Stderr = errw

	return cmd.Run()
}

func runCompilerArgv(compiler model.ID, infiles []string, outfile string) []string {
	base := []string{"-compiler", compiler.String(), "-mode", "binary", "-o", outfile}
	return append(base, infiles...)
}
