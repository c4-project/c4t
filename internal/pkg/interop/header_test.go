package interop

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

var headerDecodeCases = []struct {
	json string
	want Header
}{
	// TODO(@MattWindsor91): add more test cases
	{`{
		"name": "SBRlx",
		"locations": null,
		"init": { "x": 0, "y": 0 },
		"postcondition": "exists (0:a == 0 /\\ 1:a == 0)"
	 }`, Header{
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

// TestHeader_JsonDecode tests that the JSON tags on Header are set up to be able to parse valid headers.
func TestHeader_JsonDecode(t *testing.T) {
	for _, c := range headerDecodeCases {
		got := Header{}

		dec := json.NewDecoder(strings.NewReader(c.json))
		if err := dec.Decode(&got); err != nil {
			t.Errorf("decode failed with error (%s): %q", err, c.json)
		} else if !reflect.DeepEqual(got, c.want) {
			t.Errorf("decode got=%v; want=%v; input: %q", got, c.want, c.json)
		}
	}
}
