// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package observing_test

import (
	"testing"

	"github.com/c4-project/c4t/internal/observing"

	"github.com/c4-project/c4t/internal/helper/testhelp"
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
		in  []observer
		err error
	}{
		"ok": {
			in:  []observer{&demoObserver{}, &demoObserver{}},
			err: nil,
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
