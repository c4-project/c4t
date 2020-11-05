// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package resolver_test

import (
	"testing"

	backend2 "github.com/MattWindsor91/act-tester/internal/model/service/backend"

	"github.com/MattWindsor91/act-tester/internal/model/id"
	"github.com/MattWindsor91/act-tester/internal/serviceimpl/backend"
	"github.com/MattWindsor91/act-tester/internal/serviceimpl/backend/resolver"
	"github.com/stretchr/testify/assert"
)

// TestResolver_Capabilities tests that the standard resolver provides the correct capabilities for the known backends.
func TestResolver_Capabilities(t *testing.T) {
	t.Parallel()

	cases := map[string]backend.Capability{
		"nope":     0,
		"delitmus": backend.CanLift,
		"herd":     backend.CanRunStandalone,
		"litmus":   backend.CanLift | backend.CanRunStandalone | backend.CanProduceExecutables,
	}
	for name, c := range cases {
		name, c := name, c
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, c, resolver.Resolve.Capabilities(&backend2.Spec{Style: id.FromString(name)}))
		})
	}
}
