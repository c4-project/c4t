// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package analyse_test

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/c4-project/c4t/internal/app/analyse"
	"github.com/c4-project/c4t/internal/helper/testhelp"
	"github.com/c4-project/c4t/internal/plan"
	"github.com/c4-project/c4t/internal/stage/analyser"
	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

// TestApp_errorOnBadStatus is a top-level test for the analyse app that tests its ability to error on a bad status.
func TestApp_errorOnBadStatus(t *testing.T) {
	// The mock plan contains bad statuses!
	tpath := makeMockPlanFile(t)

	cases := map[string]struct {
		flags []string
		err   error
	}{
		"off": {flags: []string{}, err: nil},
		"on":  {flags: []string{"-" + analyse.FlagErrorOnBadStatus}, err: analyser.ErrBadStatus},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			// TODO(@MattWindsor91): can we parallelise this and the above?
			var outb, errb bytes.Buffer

			argv := append([]string{analyse.Name}, c.flags...)
			argv = append(argv, tpath)

			app := analyse.App(&outb, &errb)
			err := app.Run(argv)
			testhelp.ExpectErrorIs(t, err, c.err, "running analyser on bad status")

			assert.Empty(t, outb.Bytes(), "shouldn't have outputted anything without specific writer flags")
			assert.Empty(t, errb.Bytes(), "shouldn't have outputted anything without specific writer flags")
		})
	}
}

func makeMockPlanFile(t *testing.T) string {
	t.Helper()

	// TODO(@MattWindsor91): share this with different tests as we add them.
	// TODO(@MattWindsor91): exercise gzip too?
	tpath := filepath.Join(t.TempDir(), "plan.json")
	p := plan.Mock()
	err := p.WriteFile(tpath, plan.WriteNone)
	require.NoError(t, err, "couldn't write mock plan to temp file")
	return tpath
}
