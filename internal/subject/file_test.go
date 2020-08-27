// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package subject_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/MattWindsor91/act-tester/internal/model/id"
	"github.com/MattWindsor91/act-tester/internal/model/litmus"
	"github.com/MattWindsor91/act-tester/internal/subject"
	"github.com/MattWindsor91/act-tester/internal/subject/compilation"
)

// TestSubject_ReadCompilerLog_plain tests Subject.ReadCompilerLog with logs present as test data on the filesystem.
// accessible on the filesystem.
func TestSubject_ReadCompilerLog(t *testing.T) {
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

			cid := id.FromString("clang.3")
			cr := compilation.CompileResult{Files: compilation.CompileFileset{Log: c.path}}
			s := subject.NewOrPanic(litmus.New("foo.litmus"), subject.WithCompile(cid, cr))

			bs, err := s.ReadCompilerLog("testdata", cid)
			require.NoError(t, err, "reading plain compiler log")
			assert.Equal(t, c.want, string(bs), "incorrect compiler log contents")
		})
	}

}
