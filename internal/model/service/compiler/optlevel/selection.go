// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package optlevel

import "github.com/c4-project/c4t/internal/helper/stringhelp"

// Selection represents a piece of compiler configuration that specifies which optimisation levels to select.
type Selection struct {
	// Enabled overrides the default selection to insert optimisation levels.
	Enabled []string `toml:"enabled,omitempty"`
	// Disabled overrides the default selection to remove optimisation levels.
	Disabled []string `toml:"disabled,omitempty"`
}

// Select inserts enables from this selection into defaults, then removes disables.
// Disables take priority over enables.
// The resulting map is a copy; this function doesn't mutate defaults.
func (s Selection) Override(defaults stringhelp.Set) stringhelp.Set {
	nm := defaults.Copy()
	nm.Add(s.Enabled...)
	nm.Remove(s.Disabled...)
	return nm
}
