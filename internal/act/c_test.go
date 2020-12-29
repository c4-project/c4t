// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package act_test

import (
	"context"
	"testing"

	"github.com/c4-project/c4t/internal/model/service"
	"github.com/c4-project/c4t/internal/model/service/mocks"

	"github.com/stretchr/testify/mock"

	"github.com/c4-project/c4t/internal/act"
	"github.com/stretchr/testify/require"
)

// TestRunner_Delitmus tests the happy path of Delitmus using a mock runner.
func TestRunner_Delitmus(t *testing.T) {
	t.Parallel()

	m := new(mocks.Runner)
	m.Test(t)
	m.On("Run", mock.Anything, service.RunInfo{
		Cmd:  act.BinActC,
		Args: []string{"delitmus", "-aux-output", "aux.json", "-output", "c.json", "in.litmus"},
	}).Return(nil).Once()

	dj := act.DelitmusJob{
		InLitmus: "in.litmus",
		OutAux:   "aux.json",
		OutC:     "c.json",
	}

	a := act.Runner{Base: m}
	err := a.Delitmus(context.Background(), dj)
	require.NoError(t, err, "mocked delitmus shouldn't error")

	m.AssertExpectations(t)
}
