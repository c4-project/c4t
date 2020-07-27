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
const CurrentVer Version = 2020_07_29

// Version history since 2020_05_29:
//
// 2020_07_27: Added new (sub-)stages: Compile and Run.
// 2020_07_21: Added stage information.
// 2020_05_29: Initial version for which this comment was maintained.

// IsCurrent is true if, and only if, v is compatible with the current version.
func (v Version) IsCurrent() bool {
	return v == CurrentVer
}
