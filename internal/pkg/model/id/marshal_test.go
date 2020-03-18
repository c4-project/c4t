// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package id_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/MattWindsor91/act-tester/internal/pkg/model/id"

	"github.com/BurntSushi/toml"
)

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
			i := wrappedID{id.FromString(ids)}
			want := fmt.Sprintf("ID = %q\n", i.ID)
			if s, err := encodeToString(t, i); err != nil {
				t.Errorf("error marshalling %q: %v", i, err)
			} else if s != want {
				t.Errorf("TOML of %q=%q, want %q", i.ID, s, want)
			}
		})
	}
}

// TestID_MarshalText_roundTrip tests whether text marshalling for IDs works by means of round-trip encoding/decoding.
func TestID_MarshalText_roundTrip(t *testing.T) {
	for name, ids := range cases {
		t.Run(name, func(t *testing.T) {
			i := id.FromString(ids)

			var got wrappedID

			if str, err := encodeToString(t, wrappedID{i}); err != nil {
				t.Errorf("error marshalling %q: %v", i, err)
			} else if _, err = toml.Decode(str, &got); err != nil {
				t.Errorf("error unmarshalling %q: %v", i, err)
			} else if !i.Equal(got.ID) {
				t.Errorf("marshal roundtrip %q came back as %q", i, got.ID)
			}
		})
	}
}

// wrappedID serves to lift CompilerID into a TOMLable struct type.
type wrappedID struct {
	ID id.ID
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
