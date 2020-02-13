package model

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"github.com/BurntSushi/toml"
)

// ExampleID_Tags is a runnable example for Tags.
func ExampleID_Tags() {
	id := IDFromString("foo.bar.baz")
	for _, tag := range id.Tags() {
		fmt.Println(tag)
	}

	// Output:
	// foo
	// bar
	// baz
}

var cases = []struct {
	name string
	id   string
}{
	{"empty", ``},
	{"one-tag", `foo`},
	{"multi-tag", `foo.bar.baz`},
	{"hyphenated", `weird-hyphens.allowed`},
}

// TestID_MarshalText tests whether text marshalling for IDs works by means of TOML encoding.
func TestID_MarshalText(t *testing.T) {
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			id := wrappedID{IDFromString(c.id)}
			want := fmt.Sprintf("ID = %q\n", id.ID)
			if s, err := encodeToString(t, id); err != nil {
				t.Errorf("error marshalling %q: %v", id, err)
			} else if s != want {
				t.Errorf("TOML of %q=%q, want %q", id.ID, s, want)
			}
		})
	}
}

// TestID_MarshalText tests whether text marshalling for IDs works by means of round-trip TOML encoding and decoding.
func TestID_MarshalText_RoundTrip(t *testing.T) {
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			id := IDFromString(c.id)

			var got wrappedID

			if str, err := encodeToString(t, wrappedID{id}); err != nil {
				t.Errorf("error marshalling %q: %v", id, err)
			} else if _, err = toml.Decode(str, &got); err != nil {
				t.Errorf("error unmarshalling %q: %v", id, err)
			} else if !reflect.DeepEqual(id, got.ID) {
				t.Errorf("marshal roundtrip %q came back as %q", id, got.ID)
			}
		})
	}
}

// wrappedID serves to lift ID into a TOMLable struct type.
type wrappedID struct {
	ID ID
}

func encodeToString(t *testing.T, in wrappedID) (string, error) {
	t.Helper()

	var buf bytes.Buffer
	enc := toml.NewEncoder(&buf)
	if err := enc.Encode(in); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func TestNewID_Valid(t *testing.T) {
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
			if id, err := NewID(test.tags...); err != nil {
				t.Errorf("NewID from tags %v error: %v", test.tags, err)
			} else if id.String() != test.want {
				t.Errorf("NewID from tags %v=%s, want %s", test.tags, id.String(), test.want)
			}
		})
	}
}
