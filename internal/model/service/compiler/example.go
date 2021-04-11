// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package compiler

import (
	"github.com/c4-project/c4t/internal/id"
	"github.com/c4-project/c4t/internal/model/service/compiler/optlevel"
)

// MockPower9GCCOpt produces a GCC-compatible instance with Power9 architecture and optimisation configuration.
func MockPower9GCCOpt() Instance {
	return Instance{
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

// MockX86Gcc produces a GCC-compatible instance with X86 architecture.
func MockX86Gcc() Instance {
	return Instance{
		Compiler: Compiler{
			Style: id.CStyleGCC,
			Arch:  id.ArchX86,
		},
	}
}

// MockSet produces a mock compiler set.
func MockSet() map[id.ID]Instance {
	// These names line up with those in the example corpus.
	return map[id.ID]Instance{
		id.FromString("gcc"):   MockPower9GCCOpt(),
		id.FromString("clang"): MockX86Gcc(),
	}
}
