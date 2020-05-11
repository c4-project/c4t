// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package analysis_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/MattWindsor91/act-tester/internal/model/subject"

	"github.com/stretchr/testify/assert"

	"github.com/MattWindsor91/act-tester/internal/model/corpus"
	"github.com/MattWindsor91/act-tester/internal/model/corpus/analysis"
)

// TestAnalyse_empty tests that collating an empty corpus gives an empty collation.
func TestAnalyse_empty(t *testing.T) {
	t.Parallel()

	c, err := analysis.Analyse(context.Background(), corpus.Corpus{}, 10)
	if err != nil {
		t.Fatal("unexpected error collating empty corpus:", err)
	}
	for k, coll := range c.ByStatus {
		assert.Emptyf(t, coll, "%s not empty", k)
	}
}

// TestAnalyse_mock tests that collating an example corpus gives the expected collation.
func TestAnalyse_mock(t *testing.T) {
	t.Parallel()

	m := corpus.Mock()
	crp, err := analysis.Analyse(context.Background(), m, 10)
	if err != nil {
		t.Fatal("unexpected error collating mock corpus:", err)
	}

	cases := map[string]struct {
		subc         subject.Status
		wantSubjects []string
	}{
		"flagged":          {subc: subject.StatusFlagged, wantSubjects: []string{"baz"}},
		"run-failures":     {subc: subject.StatusRunFail, wantSubjects: []string{}},
		"run-timeouts":     {subc: subject.StatusRunTimeout, wantSubjects: []string{"barbaz"}},
		"compile-failures": {subc: subject.StatusCompileFail, wantSubjects: []string{"bar"}},
		"compile-timeouts": {subc: subject.StatusCompileTimeout, wantSubjects: []string{}},
		"successes":        {subc: subject.StatusOk, wantSubjects: []string{"foo"}},
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
