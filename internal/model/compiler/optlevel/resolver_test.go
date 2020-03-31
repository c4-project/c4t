// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package optlevel_test

import (
	"testing"

	"github.com/MattWindsor91/act-tester/internal/helper/testhelp"
	"github.com/MattWindsor91/act-tester/internal/model/compiler/optlevel"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockResolver mocks Resolver.
type mockResolver struct {
	mock.Mock
}

func (m *mockResolver) DefaultLevels() map[string]struct{} {
	return m.Called().Get(0).(map[string]struct{})
}

func (m *mockResolver) Levels() map[string]optlevel.Level {
	return m.Called().Get(0).(map[string]optlevel.Level)
}

func makeMockResolver(dls map[string]struct{}, levels map[string]optlevel.Level) *mockResolver {
	var mr mockResolver
	mr.On("DefaultLevels").Return(dls).Once()
	mr.On("Levels").Return(levels).Once()
	return &mr
}

func TestSelectLevels(t *testing.T) {
	t.Parallel()

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

	cases := map[string]struct {
		sel      optlevel.Selection
		expected map[string]struct{}
		err      error
	}{
		"defaults": {
			sel:      optlevel.Selection{},
			expected: dls,
		},
		"disable-everything": {
			sel:      optlevel.Selection{Disabled: []string{"", "size", "speed"}},
			expected: nil,
		},
		"unknown-enable": {
			sel: optlevel.Selection{Enabled: []string{"kappa"}},
			err: optlevel.ErrNoSuchLevel,
		},
	}
	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			r := makeMockResolver(dls, levels)
			ls, err := optlevel.SelectLevels(r, c.sel)
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
