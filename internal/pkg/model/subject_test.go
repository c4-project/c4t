package model

import (
	"errors"
	"fmt"
	"testing"
)

// ExampleSubject_HarnessPath is a testable example for HarnessPath.
func ExampleSubject_HarnessPath() {
	s := Subject{HarnessPaths: map[string][]string{
		"localhost:x86.64": {"foo", "bar", "baz"},
		"spikemuth:arm.7":  {"foobar", "barbaz"},
	}}
	lps, le := s.HarnessPath(IdFromString("localhost"), IdFromString("x86.64"))
	sps, se := s.HarnessPath(IdFromString("spikemuth"), IdFromString("arm.7"))

	if le != nil {
		fmt.Println(le)
	}
	for _, l := range lps {
		fmt.Println(l)
	}

	if se != nil {
		fmt.Println(se)
	}
	for _, s := range sps {
		fmt.Println(s)
	}

	// Output:
	// foo
	// bar
	// baz
	// foobar
	// barbaz
}

// TestSubject_HarnessPath_Missing checks that trying to get a harness path for a missing machine/emits pair triggers
// the appropriate error.
func TestSubject_HarnessPath_Missing(t *testing.T) {
	var s Subject
	_, err := s.HarnessPath(IdFromString("localhost"), IdFromString("x86.64"))
	if err == nil {
		t.Fatal("missing harness path gave no error")
	}
	if !errors.Is(err, ErrMissingHarness) {
		t.Fatalf("missing harness path gave wrong error: %v", err)
	}
}
