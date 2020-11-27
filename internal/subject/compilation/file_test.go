// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package compilation_test

import (
	"path/filepath"
	"testing"

	"github.com/MattWindsor91/c4t/internal/subject/compilation"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCompileFileset_ReadCompilerLog_plain tests CompileFileset.ReadLog with logs present as test data on the filesystem.
func TestCompileFileset_ReadCompilerLog(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		path, want string
	}{
		"plain file": {
			path: "plain.log",
			want: "This is an example of a plaintext compiler log file.",
		},
		"tarballed file": {
			path: "tarball/tarred.log", // slash path intentional
			want: "This is an example of a tarballed compiler log file.",
		},
	}
	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			cr := compilation.CompileFileset{Log: c.path}
			bs, err := cr.ReadLog(filepath.Join("testdata", "read_log"))
			require.NoError(t, err, "reading plain compiler log")
			assert.Equal(t, c.want, string(bs), "incorrect compiler log contents")
		})
	}

}
