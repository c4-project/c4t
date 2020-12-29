// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package obs_test

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/1set/gut/yos"
	"github.com/c4-project/c4t/internal/app/obs"
	"github.com/stretchr/testify/require"
)

// TestApp_valid tests the happy path of the obs app on input test data.
func TestApp_valid(t *testing.T) {
	t.Parallel()

	fs, err := yos.ListMatch(filepath.Join("testdata", "valid"), yos.ListIncludeFile, "*.json")
	require.NoError(t, err, "testdata listing shouldn't fail")

	for _, f := range fs {
		path := f.Path
		pprefix := strings.TrimSuffix(path, ".json")

		for _, flags := range []string{"i", "p", "ip"} {
			flags := flags
			name := strings.Join([]string{pprefix, flags}, "-")
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				fpath := name + ".txt"
				require.FileExists(t, fpath, "want-file should exist")

				want, err := ioutil.ReadFile(fpath)
				require.NoError(t, err, "want-file should be readable")

				args := []string{obs.Name, "-" + flags, path}

				var buf bytes.Buffer
				require.NoError(t, obs.App(&buf, ioutil.Discard).Run(args), "obs app should run OK")
				require.Equal(t, string(want), buf.String(), "mismatch between output")
			})
		}
	}
}
