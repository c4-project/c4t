// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package job

// Fuzzer contains information on how to fuzz a litmus file.
type Fuzzer struct {
	// Seed is the seed to use for randomising decisions made by the fuzzer.
	Seed int32

	// In is the slashpath to the file to fuzz.
	In string

	// OutLitmus is the slashpath to the litmus file that should be outputted by the fuzzer.
	OutLitmus string

	// OutTrace is the slashpath to the trace file that should be outputted by the fuzzer.
	OutTrace string
}
