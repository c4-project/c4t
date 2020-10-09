// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package coverage

import (
	"errors"

	"github.com/MattWindsor91/act-tester/internal/config"
	"github.com/MattWindsor91/act-tester/internal/model/service"
	"github.com/pelletier/go-toml"
)

// ErrConfigNil is produced when we supply a null pointer to OptionsFromConfig.
var ErrConfigNil = errors.New("supplied config is nil")

// Profile tells the coverage generator how to set up a particular coverage profile.
type Profile struct {
	// Kind specifies the type of fuzzer profile this is.
	Kind ProfileKind `toml:"kind"`

	// Run specifies, if this is a standalone profile, how to run the generator.
	Run *service.RunInfo `toml:"run"`

	// Runner specifies an overridden runner for the profile; this is basically useful only for testing.
	Runner Runner
}

// Config gathers the configuration present in coverage generator config files.
type Config struct {
	// Paths contains the input and output pathsets for the coverage generator.
	Paths config.Pathset `toml:"paths"`

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

// MakeMaker makes a maker using the configuration in cfg.
func (cfg *Config) MakeMaker(o ...Option) (*Maker, error) {
	if cfg == nil {
		return nil, ErrConfigNil
	}
	return NewMaker(cfg.Paths.OutDir, cfg.Profiles,
		OverrideQuantities(cfg.Quantities),
		AddInputs(cfg.Paths.Inputs...),
		Options(o...),
	)
}
