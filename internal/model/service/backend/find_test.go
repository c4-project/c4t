// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package backend_test

import (
	"fmt"

	"github.com/c4-project/c4t/internal/model/id"
	"github.com/c4-project/c4t/internal/model/service/backend"
)

// ExampleCriteria_String is a runnable example for Criteria.String.
func ExampleCriteria_String() {
	fmt.Println(backend.Criteria{})
	fmt.Println(backend.Criteria{IDGlob: id.FromString("litmus")})
	fmt.Println(backend.Criteria{StyleGlob: id.FromString("herdtools.*")})
	fmt.Println(backend.Criteria{IDGlob: id.FromString("litmus"), StyleGlob: id.FromString("herdtools.*")})

	// Output:
	// any
	// id=litmus
	// style=herdtools.*
	// id=litmus, style=herdtools.*
}

// ExampleCriteria_String is a runnable example for Criteria.Matches.
func ExampleCriteria_Matches() {
	spec := backend.NamedSpec{
		ID:   id.FromString("litmus"),
		Spec: backend.Spec{Style: id.FromString("herdtools.litmus")},
	}

	m1, _ := backend.Criteria{}.Matches(spec)
	fmt.Println("matches empty criteria:", m1)

	m2, _ := backend.Criteria{
		IDGlob:    id.FromString("litmus"),
		StyleGlob: id.FromString("herdtools.*"),
	}.Matches(spec)
	fmt.Println("matches first criteria:", m2)

	m3, _ := backend.Criteria{
		IDGlob:    id.FromString("litmus"),
		StyleGlob: id.FromString("herdtools.*.7"),
	}.Matches(spec)
	fmt.Println("matches second criteria:", m3)

	_, err := backend.Criteria{
		IDGlob:    id.FromString("litmus"),
		StyleGlob: id.FromString("*.herdtools.*"),
	}.Matches(spec)
	fmt.Println("error for malformed glob:", err)

	// Output:
	// matches empty criteria: true
	// matches first criteria: true
	// matches second criteria: false
	// error for malformed glob: malformed glob expression: more than one '*' character
}

// ExampleCriteria_String is a runnable example for Criteria.Find.
func ExampleCriteria_Find() {
	specs := []backend.NamedSpec{
		{ID: id.FromString("herd"), Spec: backend.Spec{Style: id.FromString("herdtools.herd")}},
		{ID: id.FromString("litmus"), Spec: backend.Spec{Style: id.FromString("herdtools.litmus")}},
		{ID: id.FromString("litmus.dev"), Spec: backend.Spec{Style: id.FromString("herdtools.litmus")}},
		{ID: id.FromString("rmem"), Spec: backend.Spec{Style: id.FromString("rmem")}},
	}

	m1, _ := backend.Criteria{}.Find(specs)
	fmt.Println("empty criteria:", m1.ID)
	m2, _ := backend.Criteria{IDGlob: id.FromString("litmus.*")}.Find(specs)
	fmt.Println("litmus criteria:", m2.ID)
	m3, _ := backend.Criteria{IDGlob: id.FromString("litmus.dev")}.Find(specs)
	fmt.Println("litmus.dev criteria:", m3.ID)
	m4, _ := backend.Criteria{StyleGlob: id.FromString("rmem.*")}.Find(specs)
	fmt.Println("rmem criteria:", m4.ID)
	_, err := backend.Criteria{IDGlob: id.FromString("litmus"), StyleGlob: id.FromString("rmem")}.Find(specs)
	fmt.Println("unmatchable criteria:", err)
	_, err = backend.Criteria{IDGlob: id.FromString("litmus"), StyleGlob: id.FromString("*.rmem.*")}.Find(specs)
	fmt.Println("malformed criteria:", err)

	// Output:
	// empty criteria: herd
	// litmus criteria: litmus
	// litmus.dev criteria: litmus.dev
	// rmem criteria: rmem
	// unmatchable criteria: no matching backend found: id=litmus, style=rmem
	// malformed criteria: malformed glob expression: more than one '*' character
}
