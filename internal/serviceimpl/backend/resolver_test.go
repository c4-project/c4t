// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package backend_test

import (
	"testing"

	"github.com/c4-project/c4t/internal/serviceimpl/backend"

	"github.com/stretchr/testify/require"

	backend2 "github.com/c4-project/c4t/internal/model/service/backend"

	"github.com/c4-project/c4t/internal/model/id"
	"github.com/stretchr/testify/assert"
)

// TestResolver_Get_capabilities tests that the standard resolver provides the correct capabilities for the known backends.
func TestResolver_Get_capabilities(t *testing.T) {
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

			r, err := backend.Resolve.Resolve(backend2.Spec{Style: id.FromString(name)})
			require.NoError(t, err, "resolution should pass")

			assert.Equal(t, c, r.Class().Metadata().Capabilities)
		})
	}
}
