// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package model

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
	"testing"

	"github.com/BurntSushi/toml"
)

// ExampleID_Less is a runnable example for Less
func ExampleID_Less() {
	ids := []ID{
		IDFromString("arm.7"),
		IDFromString("ppc.64.le"),
		IDFromString("x86.32"),
		IDFromString("x86"),
		IDFromString("arm"),
		IDFromString("ppc"),
		IDFromString("x86.64"),
		IDFromString("ppc.64"),
		IDFromString("arm.6"),
	}
	sort.Slice(ids, func(i, j int) bool {
		return ids[i].Less(ids[j])
	})
	for _, id := range ids {
		fmt.Println(id)
	}

	// Output:
	// arm
	// arm.6
	// arm.7
	// ppc
	// ppc.64
	// ppc.64.le
	// x86
	// x86.32
	// x86.64
}

// ExampleID_Join is a runnable example for Join.
func ExampleID_Join() {
	id1 := IDFromString("foo.bar")
	id2 := IDFromString("baz.barbaz")
	fmt.Println(id1.Join(id2).String())

	// empty IDs do nothing when joined
	fmt.Println(ID{}.Join(id1).String())
	fmt.Println(id2.Join(ID{}).String())

	// Output:
	// foo.bar.baz.barbaz
	// foo.bar
	// baz.barbaz
}

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

var cases = map[string]string{
	"empty":      ``,
	"one-tag":    `foo`,
	"multi-tag":  `foo.bar.baz`,
	"hyphenated": `weird-hyphens.allowed`,
}

// TestID_MarshalText tests whether text marshalling for IDs works by means of TOML encoding.
func TestID_MarshalText(t *testing.T) {
	for name, ids := range cases {
		t.Run(name, func(t *testing.T) {
			id := wrappedID{IDFromString(ids)}
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
	for name, ids := range cases {
		t.Run(name, func(t *testing.T) {
			id := IDFromString(ids)

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

// wrappedID serves to lift CompilerID into a TOMLable struct type.
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
