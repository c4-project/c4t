package model

import "fmt"

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
