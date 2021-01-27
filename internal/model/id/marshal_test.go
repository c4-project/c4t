// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package id_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/c4-project/c4t/internal/model/id"

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
	t.Parallel()
	for name, ids := range cases {
		ids := ids
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			i := wrappedID{id.FromString(ids)}
			want := fmt.Sprintf("ID = %q\n", i.ID)
			s, err := encodeToString(t, i)
			require.NoErrorf(t, err, "error marshalling %q", i)
			assert.Equalf(t, want, s, "TOML of %q=%q, want %q", i.ID, s, want)
		})
	}
}

// TestID_MarshalText_roundTrip tests whether text marshalling for IDs works by means of round-trip encoding/decoding.
func TestID_MarshalText_roundTrip(t *testing.T) {
	t.Parallel()

	for name, ids := range cases {
		ids := ids
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			i := id.FromString(ids)

			var got wrappedID

			str, err := encodeToString(t, wrappedID{i})
			require.NoErrorf(t, err, "error marshalling %q", i)
			_, err = toml.Decode(str, &got)
			require.NoErrorf(t, err, "error unmarshalling %q", i)
			assert.Truef(t, i.Equal(got.ID), "marshal roundtrip %q came back as %q", i, got.ID)
		})
	}
}

// TestID_MarshalJSON tests whether text marshalling for IDs works by means of JSON encoding.
func TestID_MarshalJSON(t *testing.T) {
	t.Parallel()

	for name, ids := range cases {
		ids := ids

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var b bytes.Buffer
			want := fmt.Sprintf("%q", ids)
			i, err := id.TryFromString(ids)
			require.NoErrorf(t, err, "error id-converting %q", ids)
			err = json.NewEncoder(&b).Encode(i)
			require.NoErrorf(t, err, "error marshalling %q", ids)
			assert.JSONEq(t, want, b.String(), "comparing baseline against marshalled")
		})
	}
}

// TestID_MarshalJSON_roundTrip tests whether JSON marshalling for IDs works by means of round-trip encoding/decoding.
func TestID_MarshalJSON_roundTrip(t *testing.T) {
	t.Parallel()

	for name, ids := range cases {
		ids := ids

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var b bytes.Buffer

			want, err := id.TryFromString(ids)
			require.NoErrorf(t, err, "error id-converting %q", ids)
			err = json.NewEncoder(&b).Encode(want)
			require.NoErrorf(t, err, "error marshalling %q", ids)

			var got id.ID
			err = json.NewDecoder(&b).Decode(&got)
			require.NoErrorf(t, err, "error unmarshalling %q", ids)

			assert.True(t, want.Equal(got), "comparing baseline against marshalled")
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
