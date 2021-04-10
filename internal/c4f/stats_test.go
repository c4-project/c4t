// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package c4f_test

import (
	"strings"
	"testing"

	"github.com/c4-project/c4t/internal/id"

	"github.com/c4-project/c4t/internal/model/litmus"

	"github.com/c4-project/c4t/internal/c4f"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStatset_Parse(t *testing.T) {
	r := strings.NewReader(`threads 3
returns 0
literals.bool 14
atomics.expression.cmpxchg 0
atomics.expression.fence 0
atomics.expression.fetch 30
atomics.expression.load 19
atomics.expression.store 0
atomics.expression.xchg 0
mem-orders.expression.memory_order_relaxed 15
mem-orders.expression.memory_order_consume 4
mem-orders.expression.memory_order_acquire 11
mem-orders.expression.memory_order_release 10
mem-orders.expression.memory_order_acq_rel 3
mem-orders.expression.memory_order_seq_cst 6
atomics.statement.cmpxchg 0
atomics.statement.fence 1
atomics.statement.fetch 2
atomics.statement.load 0
atomics.statement.store 0
atomics.statement.xchg 3
mem-orders.statement.memory_order_relaxed 8
mem-orders.statement.memory_order_consume 0
mem-orders.statement.memory_order_acquire 0
mem-orders.statement.memory_order_release 7
mem-orders.statement.memory_order_acq_rel 0
mem-orders.statement.memory_order_seq_cst 0`)

	want := litmus.Statset{
		Threads:      3,
		Returns:      0,
		LiteralBools: 14,
		AtomicExpressions: litmus.AtomicStatset{
			Types: map[id.ID]int{
				id.FromString("cmpxchg"): 0,
				id.FromString("fence"):   0,
				id.FromString("fetch"):   30,
				id.FromString("load"):    19,
				id.FromString("store"):   0,
				id.FromString("xchg"):    0,
			},
			MemOrders: map[id.ID]int{
				id.FromString("memory_order_relaxed"): 15,
				id.FromString("memory_order_consume"): 4,
				id.FromString("memory_order_acquire"): 11,
				id.FromString("memory_order_release"): 10,
				id.FromString("memory_order_acq_rel"): 3,
				id.FromString("memory_order_seq_cst"): 6,
			},
		},
		AtomicStatements: litmus.AtomicStatset{
			Types: map[id.ID]int{
				id.FromString("cmpxchg"): 0,
				id.FromString("fence"):   1,
				id.FromString("fetch"):   2,
				id.FromString("load"):    0,
				id.FromString("store"):   0,
				id.FromString("xchg"):    3,
			},
			MemOrders: map[id.ID]int{
				id.FromString("memory_order_relaxed"): 8,
				id.FromString("memory_order_consume"): 0,
				id.FromString("memory_order_acquire"): 0,
				id.FromString("memory_order_release"): 7,
				id.FromString("memory_order_acq_rel"): 0,
				id.FromString("memory_order_seq_cst"): 0,
			},
		},
	}

	var got litmus.Statset
	err := c4f.ParseStats(&got, r)
	require.NoError(t, err)
	assert.Equal(t, got, want)
}
