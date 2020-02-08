package interop

import (
	"bytes"
	"errors"
	"os"
	"strings"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"
)

// BackendFinder is the interface of things that can find backends for machines.
type BackendFinder interface {
	// FindBackend asks for a backend with the given style on any one of machines,
	// or a default machine if none have such a backend.
	FindBackend(style model.Id, machines ...model.Id) (*model.Backend, error)
}

// BinActBackend is the name of the ACT backend services binary.
const BinActBackend = "act-backend"

// ErrNoBackend occurs when no backend is reported by ACT.
var ErrNoBackend = errors.New("no backend reported")

func (a *ActRunner) FindBackend(style model.Id, machines ...model.Id) (*model.Backend, error) {
	id, err := a.runActBackend(style, machines)
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

func (a *ActRunner) runActBackend(style model.Id, machines []model.Id) (model.Id, error) {
	argv := findBackendArgv(style, machines)
	sargs := StandardArgs{Verbose: false}

	var obuf bytes.Buffer
	if err := a.Run(BinActBackend, nil, &obuf, os.Stderr, sargs, argv...); err != nil {
		return model.EmptyId, err
	}

	return model.IdFromString(strings.TrimSpace(obuf.String())), nil
}

func findBackendArgv(style model.Id, machines []model.Id) []string {
	argv := make([]string, len(machines)+2)
	argv[0] = "find"
	argv[1] = style.String()
	for i, m := range machines {
		argv[i+2] = m.String()
	}
	return argv
}
