package interop

import (
	"bytes"
	"errors"
	"strings"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"
)

// BinActBackend is the name of the ACT backend services binary.
const BinActBackend = "act-backend"

// ErrNoBackend occurs when no backend is reported by ACT.
var ErrNoBackend = errors.New("no backend reported")

// FindBackend finds a backend using ACT.
func (a ActRunner) FindBackend(style model.Id, machines ...model.Id) (*model.Backend, error) {
	id, err := a.runFindBackend(style, machines)
	if err != nil {
		return nil, err
	}

	if id.String() == "" {
		return nil, ErrNoBackend
	}

	return &model.Backend{
		Service: model.Service{Id: id, IdQualified: true, Style: style},
	}, nil
}

// runFindBackend does most of the legwork of running an ACT find-backend query.
func (a ActRunner) runFindBackend(style model.Id, machines []model.Id) (model.Id, error) {
	argv := findBackendArgv(style, machines)
	sargs := StandardArgs{Verbose: false}

	var obuf bytes.Buffer
	cmd := a.Command(BinActBackend, "find", sargs, argv...)
	cmd.Stdout = &obuf
	if err := cmd.Run(); err != nil {
		return model.EmptyId, err
	}

	return model.IdFromString(strings.TrimSpace(obuf.String())), nil
}

// findBackendArgv constructs the argv for a backend find on style and machines.
func findBackendArgv(style model.Id, machines []model.Id) []string {
	argv := make([]string, len(machines)+1)
	argv[0] = style.String()
	for i, m := range machines {
		argv[i+1] = m.String()
	}
	return argv
}

// MakeHarness makes a harness using ACT.
func (a ActRunner) MakeHarness(_ model.HarnessSpec) (outFiles []string, err error) {
	// TODO(@MattWindsor91)
	return nil, errors.New("unimplemented")
}
