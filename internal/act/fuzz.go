// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package act

import (
	"context"
	"strconv"

	"github.com/MattWindsor91/act-tester/internal/model/job"
)

// BinActFuzz is the name of the ACT fuzzer binary.
const BinActFuzz = "act-fuzz"

// FuzzSingle wraps the ACT one-file fuzzer, supplying the given seed.
func (a *Runner) Fuzz(ctx context.Context, j job.Fuzzer) error {
	sargs := StandardArgs{Verbose: false}
	seedStr := strconv.Itoa(int(j.Seed))
	args := []string{"-seed", seedStr, "-o", j.OutLitmus, "-trace-output", j.OutTrace, j.In}
	cmd := a.CommandContext(ctx, BinActFuzz, "run", sargs, args...)
	return cmd.Run()
}
