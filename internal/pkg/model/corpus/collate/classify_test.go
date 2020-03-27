// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package collate_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/MattWindsor91/act-tester/internal/pkg/model/corpus"
	"github.com/MattWindsor91/act-tester/internal/pkg/model/corpus/collate"
)

// TestCollate_empty tests that collating an empty corpus gives an empty collation.
func TestCollate_empty(t *testing.T) {
	t.Parallel()

	c, err := collate.Collate(context.Background(), corpus.Corpus{}, 10)
	if err != nil {
		t.Fatal("unexpected error collating empty corpus:", err)
	}
	for k, coll := range c.ByStatus() {
		assert.Emptyf(t, coll, "%s not empty", k)
	}
}

// TestCollate_empty tests that collating an empty corpus gives the expected collation.
func TestCollate_mock(t *testing.T) {
	t.Parallel()

	m := corpus.Mock()
	c, err := collate.Collate(context.Background(), m, 10)
	if err != nil {
		t.Fatal("unexpected error collating mock corpus:", err)
	}

	cases := map[string]struct {
		subc         corpus.Corpus
		wantSubjects []string
	}{
		"flagged":          {subc: c.Flagged, wantSubjects: []string{"baz"}},
		"run-failures":     {subc: c.Run.Failures, wantSubjects: []string{}},
		"run-timeouts":     {subc: c.Run.Timeouts, wantSubjects: []string{"barbaz"}},
		"compile-failures": {subc: c.Compile.Failures, wantSubjects: []string{"bar"}},
		"compile-timeouts": {subc: c.Compile.Timeouts, wantSubjects: []string{}},
		"successes":        {subc: c.Successes, wantSubjects: []string{"foo"}},
	}
	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := c.subc.Names()
			if !reflect.DeepEqual(got, c.wantSubjects) {
				t.Errorf("wrong subjects: got=%v; want=%v", got, c.wantSubjects)
			}
		})
	}
}
