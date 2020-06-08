// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package compile describes compilation jobs.
package compile

import (
	"github.com/MattWindsor91/act-tester/internal/model/service"
	"github.com/MattWindsor91/act-tester/internal/model/service/compiler"
)

// Compile represents a request to compile a list of files to an executable given a particular compiler.
type Compile struct {
	// Compiler describes the compiler to use for the compilation.
	Compiler *compiler.Compiler

	// In is the list of files to be sent to the compiler.
	In []string
	// Out is the file to be received from the compiler.
	Out string
}

// New is a convenience constructor for compiles.
func New(c *compiler.Compiler, out string, in ...string) *Compile {
	return &Compile{
		Compiler: c,
		In:       in,
		Out:      out,
	}
}

// CompilerRun gets the job's compiler run information if present; else, nil.
func (j *Compile) CompilerRun() *service.RunInfo {
	if j.Compiler == nil {
		return nil
	}
	return j.Compiler.Run
}

// SelectedOptName gets the name of this job's compiler's selected optimisation level, if present; else, "".
func (j *Compile) SelectedOptName() string {
	if j.Compiler == nil || j.Compiler.SelectedOpt == nil {
		return ""
	}
	return j.Compiler.SelectedOpt.Name
}

// SelectedMOptName gets the name of this job's compiler's selected machine optimisation profile, if present; else, "".
func (j *Compile) SelectedMOptName() string {
	if j.Compiler == nil {
		return ""
	}
	return j.Compiler.SelectedMOpt
}
