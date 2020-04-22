// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package mock

import (
	"context"
	"io"

	"github.com/MattWindsor91/act-tester/internal/helper/stringhelp"
	mdl "github.com/MattWindsor91/act-tester/internal/model/compiler"
	"github.com/MattWindsor91/act-tester/internal/model/compiler/optlevel"
	"github.com/MattWindsor91/act-tester/internal/model/job"
	"github.com/stretchr/testify/mock"
)

// Compiler mocks various compiler-related interfaces.
type Compiler struct {
	mock.Mock
}

// DefaultOptLevels mocks the eponymous method.
func (m *Compiler) DefaultOptLevels(c *mdl.Config) (stringhelp.Set, error) {
	args := m.Called(c)
	return args.Get(0).(stringhelp.Set), args.Error(1)
}

// OptLevels mocks the eponymous method.
func (m *Compiler) OptLevels(c *mdl.Config) (map[string]optlevel.Level, error) {
	args := m.Called(c)
	return args.Get(0).(map[string]optlevel.Level), args.Error(1)
}

// DefaultMOpts mocks the eponymous method.
func (m *Compiler) DefaultMOpts(c *mdl.Config) (stringhelp.Set, error) {
	args := m.Called(c)
	return args.Get(0).(stringhelp.Set), args.Error(1)
}

// RunCompiler mocks the eponymous method.
func (m *Compiler) RunCompiler(ctx context.Context, j job.Compile, errw io.Writer) error {
	args := m.Called(ctx, j, errw)
	return args.Error(0)
}
