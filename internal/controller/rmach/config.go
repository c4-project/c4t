// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package rmach

import (
	"context"

	"github.com/1set/gut/ystring"
	"github.com/MattWindsor91/act-tester/internal/model/plan"

	"github.com/MattWindsor91/act-tester/internal/transfer/remote"
)

// Config contains the parts of a remote-machine configuration that don't depend on the plan.
type Config struct {
	// DirLocal is the filepath to the directory to which local outcomes from this rmach run will appear.
	DirLocal string

	// Observers is the set of observers listening for file copying and remote corpus manipulations.
	Observers ObserverSet

	/// SSH tells the remote-machine invoker how to use SSH on the host machine.
	SSH *remote.Config
}

// Check makes sure all of the configuration is present and accounted-for.
func (c *Config) Check() error {
	if ystring.IsBlank(c.DirLocal) {
		return ErrDirEmpty
	}
	if c.SSH == nil {
		return ErrSSHNil
	}
	return nil
}

// Run abstracts over constructing a rmach from this config and running it on a plan.
func (c *Config) Run(ctx context.Context, p *plan.Plan) (*plan.Plan, error) {
	r, err := New(c, p)
	if err != nil {
		return nil, err
	}
	return r.Run(ctx)
}
