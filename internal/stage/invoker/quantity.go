// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package invoker

import (
	"github.com/MattWindsor91/act-tester/internal/plan"
	"github.com/MattWindsor91/act-tester/internal/quantity"
)

// PlanQuantityOverrider is the interface of things that may, or may not, override the invoker's quantity set with
// quantities according to information in a plan.
type PlanQuantityOverrider interface {
	// OverrideQuantitiesFromPlan applies to qs any quantity overrides that require information from the plan p.
	// This is part of the Factory interface because some, but not all, factories require plan information.
	OverrideQuantitiesFromPlan(p *plan.Plan, qs *quantity.MachNodeSet) error
}

// NopPlanQuantityOverrider is a PlanQuantityOverrider that doesn't actually override quantities.
type NopPlanQuantityOverrider struct{}

// OverrideQuantitiesFromPlan does nothing.
func (n NopPlanQuantityOverrider) OverrideQuantitiesFromPlan(*plan.Plan, *quantity.MachNodeSet) error {
	// Intentionally do nothing.
	return nil
}

// ConfigPlanQuantityOverrider overrides quantities by looking up the machine that the plan targets in a config map,
// then extracting the overrides from there.
type ConfigPlanQuantityOverrider struct {
	// PostPlanOverrides applies any overrides that should happen after the machine overrides.
	PostPlanOverrides quantity.MachNodeSet
}

func (c *ConfigPlanQuantityOverrider) OverrideQuantitiesFromPlan(p *plan.Plan, qs *quantity.MachNodeSet) error {
	if p == nil {
		return plan.ErrNil
	}
	// TODO(@MattWindsor91): what if qs is nil?
	if p.Machine.Quantities != nil {
		qs.Override(p.Machine.Quantities.Mach)
	}
	qs.Override(c.PostPlanOverrides)
	return nil
}
