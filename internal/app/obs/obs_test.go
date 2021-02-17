// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package obs_test

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/c4-project/c4t/internal/helper/testhelp"

	"github.com/c4-project/c4t/internal/app/obs"
	"github.com/stretchr/testify/require"
)

// TestApp_valid tests the happy path of the obs app on input test data.
func TestApp_valid(t *testing.T) {
	t.Parallel()

	dir := filepath.Join("testdata", "valid")
	testhelp.TestFilesOfExt(t, dir, ".json", func(t *testing.T, name, path string) {
		t.Parallel()

		for _, flags := range []string{"i", "p", "ip"} {
			flags := flags
			name := strings.Join([]string{name, flags}, "-")
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				fpath := filepath.Join(dir, name+".txt")
				require.FileExists(t, fpath, "want-file should exist")

				want, err := os.ReadFile(fpath)
				require.NoError(t, err, "want-file should be readable")

				args := []string{obs.Name, "-" + flags, path}

				var buf bytes.Buffer
				require.NoError(t, obs.App(&buf, io.Discard).Run(args), "obs app should run OK")
				require.Equal(t, string(want), buf.String(), "mismatch between output")
			})
		}
	})
}
