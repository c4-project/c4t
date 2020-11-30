// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package config

import (
	"path/filepath"

	"github.com/MattWindsor91/c4t/internal/helper/iohelp"
	"github.com/mitchellh/go-homedir"
)

// Pathset is a set of configuration paths (all of which are considered to be filepaths, not slashpaths).
type Pathset struct {
	// OutDir is the default output path for the test director.
	OutDir string `toml:"out_dir,omitempty,omitzero"`
	// Inputs is the default set of inputs (files and/or directories) for the test director.
	Inputs []string `toml:"inputs,omitempty,omitzero"`

	// TODO(@MattWindsor91): delete FilterFile, turn it into convention over configuration?

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

// OutPath is shorthand for getting a homedir-expanded full path for the output filename file.
func (p Pathset) OutPath(file string) (string, error) {
	return homedir.Expand(filepath.Join(p.OutDir, file))
}
