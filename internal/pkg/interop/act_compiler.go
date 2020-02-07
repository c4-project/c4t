package interop

import (
	"bytes"
	"github.com/MattWindsor91/act-tester/internal/pkg/model"
)

// CompilerFilter specifies filtering predicates used to find compilers.
type CompilerFilter struct {
	// CompPred is the compiler predicate.
	CompPred string

	// MachPred is the machine predicate.
	MachPred string
}

// ToArgv converts c to an argument vector fragment.
func (c CompilerFilter) ToArgv() []string {
	argv := []string{"list"}
	if c.CompPred != "" {
		argv = append(argv, "--compiler-predicate", c.CompPred)
	}
	if c.MachPred != "" {
		argv = append(argv, "--machine-predicate", c.MachPred)
	}
	return argv
}

// CompilerInspectable is the interface of things that can query compiler information.
type CompilerInspectable interface {
	// List asks the compiler inspector to list all available compilers given the filter f.
	List(f CompilerFilter) ([]model.Compiler, error)
}

func (a *ActRunner) List(f CompilerFilter) ([]*model.Compiler, error) {
	argv := append([]string{"list"},f.ToArgv()...)
	sargs := StandardArgs{Verbose: false}

	var obuf bytes.Buffer
	var ebuf bytes.Buffer

	if err := a.Run(BinActCompiler, nil, &obuf, &ebuf, sargs, argv...); err != nil {
		return nil, err
	}

	return model.ParseCompilerList(&obuf)
}
