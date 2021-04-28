// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package service

import (
	"fmt"
	"sort"
	"strings"

	"github.com/1set/gut/ystring"

	"github.com/c4-project/c4t/internal/id"

	"github.com/buildkite/interpolate"
	"github.com/c4-project/c4t/internal/helper/stringhelp"
)

// RunInfo gives hints as to how to run a service.
type RunInfo struct {
	// Cmd overrides the command for the service.
	Cmd string `toml:"cmd,omitzero" json:"cmd,omitempty"`

	// Args specifies (extra) arguments to supply to the service.
	Args []string `toml:"args,omitempty" json:"args,omitempty"`

	// Env defines environment variables to be set if the runner supports it.
	Env map[string]string `toml:"env,omitempty" json:"env,omitempty"`
}

// NewRunInfo programmatically creates a RunInfo using command cmd and arguments args.
func NewRunInfo(cmd string, args ...string) *RunInfo {
	return &RunInfo{Cmd: cmd, Args: args}
}

// NewIfDifferent returns a new RunInfo overriding this RunInfo's Cmd with cmd if it is different, and nil otherwise.
// This is useful for generating service configuration, if we only want to supply a specific non-default RunInfo when
// there is a real difference from the current default.
//
// This function may eventually also take arguments; if so, they will be supplied variadically.
func (r *RunInfo) NewIfDifferent(cmd string) *RunInfo {
	if cmd == r.Cmd {
		return nil
	}
	return &RunInfo{Cmd: cmd}
}

// SystematicID produces an ID based any present elements of this RunInfo.
// If the RunInfo is the zero value, the ID will be empty.
// It is not formally guaranteed to be unique, but should be close enough.
func (r *RunInfo) SystematicID() (id.ID, error) {
	// TODO(@MattWindsor91): strip - and ., perhaps

	tags := make([]string, 0, 1+len(r.Args)+len(r.Env))
	if ystring.IsNotEmpty(r.Cmd) {
		tags = append(tags, r.Cmd)
	}
	tags = append(tags, r.Args...)
	tags = append(tags, r.EnvStrings()...)

	return id.New(tags...)
}

// Invocation is Cmd appended to Args.
func (r *RunInfo) Invocation() []string {
	return append([]string{r.Cmd}, r.Args...)
}

// EnvStrings is the environment as a set of key-value pairs.
// This is compatible with exec.Cmd's Env field.
func (r *RunInfo) EnvStrings() []string {
	if len(r.Env) == 0 {
		return nil
	}

	ks, err := stringhelp.MapKeys(r.Env)
	if err != nil {
		return []string{fmt.Sprintf("(err: %s)", err)}
	}
	sort.Strings(ks)
	for i, k := range ks {
		ks[i] = fmt.Sprintf("%s=%s", k, r.Env[k])
	}

	return ks
}

// String is defined as the space-joined form of EnvString and Invocation.
func (r *RunInfo) String() string {
	parts := append(r.EnvStrings(), r.Invocation()...)
	return strings.Join(parts, " ")
}

// Override overlays this run information with that in new.
func (r *RunInfo) Override(new RunInfo) {
	r.Cmd = overrideCmd(r.Cmd, new.Cmd)
	// TODO(@MattWindsor91): we might need a way to replace arguments rather than appending to them.
	r.AppendArgs(new.Args...)
	r.overrideEnv(new.Env)
}

// AppendArgs overlays new arguments onto this run info.
func (r *RunInfo) AppendArgs(new ...string) {
	r.Args = append(r.Args, new...)
}

// Interpolate expands any interpolations in this run info's arguments and environment according to expansions.
// While it modifies this RunInfo in place, it creates a new argument slice and environment map.
func (r *RunInfo) Interpolate(expansions map[string]string) error {
	env := interpolate.NewMapEnv(expansions)
	if err := r.interpolateArgs(env); err != nil {
		return err
	}
	return r.interpolateEnv(env)
}

func (r *RunInfo) interpolateArgs(env interpolate.Env) error {
	var err error
	args := make([]string, len(r.Args))
	for i, arg := range r.Args {
		if args[i], err = interpolate.Interpolate(env, arg); err != nil {
			return err
		}
	}
	r.Args = args
	return nil
}

func (r *RunInfo) interpolateEnv(env interpolate.Env) error {
	var err error
	nenv := make(map[string]string, len(r.Env))
	for k, v := range r.Env {
		if nenv[k], err = interpolate.Interpolate(env, v); err != nil {
			return err
		}
	}
	r.Env = nenv
	return nil
}

// OverrideIfNotNil is Override if new is non-nil, and no-op otherwise.
func (r *RunInfo) OverrideIfNotNil(new *RunInfo) {
	if new != nil {
		r.Override(*new)
	}
}

func (r *RunInfo) overrideEnv(env map[string]string) {
	if len(r.Env) == 0 {
		r.Env = env
	}
	// In case we're sharing this environment by reference.
	nenv := make(map[string]string, len(r.Env))
	for k, v := range r.Env {
		nenv[k] = v
	}
	for k, v := range env {
		nenv[k] = v
	}
	r.Env = nenv
}

func overrideCmd(old, new string) string {
	if new == "" {
		return old
	}
	return new
}
