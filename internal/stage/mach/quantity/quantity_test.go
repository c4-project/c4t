// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package quantity_test

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/MattWindsor91/act-tester/internal/stage/mach/quantity"

	"github.com/stretchr/testify/assert"

	"github.com/MattWindsor91/act-tester/internal/stage/mach/timeout"
)

// ExampleQuantitySet_Log is a testable example for Log.
func ExampleQuantitySet_Log() {
	qs := quantity.Set{
		Compiler: quantity.SingleSet{
			Timeout:  timeout.Timeout(1 * time.Minute),
			NWorkers: 2,
		},
		Runner: quantity.SingleSet{
			Timeout:  timeout.Timeout(2 * time.Minute),
			NWorkers: 1,
		},
	}

	l := log.New(os.Stdout, "", 0)
	qs.Log(l)

	// Output:
	// [Compiler]
	// running across 2 workers
	// timeout at 1m0s
	// [Runner]
	// running across 1 worker
	// timeout at 2m0s
}

// TestQuantitySet_Override tests Override against some cases.
func TestQuantitySet_Override(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		old, new, want quantity.Set
	}{
		"all-old": {
			new: quantity.Set{
				Compiler: quantity.SingleSet{
					Timeout:  timeout.Timeout(4 * time.Second),
					NWorkers: 53,
				},
				Runner: quantity.SingleSet{
					Timeout:  timeout.Timeout(98 * time.Minute),
					NWorkers: 27,
				},
			},
			old: quantity.Set{},
			want: quantity.Set{
				Compiler: quantity.SingleSet{
					Timeout:  timeout.Timeout(4 * time.Second),
					NWorkers: 53,
				},
				Runner: quantity.SingleSet{
					Timeout:  timeout.Timeout(98 * time.Minute),
					NWorkers: 27,
				},
			},
		},
		"all-new": {
			old: quantity.Set{},
			new: quantity.Set{
				Compiler: quantity.SingleSet{
					Timeout:  timeout.Timeout(42 * time.Second),
					NWorkers: 27,
				},
				Runner: quantity.SingleSet{
					Timeout:  timeout.Timeout(1 * time.Minute),
					NWorkers: 42,
				},
			},
			want: quantity.Set{
				Compiler: quantity.SingleSet{
					Timeout:  timeout.Timeout(42 * time.Second),
					NWorkers: 27,
				},
				Runner: quantity.SingleSet{
					Timeout:  timeout.Timeout(1 * time.Minute),
					NWorkers: 42,
				},
			},
		},
	}

	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := c.old
			got.Override(c.new)

			assert.Equal(t, c.want, got, "quantity set override mismatch")
		})
	}
}
