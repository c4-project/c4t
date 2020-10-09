// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package job

import (
	"github.com/MattWindsor91/act-tester/internal/machine"
	fuzzer2 "github.com/MattWindsor91/act-tester/internal/model/service/fuzzer"
)

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

	// Machine is, optionally, the machine that is the target of the fuzzed output.
	Machine *machine.Machine

	// Config is the configuration for the fuzzer, if any.
	Config *fuzzer2.Configuration
}
