// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package act

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStatset_Parse(t *testing.T) {
	r := strings.NewReader(`threads 4
returns 3
literal-bools 34
atomic-cmpxchg-statements 0
atomic-fence-statements 3
atomic-fetch-statements 2
atomic-load-statements 0
atomic-store-statements 2
atomic-xchg-statements 2`)

	want := Statset{
		Threads:      4,
		Returns:      3,
		LiteralBools: 34,
		AtomicStatements: map[string]int{
			"cmpxchg": 0,
			"fence":   3,
			"fetch":   2,
			"load":    0,
			"store":   2,
			"xchg":    2,
		},
	}

	var got Statset
	err := got.Parse(r)
	require.NoError(t, err)
	assert.Equal(t, got, want)
}
