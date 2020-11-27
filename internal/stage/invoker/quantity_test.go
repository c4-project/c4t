// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package invoker_test

import (
	"testing"

	"github.com/MattWindsor91/c4t/internal/helper/testhelp"

	"github.com/MattWindsor91/c4t/internal/plan"
	"github.com/MattWindsor91/c4t/internal/quantity"
	"github.com/MattWindsor91/c4t/internal/stage/invoker"
	"github.com/stretchr/testify/require"
)

// TestNopPlanQuantityOverrider_OverrideQuantitiesFromPlan tests that overriding quantities on a nop overrider does
// nothing.
func TestNopPlanQuantityOverrider_OverrideQuantitiesFromPlan(t *testing.T) {
	t.Parallel()

	in := quantity.MachNodeSet{
		Compiler: quantity.BatchSet{
			Timeout:  1,
			NWorkers: 2,
		},
		Runner: quantity.BatchSet{
			Timeout:  3,
			NWorkers: 4,
		},
	}
	out := in
	err := invoker.NopPlanQuantityOverrider{}.OverrideQuantitiesFromPlan(plan.Mock(), &out)
	require.NoError(t, err, "nop overrider should not error")
	require.Equal(t, in, out, "nop overrider should not change quantities")
}

// TestConfigPlanQuantityOverrider_OverrideQuantitiesFromPlan checks various cases of
// ConfigPlanQuantityOverrider.OverrideQuantitiesFromPlan.
func TestConfigPlanQuantityOverrider_OverrideQuantitiesFromPlan(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		pdelta        func(*plan.Plan)
		andThen, want quantity.MachNodeSet
	}{
		"nil overrides": {
			pdelta: func(p *plan.Plan) {
				p.Machine.Quantities = nil
			},
			andThen: quantity.MachNodeSet{},
			want: quantity.MachNodeSet{
				Compiler: quantity.BatchSet{Timeout: 1, NWorkers: 2},
				Runner:   quantity.BatchSet{Timeout: 3, NWorkers: 4},
			},
		},
		"no overrides": {
			pdelta: func(p *plan.Plan) {
				p.Machine.Quantities.Mach = quantity.MachNodeSet{}
			},
			andThen: quantity.MachNodeSet{},
			want: quantity.MachNodeSet{
				Compiler: quantity.BatchSet{Timeout: 1, NWorkers: 2},
				Runner:   quantity.BatchSet{Timeout: 3, NWorkers: 4},
			},
		},
		"plan overrides only": {
			pdelta: func(p *plan.Plan) {
				p.Machine.Quantities.Mach = quantity.MachNodeSet{
					Compiler: quantity.BatchSet{Timeout: 10},
					Runner:   quantity.BatchSet{NWorkers: 40},
				}
			},
			andThen: quantity.MachNodeSet{},
			want: quantity.MachNodeSet{
				Compiler: quantity.BatchSet{Timeout: 10, NWorkers: 2},
				Runner:   quantity.BatchSet{Timeout: 3, NWorkers: 40},
			},
		},
		"post-plan overrides only": {
			pdelta: func(p *plan.Plan) {
				p.Machine.Quantities.Mach = quantity.MachNodeSet{}
			},
			andThen: quantity.MachNodeSet{
				Compiler: quantity.BatchSet{NWorkers: 20},
				Runner:   quantity.BatchSet{Timeout: 30},
			},
			want: quantity.MachNodeSet{
				Compiler: quantity.BatchSet{Timeout: 1, NWorkers: 20},
				Runner:   quantity.BatchSet{Timeout: 30, NWorkers: 4},
			},
		},
		"both overrides": {
			pdelta: func(p *plan.Plan) {
				p.Machine.Quantities.Mach = quantity.MachNodeSet{
					Compiler: quantity.BatchSet{Timeout: 10, NWorkers: 20},
				}
			},
			andThen: quantity.MachNodeSet{
				Compiler: quantity.BatchSet{NWorkers: 200},
				Runner:   quantity.BatchSet{NWorkers: 40},
			},
			want: quantity.MachNodeSet{
				Compiler: quantity.BatchSet{Timeout: 10, NWorkers: 200},
				Runner:   quantity.BatchSet{Timeout: 3, NWorkers: 40},
			},
		},
	}

	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			in := quantity.MachNodeSet{
				Compiler: quantity.BatchSet{Timeout: 1, NWorkers: 2},
				Runner:   quantity.BatchSet{Timeout: 3, NWorkers: 4},
			}

			p := plan.Mock()
			p.Machine.Quantities = &quantity.MachineSet{}
			c.pdelta(p)

			o := invoker.ConfigPlanQuantityOverrider{PostPlanOverrides: c.andThen}

			out := in
			err := o.OverrideQuantitiesFromPlan(p, &out)
			require.NoError(t, err, "override should not error")

			require.Equal(t, c.want, out, "quantity set not overridden correctly")
		})
	}
}

// TestConfigPlanQuantityOverrider_OverrideQuantitiesFromPlan_error checks whether various error cases of
// ConfigPlanQuantityOverrider.OverrideQuantitiesFromPlan work properly.
func TestConfigPlanQuantityOverrider_OverrideQuantitiesFromPlan_error(t *testing.T) {
	t.Parallel()

	o := invoker.ConfigPlanQuantityOverrider{}
	var out quantity.MachNodeSet

	err := o.OverrideQuantitiesFromPlan(nil, &out)
	testhelp.ExpectErrorIs(t, err, plan.ErrNil, "OverrideQuantitiesFromPlan with nil plan")
}
