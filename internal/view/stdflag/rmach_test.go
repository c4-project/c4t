// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package stdflag_test

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
	"time"

	runner2 "github.com/MattWindsor91/act-tester/internal/stage/invoker/runner"

	"github.com/MattWindsor91/act-tester/internal/stage/mach"
	"github.com/MattWindsor91/act-tester/internal/stage/mach/compiler"
	"github.com/MattWindsor91/act-tester/internal/stage/mach/runner"
	"github.com/MattWindsor91/act-tester/internal/stage/mach/timeout"
	"github.com/MattWindsor91/act-tester/internal/view/stdflag"
	"github.com/stretchr/testify/assert"
	c "github.com/urfave/cli/v2"
)

// ExampleMachInvoker_MachArgs is a testable example for MachArgs.
func ExampleMachInvoker_MachArgs() {
	mempty := stdflag.MachInvoker{Config: &mach.UserConfig{}}
	fmt.Println(strings.Join(mempty.MachArgs(""), ", "))

	mi := stdflag.MachInvoker{Config: &mach.UserConfig{
		OutDir:       "foo",
		SkipCompiler: true,
		SkipRunner:   false,
		Quantities: mach.QuantitySet{
			Compiler: compiler.QuantitySet{
				Timeout: timeout.Timeout(10 * time.Second),
			},
			Runner: runner.QuantitySet{
				Timeout:  timeout.Timeout(5 * time.Minute),
				NWorkers: 20,
			},
		},
	}}

	fmt.Println(strings.Join(mi.MachArgs(""), ", "))

	// Output:
	// -d, , -compiler-timeout, 0s, -run-timeout, 0s, -num-workers, 0
	// -d, foo, -compiler-timeout, 10s, -run-timeout, 5m0s, -num-workers, 20, -skip-compiler
}

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

			args := runner2.Invocation(stdflag.MachInvoker{Config: &in}, "")
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
