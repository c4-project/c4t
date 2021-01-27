// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package quantity_test

import (
	"context"
	"testing"
	"time"

	"github.com/c4-project/c4t/internal/quantity"
	"github.com/stretchr/testify/assert"
)

// TestTimeout_OnContext tests to see if Timeout.OnContext seems to be producing the right types of context.
func TestTimeout_OnContext(t *testing.T) {
	zero := quantity.Timeout(0)
	c, _ := zero.OnContext(context.Background())
	_, hasDl := c.Deadline()
	assert.False(t, hasDl, "timeout 0 should not have a deadline")

	nonzero := quantity.Timeout(1 * time.Second)
	c, _ = nonzero.OnContext(context.Background())
	_, hasDl = c.Deadline()
	assert.True(t, hasDl, "timeout 1s should have a deadline")
}
