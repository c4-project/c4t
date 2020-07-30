// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package compiler_test

import (
	"errors"
	"testing"

	"github.com/MattWindsor91/act-tester/internal/helper/stringhelp"

	"github.com/MattWindsor91/act-tester/internal/model/service/compiler"

	"github.com/MattWindsor91/act-tester/internal/helper/testhelp"
	"github.com/MattWindsor91/act-tester/internal/model/service/compiler/optlevel"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockResolver mocks Inspector.
type mockResolver struct {
	mock.Mock
}

func (m *mockResolver) DefaultOptLevels(c *compiler.Compiler) (stringhelp.Set, error) {
	args := m.Called(c)
	return args.Get(0).(stringhelp.Set), args.Error(1)
}

func (m *mockResolver) DefaultMOpts(c *compiler.Compiler) (stringhelp.Set, error) {
	args := m.Called(c)
	return args.Get(0).(stringhelp.Set), args.Error(1)
}

func (m *mockResolver) OptLevels(c *compiler.Compiler) (map[string]optlevel.Level, error) {
	args := m.Called(c)
	return args.Get(0).(map[string]optlevel.Level), args.Error(1)
}

func makeMockResolver(dls, dms stringhelp.Set, levels map[string]optlevel.Level, derr, merr, oerr error) *mockResolver {
	var mr mockResolver
	mr.On("DefaultOptLevels", mock.Anything).Return(dls, derr).Once()
	mr.On("DefaultMOpts", mock.Anything).Return(dms, merr).Once()
	mr.On("OptLevels", mock.Anything).Return(levels, oerr).Once()
	return &mr
}

// TestSelectLevels tests SelectLevels on a variety of cases.
func TestSelectLevels(t *testing.T) {
	t.Parallel()

	dls, levels := testData()
	mr := func() *mockResolver { return makeMockResolver(dls, nil, levels, nil, nil, nil) }

	err := errors.New("test error please ignore")

	cases := map[string]struct {
		conf     *compiler.Compiler
		res      func() *mockResolver
		expected stringhelp.Set
		err      error
	}{
		"defaults-nil": {
			conf:     &compiler.Compiler{Opt: nil},
			res:      mr,
			expected: dls,
		},
		"defaults": {
			conf:     &compiler.Compiler{Opt: &optlevel.Selection{}},
			res:      mr,
			expected: dls,
		},
		"disable-everything": {
			conf:     &compiler.Compiler{Opt: &optlevel.Selection{Disabled: []string{"", "size", "speed"}}},
			res:      mr,
			expected: nil,
		},
		"unknown-enable": {
			conf: &compiler.Compiler{Opt: &optlevel.Selection{Enabled: []string{"kappa"}}},
			res:  mr,
			err:  compiler.ErrNoSuchLevel,
		},
		"no-conf": {
			conf: nil,
			res:  mr,
			err:  compiler.ErrConfigNil,
		},
		"d-error": {
			conf: &compiler.Compiler{Opt: &optlevel.Selection{}},
			res:  func() *mockResolver { return makeMockResolver(nil, nil, levels, err, nil, nil) },
			err:  err,
		},
		"o-error": {
			conf: &compiler.Compiler{Opt: &optlevel.Selection{}},
			res:  func() *mockResolver { return makeMockResolver(dls, nil, nil, nil, nil, err) },
			err:  err,
		},
	}
	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ls, err := compiler.SelectLevels(c.res(), c.conf)
			if !testhelp.ExpectErrorIs(t, err, c.err, "SelectLevels") || err != nil {
				return
			}

			for n, l := range ls {
				assert.Equal(t, levels[n], l, "selected level inconsistent with input")
				assert.Contains(t, c.expected, n, "selected level not expected", n)
			}
			for n := range c.expected {
				assert.Contains(t, ls, n, "expected level not selected", n)
			}
		})
	}
}

// TestSelectMOpts tests SelectMOpts on a variety of cases.
func TestSelectMOpts(t *testing.T) {
	t.Parallel()

	dms := stringhelp.NewSet("march=native", "march=x86_64", "march=skylake")
	dmsk := dms.Copy()
	dmsk.Add("kappa")

	mr := func() *mockResolver { return makeMockResolver(nil, dms, nil, nil, nil, nil) }

	err := errors.New("test error please ignore")

	cases := map[string]struct {
		conf     *compiler.Compiler
		res      func() *mockResolver
		expected stringhelp.Set
		err      error
	}{
		"defaults-nil": {
			conf:     &compiler.Compiler{MOpt: nil},
			res:      mr,
			expected: dms,
		},
		"defaults": {
			conf:     &compiler.Compiler{MOpt: &optlevel.Selection{}},
			res:      mr,
			expected: dms,
		},
		"disable-everything": {
			conf:     &compiler.Compiler{MOpt: &optlevel.Selection{Disabled: dms.Slice()}},
			res:      mr,
			expected: nil,
		},
		"enable-new": {
			conf:     &compiler.Compiler{MOpt: &optlevel.Selection{Enabled: []string{"kappa"}}},
			res:      mr,
			expected: dmsk,
		},
		"no-conf": {
			conf: nil,
			res:  mr,
			err:  compiler.ErrConfigNil,
		},
		"error": {
			conf: &compiler.Compiler{MOpt: &optlevel.Selection{}},
			res:  func() *mockResolver { return makeMockResolver(nil, nil, nil, nil, err, nil) },
			err:  err,
		},
	}
	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ls, err := compiler.SelectMOpts(c.res(), c.conf)
			if !testhelp.ExpectErrorIs(t, err, c.err, "SelectMOpts") || err != nil {
				return
			}

			assert.ElementsMatch(t, c.expected.Slice(), ls.Slice(), "selected ops not expected")
		})
	}
}

func testData() (stringhelp.Set, map[string]optlevel.Level) {
	dls := map[string]struct{}{"": {}, "size": {}, "speed": {}}
	levels := map[string]optlevel.Level{
		"": {
			Optimises:       false,
			Bias:            optlevel.BiasDebug,
			BreaksStandards: false,
		},
		"size": {
			Optimises:       true,
			Bias:            optlevel.BiasSize,
			BreaksStandards: false,
		},
		"speed": {
			Optimises:       true,
			Bias:            optlevel.BiasSpeed,
			BreaksStandards: false,
		},
		"fast": {
			Optimises:       true,
			Bias:            optlevel.BiasSpeed,
			BreaksStandards: true,
		},
	}
	return dls, levels
}
