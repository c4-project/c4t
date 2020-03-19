// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package collate_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/MattWindsor91/act-tester/internal/pkg/corpus"
	"github.com/MattWindsor91/act-tester/internal/pkg/corpus/collate"
)

// TestCollate_empty tests that collating an empty corpus gives an empty collation.
func TestCollate_empty(t *testing.T) {
	t.Parallel()

	c, err := collate.Collate(context.Background(), corpus.Corpus{}, 10)
	if err != nil {
		t.Fatal("unexpected error collating empty corpus:", err)
	}
	if len(c.Timeouts) != 0 {
		t.Error("timeouts not empty")
	}
	if len(c.Flagged) != 0 {
		t.Error("flagged not empty")
	}
	if len(c.RunFailures) != 0 {
		t.Error("run-failures not empty")
	}
	if len(c.CompileFailures) != 0 {
		t.Error("compile-failures not empty")
	}
	if len(c.Successes) != 0 {
		t.Error("successes not empty")
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
		"timeouts":         {subc: c.Timeouts, wantSubjects: []string{"barbaz"}},
		"flagged":          {subc: c.Flagged, wantSubjects: []string{"baz"}},
		"run-failures":     {subc: c.RunFailures, wantSubjects: []string{}},
		"compile-failures": {subc: c.CompileFailures, wantSubjects: []string{"bar"}},
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
