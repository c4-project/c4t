// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package plan_test

import (
	"fmt"

	"github.com/c4-project/c4t/internal/id"

	"github.com/c4-project/c4t/internal/model/service/compiler"
	"github.com/c4-project/c4t/internal/plan"
	"github.com/c4-project/c4t/internal/subject/corpus"
)

// ExamplePlan_MaxNumRecipes is a testable example for MaxNumRecipes.
func ExamplePlan_MaxNumRecipes() {
	p := plan.Plan{
		Compilers: compiler.InstanceMap{
			id.FromString("gcc1"): compiler.MockX86Gcc(),
			id.FromString("gcc2"): compiler.MockX86Gcc(), // same architecture
			id.FromString("gcc3"): compiler.MockPower9GCCOpt(),
		},
		Corpus: corpus.New("foo", "bar", "baz"),
	}
	fmt.Println(p.MaxNumRecipes())

	// Output:
	// 6
}

// ExamplePlan_NumExpCompilations is a testable example for NumExpCompilations.
func ExamplePlan_NumExpCompilations() {
	p := plan.Plan{
		Compilers: compiler.InstanceMap{
			id.FromString("gcc1"): compiler.MockX86Gcc(),
			id.FromString("gcc2"): compiler.MockX86Gcc(),
			id.FromString("gcc3"): compiler.MockPower9GCCOpt(),
		},
		Corpus: corpus.New("foo", "bar", "baz"),
	}
	fmt.Println(p.NumExpCompilations())

	// Output:
	// 9
}
