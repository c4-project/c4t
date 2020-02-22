package interop

import (
	"context"
	"strconv"
)

// BinActFuzz is the name of the ACT fuzzer binary.
const BinActFuzz = "act-fuzz"

// FuzzSingle wraps the ACT one-file fuzzer, supplying the given seed.
func (a *ActRunner) FuzzSingle(ctx context.Context, seed int32, inPath, outPath, tracePath string) error {
	sargs := StandardArgs{Verbose: false}
	seedStr := strconv.Itoa(int(seed))
	cmd := a.CommandContext(ctx, BinActFuzz, "run", sargs, "-seed", seedStr, "-o", outPath, "-trace-output", tracePath, inPath)
	return cmd.Run()
}
