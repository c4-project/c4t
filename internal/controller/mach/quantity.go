// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package mach

import (
	"github.com/MattWindsor91/act-tester/internal/controller/mach/compiler"
	"github.com/MattWindsor91/act-tester/internal/controller/mach/runner"
)

// QuantitySet contains the tunable quantities for both batch-compiler and batch-runner.
type QuantitySet struct {
	// Compiler is the quantity set for the compiler.
	Compiler compiler.QuantitySet
	// Runner is the quantity set for the runner.
	Runner runner.QuantitySet
}
