// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package optlevel

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
func (s Selection) Override(defaults map[string]struct{}) map[string]struct{} {
	nm := make(map[string]struct{}, len(defaults)+len(s.Enabled))
	for _, o := range s.Enabled {
		nm[o] = struct{}{}
	}
	for o := range defaults {
		nm[o] = struct{}{}
	}
	for _, o := range s.Disabled {
		delete(nm, o)
	}
	return nm
}
