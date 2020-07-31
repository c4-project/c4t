// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package plan_test

import (
	"fmt"
	"time"

	"github.com/MattWindsor91/act-tester/internal/plan"
	"github.com/MattWindsor91/act-tester/internal/plan/stage"
)

// ExampleMetadata_RequireStage is a testable example for RequireStage.
func ExampleMetadata_RequireStage() {
	m := plan.NewMetadata(plan.UseDateSeed)
	fmt.Println("starts with plan stage?:", m.RequireStage(stage.Plan) == nil)
	m.ConfirmStage(stage.Plan, time.Now(), 0)
	fmt.Println("ends with plan stage?:", m.RequireStage(stage.Plan) == nil)

	// Output:
	// starts with plan stage?: false
	// ends with plan stage?: true
}
