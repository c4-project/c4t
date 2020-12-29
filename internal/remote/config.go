// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package remote

import (
	"path/filepath"

	"github.com/c4-project/c4t/internal/helper/iohelp"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

// Config is the top-level configuration for c4t's SSH support.
type Config struct {
	// SCPTimeoutMins is a timeout to apply to SCP.
	SCPTimeoutMins int `toml:"scp_timeout,omitzero"`

	// KnownHostsFilePaths is a list of raw filepaths to SSH known-hosts file.
	KnownHostsFilePaths []string `toml:"known_hosts_paths,omitempty"`
}

// knownHosts gets a known-host callback given the known-host paths in this config.
func (c *Config) knownHosts() (ssh.HostKeyCallback, error) {
	paths, err := c.knownHostsPaths()
	if err != nil {
		return nil, err
	}
	return knownhosts.New(paths...)
}

// knownHostsPaths gets a list of slash-delimited, expanded paths to check for known-hosts files.
func (c *Config) knownHostsPaths() ([]string, error) {
	fpaths := append(c.KnownHostsFilePaths, filepath.Join("~", ".ssh", "known_hosts"))
	return iohelp.ExpandMany(fpaths)
}
