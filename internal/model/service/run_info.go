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

// Override overlays this run information with that in new.
func (r *RunInfo) Override(new RunInfo) {
	r.Cmd = overrideCmd(r.Cmd, new.Cmd)
	// TODO(@MattWindsor91): we might need a way to replace arguments rather than appending to them.
	r.Args = append(r.Args, new.Args...)
}

func overrideCmd(old, new string) string {
	if new == "" {
		return old
	}
	return new
}
