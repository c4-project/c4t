// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package coverage_test

import (
	"bytes"
	"context"
	"reflect"
	"testing"

	mocks3 "github.com/MattWindsor91/act-tester/internal/model/litmus/mocks"

	"github.com/MattWindsor91/act-tester/internal/model/recipe"

	"github.com/MattWindsor91/act-tester/internal/model/id"
	"github.com/MattWindsor91/act-tester/internal/model/service"
	mocks2 "github.com/MattWindsor91/act-tester/internal/stage/lifter/mocks"

	"github.com/MattWindsor91/act-tester/internal/model/job"
	"github.com/stretchr/testify/mock"

	"github.com/MattWindsor91/act-tester/internal/model/litmus"
	"github.com/MattWindsor91/act-tester/internal/subject"

	"github.com/stretchr/testify/require"

	"github.com/MattWindsor91/act-tester/internal/coverage"
	"github.com/MattWindsor91/act-tester/internal/model/service/fuzzer"

	"github.com/MattWindsor91/act-tester/internal/stage/fuzzer/mocks"
)

// TestFuzzRunner_Run tests FuzzRunner.Run's happy path.
func TestFuzzRunner_Run(t *testing.T) {
	td := t.TempDir()

	var (
		f mocks.SingleFuzzer
		l mocks2.SingleLifter
		s mocks3.StatDumper
		b bytes.Buffer
	)
	f.Test(t)
	l.Test(t)
	s.Test(t)

	conf := fuzzer.Configuration{Params: map[string]string{"fus": "ro dah"}}
	fr := coverage.FuzzRunner{
		Fuzzer:     &f,
		Lifter:     &l,
		StatDumper: &s,
		Config:     &conf,
		Arch:       id.ArchX86,
		Backend:    &service.Backend{Style: id.FromString("litmus")},
		ErrW:       &b,
	}
	sub := subject.NewOrPanic(litmus.New("foo.litmus"))
	rc := coverage.RunContext{
		Seed:        4321,
		BucketDir:   td,
		NumInBucket: 1,
		Input:       sub,
	}

	f.On("Fuzz", mock.Anything, mock.MatchedBy(func(f job.Fuzzer) bool {
		return f.Seed == rc.Seed &&
			f.OutLitmus == rc.OutLitmus() &&
			f.Config != nil &&
			reflect.DeepEqual(conf, *(f.Config))
	})).Return(nil).Once()
	l.On("Lift", mock.Anything, mock.MatchedBy(func(l job.Lifter) bool {
		return l.Arch.Equal(fr.Arch) &&
			l.Backend == fr.Backend &&
			l.In.Filepath() == rc.OutLitmus() &&
			l.OutDir == rc.LiftOutDir()
	}), &b).Return(recipe.Recipe{}, nil).Once()
	s.On("DumpStats", mock.Anything, mock.AnythingOfType("*litmus.Statset"), rc.OutLitmus()).Return(nil).Once()

	err := fr.Run(context.Background(), rc)
	require.NoError(t, err, "mock fuzz run shouldn't error")

	f.AssertExpectations(t)
	l.AssertExpectations(t)
	s.AssertExpectations(t)
}
