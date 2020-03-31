// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package gcc

import "github.com/MattWindsor91/act-tester/internal/model/compiler/optlevel"

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

	// TODO(@MattWindsor91): use this
	// OptLevelDisabledNames contains optimisation levels that are disabled by default, as they aren't portable.
	OptLevelDisabledNames = []string{"g", "z"}
)

// DefaultLevels gets the default level set for GCC.
func (g GCC) DefaultLevels() map[string]struct{} {
	sel := optlevel.Selection{
		Enabled:  OptLevelNames,
		Disabled: OptLevelDisabledNames,
	}
	return sel.Override(nil)
}

func (_ GCC) Levels() map[string]optlevel.Level {
	return OptLevels
}
