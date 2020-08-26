// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package config

import "github.com/MattWindsor91/act-tester/internal/helper/iohelp"

// Pathset is a set of configuration paths (all of which are considered to be filepaths, not slashpaths).
type Pathset struct {
	// OutDir is the default output path for the test director.
	OutDir string `toml:"out_dir,omitempty,omitzero"`
	// Inputs is the default set of inputs (files and/or directories) for the test director.
	Inputs []string `toml:"inputs,omitempty,omitzero"`
	// FilterFile is, if present, a path pointing to a YAML file containing analysis filter definitions.
	FilterFile string `toml:"filter_file,omitempty,omitzero"`
}

// FallbackToInputs returns fs if non-empty, and the homedir-expanded version of Pathset.Inputs on p otherwise.
func (p Pathset) FallbackToInputs(fs []string) ([]string, error) {
	if len(fs) != 0 {
		return fs, nil
	}
	return iohelp.ExpandMany(p.Inputs)
}
