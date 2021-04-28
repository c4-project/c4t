// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package service

import "context"

// ExtClass contains fields common to a service class that invokes an external binary.
type ExtClass struct {
	// DefaultRunInfo captures the default run information for this class.
	DefaultRunInfo RunInfo

	// AltCommands contains possible alternative commands that can be substituted into DefaultRunInfo.
	// This is mainly for probing purposes.
	AltCommands []string
}

// ProbeByVersionCommand runs a version command for every command in this ExtClass's DefaultRunInfo and AltCommands,
// formed by appending args to the DefaultRunInfo arguments.  It returns a map from command names to versions, where
// a command name is mapped only if the version command succeeded.
func (e ExtClass) ProbeByVersionCommand(ctx context.Context, r Runner, args ...string) map[string]string {
	cmds := append(e.AltCommands, e.DefaultRunInfo.Cmd)
	versions := make(map[string]string, len(cmds))

	for _, cmd := range cmds {
		ri := e.DefaultRunInfo
		ri.Override(RunInfo{Cmd: cmd, Args: args})
		if ver, err := RunAndCaptureStdout(ctx, r, ri); err == nil {
			versions[cmd] = ver
		}
	}

	return versions
}
