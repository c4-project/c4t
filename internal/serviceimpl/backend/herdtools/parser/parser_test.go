// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package parser_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/MattWindsor91/act-tester/internal/helper/testhelp"
	"github.com/MattWindsor91/act-tester/internal/serviceimpl/backend/herdtools/herd"
	"github.com/MattWindsor91/act-tester/internal/serviceimpl/backend/herdtools/litmus"
	"github.com/MattWindsor91/act-tester/internal/serviceimpl/backend/herdtools/parser"
	"github.com/MattWindsor91/act-tester/internal/subject/obs"
)

// TestParse_error tests Parse with various error cases.
func TestParse_error(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		impl parser.Impl
		err  error
	}{
		"empty":                      {impl: herd.Herd{}, err: parser.ErrInputEmpty},
		"herd-no-states":             {impl: herd.Herd{}, err: parser.ErrBadStateCount},
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
