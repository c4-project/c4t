// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package perturber_test

import (
	"context"
	"sort"
	"testing"

	"github.com/c4-project/c4t/internal/mutation"

	"github.com/c4-project/c4t/internal/subject/corpus/builder"

	"github.com/c4-project/c4t/internal/model/service/compiler/mocks"
	mocks2 "github.com/c4-project/c4t/internal/stage/perturber/mocks"

	"github.com/c4-project/c4t/internal/model/service/compiler"
	"github.com/c4-project/c4t/internal/observing"

	"github.com/c4-project/c4t/internal/quantity"

	"github.com/c4-project/c4t/internal/helper/stringhelp"
	"github.com/c4-project/c4t/internal/model/service/compiler/optlevel"
	"github.com/stretchr/testify/mock"

	"github.com/c4-project/c4t/internal/plan"
	"github.com/c4-project/c4t/internal/stage/perturber"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPerturber_Run tests the happy path of a perturber using copious amounts of mocking.
func TestPerturber_Run(t *testing.T) {
	pm := plan.Mock()
	pm.Mutation = &mutation.Config{
		Enabled:   true,
		Selection: 42,
	}

	// This should give us a degree of sampling.
	sampleSize := len(pm.Corpus) / 2
	require.Less(t, 0, sampleSize, "sample size of mock plan is nonpositive")
	qs := quantity.PerturbSet{CorpusSize: sampleSize}

	var (
		mi mocks.Inspector
		mo mocks2.Observer
	)
	mi.Test(t)
	mo.Test(t)

	dls := stringhelp.NewSet("0", "2", "fast")
	dms := stringhelp.NewSet("march=armv7-a")

	ols := map[string]optlevel.Level{
		"0": {
			Optimises:       false,
			Bias:            optlevel.BiasDebug,
			BreaksStandards: false,
		},
		"1": {
			Optimises:       true,
			Bias:            optlevel.BiasSize,
			BreaksStandards: false,
		},
		"2": {
			Optimises:       true,
			Bias:            optlevel.BiasSpeed,
			BreaksStandards: false,
		},
		"3": {
			Optimises:       true,
			Bias:            optlevel.BiasSpeed,
			BreaksStandards: false,
		},
		"fast": {
			Optimises:       true,
			Bias:            optlevel.BiasSpeed,
			BreaksStandards: true,
		},
	}

	clen := len(pm.Compilers)
	cids, err := pm.CompilerIDs()
	require.NoError(t, err, "mock plan CompilerIDs shouldn't error")

	mi.On("DefaultMOpts", mock.Anything).Return(dms, nil).Times(clen)
	mi.On("DefaultOptLevels", mock.Anything).Return(dls, nil).Times(clen)
	mi.On("OptLevels", mock.Anything).Return(ols, nil).Times(clen)

	mockOnPerturb(&mo, perturber.KindStart, func(set *quantity.PerturbSet, i int64) bool {
		return set.CorpusSize == qs.CorpusSize
		// Remember to add more checks here
	}).Return().Once()
	mockOnPerturb(&mo, perturber.KindRandomisingOpts, func(*quantity.PerturbSet, int64) bool {
		return true
	}).Return().Once()
	mockOnPerturb(&mo, perturber.KindSamplingCorpus, func(*quantity.PerturbSet, int64) bool {
		return true
	}).Return().Once()
	mockOnPerturb(&mo, perturber.KindSeedChanged, func(*quantity.PerturbSet, int64) bool {
		return true
	}).Return().Once()

	mockOnBuild(&mo, observing.BatchStart, func(num int, name string, _ *builder.Request) bool {
		return num == sampleSize && name == "sampled"
	}).Return().Once()
	mockOnBuild(&mo, observing.BatchStep, func(num int, _ string, r *builder.Request) bool {
		// TODO(@MattWindsor91): check that the sampled item is in the corpus
		return 0 <= num && num < sampleSize && r != nil && r.Add != nil
	}).Return().Times(sampleSize)
	mockOnBuild(&mo, observing.BatchEnd, func(int, string, *builder.Request) bool {
		return true
	}).Return().Once()

	mockOnCompilerConfig(&mo, observing.BatchStart, func(n int, _ *compiler.Named) bool {
		return n == clen
	}).Return().Once()
	mockOnCompilerConfig(&mo, observing.BatchStep, func(_ int, nc *compiler.Named) bool {
		i := sort.Search(clen, func(i int) bool {
			return !cids[i].Less(nc.ID)
		})
		return i < clen && cids[i].Equal(nc.ID)
	}).Return().Times(clen)
	mockOnCompilerConfig(&mo, observing.BatchEnd, func(int, *compiler.Named) bool {
		return true
	}).Return().Once()

	pt, err := perturber.New(&mi, perturber.OverrideQuantities(qs), perturber.ObserveWith(&mo))
	require.NoError(t, err, "error when constructing perturber")
	np, err := pt.Run(context.Background(), pm)
	require.NoError(t, err, "error when perturbing")

	for n, s := range np.Corpus {
		if !assert.Contains(t, pm.Corpus, n, "sample contains spurious subject") {
			continue
		}
		assert.Equal(t, s.Source.Path, pm.Corpus[n].Source.Path, "sample has changed path")
		// TODO(@MattWindsor91): other assertions?  merge with other tests?
	}

	for n, c := range np.Compilers {
		if !assert.Contains(t, pm.Compilers, n, "compiler set contains spurious compiler") {
			continue
		}
		assert.Equal(t, c.Arch, pm.Compilers[n].Arch, "compiler randomisation changed arch")
		assert.Equal(t, c.Mutant, pm.Mutation.Selection, "compiler randomisation didn't copy mutant ID")
		// TODO(@MattWindsor91): other assertions?  merge with other tests?
	}

	mi.AssertExpectations(t)
	mo.AssertExpectations(t)
}

func mockOnBuild(mo *mocks2.Observer, kind observing.BatchKind, f func(num int, name string, request *builder.Request) bool) *mock.Call {
	return mo.On("OnBuild", mock.MatchedBy(func(m builder.Message) bool {
		if m.Kind != kind {
			return false
		}
		return f(m.Num, m.Name, m.Request)
	}))
}
func mockOnPerturb(mo *mocks2.Observer, kind perturber.Kind, f func(*quantity.PerturbSet, int64) bool) *mock.Call {
	return mo.On("OnPerturb", mock.MatchedBy(func(m perturber.Message) bool {
		if m.Kind != kind {
			return false
		}
		return f(m.Quantities, m.Seed)
	}))
}
func mockOnCompilerConfig(mo *mocks2.Observer, kind observing.BatchKind, f func(int, *compiler.Named) bool) *mock.Call {
	return mo.On("OnCompilerConfig", mock.MatchedBy(func(m compiler.Message) bool {
		if m.Kind != kind {
			return false
		}
		return f(m.Num, m.Configuration)
	}))
}
