// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package tabulator_test

import (
	"fmt"
	"os"

	"github.com/c4-project/c4t/internal/tabulator"
)

// ExampleNewTab is a runnable example for NewTab.
func ExampleNewTab() {
	w := tabulator.NewTab(os.Stdout)
	w.Header("Country", "Code")
	w.Cell("USA").Cell(1).EndRow()
	w.Cell("UK").Cell(int64(44)).EndRow()
	if err := w.Flush(); err != nil {
		fmt.Println(err)
	}

	// Output:
	// Country  Code
	// USA      1
	// UK       44
}
