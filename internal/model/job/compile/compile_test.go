// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package compile_test

import (
	"testing"

	"github.com/MattWindsor91/act-tester/internal/model/job/compile"

	"github.com/MattWindsor91/act-tester/internal/model/service"
	"github.com/MattWindsor91/act-tester/internal/model/service/compiler"
	"github.com/stretchr/testify/assert"
)

// TestCompile_CompilerRun tests the behaviour of CompilerRun on various compile jobs.
func TestCompile_CompilerRun(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		in    compile.Compile
		out   service.RunInfo
		isNil bool
	}{
		"no-compiler": {
			in:    compile.Compile{},
			isNil: true,
		},
		"no-runinfo": {
			in: compile.Compile{
				Compiler: &compiler.Configuration{},
			},
			isNil: true,
		},
		"present": {
			in: compile.Compile{
				Compiler: &compiler.Configuration{
					Config: compiler.Config{
						Run: &service.RunInfo{
							Cmd:  "foo",
							Args: []string{"bar", "baz"},
						},
					},
				},
			},
			out: service.RunInfo{
				Cmd:  "foo",
				Args: []string{"bar", "baz"},
			},
		},
	}

	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			pri := c.in.CompilerRun()
			if c.isNil {
				assert.Nil(t, pri, "compiler run expected to be nil")
			} else if assert.NotNil(t, pri, "compiler run expected to be non-nil") {
				assert.Equal(t, c.out, *pri, "compiler run info not equal")
			}
		})
	}
}
