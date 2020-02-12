package model

import (
	"io/ioutil"
)

// HarnessSpec is a specification of how to make a test harness.
type HarnessSpec struct {
	// Backend is the fully-qualified identifier of the backend to use to make this harness.
	Backend Id

	// Arch is the ID of the architecture for which a harness should be prepared.
	Arch Id

	// InFile is the path to the input litmus test file.
	InFile string

	// OutDir is the path to the output harness directory.
	OutDir string
}

// OutFiles reads s.OutDir as a directory and returns its contents as qualified paths.
func (s HarnessSpec) OutFiles() ([]string, error) {
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
