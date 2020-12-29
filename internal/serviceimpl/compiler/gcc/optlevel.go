// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package gcc

import (
	"github.com/c4-project/c4t/internal/helper/stringhelp"
	"github.com/c4-project/c4t/internal/model/service/compiler"
	"github.com/c4-project/c4t/internal/model/service/compiler/optlevel"
)

var (
	// OptLevels contains the optimisation levels known to exist on GCC, Clang, AppleClang etc.
	OptLevels = map[string]optlevel.Level{
		// no optimisation
		"0": {
			Optimises:       false,
			Bias:            optlevel.BiasDebug,
			BreaksStandards: false,
		},
		// mild optimisation
		"1": {
			Optimises:       true,
			Bias:            optlevel.BiasSpeed,
			BreaksStandards: false,
		},
		// moderate optimisation
		"2": {
			Optimises:       true,
			Bias:            optlevel.BiasSpeed,
			BreaksStandards: false,
		},
		// heavy optimisation
		"3": {
			Optimises:       true,
			Bias:            optlevel.BiasSpeed,
			BreaksStandards: false,
		},
		// standards-bending optimisation
		"fast": {
			Optimises:       true,
			Bias:            optlevel.BiasSpeed,
			BreaksStandards: true,
		},
		// optimise for size
		"s": {
			Optimises:       true,
			Bias:            optlevel.BiasSize,
			BreaksStandards: false,
		},
		// AppleClang only?
		"z": {
			Optimises:       true,
			Bias:            optlevel.BiasSize,
			BreaksStandards: false,
		},
		// debug-friendly optimisation
		"g": {
			Optimises:       true,
			Bias:            optlevel.BiasDebug,
			BreaksStandards: false,
		},
		// 'equivalent to O2'
		"": {
			Optimises:       true,
			Bias:            optlevel.BiasSpeed,
			BreaksStandards: false,
		},
	}

	// OptLevelNames is a consistently named list of the optimisation levels in OptLevels.
	OptLevelNames = []string{"", "0", "1", "2", "3", "fast", "s", "z", "g"}

	// OptLevelDisabledNames contains optimisation levels that are disabled by default, as they are redundant or non-portable.
	OptLevelDisabledNames = []string{"", "0", "g", "z"}
)

// DefaultOptLevels gets the default level set for GCC.
func (g GCC) DefaultOptLevels(_ *compiler.Compiler) (stringhelp.Set, error) {
	sel := optlevel.Selection{
		Enabled:  OptLevelNames,
		Disabled: OptLevelDisabledNames,
	}
	return sel.Override(nil), nil
}

func (_ GCC) OptLevels(_ *compiler.Compiler) (map[string]optlevel.Level, error) {
	return OptLevels, nil
}
