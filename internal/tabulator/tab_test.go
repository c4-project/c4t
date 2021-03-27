// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package tabulator_test

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/c4-project/c4t/internal/helper/iohelp"
	"github.com/stretchr/testify/assert"

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

// TestNewTab_error tests that tabbing tabulators handle errors correctly.
func TestNewTab_error(t *testing.T) {
	t.Parallel()
	err := errors.New("no u")
	w := tabulator.NewTab(iohelp.ErrWriter{Err: err})
	w.Cell("USA").Cell(1).EndRow()
	w.Cell("UK").Cell(int64(44)).EndRow()
	assert.ErrorIs(t, w.Flush(), err)
}
