// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package subject_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/MattWindsor91/act-tester/internal/pkg/model"

	"github.com/BurntSushi/toml"
	"github.com/MattWindsor91/act-tester/internal/pkg/subject"
)

// TestRun_TomlDecode tests the decoding of a test run from various TOML examples.
func TestRun_TomlDecode(t *testing.T) {
	t.Parallel()
	cases := map[string]struct {
		toml string
		want subject.Run
	}{
		"empty": {toml: "", want: subject.Run{}},
		"unsat": {
			toml: `
time = 2015-10-21T07:28:00-08:00
duration = 8675309
status = "flagged"
[obs]
  flags = "unsat"

[[obs.counter_examples]]
  "0:r0" = "1"
  x = "1"

[[obs.states]]
  "0:r0" = "1"
  x = "1"`,
			want: subject.Run{
				Time:     time.Date(2015, time.October, 21, 7, 28, 0, 0, time.FixedZone("UTC-8", -8*60*60)),
				Duration: 8675309,
				Status:   subject.StatusFlagged,
				Obs: &model.Obs{
					Flags: model.ObsUnsat,
					CounterExamples: []model.ObsState{
						{"0:r0": "1", "x": "1"},
					},
					States: []model.ObsState{
						{"0:r0": "1", "x": "1"},
					},
				},
			},
		},
	}

	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			var got subject.Run
			if _, err := toml.Decode(c.toml, &got); err != nil {
				t.Fatal("unexpected decode error:", err)
			}

			if !got.Time.Equal(c.want.Time) {
				t.Errorf("badly parsed time: got=%v, want=%v", got.Time, c.want.Time)
			}
			if got.Duration != c.want.Duration {
				t.Errorf("badly parsed duration: got=%v, want=%v", got.Duration, c.want.Duration)
			}
			if got.Status != c.want.Status {
				t.Errorf("badly parsed status: got=%q, want=%q", got.Status.String(), c.want.Status.String())
			}
			if !reflect.DeepEqual(got.Obs, c.want.Obs) {
				t.Errorf("badly parsed obs: got=%v, want=%v", got.Obs, c.want.Obs)
			}
		})
	}
}
