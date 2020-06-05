// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package fuzzer_test

import (
	"context"
	"testing"

	"github.com/MattWindsor91/act-tester/internal/model/litmus"
	mocks2 "github.com/MattWindsor91/act-tester/internal/model/litmus/mocks"

	"github.com/MattWindsor91/act-tester/internal/controller/fuzzer/mocks"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/MattWindsor91/act-tester/internal/helper/iohelp"

	"github.com/MattWindsor91/act-tester/internal/model/corpus"

	"github.com/MattWindsor91/act-tester/internal/model/subject"

	"github.com/MattWindsor91/act-tester/internal/controller/fuzzer"
	"github.com/MattWindsor91/act-tester/internal/helper/testhelp"
	"github.com/MattWindsor91/act-tester/internal/model/plan"
)

// makeConfig makes a 'valid' fuzzer config.
func makeConfig() (*fuzzer.Config, *mocks.SubjectPather, *mocks2.StatDumper) {
	mp := new(mocks.SubjectPather)
	md := new(mocks2.StatDumper)
	return &fuzzer.Config{
		Driver:     fuzzer.NopFuzzer{},
		Paths:      mp,
		StatDumper: md,
		Quantities: fuzzer.QuantitySet{
			CorpusSize:    0,
			SubjectCycles: 10,
		},
	}, mp, md
}

// makePlan makes a 'valid' plan.
func makePlan() *plan.Plan {
	return &plan.Plan{
		Metadata: plan.Header{Version: plan.CurrentVer},
		Corpus: map[string]subject.Subject{
			"foo": *subject.NewOrPanic(litmus.New("foo.litmus", litmus.WithThreads(1))),
			"bar": *subject.NewOrPanic(litmus.New("bar.litmus", litmus.WithThreads(2))),
			"baz": *subject.NewOrPanic(litmus.New("baz.litmus", litmus.WithThreads(3))),
		},
	}
}

// TestNew_error checks that New processes various error conditions correctly.
func TestNew_error(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		// cdelta modifies the configuration from a known-working value.
		cdelta func(*fuzzer.Config) *fuzzer.Config
		// pdelta modifies the plan from a known-working value.
		pdelta func(*plan.Plan) *plan.Plan
		// err is any error expected to occur on constructing with the modified plan and configuraiton.
		err error
	}{
		"ok": {
			err: nil,
		},
		"nil-config": {
			cdelta: func(c *fuzzer.Config) *fuzzer.Config {
				return nil
			},
			err: fuzzer.ErrConfigNil,
		},
		"nil-driver": {
			cdelta: func(c *fuzzer.Config) *fuzzer.Config {
				c.Driver = nil
				return c
			},
			err: fuzzer.ErrDriverNil,
		},
		"nil-paths": {
			cdelta: func(c *fuzzer.Config) *fuzzer.Config {
				c.Paths = nil
				return c
			},
			err: iohelp.ErrPathsetNil,
		},
		"nil-plan": {
			pdelta: func(p *plan.Plan) *plan.Plan {
				return nil
			},
			err: plan.ErrNil,
		},
		"no-corpus": {
			pdelta: func(p *plan.Plan) *plan.Plan {
				p.Corpus = nil
				return p
			},
			err: corpus.ErrNone,
		},
		"small-corpus": {
			cdelta: func(c *fuzzer.Config) *fuzzer.Config {
				c.Quantities.CorpusSize = 255
				return c
			},
			err: corpus.ErrSmall,
		},
		"bad-cycles": {
			cdelta: func(c *fuzzer.Config) *fuzzer.Config {
				c.Quantities.SubjectCycles = 0
				return c
			},
			err: corpus.ErrSmall,
		},
	}

	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			cfg, _, _ := makeConfig()

			if f := c.cdelta; f != nil {
				cfg = f(cfg)
			}

			p := makePlan()
			if f := c.pdelta; f != nil {
				p = f(p)
			}

			_, err := fuzzer.New(cfg, p)
			testhelp.ExpectErrorIs(t, err, c.err, "in New()")
		})
	}
}

// TestFuzzer_Fuzz_nop tests the happy path of running the Fuzzer with a driver that doesn't do anything.
func TestFuzzer_Fuzz_nop(t *testing.T) {
	t.Parallel()

	cfg, mp, md := makeConfig()
	mp.On("Prepare").Return(nil).Once()
	mp.On("SubjectLitmus", mock.Anything).Return("fuzz.litmus")
	mp.On("SubjectTrace", mock.Anything).Return("fuzz.trace.txt")
	md.On("DumpStats", mock.Anything, mock.Anything, "fuzz.litmus").Return(nil)
	// TODO(@MattWindsor91): md mocks

	p := makePlan()

	f, err := fuzzer.New(cfg, p)
	require.NoError(t, err, "unexpected error in New")
	p2, err := f.Fuzz(context.Background())
	require.NoError(t, err, "unexpected error in Fuzz")

	for name, s := range p2.Corpus {
		sc, err := fuzzer.ParseSubjectCycle(name)
		require.NoError(t, err, "name of fuzzer output not a subject-cycle name:", name)

		sf, ok := p.Corpus[sc.Name]
		require.Truef(t, ok, "subject %s in fuzzer output has no corresponding input", name)

		// This isn't exhaustive, but should be enough to catch out issues.
		//assert.Equal(t, sf.Stats, s.Stats, "stats mismatch")
		assert.Equal(t, sf.Source, s.Source, "litmus file mismatch")
	}

	mp.AssertExpectations(t)
	md.AssertExpectations(t)
}
