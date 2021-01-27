// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package pretty_test

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/c4-project/c4t/internal/helper/testhelp"
	"github.com/stretchr/testify/require"

	"github.com/c4-project/c4t/internal/plan"
	"github.com/c4-project/c4t/internal/plan/analysis"
	"github.com/c4-project/c4t/internal/stage/analyser/pretty"
)

// ExamplePrinter_OnAnalysis is a testable example for Printer.OnAnalysis.
func ExamplePrinter_OnAnalysis() {
	p := plan.Mock()
	a, err := analysis.Analyse(context.Background(), p)
	if err != nil {
		fmt.Println("analysis error:", err)
		return
	}
	pw, err := pretty.NewPrinter(pretty.ShowCompilers(true))
	if err != nil {
		fmt.Println("printer init error:", err)
		return
	}
	pw.OnAnalysis(*a)

	// Output:
	// # Compilers
	//   ## clang
	//     - style: gcc
	//     - arch: x86
	//     - opt: none
	//     - mopt: none
	//     ### Times (sec)
	//       - compile: Min 200 Avg 200 Max 200
	//       - run: Min 0 Avg 0 Max 0
	//     ### Results
	//       - Ok: 1 subject(s)
	//   ## gcc
	//     - style: gcc
	//     - arch: ppc.64le.power9
	//     - opt: none
	//     - mopt: none
	//     ### Times (sec)
	//       - compile: Min 200 Avg 200 Max 200
	//       - run: Min 0 Avg 0 Max 0
	//     ### Results
	//       - Flagged: 1 subject(s)
	//       - CompileFail: 1 subject(s)
}

// TestPrinter_OnAnalysis_regress performs regression testing on the output of the pretty printer.
//
// It reads every json file in testdata/, and tests output against a few profiles of interest.
func TestPrinter_OnAnalysis_regress(t *testing.T) {
	t.Parallel()

	testhelp.TestFilesOfExt(t, "testdata", ".json", func(t *testing.T, name, path string) {
		var p plan.Plan
		require.NoError(t, plan.ReadFile(path, &p), "couldn't read plan")

		an, err := analysis.Analyse(context.Background(), &p)
		require.NoError(t, err, "couldn't analyse plan")

		cases := map[string]pretty.Option{
			"cp":  pretty.Options(pretty.ShowCompilers(true), pretty.ShowPlanInfo(true)),
			"cs":  pretty.Options(pretty.ShowCompilers(true), pretty.ShowSubjects(true)),
			"ps":  pretty.Options(pretty.ShowPlanInfo(true), pretty.ShowSubjects(true)),
			"cps": pretty.Options(pretty.ShowCompilers(true), pretty.ShowPlanInfo(true), pretty.ShowSubjects(true)),
			"mut": pretty.Options(pretty.ShowMutation(true)),
		}

		var gotw bytes.Buffer
		for cname, popt := range cases {
			gotw.Reset()
			cfile := strings.Join([]string{name, cname, "txt"}, ".")

			t.Run(cfile, func(t *testing.T) {
				wbytes, err := ioutil.ReadFile(filepath.Join("testdata", cfile))
				require.NoError(t, err, "couldn't load case file", cfile)
				want := string(wbytes)

				pp, err := pretty.NewPrinter(popt, pretty.WriteTo(&gotw))
				require.NoError(t, err, "couldn't set up pretty printer for profile", cname)
				require.NoError(t, pp.Write(*an), "couldn't write analysis for profile", cname)
				require.Equal(t, want, gotw.String())
			})
		}
	})
}
