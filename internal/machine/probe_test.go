// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package machine_test

import (
	"testing"

	"github.com/c4-project/c4t/internal/id"

	"github.com/c4-project/c4t/internal/machine"
	"github.com/c4-project/c4t/internal/machine/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestProbe tests the happy path of Probe.
func TestProbe(t *testing.T) {
	t.Parallel()

	var m mocks.Prober
	m.Test(t)

	ncores := 4
	arch := id.ArchAArch64

	m.On("NCores").Return(ncores, nil).Once()
	m.On("Arch").Return(arch, nil).Once()
	// TODO(@MattWindsor91): add other probing

	var got machine.Config
	require.NoError(t, got.Probe(&m), "probe shouldn't error")

	assert.Equal(t, ncores, got.Cores, "probed cores not set correctly")
	assert.Equal(t, arch, got.Arch, "probed arch not set correctly")

	m.AssertExpectations(t)
}
