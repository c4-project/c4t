// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml"

	"github.com/1set/gut/yos"
	"github.com/1set/gut/ystring"
)

const (
	// dirConfig is the subdirectory under the user config directory in which c4t will check for a config file.
	dirConfig = "c4t"
	// fileConfig is the default name that c4t will use when looking for a config file.
	fileConfig = "tester.toml"
)

// NoConfigFileError wraps an error when trying to resolve a config file through FilePath.
type NoConfigFileError struct {
	// Tried is the list of config files that have been tried.
	Tried []string
}

// Error returns the error string for this error.
func (n *NoConfigFileError) Error() string {
	return fmt.Sprintf("no config file found (tried %s)", strings.Join(n.Tried, ", "))
}

// Unwrap returns os.ErrNotExist, as this is a particular kind of not-exist error.
func (n *NoConfigFileError) Unwrap() error {
	return os.ErrNotExist
}

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
	t, err := toml.LoadFile(f)
	if err != nil {
		return nil, err
	}
	var c Config
	return &c, t.Unmarshal(&c)
}

// FilePath resolves the config file path file using hooks.
// The first hook to successfully find a valid file wins.
func FilePath(file string, hooks ...func(string) (string, error)) (string, error) {
	tried := make([]string, 0, len(hooks))

	for _, h := range hooks {
		path, err := h(file)
		if err != nil {
			return "", err
		}
		if ystring.IsBlank(path) {
			continue
		}
		if yos.ExistFile(path) {
			return path, nil
		}
		tried = append(tried, path)
	}
	return "", &NoConfigFileError{Tried: tried}
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
