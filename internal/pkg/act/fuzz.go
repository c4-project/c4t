// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package act

import (
	"context"
	"strconv"

	"github.com/MattWindsor91/act-tester/internal/pkg/subject"
)

// BinActFuzz is the name of the ACT fuzzer binary.
const BinActFuzz = "act-fuzz"

// FuzzSingle wraps the ACT one-file fuzzer, supplying the given seed.
func (a *Runner) FuzzSingle(ctx context.Context, seed int32, inPath string, outPaths subject.FuzzFileset) error {
	sargs := StandardArgs{Verbose: false}
	seedStr := strconv.Itoa(int(seed))
	args := []string{"-seed", seedStr, "-o", outPaths.Litmus, "-trace-output", outPaths.Trace, inPath}
	cmd := a.CommandContext(ctx, BinActFuzz, "run", sargs, args...)
	return cmd.Run()
}
