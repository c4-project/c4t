// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package parser_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/MattWindsor91/c4t/internal/serviceimpl/backend/herdstyle/rmem"

	"github.com/stretchr/testify/require"

	"github.com/MattWindsor91/c4t/internal/helper/testhelp"
	"github.com/MattWindsor91/c4t/internal/serviceimpl/backend/herdstyle/herd"
	"github.com/MattWindsor91/c4t/internal/serviceimpl/backend/herdstyle/litmus"
	"github.com/MattWindsor91/c4t/internal/serviceimpl/backend/herdstyle/parser"
	"github.com/MattWindsor91/c4t/internal/subject/obs"
)

// TestParse_error tests Parse with various error cases.
func TestParse_error(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		impl parser.Impl
		err  error
	}{
		"empty":                      {impl: herd.Herd{}, err: parser.ErrInputEmpty},
		"herd-no-states":             {impl: herd.Herd{}, err: parser.ErrNoStates},
		"herd-no-states-end":         {impl: herd.Herd{}, err: parser.ErrNoStates},
		"herd-not-enough-states":     {impl: herd.Herd{}, err: parser.ErrBadStateLine},
		"herd-not-enough-states-end": {impl: herd.Herd{}, err: parser.ErrNotEnoughStates},
		"herd-too-many-states":       {impl: herd.Herd{}, err: parser.ErrBadSummary},
		"herd-no-summary-end":        {impl: herd.Herd{}, err: parser.ErrNoSummary},
		"litmus-no-test":             {impl: litmus.Litmus{}, err: parser.ErrNoTest},
	}

	for name, c := range cases {
		name, c := name, c
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			f, err := os.Open(filepath.Join("testdata", "error-in", name+".txt"))
			require.NoError(t, err, "missing test file")
			defer func() { _ = f.Close() }()

			o := new(obs.Obs)
			err = parser.Parse(c.impl, f, o)

			testhelp.ExpectErrorIs(t, err, c.err, "parsing bad input")
		})
	}
}

// TestParse_valid tests Parse with various valid, or should-be-valid, cases.
func TestParse_valid(t *testing.T) {
	t.Parallel()

	cases := map[string]parser.Impl{
		"herd-ok-small":         herd.Herd{},
		"rmem-ok-unsat-partial": rmem.Rmem{},
	}

	for name, c := range cases {
		name, c := name, c
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			f, err := os.Open(filepath.Join("testdata", "valid", name+".txt"))
			require.NoError(t, err, "missing input file")
			defer func() { _ = f.Close() }()

			out, err := ioutil.ReadFile(filepath.Join("testdata", "valid", name+".json"))
			require.NoError(t, err, "missing output file")

			o := new(obs.Obs)
			err = parser.Parse(c, f, o)
			require.NoError(t, err, "parse should not error")

			var b bytes.Buffer
			err = json.NewEncoder(&b).Encode(o)
			require.NoError(t, err, "observation should encode into buffer")

			require.JSONEq(t, string(out), b.String(), "observation JSON not as expected")
		})
	}
}
