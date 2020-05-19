// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package analyse_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/MattWindsor91/act-tester/internal/controller/analyse"

	"github.com/MattWindsor91/act-tester/internal/model/status"

	"github.com/MattWindsor91/act-tester/internal/helper/testhelp"
	"github.com/MattWindsor91/act-tester/internal/model/corpus"

	"github.com/MattWindsor91/act-tester/internal/model/plan"
)

// TestNewAnalyser_empty tests that analysing an empty corpus gives an error.
func TestNewAnalyser_empty(t *testing.T) {
	t.Parallel()

	_, err := analyse.NewAnalyser(&plan.Plan{Metadata: plan.Header{Version: plan.CurrentVer}}, 10)
	testhelp.ExpectErrorIs(t, err, corpus.ErrNone, "analysing empty plan")
}

// TestAnalyser_Analyse_mock tests that analysing an example corpus gives the expected collation.
func TestAnalyser_Analyse_mock(t *testing.T) {
	t.Parallel()

	m := plan.Mock()
	a, err := analyse.NewAnalyser(m, 10)
	require.NoError(t, err, "unexpected error initialising analyser")
	crp, err := a.Analyse(context.Background())
	require.NoError(t, err, "unexpected error analysing")

	cases := map[string]struct {
		subc         status.Status
		wantSubjects []string
	}{
		"flagged":          {subc: status.Flagged, wantSubjects: []string{"baz"}},
		"run-failures":     {subc: status.RunFail, wantSubjects: []string{}},
		"run-timeouts":     {subc: status.RunTimeout, wantSubjects: []string{"barbaz"}},
		"compile-failures": {subc: status.CompileFail, wantSubjects: []string{"bar"}},
		"compile-timeouts": {subc: status.CompileTimeout, wantSubjects: []string{}},
		"successes":        {subc: status.Ok, wantSubjects: []string{"foo"}},
	}
	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := crp.ByStatus[c.subc].Names()
			if !reflect.DeepEqual(got, c.wantSubjects) {
				t.Errorf("wrong subjects: got=%v; want=%v", got, c.wantSubjects)
			}
		})
	}
}
