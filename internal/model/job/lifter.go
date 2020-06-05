// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package job

import (
	"io/ioutil"

	"github.com/MattWindsor91/act-tester/internal/model/litmus"

	"github.com/MattWindsor91/act-tester/internal/model/service"

	"github.com/MattWindsor91/act-tester/internal/model/id"
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

// OutFiles reads s.OutDir as a directory and returns its contents as qualified paths.
// This is useful for using a recipe job to feed a compiler job.
func (s Lifter) OutFiles() ([]string, error) {
	fs, err := ioutil.ReadDir(s.OutDir)
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
