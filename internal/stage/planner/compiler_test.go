// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package planner_test

import (
	"testing"

	"github.com/c4-project/c4t/internal/stage/planner/mocks"

	"github.com/c4-project/c4t/internal/observing"

	"github.com/1set/gut/ystring"

	"github.com/stretchr/testify/assert"

	"github.com/c4-project/c4t/internal/model/service"

	"github.com/stretchr/testify/require"

	"github.com/c4-project/c4t/internal/stage/planner"

	"github.com/c4-project/c4t/internal/helper/stringhelp"
	"github.com/c4-project/c4t/internal/id"
	"github.com/c4-project/c4t/internal/model/service/compiler"
	cmocks "github.com/c4-project/c4t/internal/model/service/compiler/mocks"
	"github.com/c4-project/c4t/internal/model/service/compiler/optlevel"
	"github.com/stretchr/testify/mock"
)

// TestCompilerPlanner_Plan tests the happy path of a compiler planner using copious amounts of mocking.
func TestCompilerPlanner_Plan(t *testing.T) {
	var (
		ml mocks.CompilerLister
		mo cmocks.Observer
	)
	ml.Test(t)
	mo.Test(t)

	cfgs := map[id.ID]compiler.Compiler{
		id.FromString("gcc"): {
			Style: id.CStyleGCC,
			Arch:  id.ArchArmCortexA72,
		},
		id.FromString("gccnt"): {
			Disabled: true,
			Style:    id.CStyleGCC,
			Arch:     id.ArchArmCortexA72,
		},
		id.FromString("clang"): {
			Style: id.CStyleGCC,
			Arch:  id.ArchArm8,
			Run: &service.RunInfo{
				Cmd: "clang",
			},
			MOpt: &optlevel.Selection{
				Enabled:  []string{"march=armv8-a"},
				Disabled: []string{"march=armv7-a"},
			},
			Opt: &optlevel.Selection{
				Enabled:  []string{"1"},
				Disabled: []string{"fast"},
			},
		},
	}
	ncfgs := len(cfgs)

	dls := stringhelp.NewSet("0", "2", "fast")
	dms := stringhelp.NewSet("march=armv7-a")

	ml.On("Compilers").Return(cfgs, nil).Once()

	keys, _ := id.MapKeys(cfgs)

	mockOnCompilerConfig(&mo, observing.BatchStart, func(n int, _ *compiler.Named) bool {
		return n == ncfgs-1
	}).Return().Once()
	mockOnCompilerConfig(&mo, observing.BatchStep, func(_ int, nc *compiler.Named) bool {
		cs := nc.ID
		i := id.SearchSlice(keys, cs)
		return i < ncfgs && keys[i] == cs && !nc.Disabled
	}).Return().Times(ncfgs - 1)
	mockOnCompilerConfig(&mo, observing.BatchEnd, func(int, *compiler.Named) bool {
		return true
	}).Return().Once()

	cp := planner.CompilerPlanner{
		Lister:    &ml,
		Observers: []compiler.Observer{&mo},
	}

	cs, err := cp.Plan()
	require.NoError(t, err)

	ml.AssertExpectations(t)
	mo.AssertExpectations(t)

	for n, c := range cs {
		assert.Equalf(t, cfgs[n], c.Compiler, "config not passed through correctly for %s", n)

		if !ystring.IsBlank(c.SelectedMOpt) {
			checkSelection(t, "MOpt", n.String(), c.SelectedMOpt, dms.Slice(), c.MOpt)
		}
		if c.SelectedOpt != nil {
			checkSelection(t, "Opt", n.String(), c.SelectedOpt.Name, dls.Slice(), c.Opt)
		}
		assert.Falsef(t, c.Disabled, "picked up disabled compiler %s", n)
	}
}

func mockOnCompilerConfig(mo *cmocks.Observer, kind observing.BatchKind, f func(int, *compiler.Named) bool) *mock.Call {
	return mo.On("OnCompilerConfig", mock.MatchedBy(func(m compiler.Message) bool {
		if m.Kind != kind {
			return false
		}
		return f(m.Num, m.Configuration)
	}))
}

func checkSelection(t *testing.T, ty, n, chosen string, defaults []string, sel *optlevel.Selection) {
	t.Helper()

	allowed := defaults
	if sel != nil {
		allowed = append(allowed, sel.Enabled...)
		assert.NotContainsf(t, sel.Disabled, chosen, "selected %s for %s (%s) disabled", ty, n, chosen)
	}
	assert.Containsf(t, allowed, chosen, "selected %s for %s (%s) not allowed", ty, n, chosen)
}
