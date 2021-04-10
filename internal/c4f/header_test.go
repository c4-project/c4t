// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package c4f_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/c4-project/c4t/internal/c4f"
)

var headerDecodeCases = map[string]struct {
	json string
	want c4f.Header
}{
	// TODO(@MattWindsor91): add more test cases
	"sbrlx": {
		`{
		"name": "SBRlx",
		"locations": null,
		"init": { "x": 0, "y": 0 },
		"postcondition": "exists (0:a == 0 /\\ 1:a == 0)"
	 }`, c4f.Header{
			Name:      "SBRlx",
			Locations: nil,
			Init: map[string]int{
				"x": 0,
				"y": 0,
			},
			Postcondition: `exists (0:a == 0 /\ 1:a == 0)`,
		},
	},
}

// TestReadHeader tests that we can read headers properly from JSON.
func TestReadHeader(t *testing.T) {
	for name, c := range headerDecodeCases {
		t.Run(name, func(t *testing.T) {
			rd := strings.NewReader(c.json)
			var got c4f.Header
			if err := got.Read(rd); err != nil {
				t.Errorf("decode failed with error (%s): %q", err, c.json)
			} else if !reflect.DeepEqual(got, c.want) {
				t.Errorf("decode got=%v; want=%v; input: %q", got, c.want, c.json)
			}
		})
	}
}
