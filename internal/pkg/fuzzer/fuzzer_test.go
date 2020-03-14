// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package fuzzer_test

import (
	"context"
	"testing"

	"github.com/MattWindsor91/act-tester/internal/pkg/iohelp"

	"github.com/MattWindsor91/act-tester/internal/pkg/corpus"

	"github.com/MattWindsor91/act-tester/internal/pkg/subject"

	"github.com/MattWindsor91/act-tester/internal/pkg/fuzzer"
	"github.com/MattWindsor91/act-tester/internal/pkg/plan"
	"github.com/MattWindsor91/act-tester/internal/pkg/testhelp"
)

// makeConfig makes a 'valid' fuzzer config.
func makeConfig() *fuzzer.Config {
	return &fuzzer.Config{
		Driver: fuzzer.NopFuzzer{},
		Paths:  &fuzzer.MockPathset{},
		Quantities: fuzzer.QuantitySet{
			CorpusSize:    0,
			SubjectCycles: 10,
		},
	}
}

// makePlan makes a 'valid' plan.
func makePlan() *plan.Plan {
	return &plan.Plan{
		Corpus: map[string]subject.Subject{
			"foo": {
				Threads: 1,
				Litmus:  "foo.litmus",
			},
			"bar": {
				Threads: 2,
				Litmus:  "bar.litmus",
			},
			"baz": {
				Threads: 3,
				Litmus:  "baz.litmus",
			},
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

			cfg := makeConfig()
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

	cfg := makeConfig()
	p := makePlan()

	f, err := fuzzer.New(cfg, p)
	if err != nil {
		t.Fatal("unexpected error in New:", err)
	}
	p2, err := f.Fuzz(context.Background())
	if err != nil {
		t.Fatal("unexpected error in Fuzz:", err)
	}

	for name, s := range p2.Corpus {
		sc, err := fuzzer.ParseSubjectCycle(name)
		if err != nil {
			t.Fatal("name of fuzzer output not a subject-cycle name:", name)
		}

		sf, ok := p.Corpus[sc.Name]
		if !ok {
			t.Fatalf("subject %s in fuzzer output has no corresponding input", name)
		}

		// This isn't exhaustive, but should be enough to catch out issues.
		if s.Threads != sf.Threads {
			t.Errorf("thread mismatch: orig=%d, fuzz=%d", s.Threads, sf.Threads)
		}
		if s.Litmus != sf.Litmus {
			t.Errorf("litmus mismatch: orig=%q, fuzz=%q", s.Litmus, sf.Litmus)
		}
	}
}
