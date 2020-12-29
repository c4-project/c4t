// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package coverage_test

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/c4-project/c4t/internal/observing"
	"github.com/stretchr/testify/mock"

	"github.com/c4-project/c4t/internal/coverage"

	"github.com/stretchr/testify/require"

	"github.com/c4-project/c4t/internal/coverage/mocks"
)

// TestMaker_Run tests the happy path of Maker.Run using mocks.
func TestMaker_Run(t *testing.T) {
	td := t.TempDir()

	var (
		mr1, mr2 mocks.Runner
		mo       mocks.Observer
	)
	mr1.Test(t)
	mr2.Test(t)
	mo.Test(t)

	p1 := coverage.Profile{Kind: coverage.Known, Runner: &mr1}
	p2 := coverage.Profile{Kind: coverage.Standalone, Runner: &mr2}

	mk, err := coverage.NewMaker(td,
		map[string]coverage.Profile{"p1": p1, "p2": p2},
		coverage.OverrideQuantities(coverage.QuantitySet{
			Count:     12,
			Divisions: []int{3, 2},
		}),
		coverage.ObserveWith(&mo),
	)
	require.NoError(t, err, "shouldn't error on creating the maker")

	mockObs(&mo, "p1", 12)
	mockObs(&mo, "p2", 12)

	for i := 1; i <= 2; i++ {
		mockRun(&mr1, "p1", td, fmt.Sprintf("1_%d", i), 2)
		mockRun(&mr2, "p2", td, fmt.Sprintf("1_%d", i), 2)
	}
	for i := 2; i <= 3; i++ {
		mockRun(&mr1, "p1", td, fmt.Sprintf("%d", i), 4)
		mockRun(&mr2, "p2", td, fmt.Sprintf("%d", i), 4)
	}

	err = mk.Run(context.Background())
	require.NoError(t, err, "shouldn't error on running the maker")

	mr1.AssertExpectations(t)
	mr2.AssertExpectations(t)
	mo.AssertExpectations(t)
}

func mockRun(mr *mocks.Runner, pname, rootdir, bdir string, bsize int) {
	for i := 0; i < bsize; i++ {
		i := i
		mr.On("Run", mock.Anything, mock.MatchedBy(func(mr coverage.RunContext) bool {
			return mr.NumInBucket == i && mr.BucketDir == filepath.Join(rootdir, pname, bdir)
		})).Return(nil).Once()
	}
}

func mockObs(mo *mocks.Observer, pname string, total int) {
	mockOnCoverageRun(mo, observing.BatchStart, func(pname2 string, i int) bool {
		return pname == pname2 && i == total
	}).Return().Once()
	for i := 1; i <= total; i++ {
		i := i
		// TODO(@MattWindsor91): check buckets
		mockOnCoverageRun(mo, observing.BatchStep, func(_ string, j int) bool {
			return i == j
		}).Return().Once()
	}
	mockOnCoverageRun(mo, observing.BatchEnd, func(pname2 string, _ int) bool {
		return pname == pname2
	}).Return().Once()
}

func mockOnCoverageRun(mo *mocks.Observer, kind observing.BatchKind, f func(pname string, i int) bool) *mock.Call {
	return mo.On("OnCoverageRun",
		mock.MatchedBy(func(m coverage.RunMessage) bool {
			return m.Kind == kind && f(m.ProfileName, m.Num)
		}))
}
