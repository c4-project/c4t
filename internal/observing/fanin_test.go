// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package observing_test

import (
	"context"
	"sync"
	"testing"

	"github.com/MattWindsor91/act-tester/internal/helper/testhelp"

	"github.com/stretchr/testify/assert"

	"github.com/MattWindsor91/act-tester/internal/observing"
	"github.com/stretchr/testify/require"
)

// TestFanIn_Run_empty tests that trying to run a fan-in with no channels terminates without errors.
func TestFanIn_Run_empty(t *testing.T) {
	fi := observing.NewFanIn(func(int, interface{}) error {
		return nil
	}, 0)
	err := fi.Run(context.Background())
	require.NoError(t, err, "should terminate with no errors")
}

// TestFanIn_Run_noCancel tests that trying to run a fan-in with no cancelling terminates without errors.
func TestFanIn_Run_noCancel(t *testing.T) {
	var mp [10]int

	fi := observing.NewFanIn(func(i int, v interface{}) error {
		mp[i] = v.(int)
		return nil
	}, len(mp))

	for i := 0; i < len(mp); i++ {
		i := i
		ch := make(chan int)
		go func() {
			ch <- i + 1
			close(ch)
		}()
		fi.Add(ch)
	}
	err := fi.Run(context.Background())
	require.NoError(t, err, "should terminate with no errors")

	for i := 0; i < len(mp); i++ {
		assert.Equal(t, i+1, mp[i], "didn't receive this message")
	}
}

// TestFanIn_Run_instantCancel tests that trying to run a fan-in with instant cancelling works properly
func TestFanIn_Run_instantCancel(t *testing.T) {
	var mp [10]int

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	fi := observing.NewFanIn(func(i int, v interface{}) error {
		mp[i] = v.(int)
		return nil
	}, len(mp))

	var wg sync.WaitGroup
	wg.Add(len(mp))
	for i := 0; i < len(mp); i++ {
		i := i
		ch := make(chan int)
		go func() {
			ch <- i + 1
			close(ch)
			wg.Done()
		}()
		fi.Add(ch)
	}
	err := fi.Run(ctx)
	testhelp.ExpectErrorIs(t, err, context.Canceled, "should propagate cancel")

	wg.Wait()
	for i := 0; i < len(mp); i++ {
		assert.Equal(t, 0, mp[i], "should not have received this message")
	}

}
