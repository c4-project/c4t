// Copyright (c) 2021 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package mutation_test

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/c4-project/c4t/internal/mutation"
	"github.com/stretchr/testify/require"
)

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
		assert.Equal(t, want, got, "mutant wrong at position", i)
		kch <- struct{}{}
	}

	got = <-mch
	assert.Equal(t, wants[0], got, "mutants didn't wrap around")

	cancel()
	close(kch)
	wg.Wait()
}
