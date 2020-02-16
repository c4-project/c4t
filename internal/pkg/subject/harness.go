package subject

import "path"

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
