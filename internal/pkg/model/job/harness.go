// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package job

import (
	"io/ioutil"

	"github.com/MattWindsor91/act-tester/internal/pkg/model/service"

	"github.com/MattWindsor91/act-tester/internal/pkg/model/id"
)

// Harness is a specification of how to make a test harness.
type Harness struct {
	// Backend is the backend to use to make this harness.
	Backend *service.Backend

	// Arch is the ID of the architecture for which a harness should be prepared.
	Arch id.ID

	// InFile is the path to the input litmus test file.
	InFile string

	// OutDir is the path to the output harness directory.
	OutDir string
}

// OutFiles reads s.OutDir as a directory and returns its contents as qualified paths.
// This is useful for using a harness job to feed a compiler job.
func (s Harness) OutFiles() ([]string, error) {
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
