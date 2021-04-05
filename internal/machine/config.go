// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package machine

import (
	"fmt"

	"github.com/c4-project/c4t/internal/id"
	"github.com/c4-project/c4t/internal/model/litmus"
	"github.com/c4-project/c4t/internal/model/service/compiler"
	"github.com/c4-project/c4t/internal/mutation"
)

// Config is a config record for a particular machine.
//
// The difference between a Machine and a Config is that the latter contains raw configuration data for things that get
// mapped into expanded forms in the plan, for instance compilers.
type Config struct {
	Machine

	// Arch is the default architecture for compilers in this configuration.
	Arch id.ID `toml:"arch,omitempty"`

	// RawCompilers contains raw information about the compilers attached to this machine.
	//
	// The keys are strings that will become IDs in the fully processed compiler map.
	// We don't store them as IDs here, to avoid issues with TOML serialisation.
	//
	// This doesn't contain machine-level defaults; use Compilers() to get a fully resolved version.
	RawCompilers map[string]compiler.Config `toml:"compilers,omitempty"`

	// Mutation contains information about how to mutation-test on this machine.
	Mutation *mutation.Config `toml:"mutation,omitempty"`
}

// Compilers prepares a fully resolved compiler map, with any machine defaults filled in.
// It errors if there are missing parts of compiler configuration for a particular compiler,
// or if there is a problem with the compiler ID.
//
// This is always a separate map from RawCompilers, even when no defaults exist.
func (c *Config) Compilers() (map[id.ID]compiler.Compiler, error) {
	cs := make(map[id.ID]compiler.Compiler, len(c.RawCompilers))
	for n, raw := range c.RawCompilers {
		if err := c.addAndExpandCompiler(cs, n, raw); err != nil {
			return nil, err
		}
	}
	return cs, nil
}

func (c *Config) addAndExpandCompiler(cs map[id.ID]compiler.Compiler, n string, raw compiler.Config) error {
	cid, err := id.TryFromString(n)
	if err != nil {
		return fmt.Errorf("name of compiler %s: %w", n, err)
	}
	if cs[cid], err = c.prepareCompiler(raw); err != nil {
		return fmt.Errorf("compiler %s: %w", n, err)
	}
	return nil
}

// prepareCompiler expands a compiler by applying machine defaults where needed.
func (c *Config) prepareCompiler(raw compiler.Config) (compiler.Compiler, error) {
	prep := (compiler.Compiler)(raw)
	if prep.Arch.IsEmpty() {
		prep.Arch = c.Arch
	}
	if prep.Arch.IsEmpty() {
		// TODO(@MattWindsor91): error should be moved.
		return prep, litmus.ErrEmptyArch
	}
	return prep, nil
}
