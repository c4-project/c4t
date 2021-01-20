// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

// Package compiler contains types for compilers, which are a particular type of service.
package compiler

import (
	"fmt"
	"strings"
	"time"

	"github.com/1set/gut/ystring"

	"github.com/c4-project/c4t/internal/model/id"
	"github.com/c4-project/c4t/internal/model/service"
	"github.com/c4-project/c4t/internal/model/service/compiler/optlevel"
)

// Compiler collects the part of a compiler's specification that comes from the c4t configuration.
type Compiler struct {
	// Disabled specifies whether this compiler has been disabled.
	Disabled bool `toml:"disabled,omitempty" json:"disabled,omitempty"`

	// Style is the declared style of the compile.
	Style id.ID `toml:"style" json:"style"`

	// Arch is the architecture (or 'emits') ID for the compiler.
	Arch id.ID `toml:"arch" json:"arch"`

	// Run contains information on how to run the compiler.
	Run *service.RunInfo `toml:"run,omitempty" json:"run,omitempty"`

	// MOpt contains information on the 'mopt' (compiler architecture tuning) levels to select for the compiler.
	MOpt *optlevel.Selection `toml:"mopt,optempty" json:"mopt,omitempty"`

	// Opt contains information on the optimisation levels to select for the compiler.
	Opt *optlevel.Selection `toml:"opt,omitempty" json:"opt,omitempty"`
}

// Configuration collects all test-relevant information about a compiler.
type Configuration struct {
	// SelectedMOpt refers to an architecture tuning level chosen using the compiler's configured march selection.
	SelectedMOpt string `toml:"selected_mopt,optempty" json:"selected_mopt,omitempty"`
	// SelectedOpt refers to an optimisation level chosen using the compiler's configured optimisation selection.
	SelectedOpt *optlevel.Named `toml:"selected_opt,omitempty" json:"selected_opt,omitempty"`
	// ConfigTime captures the time at which this compiler configuration was generated.
	//
	// An example of when this may be useful is when using a compiler with run-time mutations enabled; we can use the
	// configuration time as a seed (by interpolating it out into the arguments or environment variables) to choose
	// mutations.
	ConfigTime time.Time `toml:"config_time,omitempty" json:"config_time,omitempty"`

	Compiler
}

// SelectedOptName returns the name of the selected optimisation level, or the empty string if there isn't one.
func (c Configuration) SelectedOptName() string {
	if c.SelectedOpt == nil {
		return ""
	}
	return c.SelectedOpt.Name
}

// String outputs a human-readable but machine-separable summary of this compiler configuration.
func (c Configuration) String() string {
	s, err := c.stringErr()
	if err != nil {
		return fmt.Sprintf("error: %s", err)
	}
	return s
}

func (c Configuration) stringErr() (string, error) {
	var sb strings.Builder
	if _, err := fmt.Fprintf(&sb, "%s@%s", c.Style, c.Arch); err != nil {
		return "", err
	}
	if c.Run != nil {
		if _, err := fmt.Fprintf(&sb, " (%s)", c.Run); err != nil {
			return "", err
		}
	}
	oname := c.SelectedOptName()
	if !ystring.IsBlank(oname) {
		if _, err := fmt.Fprintf(&sb, " opt %q", oname); err != nil {
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
