// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package act

import (
	"bytes"
	"context"
	"fmt"

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
	sargs := StandardArgs{Verbose: false}

	cmd := a.CommandContext(ctx, BinActC, "dump-header", sargs, path)
	cmd.Stdout = &obuf

	if err := cmd.Run(); err != nil {
		return err
	}

	return h.Read(&obuf)
}

// DumpStats runs act-c dump-stats on the subject at path, writing the stats to s.
func (a *Runner) DumpStats(ctx context.Context, s *litmus.Statset, path string) error {
	var obuf bytes.Buffer
	sargs := StandardArgs{Verbose: false}

	cmd := a.CommandContext(ctx, BinActC, "dump-stats", sargs, path)
	cmd.Stdout = &obuf

	if err := cmd.Run(); err != nil {
		return err
	}

	return ParseStats(s, &obuf)
}
