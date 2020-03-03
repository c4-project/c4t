// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package testhelp contains test helpers.
package testhelp

import (
	"bytes"
	"errors"
	"reflect"
	"testing"

	"github.com/BurntSushi/toml"
)

// ExpectErrorIs checks whether got has an 'Is' relation to want.
// If not, it fails the test with a message mentioning context.
func ExpectErrorIs(t *testing.T, got, want error, context string) {
	if got == nil {
		t.Helper()
		t.Errorf("%s: error nil; want=%q", context, want)
	} else if !errors.Is(got, want) {
		t.Helper()
		t.Errorf("%s: error=%q; want=%q", context, got, want)
	}
}

// TestTomlRoundTrip tests that encoding want to TOML then decoding it produces a value that is deep-equal to want.
func TestTomlRoundTrip(t *testing.T, want interface{}, context string) {
	var b bytes.Buffer

	if err := toml.NewEncoder(&b).Encode(want); err != nil {
		t.Helper()
		t.Errorf("%s: unexpected encode error: %v", context, err)
		return
	}

	// We have to make sure that 'got' has the same type as 'want', even though at compile-time we don't know that type.
	// If we don't, then the TOML decoder can't do type-driven decoding, and instead gives us the wrong type back.
	vgot := reflect.New(reflect.TypeOf(want))

	if _, err := toml.DecodeReader(&b, vgot.Interface()); err != nil {
		t.Helper()
		t.Errorf("%s: unexpected decode error: %v", context, err)
		return
	}

	got := reflect.Indirect(vgot).Interface()
	if !reflect.DeepEqual(got, want) {
		t.Helper()
		t.Errorf("%s: got=%v; want=%v", context, got, want)
	}
}
