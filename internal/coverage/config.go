// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package coverage

import (
	"errors"

	backend2 "github.com/MattWindsor91/c4t/internal/model/service/backend"

	"github.com/MattWindsor91/c4t/internal/model/litmus"

	"github.com/MattWindsor91/c4t/internal/stage/lifter"
	"github.com/mitchellh/go-homedir"

	"github.com/MattWindsor91/c4t/internal/model/id"

	fuzzer2 "github.com/MattWindsor91/c4t/internal/model/service/fuzzer"

	"github.com/MattWindsor91/c4t/internal/stage/fuzzer"

	"github.com/MattWindsor91/c4t/internal/config"
	"github.com/MattWindsor91/c4t/internal/model/service"
	"github.com/pelletier/go-toml"
)

// ErrConfigNil is produced when we supply a null pointer to OptionsFromConfig.
var ErrConfigNil = errors.New("supplied config is nil")

// Profile tells the coverage generator how to set up a particular coverage profile.
type Profile struct {
	// Kind specifies the type of fuzzer profile this is.
	Kind ProfileKind `toml:"kind"`

	// Arch is the target architecture for the profile, if it uses one.
	Arch id.ID `toml:"arch"`

	// Backend directly feeds in the target backend for the profile, if it uses one.
	Backend *backend2.Spec `toml:"backend"`

	// Run specifies, if this is a standalone profile, how to run the generator.
	Run *service.RunInfo `toml:"run"`

	// Fuzz specifies a fuzzer configuration to use if this is an known-fuzzer profile.
	Fuzz *fuzzer2.Configuration `toml:"fuzz"`

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

func OptionsFromConfig(cfg *Config) Option {
	return Options(
		OverrideQuantities(cfg.Quantities),
		AddInputs(cfg.Paths.Inputs...),
	)
}

// UseFuzzer adds support for f as a 'known' fuzzer.
func UseFuzzer(f fuzzer.SingleFuzzer) Option {
	return func(maker *Maker) error {
		// TODO(@MattWindsor91): multiple known fuzzers
		maker.fuzz = f
		return nil
	}
}

// UseStatDumper adds support for d as the statistics dumper.
func UseStatDumper(d litmus.StatDumper) Option {
	return func(maker *Maker) error {
		maker.sdump = d
		return nil
	}
}

// UseLifter adds support for l as the source of lifters.
func UseLifter(l lifter.SingleLifter) Option {
	return func(maker *Maker) error {
		maker.lift = l
		return nil
	}
}

// MakeMaker makes a maker using the configuration in cfg.
func (cfg *Config) MakeMaker(o ...Option) (*Maker, error) {
	if cfg == nil {
		return nil, ErrConfigNil
	}
	od, err := homedir.Expand(cfg.Paths.OutDir)
	if err != nil {
		return nil, err
	}
	return NewMaker(od, cfg.Profiles,
		OptionsFromConfig(cfg),
		Options(o...),
	)
}
