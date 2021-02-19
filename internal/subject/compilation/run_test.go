// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package compilation_test

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/c4-project/c4t/internal/timing"
	"github.com/stretchr/testify/assert"

	"github.com/c4-project/c4t/internal/subject/compilation"

	"github.com/c4-project/c4t/internal/subject/status"

	"github.com/c4-project/c4t/internal/subject/obs"
)

// TestRun_JSONDecode tests the decoding of a test run from various JSON examples.
func TestRun_JSONDecode(t *testing.T) {
	t.Parallel()
	cases := map[string]struct {
		json string
		want compilation.RunResult
	}{
		"empty": {json: "{}", want: compilation.RunResult{}},
		"unsat": {
			json: `{
"time_span": {
  "start": "2015-10-21T07:28:00-08:00",
  "end": "2015-10-21T08:28:00-08:00"
},
"status": "flagged",
"obs": {
  "flags": ["unsat"],
  "states": [
    { "tag": "counter",
      "values": {"0:r0": "1", "x": "1"} } ] }}
`,
			want: compilation.RunResult{
				Result: compilation.Result{
					Timespan: timing.Span{
						Start: time.Date(2015, time.October, 21, 7, 28, 0, 0, time.FixedZone("UTC-8", -8*60*60)),
						End:   time.Date(2015, time.October, 21, 8, 28, 0, 0, time.FixedZone("UTC-8", -8*60*60)),
					},
					Status: status.Flagged,
				},
				Obs: &obs.Obs{
					Flags: obs.Unsat,
					States: []obs.State{
						{Values: obs.Valuation{"0:r0": "1", "x": "1"}, Tag: obs.TagCounter},
					},
				},
			},
		},
	}

	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			var got compilation.RunResult
			if err := json.NewDecoder(strings.NewReader(c.json)).Decode(&got); err != nil {
				t.Fatal("unexpected decode error:", err)
			}

			assert.True(t, c.want.Timespan.Start.Equal(got.Timespan.Start))
			assert.True(t, c.want.Timespan.End.Equal(got.Timespan.End))
			assert.Equal(t, c.want.Status, got.Status)
			assert.Equal(t, c.want.Obs, got.Obs)
		})
	}
}
