// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package herdtools_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/MattWindsor91/act-tester/internal/serviceimpl/backend/herdtools/litmus"

	"github.com/MattWindsor91/act-tester/internal/serviceimpl/backend/herdtools/herd"

	"github.com/MattWindsor91/act-tester/internal/helper/iohelp"
	"github.com/MattWindsor91/act-tester/internal/model/obs"
	"github.com/MattWindsor91/act-tester/internal/serviceimpl/backend/herdtools"
	"github.com/stretchr/testify/assert"
)

// TestBackend_ParseObs tests Herdtools observation parsing on various sample files.
func TestBackend_ParseObs(t *testing.T) {
	t.Parallel()

	impls := map[string]herdtools.BackendImpl{
		"herd":   herd.Herd{},
		"litmus": litmus.Litmus{},
	}

	for name, i := range impls {
		name, i := name, i
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			b := herdtools.Backend{Impl: i}

			indir := filepath.Join("testdata", name, "in")
			fs, err := ioutil.ReadDir(indir)
			if !assert.NoError(t, err, "reading test data directory", indir) {
				return
			}
			for _, f := range fs {
				fname := f.Name()
				ename := iohelp.ExtlessFile(fname)
				t.Run(ename, func(t *testing.T) {
					t.Parallel()

					file, err := os.Open(filepath.Join(indir, fname))
					if !assert.NoError(t, err, "opening test file", fname) {
						return
					}
					var o obs.Obs
					err = b.ParseObs(context.Background(), nil, file, &o)
					_ = file.Close()
					if !assert.NoError(t, err, "parsing test file", fname) {
						return
					}
					inJson, ok := obsToJsonString(t, &o)
					if !ok {
						return
					}
					outname := filepath.Join("testdata", name, "out", ename+".json")
					outJson, err := ioutil.ReadFile(outname)
					if assert.NoError(t, err, "opening expected output file", outname) {
						assert.JSONEq(t, string(outJson), string(inJson), "JSON for observations didn't match")
					}
				})
			}
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
