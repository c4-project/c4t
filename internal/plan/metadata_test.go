// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package plan_test

import (
	"fmt"
	"time"

	"github.com/c4-project/c4t/internal/plan"
	"github.com/c4-project/c4t/internal/plan/stage"
)

// ExampleMetadata_RequireStage is a testable example for Metadata.RequireStage.
func ExampleMetadata_RequireStage() {
	m := plan.NewMetadata(plan.UseDateSeed)
	fmt.Println("starts with plan stage?:", m.RequireStage(stage.Plan) == nil)
	m.ConfirmStage(stage.Plan, time.Now(), 0)
	fmt.Println("ends with plan stage?:", m.RequireStage(stage.Plan) == nil)

	// Output:
	// starts with plan stage?: false
	// ends with plan stage?: true
}

// ExampleMetadata_ForbidStage is a testable example for Metadata.ForbidStage.
func ExampleMetadata_ForbidStage() {
	m := plan.NewMetadata(plan.UseDateSeed)
	fmt.Println("starts without plan stage?:", m.ForbidStage(stage.Plan) == nil)
	m.ConfirmStage(stage.Plan, time.Now(), 0)
	fmt.Println("ends without plan stage?:", m.ForbidStage(stage.Plan) == nil)

	// Output:
	// starts without plan stage?: true
	// ends without plan stage?: false
}
