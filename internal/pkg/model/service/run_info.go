// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package service

// RunInfo gives hints as to how to run a service.
type RunInfo struct {
	// Cmd overrides the command for the service.
	Cmd string `toml:"cmd,omitzero"`

	// Args specifies (extra) arguments to supply to the service.
	Args []string `toml:"args,omitempty"`
}

// Override creates run information by overlaying this run information with that in new.
func (r RunInfo) Override(new *RunInfo) RunInfo {
	if new == nil {
		return r
	}
	r.Cmd = overrideCmd(r.Cmd, new.Cmd)
	r.Args = append(r.Args, new.Args...)
	return r
}

func overrideCmd(old, new string) string {
	if new == "" {
		return old
	}
	return new
}
