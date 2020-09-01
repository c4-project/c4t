// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/1set/gut/yos"
	"github.com/1set/gut/ystring"

	"github.com/BurntSushi/toml"
)

const (
	// dirConfig is the subdirectory under the user config directory in which act-tester will check for a config file.
	dirConfig = "act"
	// fileConfig is the default name that act-tester will use when looking for a config file.
	fileConfig = "tester.toml"
)

// Load tries to load a tester config from various places.
// If override is non-empty, it tries there.
// Else, it first tries the current working directory, and then tries the user config directory.
func Load(override string) (*Config, error) {
	path, err := FilePath(fileConfig, StandardHooks(override)...)
	if err != nil {
		return nil, err
	}
	return tryLoad(path)
}

func tryLoad(f string) (*Config, error) {
	var c Config
	_, derr := toml.DecodeFile(f, &c)
	return &c, derr
}

// FilePath resolves the config file path file using hooks.
// The first hook to successfully find a valid file wins.
func FilePath(file string, hooks ...func(string) (string, error)) (string, error) {
	for _, h := range hooks {
		path, err := h(file)
		if err != nil {
			return "", err
		}
		if ystring.IsNotBlank(path) && yos.ExistFile(path) {
			return path, nil
		}
	}
	return "", fmt.Errorf("%w: config file %s", os.ErrNotExist, file)
}

// StandardHooks is the default set of hooks to use when getting a config file.
// These hooks are as follows:
// - Try the literal path override;
// - Try looking in the current working directory;
// - Try looking in the Go user config directory.
func StandardHooks(override string) []func(string) (string, error) {
	return []func(string) (string, error){
		func(string) (string, error) { return override, nil },
		configCWD,
		configUCD,
	}
}

func configCWD(path string) (string, error) {
	return path, nil
}

func configUCD(path string) (string, error) {
	cdir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(cdir, dirConfig, path), nil
}
