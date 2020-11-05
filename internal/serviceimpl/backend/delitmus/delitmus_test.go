// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package delitmus_test

import (
	"bytes"
	"context"
	"io"
	"path"
	"path/filepath"
	"testing"

	"github.com/MattWindsor91/act-tester/internal/model/service"
	"github.com/MattWindsor91/act-tester/internal/model/service/mocks"

	"github.com/MattWindsor91/act-tester/internal/model/service/backend"

	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/MattWindsor91/act-tester/internal/act"
	"github.com/MattWindsor91/act-tester/internal/model/litmus"
	"github.com/MattWindsor91/act-tester/internal/serviceimpl/backend/delitmus"
)

// TestDelitmus_Lift tests the happy path of Delitmus.Lift.
func TestDelitmus_Lift(t *testing.T) {
	cr := new(mocks.Runner)
	cr.Test(t)

	j := backend.LiftJob{
		In:     *litmus.New(path.Join("in", "foo.litmus")),
		OutDir: "out",
	}

	// We don't actually use this, but it helps us check the runner construction.
	errw := new(bytes.Buffer)

	cr.On("Run", mock.Anything, service.RunInfo{
		Cmd: act.BinActC,
		Args: []string{
			"delitmus",
			"-aux-output", filepath.Join("out", "aux.json"),
			"-output", filepath.Join("out", "delitmus.c"),
			j.In.Filepath(),
		},
	}).Return(nil).Once()

	dl := delitmus.Delitmus{BaseRunner: act.Runner{RunnerFactory: func(io.Writer, io.Writer) service.Runner { return cr }}}
	recipe, err := dl.Lift(context.Background(), j, errw)
	require.NoError(t, err, "lifting with mock delitmus run")

	assert.Equal(t, j.OutDir, recipe.Dir, "recipe should output to job output directory")
	assert.Nil(t, dl.BaseRunner.Stderr, "should not have overwritten the base runner")

	cr.AssertExpectations(t)
}
