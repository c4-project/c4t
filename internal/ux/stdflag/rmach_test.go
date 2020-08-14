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

	"github.com/MattWindsor91/act-tester/internal/quantity"

	runner2 "github.com/MattWindsor91/act-tester/internal/stage/invoker/runner"

	"github.com/MattWindsor91/act-tester/internal/stage/mach"
	"github.com/MattWindsor91/act-tester/internal/ux/stdflag"
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
		Quantities: quantity.MachNodeSet{
			Compiler: quantity.BatchSet{
				Timeout:  quantity.Timeout(10 * time.Second),
				NWorkers: 10,
			},
			Runner: quantity.BatchSet{
				Timeout:  quantity.Timeout(5 * time.Minute),
				NWorkers: 20,
			},
		},
	}}

	fmt.Println(strings.Join(mi.MachArgs(""), ", "))

	// Output:
	// -d, , -compiler-timeout, 0s, -run-timeout, 0s, -num-compiler-workers, 0, -num-run-workers, 0
	// -d, foo, -compiler-timeout, 10s, -run-timeout, 5m0s, -num-compiler-workers, 10, -num-run-workers, 20, -skip-compiler
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
			Quantities: quantity.MachNodeSet{
				Compiler: quantity.BatchSet{
					Timeout:  quantity.Timeout(27 * time.Second),
					NWorkers: 64,
				},
				Runner: quantity.BatchSet{
					Timeout:  quantity.Timeout(53 * time.Second),
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
					out := stdflag.MachConfigFromCli(ctx, quantity.MachNodeSet{})
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
