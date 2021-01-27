// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package litmus_test

import (
	"bytes"
	"strings"
	"testing"

	litmus2 "github.com/c4-project/c4t/internal/serviceimpl/backend/herdstyle/litmus"
)

const mainFileExample = `// File example

/* Includes */
#include <stdio.h>

int
main(int argc, char **argv)
{
  x = (_Atomic int)y;
  fprintf(fhist, "x=%i", (_Atomic int)y);
}
`

func TestFixset_PatchMainFile(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		fixset litmus2.Fixset
		want   string
	}{
		"no-fixes": {
			fixset: litmus2.Fixset{},
			want:   mainFileExample,
		},
		"stdbool": {
			fixset: litmus2.Fixset{InjectStdbool: true},
			want: `// File example

/* Includes */
#include <stdbool.h>
#include <stdio.h>

int
main(int argc, char **argv)
{
  x = (_Atomic int)y;
  fprintf(fhist, "x=%i", (_Atomic int)y);
}
`,
		},
		"casts": {
			fixset: litmus2.Fixset{RemoveAtomicCasts: true},
			want: `// File example

/* Includes */
#include <stdio.h>

int
main(int argc, char **argv)
{
  x = (_Atomic int)y;
  fprintf(fhist, "x=%i", y);
}
`,
		},
	}

	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer

			r := strings.NewReader(mainFileExample)

			if err := c.fixset.PatchMainFile(r, &buf); err != nil {
				t.Fatalf("unexpected write error: %v", err)
			}

			got := buf.String()
			if got != c.want {
				t.Fatalf("patch: got %q; want %q", got, c.want)
			}
		})
	}
}
