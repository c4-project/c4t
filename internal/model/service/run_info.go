// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package service

import (
	"context"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"

	"github.com/c4-project/c4t/internal/helper/stringhelp"
)

// Runner is the interface of things that can run, or pretend to run, services.
type Runner interface {
	// WithStdout should return a new runner with the standard output overridden to w.
	WithStdout(w io.Writer) Runner

	// WithStderr should return a new runner with the standard error overridden to w.
	WithStderr(w io.Writer) Runner

	// WithGrace should return a new runner with the timeout grace period set to d.
	WithGrace(d time.Duration) Runner

	// Run runs r using context ctx.
	Run(ctx context.Context, r RunInfo) error
}

//go:generate mockery --name=Runner

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
	for k, v := range env {
		r.Env[k] = v
	}
}

func overrideCmd(old, new string) string {
	if new == "" {
		return old
	}
	return new
}
