package model

import (
	"errors"
	"fmt"
	"testing"
)

// ExampleSubject_Harness is a testable example for Harness.
func ExampleSubject_Harness() {
	s := Subject{Harnesses: map[string]Harness{
		"localhost:x86.64": Harness{Dir: "foo", Files: []string{"bar", "baz"}},
		"spikemuth:arm.7":  Harness{Dir: "foobar", Files: []string{"barbaz"}},
	}}
	lps, le := s.Harness(IDFromString("localhost"), IDFromString("x86.64"))
	sps, se := s.Harness(IDFromString("spikemuth"), IDFromString("arm.7"))

	if le != nil {
		fmt.Println(le)
	}
	for _, l := range lps.Files {
		fmt.Println(l)
	}

	if se != nil {
		fmt.Println(se)
	}
	for _, s := range sps.Files {
		fmt.Println(s)
	}

	// Output:
	// bar
	// baz
	// barbaz
}

// TestSubject_Harness_Missing checks that trying to get a harness path for a missing machine/emits pair triggers
// the appropriate error.
func TestSubject_Harness_Missing(t *testing.T) {
	var s Subject
	_, err := s.Harness(IDFromString("localhost"), IDFromString("x86.64"))
	if err == nil {
		t.Fatal("missing harness path gave no error")
	}
	if !errors.Is(err, ErrMissingHarness) {
		t.Fatalf("missing harness path gave wrong error: %v", err)
	}
}
