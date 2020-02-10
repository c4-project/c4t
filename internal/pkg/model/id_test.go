package model

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"github.com/BurntSushi/toml"
)

var cases = []string{
	``,
	`foo`,
	`foo.bar.baz`,
	`weird-hyphens.allowed`,
}

func TestId_MarshalText(t *testing.T) {
	for _, c := range cases {
		id := wrappedId{IdFromString(c)}
		want := fmt.Sprintf("Id = %q\n", id.Id)
		if j, err := encodeToString(t, id); err != nil {
			t.Errorf("error marshalling %q: %v", id, err)
		} else if string(j) != want {
			t.Errorf("TOML of %q=%q, want %q", id.Id, string(j), want)
		}
	}
}

func TestId_MarshalText_RoundTrip(t *testing.T) {
	for _, c := range cases {
		id := IdFromString(c)

		var got wrappedId

		if str, err := encodeToString(t, wrappedId{id}); err != nil {
			t.Errorf("error marshalling %q: %v", id, err)
		} else if _, err = toml.Decode(str, &got); err != nil {
			t.Errorf("error unmarshalling %q: %v", id, err)
		} else if !reflect.DeepEqual(id, got.Id) {
			t.Errorf("marshal roundtrip %q came back as %q", id, got.Id)
		}
	}
}

type wrappedId struct {
	Id Id
}

func encodeToString(t *testing.T, in wrappedId) (string, error) {
	t.Helper()

	var buf bytes.Buffer
	enc := toml.NewEncoder(&buf)
	if err := enc.Encode(in); err != nil {
		return "", err
	}
	return buf.String(), nil
}
