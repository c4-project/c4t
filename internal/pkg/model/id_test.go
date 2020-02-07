package model

import (
	"encoding/json"
	"reflect"
	"strconv"
	"testing"
)

var cases = []string{
	``,
	`foo`,
	`foo.bar.baz`,
	`weird-hyphens.allowed`,
}

func TestId_MarshalJSON(t *testing.T) {
	for _, c := range cases {
		id := IdFromString(c)
		want := strconv.Quote(c)
		if j, err := json.Marshal(id); err != nil {
			t.Errorf("error marshalling %q: %v", id, err)
		} else if string(j) != want {
			t.Errorf("json of %q=%q, want %q", id, string(j), want)
		}
	}
}

func TestId_MarshalJSON_RoundTrip(t *testing.T) {
	for _, c := range cases {
		id := IdFromString(c)

		var got *Id

		if j, err := json.Marshal(id); err != nil {
			t.Errorf("error marshalling %q: %v", id, err)
		} else if err = json.Unmarshal(j, &got); err != nil {
			t.Errorf("error unmarshalling %q: %v", id, err)
		} else if !reflect.DeepEqual(id, *got) {
			t.Errorf("marshal roundtrip %q came back as %q", id, got)
		}
	}
}
