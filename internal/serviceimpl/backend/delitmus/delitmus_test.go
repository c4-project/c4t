// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package delitmus_test

import (
	"context"
	"path"
	"path/filepath"
	"testing"

	"github.com/c4-project/c4t/internal/model/id"

	"github.com/c4-project/c4t/internal/model/service"
	"github.com/c4-project/c4t/internal/model/service/mocks"

	"github.com/c4-project/c4t/internal/model/service/backend"

	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/c4-project/c4t/internal/act"
	"github.com/c4-project/c4t/internal/model/litmus"
	"github.com/c4-project/c4t/internal/serviceimpl/backend/delitmus"
)

// TestDelitmus_Lift tests the happy path of Delitmus.Lift.
func TestDelitmus_Lift(t *testing.T) {
	cr := new(mocks.Runner)
	cr.Test(t)

	j := backend.LiftJob{
		In: backend.LiftLitmusInput(
			litmus.NewOrPanic(path.Join("in", "foo.litmus"), litmus.WithArch(id.ArchC)),
		),
		Out: backend.LiftOutput{Dir: "out", Target: backend.ToDefault},
	}

	cr.On("Run", mock.Anything, service.RunInfo{
		Cmd: act.BinActC,
		Args: []string{
			"delitmus",
			"-aux-output", filepath.Join("out", "aux.json"),
			"-output", filepath.Join("out", "delitmus.c"),
			j.In.Litmus.Filepath(),
		},
	}).Return(nil).Once()

	dl := delitmus.Delitmus{BaseRunner: act.Runner{}}
	recipe, err := dl.Lift(context.Background(), j, cr)
	require.NoError(t, err, "lifting with mock delitmus run")

	assert.Equal(t, j.Out.Dir, recipe.Dir, "recipe should output to job output directory")
	assert.Nil(t, dl.BaseRunner.Base, "should not have changed base of original runner")

	cr.AssertExpectations(t)
}
