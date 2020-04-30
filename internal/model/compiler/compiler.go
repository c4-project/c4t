// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package compiler contains types for compilers, which are a particular type of service.
package compiler

import (
	"fmt"
	"strings"

	"github.com/1set/gut/ystring"

	"github.com/MattWindsor91/act-tester/internal/model/compiler/optlevel"
	"github.com/MattWindsor91/act-tester/internal/model/id"
	"github.com/MattWindsor91/act-tester/internal/model/service"
)

// Config collects the part of a compiler's specification that comes from the act-tester configuration.
type Config struct {
	// Disabled specifies whether this compiler has been disabled.
	Disabled bool `toml:"disabled,omitempty"`

	// Style is the declared style of the compile.
	Style id.ID `toml:"style"`

	// Arch is the architecture (or 'emits') ID for the compiler.
	Arch id.ID `toml:"arch"`

	// Run contains information on how to run the compiler.
	Run *service.RunInfo `toml:"run,omitempty"`

	// MOpt contains information on the 'mopt' (compiler architecture tuning) levels to select for the compiler.
	MOpt *optlevel.Selection `toml:"mopt,optempty"`

	// Opt contains information on the optimisation levels to select for the compiler.
	Opt *optlevel.Selection `toml:"opt,omitempty"`
}

// Compiler collects all test-relevant information about a compiler.
type Compiler struct {
	// SelectedMOpt refers to an architecture tuning level chosen using the compiler's configured march selection.
	SelectedMOpt string `toml:"selected_mopt,optempty"`
	// SelectedOpt refers to an optimisation level chosen using the compiler's configured optimisation selection.
	SelectedOpt *optlevel.Named `toml:"selected_opt,omitempty"`

	Config
}

// String outputs a human-readable but machine-separable summary of this compiler.
func (c Compiler) String() string {
	s, err := c.stringErr()
	if err != nil {
		return fmt.Sprintf("error: %s", err)
	}
	return s
}

func (c Compiler) stringErr() (string, error) {
	var sb strings.Builder
	if _, err := fmt.Fprintf(&sb, "%s@%s", c.Style, c.Arch); err != nil {
		return "", err
	}
	if c.Run != nil {
		if _, err := fmt.Fprintf(&sb, " (%s)", c.Run); err != nil {
			return "", err
		}
	}
	if c.SelectedOpt != nil {
		if _, err := fmt.Fprintf(&sb, " opt %q", c.SelectedOpt.Name); err != nil {
			return "", err
		}
	}
	if !ystring.IsBlank(c.SelectedMOpt) {
		if _, err := fmt.Fprintf(&sb, " march %q", c.SelectedMOpt); err != nil {
			return "", err
		}
	}

	return sb.String(), nil
}
