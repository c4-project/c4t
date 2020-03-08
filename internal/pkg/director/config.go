// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package director

import (
	"errors"
	"log"

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
}

// ConfigFromGlobal extracts the parts of a global config file relevant to a director, and builds a config from them.
func ConfigFromGlobal(g *config.Config, l *log.Logger) (*Config, error) {
	if g == nil {
		return nil, config.ErrNil
	}
	if g.Backend == nil {
		return nil, errors.New("config has no backend")
	}
	if g.OutDir == "" {
		return nil, ErrNoOutDir
	}

	return &Config{Logger: l, Paths: NewPathset(g.OutDir), Machines: g.Machines}, nil
}
