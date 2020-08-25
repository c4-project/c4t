// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package plan

// Type of plan version numbers.
type Version uint32

// CurrentVer is the current plan version.
// It changes when the interface between various bits of the tester (generally manifested within the plan version)
// changes.
const CurrentVer Version = 2020_08_25

// Version history since 2020_05_29:
//
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
