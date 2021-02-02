// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package pretty_test

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/c4-project/c4t/internal/stat"
	"github.com/c4-project/c4t/internal/stat/pretty"

	"github.com/c4-project/c4t/internal/helper/testhelp"
	"github.com/stretchr/testify/require"
)

// TestPrinter_OnAnalysis_regress performs regression testing on the output of the pretty printer.
//
// It reads every json file in testdata/, and tests output against a few profiles of interest.
func TestPrinter_OnAnalysis_regress(t *testing.T) {
	t.Parallel()

	testhelp.TestFilesOfExt(t, "testdata", ".json", func(t *testing.T, name, path string) {
		var s stat.Set
		require.NoError(t, s.LoadFile(path), "couldn't read stat set")

		cases := map[string]pretty.Option{
			// TODO(@MattWindsor91): totals
			"mk": pretty.Options(pretty.ShowMutants(stat.FilterKilledMutants), pretty.UseTotals(false)),
			"ma": pretty.Options(pretty.ShowMutants(stat.FilterAllMutants), pretty.UseTotals(false)),
		}

		var gotw bytes.Buffer
		for cname, popt := range cases {
			gotw.Reset()
			cfile := strings.Join([]string{name, cname, "txt"}, ".")

			t.Run(cfile, func(t *testing.T) {
				wbytes, err := ioutil.ReadFile(filepath.Join("testdata", cfile))
				require.NoErrorf(t, err, "couldn't load case file: %s", cfile)
				want := string(wbytes)

				pp, err := pretty.NewPrinter(popt, pretty.WriteTo(&gotw))
				require.NoErrorf(t, err, "couldn't set up pretty printer for profile: %s", cname)
				require.NoErrorf(t, pp.Write(s), "couldn't write stats for profile: %s", cname)
				require.Equal(t, want, gotw.String())
			})
		}
	})
}
