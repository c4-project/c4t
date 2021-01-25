// Copyright (c) 2021 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package stat_test

import (
	"encoding/csv"
	"fmt"
	"os"

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
						Selections: map[mutation.Mutant]uint64{2: 10, 3: 4},
						Hits:       map[mutation.Mutant]uint64{2: 1000, 3: 999},
						Kills:      map[mutation.Mutant]uint64{3: 2},
					},
				},
				Total: stat.MachineSpan{
					Mutation: mutation.Statset{
						Selections: map[mutation.Mutant]uint64{1: 1000, 2: 100, 3: 40},
						Hits:       map[mutation.Mutant]uint64{1: 3000, 2: 2000, 3: 5000},
						Kills:      map[mutation.Mutant]uint64{3: 20},
					},
				},
			},
			"bar": {
				Total: stat.MachineSpan{
					Mutation: mutation.Statset{
						Selections: map[mutation.Mutant]uint64{1: 500},
						Hits:       map[mutation.Mutant]uint64{},
						Kills:      map[mutation.Mutant]uint64{},
					},
				},
			},
		},
	}
	_ = s.DumpMutationCSV(csv.NewWriter(os.Stdout), false)
	fmt.Println("--")
	_ = s.DumpMutationCSV(csv.NewWriter(os.Stdout), true)

	// Output:
	// foo,2,10,1000,0
	// foo,3,4,999,2
	// --
	// bar,1,500,0,0
	// foo,1,1000,3000,0
	// foo,2,100,2000,0
	// foo,3,40,5000,20
}
