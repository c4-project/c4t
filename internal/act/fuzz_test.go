// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package act_test

import (
	"context"
	"io"
	"strconv"
	"testing"

	"github.com/MattWindsor91/act-tester/internal/model/service"
	"github.com/MattWindsor91/act-tester/internal/model/service/mocks"

	"github.com/1set/gut/ystring"

	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/require"

	"github.com/MattWindsor91/act-tester/internal/act"
	"github.com/MattWindsor91/act-tester/internal/machine"
	"github.com/MattWindsor91/act-tester/internal/model/service/fuzzer"
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

			a := act.Runner{RunnerFactory: func(io.Writer, io.Writer) service.Runner { return cr }}

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
	return c.Cmd == act.BinActFuzz &&
		c.Args[0] == "run" &&
		c.Args[1] == "-config" &&
		// TODO(@MattWindsor91): check config file?
		c.Args[3] == "-seed" &&
		c.Args[4] == strconv.Itoa(int(j.Seed)) &&
		c.Args[5] == "-o" &&
		c.Args[6] == j.OutLitmus &&
		c.Args[wantLen-1] == j.In
}