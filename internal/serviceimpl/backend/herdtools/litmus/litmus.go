// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package litmus contains the parts of a Herdtools backend specific to herd7.
package litmus

import (
	"fmt"

	"github.com/MattWindsor91/act-tester/internal/model/job"
	"github.com/MattWindsor91/act-tester/internal/model/service"
)

// Litmus describes the parts of a Litmus invocation that are specific to Herd.
type Litmus struct{}

// Args deduces the appropriate arguments for running Litmus on job j, with the merged run information r.
func (l Litmus) Args(j job.Lifter, r service.RunInfo) ([]string, error) {
	larch, err := lookupArch(j.Arch)
	if err != nil {
		return nil, fmt.Errorf("when looking up -carch: %w", err)
	}
	args := []string{
		"-o", j.OutDir,
		"-carch", larch,
		"-c11", "true",
	}
	args = append(args, r.Args...)
	args = append(args, j.InFile)
	return args, nil
}
