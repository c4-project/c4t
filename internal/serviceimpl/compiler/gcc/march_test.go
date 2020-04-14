// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package gcc_test

import (
	"testing"

	"github.com/MattWindsor91/act-tester/internal/serviceimpl/compiler/gcc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/MattWindsor91/act-tester/internal/model/id"
)

// TestMOpts tests the mopt calculation for various platforms.
func TestMOpts(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		in  id.ID
		out []string
	}{
		"skylake":      {in: id.ArchX86Skylake, out: []string{"", "arch=native", "arch=x86_64", "arch=skylake"}},
		"arm7":         {in: id.ArchArm7, out: []string{"arch=armv7-a"}},
		"arm8":         {in: id.ArchArm8, out: []string{"arch=armv8-a", "arch=armv7-a"}},
		"armcortexa72": {in: id.ArchArmCortexA72, out: []string{"cpu=cortex-a72", "arch=armv8-a", "arch=armv7-a"}},
		"ppc64le":      {in: id.ArchPPC64LE, out: []string{"", "cpu=native", "cpu=powerpc64le"}},
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
