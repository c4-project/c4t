// Copyright (c) 2021 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package compiler

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/c4-project/c4t/internal/mutation"

	"github.com/1set/gut/ystring"
	"github.com/c4-project/c4t/internal/model/service/compiler/optlevel"
)

const (
	// varConfigTime is the interpolation variable for config time (UNIX timestamp).
	varConfigTime = "config_time"
	// varMutant is the interpolation variable for mutant IDs.
	varMutant = "mutant"
)

// Instance represents a fully configured instance of a compiler.
type Instance struct {
	// SelectedMOpt refers to an architecture tuning level chosen using the compiler's configured march selection.
	SelectedMOpt string `json:"selected_mopt,omitempty"`
	// SelectedOpt refers to an optimisation level chosen using the compiler's configured optimisation selection.
	SelectedOpt *optlevel.Named `json:"selected_opt,omitempty"`
	// ConfigTime captures the time at which this compiler configuration was generated.
	//
	// An example of when this may be useful is when using a compiler with run-time mutations enabled; we can use the
	// configuration time as a seed (by interpolating it out into the arguments or environment variables) to choose
	// mutations.
	ConfigTime time.Time `json:"config_time,omitempty"`
	// Mutant captures any mutant ID attached to this compiler instance.
	Mutant mutation.Mutant `json:"mutant,omitempty"`
	Compiler
}

// SelectedOptName returns the name of the selected optimisation level, or the empty string if there isn't one.
func (c Instance) SelectedOptName() string {
	if c.SelectedOpt == nil {
		return ""
	}
	return c.SelectedOpt.Name
}

// String outputs a human-readable but machine-separable summary of this compiler configuration.
func (c Instance) String() string {
	s, err := c.stringErr()
	if err != nil {
		return fmt.Sprintf("error: %s", err)
	}
	return s
}

func (c Instance) stringErr() (string, error) {
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

// Interpolations gets a map of variable interpolations that should be used in any job constructed from this instance.
func (c Instance) Interpolations() map[string]string {
	return map[string]string{
		varConfigTime: c.unixTimeString(),
		varMutant:     strconv.FormatUint(c.Mutant, 10),
	}
}

func (c Instance) unixTimeString() string {
	return strconv.FormatInt(c.ConfigTime.Unix(), 10)
}
