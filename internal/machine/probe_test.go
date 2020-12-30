// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package machine_test

import (
	"github.com/c4-project/c4t/internal/machine"
	"github.com/c4-project/c4t/internal/machine/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// TestProbe tests the happy path of Probe.
func TestProbe(t *testing.T) {
	t.Parallel()

	var m mocks.Prober
	m.Test(t)

	ncores := 4

	m.On("NCores").Return(ncores, nil).Once()
	// TODO(@MattWindsor91): add other probing

	var got machine.Machine
	require.NoError(t, machine.Probe(&got, &m), "probe shouldn't error")

	assert.Equal(t, ncores, got.Cores, "probed cores not set correctly")
}
