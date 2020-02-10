package interop

import (
	"os"
	"strconv"
)

// BinActFuzz is the name of the ACT fuzzer binary.
const BinActFuzz = "act-fuzz"

// FuzzSingle wraps the ACT one-file fuzzer, supplying the given seed.
func (a ActRunner) FuzzSingle(seed int, inPath, outPath string) error {
	sargs := StandardArgs{Verbose: false}
	seedStr := strconv.Itoa(seed)
	return a.Run(BinActCompiler, nil, nil, os.Stderr, sargs, "-seed", seedStr, "-o", outPath, inPath)
}
