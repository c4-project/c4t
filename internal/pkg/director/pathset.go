// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package director

import (
	"path"
	"strings"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"

	"github.com/MattWindsor91/act-tester/internal/pkg/iohelp"
)

const (
	segFuzz    = "fuzz"
	segLift    = "lift"
	segPlan    = "plan"
	segRun     = "run"
	segSaved   = "saved"
	segScratch = "scratch"
)

// Pathset contains the pre-computed paths used by the director.
type Pathset struct {
	// DirSaved is the directory into which saved runs get copied.
	DirSaved string

	// DirScratch is the directory that the director uses for ephemeral run data.
	DirScratch string
}

// NewPathset constructs a new pathset from the directory root.
func NewPathset(root string) *Pathset {
	return &Pathset{
		DirSaved:   path.Join(root, segSaved),
		DirScratch: path.Join(root, segScratch),
	}
}

// Prepare prepares this pathset by making its directories.
func (p *Pathset) Prepare() error {
	return iohelp.Mkdirs(p.DirSaved, p.DirScratch)
}

// MachineScratch gets the scratch pathset for a machine with
func (p *Pathset) MachineScratch(mid model.ID) *MachinePathset {
	segs := append([]string{p.DirScratch}, mid.Tags()...)
	return NewMachinePathset(path.Join(segs...))
}

// MachinePathset contains the pre-computed paths for a machine run.
type MachinePathset struct {
	// DirFuzz is the directory to which fuzzed subjects will be output.
	DirFuzz string
	// DirLift is the directory to which lifted harnesses will be output.
	DirLift string
	// DirPlan is the directory to which plans will be written.
	DirPlan string
	// DirRun is the directory into which act-tester-mach output will go.
	DirRun string
}

func NewMachinePathset(root string) *MachinePathset {
	return &MachinePathset{
		DirFuzz: path.Join(root, segFuzz),
		DirLift: path.Join(root, segLift),
		DirPlan: path.Join(root, segPlan),
		DirRun:  path.Join(root, segRun),
	}
}

// Prepare prepares this pathset by making its directories.
func (p *MachinePathset) Prepare() error {
	return iohelp.Mkdirs(p.DirPlan, p.DirFuzz, p.DirLift, p.DirRun)
}

// PlanForStage gets the path to the plan file for stage stage.
// Note that neither Prepare nor this method create or otherwise access the plan file.
func (p *MachinePathset) PlanForStage(stage string) string {
	file := strings.Join([]string{"plan", stage, "toml"}, ".")
	return path.Join(p.DirPlan, file)
}
