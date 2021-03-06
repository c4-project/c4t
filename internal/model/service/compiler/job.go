// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package compiler

import (
	"github.com/c4-project/c4t/internal/model/service"
)

// Job represents a request to compile a list of files to a particular target given a particular compiler.
type Job struct {
	// Compiler describes the compiler to use for the compilation.
	Compiler *Instance

	// In is the list of files to be sent to the compiler.
	In []string
	// Out is the file to be received from the compiler.
	Out string

	// Kind is the kind of file being produced by this compile.
	Kind Target
}

// NewJob is a convenience constructor for compiles.
func NewJob(k Target, c *Instance, out string, in ...string) *Job {
	return &Job{
		Kind:     k,
		Compiler: c,
		In:       in,
		Out:      out,
	}
}

// CompilerRun gets the job's compiler run information if present; else, nil.
func (j *Job) CompilerRun() *service.RunInfo {
	if j.Compiler == nil {
		return nil
	}
	return j.Compiler.Run
}

// SelectedOptName gets the name of this job's compiler's selected optimisation level, if present; else, "".
func (j *Job) SelectedOptName() string {
	if j.Compiler == nil || j.Compiler.SelectedOpt == nil {
		return ""
	}
	return j.Compiler.SelectedOpt.Name
}

// SelectedMOptName gets the name of this job's compiler's selected machine optimisation profile, if present; else, "".
func (j *Job) SelectedMOptName() string {
	if j.Compiler == nil {
		return ""
	}
	return j.Compiler.SelectedMOpt
}
