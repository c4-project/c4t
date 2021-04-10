// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package obs_test

import (
	"testing"

	"github.com/c4-project/c4t/internal/subject/obs"

	"github.com/c4-project/c4t/internal/helper/testhelp"
)

// TestTag_MarshalText_roundTrip tests Tag.MarshalText by round-tripping a JSON encoding.
func TestTag_MarshalText_roundTrip(t *testing.T) {
	t.Parallel()
	for i := obs.TagUnknown; i <= obs.TagLast; i++ {
		i := i
		t.Run(i.String(), func(t *testing.T) {
			t.Parallel()
			testhelp.TestJSONRoundTrip(t, i, "round-tripping tag")
		})
	}
}
