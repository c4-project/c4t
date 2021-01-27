// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package stat_test

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/c4-project/c4t/internal/subject/status"

	"github.com/c4-project/c4t/internal/mutation"
	"github.com/c4-project/c4t/internal/stat"
)

// ExampleSet_DumpMutationCSV is a runnable example for Set.DumpMutationCSV.
func ExampleSet_DumpMutationCSV() {
	s := stat.Set{
		Machines: map[string]stat.Machine{
			"foo": {
				Session: stat.MachineSpan{
					Mutation: mutation.Statset{
						ByMutant: map[mutation.Mutant]mutation.MutantStatset{
							2:  {Selections: 1, Hits: 0, Kills: 0, Statuses: map[status.Status]uint64{status.Filtered: 1}},
							42: {Selections: 10, Hits: 1, Kills: 0, Statuses: map[status.Status]uint64{status.Ok: 9, status.CompileTimeout: 1}},
							53: {Selections: 20, Hits: 400, Kills: 15, Statuses: map[status.Status]uint64{status.Flagged: 15, status.CompileFail: 3, status.RunFail: 2}},
						},
					},
				},
				Total: stat.MachineSpan{
					Mutation: mutation.Statset{
						ByMutant: map[mutation.Mutant]mutation.MutantStatset{
							2:  {Selections: 41, Hits: 5000, Kills: 40, Statuses: map[status.Status]uint64{status.Flagged: 40, status.Filtered: 1}},
							42: {Selections: 100, Hits: 1, Kills: 0, Statuses: map[status.Status]uint64{status.Ok: 99, status.CompileTimeout: 1}},
							53: {Selections: 20, Hits: 400, Kills: 15, Statuses: map[status.Status]uint64{status.Flagged: 15, status.CompileFail: 3, status.RunFail: 2}},
						},
					},
				},
			},
			"bar": {
				Total: stat.MachineSpan{
					Mutation: mutation.Statset{
						ByMutant: map[mutation.Mutant]mutation.MutantStatset{
							1: {Selections: 500, Hits: 0, Kills: 0, Statuses: map[status.Status]uint64{status.Ok: 500}},
						},
					},
				},
			},
		},
	}

	w := csv.NewWriter(os.Stdout)
	_ = s.DumpMutationCSVHeader(w)
	_ = s.DumpMutationCSV(w, false)
	fmt.Println("--")
	_ = s.DumpMutationCSV(w, true)

	// Output:
	// Machine,Mutant,Selections,Hits,Kills,Ok,Filtered,Flagged,CompileFail,CompileTimeout,RunFail,RunTimeout
	// foo,2,1,0,0,0,1,0,0,0,0,0
	// foo,42,10,1,0,9,0,0,0,1,0,0
	// foo,53,20,400,15,0,0,15,3,0,2,0
	// --
	// bar,1,500,0,0,500,0,0,0,0,0,0
	// foo,2,41,5000,40,0,1,40,0,0,0,0
	// foo,42,100,1,0,99,0,0,0,1,0,0
	// foo,53,20,400,15,0,0,15,3,0,2,0
}
