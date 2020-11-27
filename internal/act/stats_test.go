// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package act_test

import (
	"strings"
	"testing"

	"github.com/MattWindsor91/c4t/internal/model/litmus"

	"github.com/MattWindsor91/c4t/internal/act"

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
			Types: map[string]int{
				"cmpxchg": 0,
				"fence":   0,
				"fetch":   30,
				"load":    19,
				"store":   0,
				"xchg":    0,
			},
			MemOrders: map[string]int{
				"memory_order_relaxed": 15,
				"memory_order_consume": 4,
				"memory_order_acquire": 11,
				"memory_order_release": 10,
				"memory_order_acq_rel": 3,
				"memory_order_seq_cst": 6,
			},
		},
		AtomicStatements: litmus.AtomicStatset{
			Types: map[string]int{
				"cmpxchg": 0,
				"fence":   1,
				"fetch":   2,
				"load":    0,
				"store":   0,
				"xchg":    3,
			},
			MemOrders: map[string]int{
				"memory_order_relaxed": 8,
				"memory_order_consume": 0,
				"memory_order_acquire": 0,
				"memory_order_release": 7,
				"memory_order_acq_rel": 0,
				"memory_order_seq_cst": 0,
			},
		},
	}

	var got litmus.Statset
	err := act.ParseStats(&got, r)
	require.NoError(t, err)
	assert.Equal(t, got, want)
}
