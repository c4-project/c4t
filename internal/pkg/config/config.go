// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package config describes the top-level tester configuration.

// TODO(@MattWindsor91): slowly wrest control of the configuration from OCaml act.

package config

import (
	"github.com/MattWindsor91/act-tester/internal/pkg/model"
)

// Config is a top-level tester config struct.
type Config struct {
	// Backend contains information about the backend being used to generate test harnesses.
	Backend *model.Backend `toml:"backend,omitempty"`

	// Machines enumerates the machines available for testing.
	Machines map[string]Machine `toml:"machines,omitempty"`

	// OutDir is the output directory for fully directed test runs.
	OutDir string `toml:"out_dir"`
}
