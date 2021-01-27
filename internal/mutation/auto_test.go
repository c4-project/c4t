// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package mutation_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/c4-project/c4t/internal/mutation/mocks"
	"github.com/c4-project/c4t/internal/quantity"

	"github.com/stretchr/testify/assert"

	"github.com/c4-project/c4t/internal/mutation"
	"github.com/stretchr/testify/require"
)

// TestNewAutomator_error makes sure NewAutomator fails if we try to automate an inactive mutation test configuration.
func TestNewAutomator_error(t *testing.T) {
	_, err := mutation.NewAutomator(mutation.AutoConfig{})
	require.ErrorIs(t, err, mutation.ErrNotActive, "shouldn't be able to construct this")
}

// TestAutomator_Run_killOnly is a test run of Automator.Run when kill signalling is active but time-slicing isn't.
func TestAutomator_Run_killOnly(t *testing.T) {
	t.Parallel()

	cfg := mutation.AutoConfig{
		Ranges:       []mutation.Range{{Start: 1, End: 3}, {Start: 5, End: 6}},
		ChangeKilled: true,
	}

	a, err := mutation.NewAutomator(cfg)
	require.NoError(t, err, "automator should be constructible")

	kch := a.KillCh()
	require.NotNil(t, kch, "kill channel should be non-nil")

	mch := a.MutantCh()
	require.NotNil(t, mch, "mutant channel should be non-nil")

	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup
	wg.Add(1)
	go func() { a.Run(ctx); wg.Done() }()

	wants := cfg.Mutants()

	var got mutation.Mutant
	for i, want := range wants {
		got = <-mch
		assert.Equalf(t, want, got, "mutant wrong at position %d", i)
		kch <- struct{}{}
	}

	got = <-mch
	assert.Equal(t, wants[0], got, "mutants didn't wrap around")

	cancel()
	// Test draining of spurious kills.
	kch <- struct{}{}
	close(kch)
	wg.Wait()
}

// TestAutomator_Run_timeOnly is a test run of Automator.Run when kill signalling is active but time-slicing isn't.
func TestAutomator_Run_timeOnly(t *testing.T) {
	t.Parallel()

	ch := make(chan time.Time)
	ticker := new(mocks.Ticker)
	ticker.Test(t)
	// The ticker should be stopped before the end of the run.
	ticker.On("Stop").Return().Once()

	cfg := mutation.AutoConfig{
		Ranges:      []mutation.Range{{Start: 1, End: 3}, {Start: 5, End: 6}},
		ChangeAfter: quantity.Timeout(10 * time.Minute),
	}

	a, err := mutation.NewAutomator(cfg)
	require.NoError(t, err, "automator should be constructible")
	a.TickerF = func(duration time.Duration) (<-chan time.Time, mutation.Ticker) {
		return ch, ticker
	}

	kch := a.KillCh()
	require.Nil(t, kch, "kill channel should be nil")

	mch := a.MutantCh()
	require.NotNil(t, mch, "mutant channel should be non-nil")

	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup
	wg.Add(1)
	go func() { a.Run(ctx); wg.Done() }()

	wants := cfg.Mutants()

	var got mutation.Mutant
	for i, want := range wants {
		got = <-mch
		assert.Equal(t, want, got, "mutant wrong at position", i)
		ch <- time.Now()
	}

	got = <-mch
	assert.Equal(t, wants[0], got, "mutants didn't wrap around")

	cancel()
	wg.Wait()

	ticker.AssertExpectations(t)
}

// TestAutoPool_Mutant checks whether AutoPool.Mutant returns what we expect to be the right mutants.
func TestAutoPool_Mutant(t *testing.T) {
	t.Parallel()
	asrt := assert.New(t)

	orig := []mutation.Mutant{1, 2, 4, 5, 10, 11, 12}
	var a mutation.AutoPool

	a.Init(orig)

	for i, m := range orig {
		asrt.Equalf(m, a.Mutant(), "mutant at position %d unexpected", i)
		if i%2 == 0 {
			a.Kill()
		} else {
			a.Advance()
		}
	}

	// This should leave the odd mutants.
	for i := 1; i < len(orig); i += 2 {
		asrt.Equalf(orig[i], a.Mutant(), "mutant at position %d unexpected (pass 2)", i)
		a.Kill()
	}

	// We should now have killed all the mutants, testing for wraparound now.
	for i, m := range orig {
		asrt.Equalf(m, a.Mutant(), "mutant at position %d unexpected (pass 3)", i)
		if 0 < i {
			a.Kill()
		} else {
			a.Advance()
		}
	}

	// We left exactly one mutant live.
	for i := 0; i < 100; i++ {
		asrt.Equalf(orig[0], a.Mutant(), "mutant unexpected (pass 4)")
		a.Advance()
	}
}
