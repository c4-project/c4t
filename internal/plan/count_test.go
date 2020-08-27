// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package plan_test

import (
	"fmt"

	"github.com/MattWindsor91/act-tester/internal/model/service/compiler"
	"github.com/MattWindsor91/act-tester/internal/plan"
	"github.com/MattWindsor91/act-tester/internal/subject/corpus"
)

// ExamplePlan_MaxNumRecipes is a testable example for MaxNumRecipes.
func ExamplePlan_MaxNumRecipes() {
	p := plan.Plan{
		Compilers: map[string]compiler.Configuration{
			"gcc1": compiler.MockX86Gcc(),
			"gcc2": compiler.MockX86Gcc(), // same architecture
			"gcc3": compiler.MockPower9GCCOpt(),
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
		Compilers: map[string]compiler.Configuration{
			"gcc1": compiler.MockX86Gcc(),
			"gcc2": compiler.MockX86Gcc(),
			"gcc3": compiler.MockPower9GCCOpt(),
		},
		Corpus: corpus.New("foo", "bar", "baz"),
	}
	fmt.Println(p.NumExpCompilations())

	// Output:
	// 9
}
