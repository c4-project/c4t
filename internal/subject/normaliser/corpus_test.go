// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package normaliser_test

import (
	"testing"

	"github.com/c4-project/c4t/internal/subject/normaliser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/c4-project/c4t/internal/subject/corpus"
)

// TestCorpus_Normalise tests the normaliser on a corpus stitched together from all of the individual subject
// cases.
func TestCorpus_Normalise(t *testing.T) {
	t.Parallel()

	in := make(corpus.Corpus, len(testSubjects))
	want := make(corpus.Corpus, len(testSubjects))
	maps := make(map[string]normaliser.Map, len(testSubjects))
	for n, v := range testSubjects {
		v := v(n)
		in[n] = v.in
		want[n] = v.out
		maps[n] = v.maps
	}

	norm := normaliser.NewCorpus("")
	got, err := norm.Normalise(in)
	require.NoError(t, err, "normalising a corpus")

	assert.Equal(t, want, got, "comparing result corpora")

	for n, m := range maps {
		assert.Equalf(t, m, norm.BySubject[n].Mappings, "comparing mappings for subject %q", n)
	}
}
