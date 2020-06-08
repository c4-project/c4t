// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package job

import (
	"errors"
	"io/ioutil"

	"github.com/1set/gut/ystring"

	"github.com/MattWindsor91/act-tester/internal/model/litmus"

	"github.com/MattWindsor91/act-tester/internal/model/service"

	"github.com/MattWindsor91/act-tester/internal/model/id"
)

var (
	// ErrInLitmusBlank occurs when the input file of al ifter job is checked and found to be blank.
	ErrInLitmusBlank = errors.New("input litmus file path blank")
	// ErrOutDirBlank occurs when the output directory of a lifter job is checked and found to be blank.
	ErrOutDirBlank = errors.New("output directory path blank")
)

// Lifter is a specification of how to lift a test into a compilable recipe.
type Lifter struct {
	// Backend is the backend to use to perform the lifting.
	Backend *service.Backend

	// Arch is the ID of the architecture for which a recipe should be prepared.
	Arch id.ID

	// In is the input litmus test file and its associated data from fuzzing and/or planning.
	In litmus.Litmus

	// OutDir is the path to the output directory.
	OutDir string
}

// Check performs several in-flight checks on a lifter job.
func (l *Lifter) Check() error {
	if !l.In.HasPath() {
		return ErrInLitmusBlank
	}
	if ystring.IsBlank(l.OutDir) {
		return ErrOutDirBlank
	}
	return nil
}

// OutFiles reads s.OutDir as a directory and returns its contents as qualified paths.
// This is useful for using a recipe job to feed a compiler job.
func (l *Lifter) OutFiles() ([]string, error) {
	fs, err := ioutil.ReadDir(l.OutDir)
	if err != nil {
		return nil, err
	}

	ps := make([]string, len(fs))
	i := 0
	for _, f := range fs {
		if f.IsDir() {
			continue
		}
		ps[i] = f.Name()
		i++
	}
	return ps[:i], nil
}
