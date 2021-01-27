// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package c4f_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/c4-project/c4t/internal/model/service"
	"github.com/c4-project/c4t/internal/model/service/mocks"

	"github.com/1set/gut/ystring"

	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/require"

	"github.com/c4-project/c4t/internal/c4f"
	"github.com/c4-project/c4t/internal/machine"
	"github.com/c4-project/c4t/internal/model/service/fuzzer"
)

// TestRunner_Fuzz tests the happy path of Runner.Fuzz using a mock command runner.
func TestRunner_Fuzz(t *testing.T) {
	// TODO(@MattWindsor91): This affects the filesystem, so I'm unsure as to whether it can be parallel
	j1 := fuzzer.Job{
		Seed:      8675309,
		In:        "foo.litmus",
		OutLitmus: "foo.fuzz.litmus",
		OutTrace:  "foo.trace",
		Machine:   &machine.Machine{Cores: 16},
		Config: &fuzzer.Configuration{Params: map[string]string{
			"action.var.make":               "10",
			"bool.mem.unsafe-weaken-orders": "true",
			"int.action.cap.upper":          "1000",
		}},
	}
	// Also test without the trace enabled.
	j2 := j1
	j2.OutLitmus = ""

	cases := map[string]fuzzer.Job{
		"with-trace":    j1,
		"without-trace": j2,
	}
	for name, j := range cases {
		t.Run(name, func(t *testing.T) {
			// TODO(@MattWindsor91): see above wrt Parallel
			cr := new(mocks.Runner)
			cr.Test(t)

			cr.On("Run", mock.Anything, mock.MatchedBy(func(c service.RunInfo) bool {
				return checkFuzzCmdSpec(c, j)
			})).Return(nil).Once()

			a := c4f.Runner{Base: cr}

			err := a.Fuzz(context.Background(), j)
			require.NoError(t, err, "mocked fuzzing should succeed")

			cr.AssertExpectations(t)
		})
	}
}

func checkFuzzCmdSpec(c service.RunInfo, j fuzzer.Job) bool {
	haveTrace := ystring.IsNotEmpty(j.OutTrace)
	wantLen := 8
	if haveTrace {
		wantLen += 2
	}
	if len(c.Args) != wantLen {
		return false
	}
	if haveTrace && (c.Args[7] != "-trace-output" || c.Args[8] != j.OutTrace) {
		return false
	}
	return c.Cmd == c4f.BinC4fFuzz &&
		c.Args[0] == "run" &&
		c.Args[1] == "-config" &&
		// TODO(@MattWindsor91): check config file?
		c.Args[3] == "-seed" &&
		c.Args[4] == strconv.Itoa(int(j.Seed)) &&
		c.Args[5] == "-o" &&
		c.Args[6] == j.OutLitmus &&
		c.Args[wantLen-1] == j.In
}
