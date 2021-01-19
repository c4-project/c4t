// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package gcc_test

import (
	"testing"

	"github.com/c4-project/c4t/internal/serviceimpl/compiler/gcc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/c4-project/c4t/internal/model/id"
)

// TestDefaultMOpts tests the mopt calculation for various platforms.
func TestDefaultMOpts(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		in  id.ID
		out []string
	}{
		"skylake":      {in: id.ArchX86Skylake, out: []string{"", "arch=native", "arch=x86-64", "arch=broadwell", "arch=skylake"}},
		"broadwell":    {in: id.ArchX86Broadwell, out: []string{"", "arch=native", "arch=x86-64", "arch=broadwell"}},
		"aarch648.1":   {in: id.ArchAArch6481, out: []string{"", "cpu=native", "cpu=generic", "arch=armv8-a", "arch=armv8.1-a"}},
		"aarch648":     {in: id.ArchAArch648, out: []string{"", "cpu=native", "cpu=generic", "arch=armv8-a"}},
		"aarch64":      {in: id.ArchAArch64, out: []string{"", "cpu=native", "cpu=generic"}},
		"armcortexa72": {in: id.ArchArmCortexA72, out: []string{"cpu=cortex-a72", "arch=armv8-a", "arch=armv7-a"}},
		"arm8":         {in: id.ArchArm8, out: []string{"arch=armv8-a", "arch=armv7-a"}},
		"arm7":         {in: id.ArchArm7, out: []string{"arch=armv7-a"}},
		"ppc64le":      {in: id.ArchPPC64LE, out: []string{"", "cpu=native", "cpu=powerpc64le"}},
		"power9":       {in: id.ArchPPCPOWER9, out: []string{"", "cpu=native", "cpu=powerpc64le", "cpu=power7", "cpu=power8", "cpu=power9"}},
		"power8":       {in: id.ArchPPCPOWER8, out: []string{"", "cpu=native", "cpu=powerpc64le", "cpu=power7", "cpu=power8"}},
		"power7":       {in: id.ArchPPCPOWER7, out: []string{"", "cpu=native", "cpu=powerpc64le", "cpu=power7"}},
	}

	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			want := c.out
			gset, err := gcc.MOpts(c.in)
			require.NoErrorf(t, err, "MOpts(%v)", c.in)
			got := gset.Slice()
			assert.ElementsMatchf(t, want, got, "MOpts(%v)", c.in)
		})
	}
}
