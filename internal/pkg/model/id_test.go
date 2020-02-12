package model

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"github.com/BurntSushi/toml"
)

var cases = []struct{
	name string
	id   string
}{
	{"empty", ``},
	{"one-tag", `foo`},
	{"multi-tag", `foo.bar.baz`},
	{"hyphenated", `weird-hyphens.allowed`},
}

// TestId_MarshalText tests whether text marshalling for IDs works by means of TOML encoding.
func TestId_MarshalText(t *testing.T) {
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			id := wrappedId{IdFromString(c.id)}
			want := fmt.Sprintf("Id = %q\n", id.Id)
			if s, err := encodeToString(t, id); err != nil {
				t.Errorf("error marshalling %q: %v", id, err)
			} else if s != want {
				t.Errorf("TOML of %q=%q, want %q", id.Id, s, want)
			}
		})
	}
}

// TestId_MarshalText tests whether text marshalling for IDs works by means of round-trip TOML encoding and decoding.
func TestId_MarshalText_RoundTrip(t *testing.T) {
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			id := IdFromString(c.id)

			var got wrappedId

			if str, err := encodeToString(t, wrappedId{id}); err != nil {
				t.Errorf("error marshalling %q: %v", id, err)
			} else if _, err = toml.Decode(str, &got); err != nil {
				t.Errorf("error unmarshalling %q: %v", id, err)
			} else if !reflect.DeepEqual(id, got.Id) {
				t.Errorf("marshal roundtrip %q came back as %q", id, got.Id)
			}
		})
	}
}

// wrappedID serves to lift Id into a TOMLable struct type.
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

func TestNewId_Valid(t *testing.T) {
	tests := []struct {
		name string
		tags []string
		want string
	}{
		{"empty", []string{""}, ""},
		{"one-tag", []string{"foo"}, "foo"},
		{"multi-tag", []string{"foo", "bar", "baz"}, "foo.bar.baz"},
		{"hyphenated", []string{"weird-hyphens", "allowed"}, "weird-hyphens.allowed"},	
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if id, err := NewId(test.tags...); err != nil {
				t.Errorf("NewId from tags %v error: %v", test.tags, err)
			} else if id.String() != test.want {
				t.Errorf("NewId from tags %v=%s, want %s", test.tags, id.String(), test.want)
			}
		})
	}
}
