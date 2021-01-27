// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package testhelp

import (
	"bytes"
	"encoding/json"
	"io"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/BurntSushi/toml"
)

// TestJSONRoundTrip tests that encoding want to TOML then decoding it produces a value that is deep-equal to want.
func TestJSONRoundTrip(t *testing.T, want interface{}, context string) {
	t.Helper()
	testRoundTrip(t, want, func(w io.Writer, i interface{}) error {
		return json.NewEncoder(w).Encode(i)
	}, func(r io.Reader, i interface{}) error {
		return json.NewDecoder(r).Decode(i)
	}, context)
}

// TestTomlRoundTrip tests that encoding want to TOML then decoding it produces a value that is deep-equal to want.
func TestTomlRoundTrip(t *testing.T, want interface{}, context string) {
	t.Helper()
	testRoundTrip(t, want, func(w io.Writer, i interface{}) error {
		return toml.NewEncoder(w).Encode(i)
	}, func(r io.Reader, i interface{}) error {
		_, err := toml.DecodeReader(r, i)
		return err
	}, context)
}

func testRoundTrip(t *testing.T, want interface{}, to func(io.Writer, interface{}) error, fro func(io.Reader, interface{}) error, context string) {
	t.Helper()

	var b bytes.Buffer

	err := to(&b, want)
	if !assert.NoErrorf(t, err, "%s: unexpected encode error", context) {
		return
	}

	// We have to make sure that 'got' has the same type as 'want', even though at compile-time we don't know that type.
	// If we don't, then the decoder can't do type-driven decoding, and instead gives us the wrong type back.
	vgot := reflect.New(reflect.TypeOf(want))

	err = fro(&b, vgot.Interface())
	if !assert.NoErrorf(t, err, "%s: unexpected decode error", context) {
		return
	}

	got := reflect.Indirect(vgot).Interface()
	assert.Equalf(t, want, got, "%s: round-trip didn't match", context)
}
