// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package litmus_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/MattWindsor91/act-tester/internal/tool/litmus"
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
		fixset litmus.Fixset
		want   string
	}{
		"no-fixes": {
			fixset: litmus.Fixset{},
			want:   mainFileExample,
		},
		"stdbool": {
			fixset: litmus.Fixset{InjectStdbool: true},
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
			fixset: litmus.Fixset{RemoveAtomicCasts: true},
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
