// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package remote

import (
	"fmt"
	"net"
	"os"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

// MachineConfig is SSH configuration for a remote machine.
type MachineConfig struct {
	// The host to use when dialing into the machine.
	Host string `toml:"host"`
	// The user to use when dialing into the machine.
	User string `toml:"user,omitzero"`
	// The port to use when dialing into the machine.
	// If zero, defaults to 22.
	Port int `toml:"port,omitzero"`
	// The directory to which we shall copy intermediate files.
	DirCopy string `toml:"copy_dir"`
}

// MachineRunner encapsulates information about how to run jobs remotely through SSH.
type MachineRunner struct {
	// Config points to the machine configuration that was used to create this runner.
	Config *MachineConfig
	ssh    *ssh.ClientConfig
	cli    *ssh.Client
}

// NewSession opens a new SSH session.
func (r *MachineRunner) NewSession() (*ssh.Session, error) {
	return r.cli.NewSession()
}

// NewSFTP opens a new SFTP session.
func (r *MachineRunner) NewSFTP() (*sftp.Client, error) {
	return sftp.NewClient(r.cli)
}

// Close closes this MachineRunner's underlying SSH connection.
func (r *MachineRunner) Close() error {
	if r.cli == nil {
		return nil
	}
	return r.cli.Close()
}

// Runner gets a SSH runner for this machine, given the configuration in c.
func (m *MachineConfig) MachineRunner(c *Config) (*MachineRunner, error) {
	s, err := m.clientConfig(c)
	if err != nil {
		return nil, err
	}
	cli, err := ssh.Dial("tcp", m.hostPort(), s)
	if err != nil {
		return nil, err
	}
	mr := MachineRunner{
		Config: m,
		cli:    cli,
		ssh:    s,
	}
	return &mr, nil
}

func (m *MachineConfig) hostPort() string {
	return fmt.Sprintf("%s:%d", m.Host, m.portOrDefault())
}

func (m *MachineConfig) portOrDefault() int {
	if m.Port == 0 {
		return 22
	}
	return m.Port
}

// clientConfig gets the SSH config for this machine, given the global configuration c.
func (m *MachineConfig) clientConfig(c *Config) (*ssh.ClientConfig, error) {
	// Fall back to defaults if c is nil.
	if c == nil {
		c = &Config{}
	}

	auths, err := authMethods()
	if err != nil {
		return nil, fmt.Errorf("while getting auth methods: %w", err)
	}
	kh, err := c.knownHosts()
	if err != nil {
		return nil, fmt.Errorf("while getting known-hosts: %w", err)
	}

	cfg := ssh.ClientConfig{
		User:            m.User,
		Auth:            auths,
		HostKeyCallback: kh,
	}

	return &cfg, nil
}

// authMethods gets together a set of authentication methods.
func authMethods() ([]ssh.AuthMethod, error) {
	// TODO(@MattWindsor91): alternative arrangements for when we don't have an SSH agent
	agentClient, err := getAgent()
	if err != nil {
		return nil, err
	}
	return []ssh.AuthMethod{ssh.PublicKeysCallback(agentClient.Signers)}, nil
}

func getAgent() (agent.ExtendedAgent, error) {
	socket := os.Getenv("SSH_AUTH_SOCK")
	conn, err := net.Dial("unix", socket)
	if err != nil {
		return nil, fmt.Errorf("failed to open SSH_AUTH_SOCK: %w", err)
	}
	return agent.NewClient(conn), nil
}
