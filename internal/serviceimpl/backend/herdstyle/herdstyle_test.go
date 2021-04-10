// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package herdstyle_test

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/c4-project/c4t/internal/helper/testhelp"

	"github.com/c4-project/c4t/internal/model/service/backend"

	"github.com/c4-project/c4t/internal/serviceimpl/backend/herdstyle/litmus"

	"github.com/c4-project/c4t/internal/serviceimpl/backend/herdstyle/herd"

	"github.com/c4-project/c4t/internal/serviceimpl/backend/herdstyle"
	"github.com/c4-project/c4t/internal/subject/obs"
	"github.com/stretchr/testify/assert"
)

// TestBackend_ParseObs tests Herdtools observation parsing on various sample files.
func TestBackend_ParseObs(t *testing.T) {
	t.Parallel()

	impls := map[string]herdstyle.BackendImpl{
		"herd":   herd.Herd{},
		"litmus": litmus.Litmus{},
	}

	for name, i := range impls {
		name, i := name, i
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			b := herdstyle.Class{Impl: i}.Instantiate(backend.Spec{})

			indir := filepath.Join("testdata", name, "in")
			testhelp.TestFilesOfExt(t, indir, ".txt", func(t *testing.T, fname, path string) {
				t.Helper()
				t.Parallel()

				file, err := os.Open(path)
				if !assert.NoErrorf(t, err, "opening test file %q", fname) {
					return
				}
				var o obs.Obs
				err = b.ParseObs(context.Background(), file, &o)
				_ = file.Close()
				if !assert.NoErrorf(t, err, "parsing test file %q", fname) {
					return
				}
				inJson, ok := obsToJsonString(t, &o)
				if !ok {
					return
				}
				outname := filepath.Join("testdata", name, "out", fname+".json")
				outJson, err := os.ReadFile(outname)
				if assert.NoErrorf(t, err, "opening expected output file %q", outname) {
					assert.JSONEq(t, string(outJson), string(inJson), "JSON for observations didn't match")
				}
			})
		})
	}
}

func obsToJsonString(t *testing.T, o *obs.Obs) ([]byte, bool) {
	t.Helper()

	var b bytes.Buffer
	e := json.NewEncoder(&b)
	err := e.Encode(o)
	ok := assert.NoError(t, err, "couldn't encode observation to JSON")
	return b.Bytes(), ok
}
