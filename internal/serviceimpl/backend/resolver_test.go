// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package backend_test

import (
	"testing"

	"github.com/c4-project/c4t/internal/serviceimpl/backend"

	"github.com/stretchr/testify/require"

	backend2 "github.com/c4-project/c4t/internal/model/service/backend"

	"github.com/c4-project/c4t/internal/id"
	"github.com/stretchr/testify/assert"
)

// TestResolver_Resolve_capabilities tests that the standard resolver provides the correct capabilities for the known backends.
func TestResolver_Resolve_capabilities(t *testing.T) {
	t.Parallel()

	cases := map[string]backend2.Capability{
		"delitmus":         backend2.CanLiftLitmus | backend2.CanProduceObj,
		"herdtools.herd":   backend2.CanLiftLitmus | backend2.CanRunStandalone,
		"herdtools.litmus": backend2.CanLiftLitmus | backend2.CanRunStandalone | backend2.CanProduceExe,
		"rmem":             backend2.CanLiftLitmus | backend2.CanRunStandalone,
	}
	for name, c := range cases {
		name, c := name, c
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			r, err := backend.Resolve.Resolve(id.FromString(name))
			require.NoError(t, err, "resolution should pass")

			assert.Equal(t, c, r.Metadata().Capabilities)
		})
	}
}
