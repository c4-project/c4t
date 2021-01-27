// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package c4f_test

import (
	"context"
	"io"
	"testing"

	"github.com/c4-project/c4t/internal/model/service"
	"github.com/c4-project/c4t/internal/model/service/mocks"

	"github.com/stretchr/testify/mock"

	"github.com/c4-project/c4t/internal/c4f"
	"github.com/stretchr/testify/require"
)

// TestRunner_Delitmus tests the happy path of Delitmus using a mock runner.
func TestRunner_Delitmus(t *testing.T) {
	t.Parallel()

	m := new(mocks.Runner)
	m.Test(t)
	m.On("Run", mock.Anything, service.RunInfo{
		Cmd:  c4f.BinC4fC,
		Args: []string{"delitmus", "-aux-output", "aux.json", "-output", "c.json", "in.litmus"},
	}).Return(nil).Once()

	dj := c4f.DelitmusJob{
		InLitmus: "in.litmus",
		OutAux:   "aux.json",
		OutC:     "c.json",
	}

	a := c4f.Runner{Base: m}
	err := a.Delitmus(context.Background(), dj)
	require.NoError(t, err, "mocked delitmus shouldn't error")

	m.AssertExpectations(t)
}

// TestRunner_CVersion tests the happy path of CVersion using a mock runner.
func TestRunner_CVersion(t *testing.T) {
	t.Parallel()

	var w io.Writer
	want := "test-version\n"

	m := new(mocks.Runner)
	m.Test(t)
	m.On("WithStdout", mock.Anything).Return(m).Run(func(args mock.Arguments) {
		w = args.Get(0).(io.Writer)
	}).Once()
	m.On("Run", mock.Anything, service.RunInfo{
		Cmd:  c4f.BinC4fC,
		Args: []string{"version", "-version"},
	}).Run(func(mock.Arguments) {
		_, _ = w.Write([]byte(want))
	}).Return(nil).Once()

	got, err := (&c4f.Runner{Base: m}).CVersion(context.Background())
	require.NoError(t, err, "mocked c-version shouldn't error")
	require.Equal(t, want, got, "c-version didn't match")

	m.AssertExpectations(t)
}
