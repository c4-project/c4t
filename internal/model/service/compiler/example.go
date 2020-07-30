// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package compiler

import (
	"github.com/MattWindsor91/act-tester/internal/model/id"
	"github.com/MattWindsor91/act-tester/internal/model/service/compiler/optlevel"
)

// MockPower9GCCOpt produces a GCC-compatible power entry with Power9 architecture and optimisation configuration.
func MockPower9GCCOpt() Configuration {
	return Configuration{
		Compiler: Compiler{
			Style: id.CStyleGCC,
			Arch:  id.ArchPPCPOWER9,
			Opt: &optlevel.Selection{
				Enabled:  []string{"1", "2", "3"},
				Disabled: []string{"fast"},
			},
		},
	}
}

// MockPower9GCCOpt produces a GCC-compatible power entry with X86 architecture.
func MockX86Gcc() Configuration {
	return Configuration{
		Compiler: Compiler{
			Style: id.CStyleGCC,
			Arch:  id.ArchX86,
		},
	}
}

// MockSet produces a mock compiler set.
func MockSet() map[string]Configuration {
	// These names line up with those in the example corpus.
	return map[string]Configuration{
		"gcc":   MockPower9GCCOpt(),
		"clang": MockX86Gcc(),
	}
}
