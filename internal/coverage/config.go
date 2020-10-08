// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package coverage

import (
	"github.com/MattWindsor91/act-tester/internal/model/service"
	toml "github.com/pelletier/go-toml"
)

// Profile tells the coverage generator how to set up a particular coverage profile.
type Profile struct {
	// Kind specifies the type of fuzzer profile this is.
	Kind ProfileKind `toml:"kind"`

	// Run specifies, if this is a standalone profile, how to run the generator.
	Run *service.RunInfo `toml:"run"`
}

// QuantitySet contains the quantities tracked by the coverage generator.
type QuantitySet struct {
	// Count is the number of subjects to fuzz for each profile.
	Count int `toml:"count"`

	// Divisions specifies how to divide Count subjects into buckets.
	// Divisions behave recursively: each subsequent level of division gets applied to the first bucket in the
	// previous level.
	Divisions []int `toml:"divisions"`
}

// Config gathers the configuration present in coverage generator config files.
type Config struct {
	// Quantities contains quantities for the coverage generator.
	Quantities QuantitySet `toml:"quantities"`

	// Profiles contains the list of coverage profiles to use.
	Profiles map[string]Profile `toml:"profiles"`
}

// LoadConfigFromFile loads a coverage configuration from the filepath path.
func LoadConfigFromFile(path string) (*Config, error) {
	tree, err := toml.LoadFile(path)
	if err != nil {
		return nil, err
	}
	var c Config
	err = tree.Unmarshal(&c)
	return &c, err
}
