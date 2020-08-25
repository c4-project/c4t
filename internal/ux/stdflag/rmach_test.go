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

	"github.com/MattWindsor91/act-tester/internal/ux/stdflag"
	"github.com/stretchr/testify/assert"
	c "github.com/urfave/cli/v2"
)

// ExampleMachArgs is a testable example for MachArgs.
func ExampleMachArgs() {
	qempty := quantity.MachNodeSet{}
	fmt.Println(strings.Join(stdflag.MachArgs("", qempty), ", "))

	qi := quantity.MachNodeSet{
		Compiler: quantity.BatchSet{
			Timeout:  quantity.Timeout(10 * time.Second),
			NWorkers: 10,
		},
		Runner: quantity.BatchSet{
			Timeout:  quantity.Timeout(5 * time.Minute),
			NWorkers: 20,
		},
	}

	fmt.Println(strings.Join(stdflag.MachArgs("foo", qi), ", "))

	// Output:
	// -d, , -compiler-timeout, 0s, -run-timeout, 0s, -num-compiler-workers, 0, -num-run-workers, 0
	// -d, foo, -compiler-timeout, 10s, -run-timeout, 5m0s, -num-compiler-workers, 10, -num-run-workers, 20
}

// TestMachConfigFromCli_roundTrip tests that sending a local config through CLI flags works properly.
func TestMachConfigFromCli_roundTrip(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		dir string
		qs  quantity.MachNodeSet
	}{
		"empty": {},
		"quantities": {
			dir: "foo",
			qs: quantity.MachNodeSet{
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

			args := stdflag.MachInvocation(in.dir, in.qs)
			a := testApp(
				func(ctx *c.Context) error {
					t.Helper()
					odir := stdflag.OutDirFromCli(ctx)
					assert.Equal(t, in.dir, odir, "directories didn't match")
					oqs := stdflag.MachNodeQuantitySetFromCli(ctx)
					assert.Equal(t, in.qs, oqs, "quantities didn't match")

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
