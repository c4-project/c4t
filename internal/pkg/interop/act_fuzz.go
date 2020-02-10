package interop

import (
	"os"
	"strconv"
)

// BinActFuzz is the name of the ACT fuzzer binary.
const BinActFuzz = "act-fuzz"

// FuzzSingle wraps the ACT one-file fuzzer, supplying the given seed.
func (a ActRunner) FuzzSingle(seed int32, inPath, outPath, tracePath string) error {
	sargs := StandardArgs{Verbose: false}
	seedStr := strconv.Itoa(int(seed))
	return a.Run(BinActFuzz, nil, nil, os.Stderr, sargs, "run", "-seed", seedStr, "-o", outPath,
		"-trace-output", tracePath, inPath)
}
