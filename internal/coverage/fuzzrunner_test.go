// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package coverage_test

import (
	"context"
	"fmt"
	"path/filepath"
	"reflect"
	"testing"

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

	var m mocks.SingleFuzzer

	conf := fuzzer.Configuration{Params: map[string]string{"fus": "ro dah"}}
	fr := coverage.FuzzRunner{Fuzzer: &m, Config: &conf}
	sub := subject.NewOrPanic(litmus.New("foo.litmus"))
	rc := coverage.RunnerContext{
		Seed:        4321,
		BucketDir:   td,
		NumInBucket: 1,
		Input:       sub,
	}

	m.On("Fuzz", mock.Anything, mock.MatchedBy(func(f job.Fuzzer) bool {
		return f.Seed == rc.Seed &&
			f.OutLitmus == filepath.Join(td, fmt.Sprintf("%d.litmus", rc.NumInBucket)) &&
			f.Config != nil &&
			reflect.DeepEqual(conf, *(f.Config))
	})).Return(nil).Once()

	err := fr.Run(context.Background(), rc)
	require.NoError(t, err, "mock fuzz run shouldn't error")

	m.AssertExpectations(t)
}
