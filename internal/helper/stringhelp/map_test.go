// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package stringhelp_test

import (
	"testing"

	"github.com/c4-project/c4t/internal/helper/stringhelp"
	"github.com/c4-project/c4t/internal/helper/testhelp"
)

// TestMapKeys_notStringMaps makes sure MapKeys does the right thing when given things that aren't string maps.
func TestMapKeys_notStringMaps(t *testing.T) {
	t.Parallel()

	cases := map[string]interface{}{
		"int":          5,
		"string":       "foo",
		"string-slice": []string{"foo", "bar"},
		"int-map":      map[int]string{1: "foo", 2: "bar"},
	}

	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			_, err := stringhelp.MapKeys(c)
			testhelp.ExpectErrorIs(t, err, stringhelp.ErrNotMap, "MapKeys")
		})
	}
}
