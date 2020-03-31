// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package config

import (
	"os"
	"path"

	"github.com/BurntSushi/toml"
)

const (
	// dirConfig is the subdirectory under the user config directory in which act-tester will check for a config file.
	dirConfig = "act"
	// fileConfig is the default name that act-tester will use when looking for a config file.
	fileConfig = "tester.toml"
)

// Load tries to load a tester config from various places.
// If f is non-empty, it tries there.
// Else, it first tries the current working directory, and then tries the user config directory.
func Load(f string) (*Config, error) {
	if f != "" {
		return tryLoad(f)
	}

	c, err := loadConfigCWD()
	if err != nil {
		return loadConfigUCD()
	}
	return c, err
}

func loadConfigCWD() (*Config, error) {
	return tryLoad(fileConfig)
}

func loadConfigUCD() (*Config, error) {
	cdir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}
	return tryLoad(path.Join(cdir, dirConfig, fileConfig))
}

func tryLoad(f string) (*Config, error) {
	var c Config
	_, derr := toml.DecodeFile(f, &c)
	return &c, derr
}
