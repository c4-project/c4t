// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package director

import (
	"errors"
	"log"
	"path/filepath"

	"github.com/mitchellh/go-homedir"

	"github.com/MattWindsor91/act-tester/internal/pkg/planner"

	"github.com/MattWindsor91/act-tester/internal/pkg/config"
)

var (
	// ErrConfigNil occurs when we try to build a director from a nil config.
	ErrConfigNil = errors.New("config nil")

	// ErrNoMachines occurs when we try to build a director from a config with no machines defined.
	ErrNoMachines = errors.New("no machines defined in config")

	// ErrNoOutDir occurs when we try to build a Director with no output directory specified in the config.
	ErrNoOutDir = errors.New("no output directory specified in config")
)

// Config groups together the various bits of configuration needed to create a director.
type Config struct {
	// Logger is the logger to which the director should log.
	// If nil, logging will proceed silently.
	Logger *log.Logger

	// Paths provides path resolving functionality for the director.
	Paths *Pathset

	// Machines contains the machines that will be used in the test.
	Machines map[string]config.Machine

	// Env groups together the bits of configuration that pertain to dealing with the environment.
	Env Env
}

// Env groups together the bits of configuration that pertain to dealing with the environment.
type Env struct {
	// Planner instructs any planners built for this director as to how to acquire information about compilers, etc.
	Planner planner.Source
}

// ConfigFromGlobal extracts the parts of a global config file relevant to a director, and builds a config from them.
func ConfigFromGlobal(g *config.Config, l *log.Logger, e Env) (*Config, error) {
	if g == nil {
		return nil, config.ErrNil
	}
	if g.Backend == nil {
		return nil, errors.New("config has no backend")
	}
	if g.OutDir == "" {
		return nil, ErrNoOutDir
	}

	edir, err := homedir.Expand(g.OutDir)
	if err != nil {
		return nil, err
	}
	odir := filepath.ToSlash(edir)

	return &Config{Logger: l, Env: e, Paths: NewPathset(odir), Machines: g.Machines}, nil
}
