// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package status_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/c4-project/c4t/internal/subject/status"
)

// TestStatus_MarshalJSON_roundTrip checks whether the JSON (un)marshalling of statuses works appropriately.
func TestStatus_MarshalJSON_roundTrip(t *testing.T) {
	t.Parallel()

	for i := status.Ok; i <= status.Last; i++ {
		i := i
		t.Run(i.String(), func(t *testing.T) {
			var b bytes.Buffer

			err := json.NewEncoder(&b).Encode(i)
			require.NoError(t, err, "encoding")

			var d status.Status
			err = json.NewDecoder(&b).Decode(&d)
			require.NoError(t, err, "decoding")
			require.Equal(t, i, d, "decoded value didn't match")
		})
	}
}
