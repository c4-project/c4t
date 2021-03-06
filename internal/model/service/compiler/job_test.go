// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package compiler_test

import (
	"testing"

	"github.com/c4-project/c4t/internal/model/service/compiler"

	"github.com/c4-project/c4t/internal/model/service"
	"github.com/stretchr/testify/assert"
)

// TestCompile_CompilerRun tests the behaviour of CompilerRun on various compile jobs.
func TestCompile_CompilerRun(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		in    compiler.Job
		out   service.RunInfo
		isNil bool
	}{
		"no-compiler": {
			in:    compiler.Job{},
			isNil: true,
		},
		"no-runinfo": {
			in: compiler.Job{
				Compiler: &compiler.Instance{},
			},
			isNil: true,
		},
		"present": {
			in: compiler.Job{
				Compiler: &compiler.Instance{
					Compiler: compiler.Compiler{
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
