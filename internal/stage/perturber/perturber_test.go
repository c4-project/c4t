// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package perturber_test

import (
	"context"
	"testing"

	"github.com/MattWindsor91/act-tester/internal/model/service/compiler/mocks"

	"github.com/MattWindsor91/act-tester/internal/helper/stringhelp"
	"github.com/MattWindsor91/act-tester/internal/model/service/compiler/optlevel"
	"github.com/stretchr/testify/mock"

	"github.com/MattWindsor91/act-tester/internal/plan"
	"github.com/MattWindsor91/act-tester/internal/stage/perturber"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPerturber_Run(t *testing.T) {
	pm := plan.Mock()

	// This should give us a degree of sampling.
	sampleSize := len(pm.Corpus) / 2
	require.Less(t, 0, sampleSize, "sample size of mock plan is nonpositive")
	qs := perturber.QuantitySet{CorpusSize: sampleSize}

	var mi mocks.Inspector

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
	mi.On("DefaultMOpts", mock.Anything).Return(dms, nil).Times(clen)
	mi.On("DefaultOptLevels", mock.Anything).Return(dls, nil).Times(clen)
	mi.On("OptLevels", mock.Anything).Return(ols, nil).Times(clen)

	// TODO(@MattWindsor91): test observers
	pt, err := perturber.New(&mi, perturber.OverrideQuantities(qs))
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
		// TODO(@MattWindsor91): other assertions?  merge with other tests?
	}

	mi.AssertExpectations(t)
}
