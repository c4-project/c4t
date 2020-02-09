package interop

import (
	"os"
	"strconv"
)

// BinActFuzz is the name of the ACT fuzzer binary.
const BinActFuzz = "act-fuzz"

// Fuzzer represents types that can commune with a C litmus test fuzzer.
type Fuzzer interface {
	// Fuzz fuzzes the test at path inPath using the given seed, outputting to path outPath.
	Fuzz(seed int, inPath, outPath string) error
}

func (a ActRunner) Fuzz(seed int, inPath, outPath string) error {
	argv := []string{"-seed", strconv.Itoa(seed), "-o", outPath, inPath}
	sargs := StandardArgs{Verbose: false}

	return a.Run(BinActCompiler, nil, nil, os.Stderr, sargs, argv...)
}
