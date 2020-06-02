// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package pathset

import (
	"fmt"
	"path/filepath"

	"github.com/MattWindsor91/act-tester/internal/model/plan"
)

// Scratch contains the pre-computed paths for a machine run.
type Scratch struct {
	// DirFuzz is the directory to which fuzzed subjects will be output.
	DirFuzz string
	// DirLift is the directory to which lifter outputs will be written.
	DirLift string
	// DirPlan is the directory to which plans will be written.
	DirPlan string
	// DirRun is the directory into which act-tester-mach output will go.
	DirRun string
}

// NewScratch creates a machine pathset rooted at root.
func NewScratch(root string) *Scratch {
	return &Scratch{
		DirFuzz: filepath.Join(root, segFuzz),
		DirLift: filepath.Join(root, segLift),
		DirPlan: filepath.Join(root, segPlan),
		DirRun:  filepath.Join(root, segRun),
	}
}

// PlanForStage gets the path to the plan file for stage stage.
// Note that neither Prepare nor this method create or otherwise access the plan file.
func (p *Scratch) PlanForStage(stage string) string {
	file := fmt.Sprintf("plan.%s%s", stage, plan.Ext)
	return filepath.Join(p.DirPlan, file)
}
