package litmus_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/MattWindsor91/act-tester/internal/pkg/litmus"
)

const mainFileExample = `#include <stdio.h>

int
main(int argc, char **argv)
{
  printf("hello, world\n");
}
`

func TestFixset_PatchMainFile(t *testing.T) {
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
			want:   strings.Join([]string{litmus.IncludeStdbool, mainFileExample}, "\n"),
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
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
