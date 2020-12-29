// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package quantity_test

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/c4-project/c4t/internal/quantity"

	"github.com/stretchr/testify/assert"
)

// ExampleMachNodeSet_Log is a testable example for MachNodeSet.Log.
func ExampleMachNodeSet_Log() {
	qs := quantity.MachNodeSet{
		Compiler: quantity.BatchSet{
			Timeout:  quantity.Timeout(1 * time.Minute),
			NWorkers: 2,
		},
		Runner: quantity.BatchSet{
			Timeout:  quantity.Timeout(2 * time.Minute),
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

// TestMachNodeSet_Override tests MachNodeSet.Override against some cases.
func TestMachNodeSet_Override(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		old, new, want quantity.MachNodeSet
	}{
		"all-old": {
			new: quantity.MachNodeSet{
				Compiler: quantity.BatchSet{
					Timeout:  quantity.Timeout(4 * time.Second),
					NWorkers: 53,
				},
				Runner: quantity.BatchSet{
					Timeout:  quantity.Timeout(98 * time.Minute),
					NWorkers: 27,
				},
			},
			old: quantity.MachNodeSet{},
			want: quantity.MachNodeSet{
				Compiler: quantity.BatchSet{
					Timeout:  quantity.Timeout(4 * time.Second),
					NWorkers: 53,
				},
				Runner: quantity.BatchSet{
					Timeout:  quantity.Timeout(98 * time.Minute),
					NWorkers: 27,
				},
			},
		},
		"all-new": {
			old: quantity.MachNodeSet{},
			new: quantity.MachNodeSet{
				Compiler: quantity.BatchSet{
					Timeout:  quantity.Timeout(42 * time.Second),
					NWorkers: 27,
				},
				Runner: quantity.BatchSet{
					Timeout:  quantity.Timeout(1 * time.Minute),
					NWorkers: 42,
				},
			},
			want: quantity.MachNodeSet{
				Compiler: quantity.BatchSet{
					Timeout:  quantity.Timeout(42 * time.Second),
					NWorkers: 27,
				},
				Runner: quantity.BatchSet{
					Timeout:  quantity.Timeout(1 * time.Minute),
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
