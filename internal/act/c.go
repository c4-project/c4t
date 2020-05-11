// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package act

import (
	"bytes"
	"context"
	"fmt"

	"github.com/MattWindsor91/act-tester/internal/model"

	"github.com/MattWindsor91/act-tester/internal/model/subject"
)

// BinActC is the name of the ACT C services binary.
const BinActC = "act-c"

// ProbeSubject probes the litmus test at path litmus, returning a named subject record.
func (a *Runner) ProbeSubject(ctx context.Context, litmus string) (subject.Named, error) {
	var h Header
	if err := a.DumpHeader(ctx, &h, litmus); err != nil {
		return subject.Named{}, fmt.Errorf("header read on %s failed: %w", litmus, err)
	}

	s := subject.Named{
		Name: h.Name,
		Subject: subject.Subject{
			OrigLitmus: litmus,
		},
	}

	if err := a.DumpStats(ctx, &s.Stats, litmus); err != nil {
		return subject.Named{}, fmt.Errorf("stats read on %s failed: %w", litmus, err)
	}

	return s, nil
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
func (a *Runner) DumpStats(ctx context.Context, s *model.Statset, path string) error {
	var obuf bytes.Buffer
	sargs := StandardArgs{Verbose: false}

	cmd := a.CommandContext(ctx, BinActC, "dump-stats", sargs, path)
	cmd.Stdout = &obuf

	if err := cmd.Run(); err != nil {
		return err
	}

	return ParseStats(s, &obuf)
}
