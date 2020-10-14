// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package act

import (
	"bytes"
	"context"
	"fmt"

	"github.com/1set/gut/ystring"

	"github.com/MattWindsor91/act-tester/internal/model/litmus"

	"github.com/MattWindsor91/act-tester/internal/subject"
)

// BinActC is the name of the ACT C services binary.
const BinActC = "act-c"

// ProbeSubject probes the litmus test at path, returning a named subject record.
func (a *Runner) ProbeSubject(ctx context.Context, path string) (*subject.Named, error) {
	// TODO(@MattWindsor91): stat dumping and subject probing should likely be two separate things.
	var h Header
	if err := a.DumpHeader(ctx, &h, path); err != nil {
		return nil, fmt.Errorf("header read on %s failed: %w", path, err)
	}

	l, err := litmus.NewWithStats(ctx, path, a)
	if err != nil {
		return nil, fmt.Errorf("stats read on %s failed: %w", path, err)
	}
	s, err := subject.New(l)
	if err != nil {
		return nil, err
	}
	return s.AddName(h.Name), nil
}

// DumpHeader runs act-c dump-header on the subject at path, writing the results to h.
func (a *Runner) DumpHeader(ctx context.Context, h *Header, path string) error {
	var obuf bytes.Buffer
	cs := CmdSpec{
		Cmd:    BinActC,
		Subcmd: "dump-header",
		Args:   []string{path},
		Stdout: &obuf,
	}
	if err := a.Run(ctx, cs); err != nil {
		return err
	}
	return h.Read(&obuf)
}

// DumpStats runs act-c dump-stats on the subject at path, writing the stats to s.
func (a *Runner) DumpStats(ctx context.Context, s *litmus.Statset, path string) error {
	var obuf bytes.Buffer
	cs := CmdSpec{
		Cmd:    BinActC,
		Subcmd: "dump-stats",
		Args:   []string{path},
		Stdout: &obuf,
	}
	if err := a.Run(ctx, cs); err != nil {
		return err
	}
	return ParseStats(s, &obuf)
}

// DelitmusJob holds information about a single delitmus run.
type DelitmusJob struct {
	// InLitmus is the filepath of the input litmus file.
	InLitmus string
	// OutAux is the filepath to which the delitmusifier should write the auxiliary file.
	OutAux string
	// OutC is the filepath to which the delitmusifier should write the output file.
	OutC string
	// TODO(@MattWindsor91): impl-suffix, no-qualify-locals, style, etc.
}

// Args gets the argument vector for DelitmusJob.
func (d DelitmusJob) Args() []string {
	// TODO(@MattWindsor91): hook up style etc.
	var args []string
	if ystring.IsNotBlank(d.OutAux) {
		args = append(args, "-aux-output", d.OutAux)
	}
	if ystring.IsNotBlank(d.OutC) {
		args = append(args, "-output", d.OutC)
	}
	return append(args, d.InLitmus)
}

// Delitmus runs act-c delitmus as directed by d.
func (a *Runner) Delitmus(ctx context.Context, d DelitmusJob) error {
	cs := CmdSpec{
		Cmd:    BinActC,
		Subcmd: "delitmus",
		Args:   d.Args(),
	}
	return a.Run(ctx, cs)
}
