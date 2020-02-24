// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package interop

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"strings"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"
)

// BinActBackend is the name of the ACT backend services binary.
const BinActBackend = "act-backend"

// ErrNoBackend occurs when no backend is reported by ACT.
var ErrNoBackend = errors.New("no backend reported")

// FindBackend finds a backend using ACT.
func (a *ActRunner) FindBackend(ctx context.Context, style model.ID, machines ...model.ID) (*model.Backend, error) {
	id, err := a.runFindBackend(ctx, style, machines)
	if err != nil {
		return nil, err
	}

	if id.String() == "" {
		return nil, ErrNoBackend
	}

	return &model.Backend{
		ID: id, IDQualified: true, Style: style,
	}, nil
}

// runFindBackend does most of the legwork of running an ACT find-backend query.
func (a *ActRunner) runFindBackend(ctx context.Context, style model.ID, machines []model.ID) (model.ID, error) {
	argv := findBackendArgv(style, machines)
	sargs := StandardArgs{Verbose: false}

	var obuf bytes.Buffer
	cmd := a.CommandContext(ctx, BinActBackend, "find", sargs, argv...)
	cmd.Stdout = &obuf
	if err := cmd.Run(); err != nil {
		return model.ID{}, err
	}

	return model.IDFromString(strings.TrimSpace(obuf.String())), nil
}

// findBackendArgv constructs the argv for a backend find on style and machines.
func findBackendArgv(style model.ID, machines []model.ID) []string {
	argv := make([]string, len(machines)+1)
	argv[0] = style.String()
	for i, m := range machines {
		argv[i+1] = m.String()
	}
	return argv
}

// MakeHarness makes a harness using ACT.
func (a *ActRunner) MakeHarness(ctx context.Context, s model.HarnessSpec) (outFiles []string, err error) {
	argv := makeHarnessArgv(s)
	sargs := StandardArgs{Verbose: false}

	cmd := a.CommandContext(ctx, BinActBackend, "make-harness", sargs, argv...)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	return s.OutFiles()
}

// makeHarnessArgv creates the appropriate harness making argv for s.
func makeHarnessArgv(s model.HarnessSpec) []string {
	return []string{
		"-backend",
		s.Backend.String(),
		"-carch",
		s.Arch.String(),
		"-o",
		s.OutDir,
		s.InFile,
	}
}

// ParseObs uses act-backend to parse the observation coming in from r into o according to b.
func (a *ActRunner) ParseObs(ctx context.Context, b model.Backend, r io.Reader, o *model.Obs) error {
	cmd := a.CommandContext(ctx, BinActBackend, "parse", StandardArgs{}, "-backend", b.ID.String())
	cmd.Stdin = r
	// TODO(@MattWindsor91): do something useful with this
	cmd.Stderr = os.Stderr

	obsr, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return nil
	}

	if err := json.NewDecoder(obsr).Decode(o); err != nil {
		_ = cmd.Wait()
		return err
	}
	return cmd.Wait()
}
