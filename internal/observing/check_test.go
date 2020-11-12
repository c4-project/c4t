// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package observing_test

import (
	"testing"

	"github.com/MattWindsor91/act-tester/internal/observing"

	"github.com/MattWindsor91/act-tester/internal/helper/testhelp"
)

type observer interface {
	Foo()
}

type demoObserver struct{}

func (*demoObserver) Foo() {}

// TestCheckObservers tests various facets of CheckObservers.
func TestCheckObservers(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		in  interface{}
		err error
	}{
		"ok": {
			in:  []observer{&demoObserver{}, &demoObserver{}},
			err: nil,
		},
		"not-a-slice": {
			in:  3,
			err: observing.ErrNotSlice,
		},
		"nil-observer": {
			in:  []observer{nil},
			err: observing.ErrObserverNil,
		},
		"empty": {
			in:  []observer{},
			err: nil,
		},
	}

	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			testhelp.ExpectErrorIs(t, observing.CheckObservers(c.in), c.err, "CheckObservers")
		})
	}
}
