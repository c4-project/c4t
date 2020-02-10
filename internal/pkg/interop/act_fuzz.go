package interop

import (
	"os"
	"strconv"
)

// BinActFuzz is the name of the ACT fuzzer binary.
const BinActFuzz = "act-fuzz"

// FuzzSingle wraps the ACT one-file fuzzer, supplying the given seed.
func (a ActRunner) FuzzSingle(seed int, inPath, outPath string) error {
	argv := []string{"-seed", strconv.Itoa(seed), "-o", outPath, inPath}
	sargs := StandardArgs{Verbose: false}

	return a.Run(BinActCompiler, nil, nil, os.Stderr, sargs, argv...)
}
