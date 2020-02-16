package model

import (
	"io/ioutil"
	"path"
)

// HarnessSpec is a specification of how to make a test harness.
type HarnessSpec struct {
	// Backend is the fully-qualified identifier of the backend to use to make this harness.
	Backend ID

	// Arch is the CompilerID of the architecture for which a harness should be prepared.
	Arch ID

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

// Harness represents information about a lifted test harness.
type Harness struct {
	// Dir is the root directory of the harness.
	Dir string `toml:"dir"`

	// Files is a list of files in the harness.
	Files []string `toml:"files"`
}

// Paths retrieves the joined dir/file paths for each file in the harness.
func (h Harness) Paths() []string {
	paths := make([]string, len(h.Files))
	for i, f := range h.Files {
		paths[i] = path.Join(h.Dir, f)
	}
	return paths
}
