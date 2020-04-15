// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package stdflag_test

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/MattWindsor91/act-tester/internal/controller/mach"
	"github.com/MattWindsor91/act-tester/internal/controller/mach/compiler"
	"github.com/MattWindsor91/act-tester/internal/controller/mach/runner"
	"github.com/MattWindsor91/act-tester/internal/controller/mach/timeout"
	"github.com/MattWindsor91/act-tester/internal/view/stdflag"
	"github.com/stretchr/testify/assert"
	c "github.com/urfave/cli/v2"
)

// TestMachConfigFromCli_roundTrip tests that sending a local config through CLI flags works properly.
func TestMachConfigFromCli_roundTrip(t *testing.T) {
	t.Parallel()

	cases := map[string]mach.UserConfig{
		"empty": {},
		"skip-compiler": {
			SkipCompiler: true,
		},
		"skip-runner": {
			SkipRunner: true,
		},
		"quantities": {
			OutDir: "foo",
			Quantities: mach.QuantitySet{
				Compiler: compiler.QuantitySet{
					Timeout: timeout.Timeout(27 * time.Second),
				},
				Runner: runner.QuantitySet{
					Timeout:  timeout.Timeout(53 * time.Second),
					NWorkers: 42,
				},
			},
		},
	}

	for name, in := range cases {
		in := in
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			args := stdflag.MachInvoker{Config: &in}.MachArgs("")
			a := testApp(
				func(ctx *c.Context) error {
					t.Helper()
					out := stdflag.MachConfigFromCli(ctx, mach.QuantitySet{})
					assert.Equal(t, in, out, "user config didn't match")

					return nil
				})
			err := a.Run(args)
			assert.NoError(t, err)
		})
	}
}

func testApp(action func(*c.Context) error) *c.App {
	a := c.NewApp()
	a.Flags = stdflag.MachCliFlags()
	a.Writer = ioutil.Discard
	a.Action = action
	return a
}
