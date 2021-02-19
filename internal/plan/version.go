// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package plan

// Type of plan version numbers.
type Version uint32

// CurrentVer is the current plan version.
// It changes when the interface between various bits of the tester (generally manifested within the plan version)
// changes.
const CurrentVer Version = 2021_02_19

// Version history since 2020_05_29:
//
// 2021_02_19: Everything tracking time plus duration has been standardised to take a "time_span" key; this contains a
//             "start" time and an "end" time.  Currently, both can be in different timezones.
//             Mutation analysis has changed to associate selections, hits, and kills with such timespans.
// 2021_01_27: Observation structures completely refactored.  Instead of "counter_examples" and "witnesses" keys, each
//             state has a "tag" key that contains "counter", "witness", "unknown", or is empty.  New "occurrences" key
//             in each state tracks number of occurrences where available; key-value maps are now inside a "values" key
//             inside the state.
// 2021_01_26: New stages: SetCompiler and Mach, resulting from all plan manipulating tools becoming stages.
// 2021_01_24: Fixes to keys.
// 2021_01_22: Changes to 'mutation' key; auto-incrementing config is now in an 'auto' subkey.
// 2021_01_21: Added 'mutation' key to plan, containing information about mutation testing.
// 2021_01_20: Service run info now contains an 'env' key mapping environment variable names to values.  If 'env' is
//             not provided, the parent process's environment is used.  On compilers only, we now have shell-style
//             variable interpolation in args and env: the only variable available so far is '${time}' which expands to
//             the UNIX timestamp at which the compiler was configured.
// 2020_12_10: Plan backends are now NamedSpecs rather than Specs, and so contain the backend ID.
// 2020_11_12: Compile results now hold the ID of the recipe used.  Recipes have extended information about targets.
// 2020_09_24: SSH configuration now uses camel_case tags.  The 'compiles' and 'runs' maps have become one
//             'compilations' map, with optional 'compile' and 'run' subkeys (these currently contain exactly the same
//             data as the original map entries).
// 2020_08_25: Machine configuration in plans now carries machine-specific quantity overrides.
// 2020_07_30: New 'perturb' stage.  Some changes to observations that may alter the interface with the machine node.
// 2020_07_28: No changes to the plan per se, but the machine node no longer supports human-readable output, and the
//             JSON mode flag has been removed.
// 2020_07_27: Added new (sub-)stages: Compile and Run.  (A typo meant that this version got stored as 2020_07_29,
//             but at this point versions are compared for equality rather than ranges, so there is no practical
//             issue.)
// 2020_07_21: Added stage information.
// 2020_05_29: Initial version for which this comment was maintained.

// IsCurrent is true if, and only if, v is compatible with the current version.
func (v Version) IsCurrent() bool {
	return v == CurrentVer
}
