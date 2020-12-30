// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package lifter_test

import (
	"errors"
	"testing"

	"github.com/c4-project/c4t/internal/model/service/backend"

	mocks2 "github.com/c4-project/c4t/internal/model/service/backend/mocks"

	"github.com/c4-project/c4t/internal/stage/lifter/mocks"

	"github.com/c4-project/c4t/internal/helper/iohelp"

	"github.com/c4-project/c4t/internal/helper/testhelp"

	"github.com/c4-project/c4t/internal/stage/lifter"
)

// TestNew_errors tests the error result of New in various situations.
func TestNew_errors(t *testing.T) {
	t.Parallel()

	opterr := errors.New("oopsie")

	cases := map[string]struct {
		// ddelta modifies the driver from a known-working value.
		ddelta func(*mocks2.Resolver) backend.Resolver
		// padelta modifies the pather from a known-working value.
		pdelta func(*mocks.Pather) lifter.Pather
		// os adds options to the constructor.
		os []lifter.Option
		// err is any error expected to occur on constructing with the modified plan and configuraiton.
		err error
	}{
		"ok": {
			err: nil,
		},
		"nil-driver": {
			ddelta: func(l *mocks2.Resolver) backend.Resolver {
				return nil
			},
			err: lifter.ErrDriverNil,
		},
		"nil-paths": {
			pdelta: func(p *mocks.Pather) lifter.Pather {
				return nil
			},
			err: iohelp.ErrPathsetNil,
		},
		"opt-err": {
			os: []lifter.Option{func(*lifter.Lifter) error {
				return opterr
			}},
			err: opterr,
		},
	}

	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var (
				mr  mocks2.Resolver
				r   backend.Resolver
				mpl mocks.Pather
				pl  lifter.Pather
			)
			mr.Test(t)
			mpl.Test(t)

			if f := c.ddelta; f != nil {
				r = f(&mr)
			} else {
				r = &mr
			}
			if f := c.pdelta; f != nil {
				pl = f(&mpl)
			} else {
				pl = &mpl
			}

			_, err := lifter.New(r, pl, c.os...)
			testhelp.ExpectErrorIs(t, err, c.err, "in New()")
			mr.AssertExpectations(t)
			mpl.AssertExpectations(t)
		})
	}
}
