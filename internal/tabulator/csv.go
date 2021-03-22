// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package tabulator

import (
	"encoding/csv"
	"io"
)

// NewCsv constructs a Tabulator that outputs CSV to w.
func NewCsv(w io.Writer) Tabulator {
	return NewLineTabulator((*CsvLineWriter)(csv.NewWriter(w)))
}

// CsvLineWriter wraps a csv.Writer to make it a LineWriter.
type CsvLineWriter csv.Writer

func (c *CsvLineWriter) Write(cells []string) error {
	return (*csv.Writer)(c).Write(cells)
}

func (c *CsvLineWriter) Flush() error {
	(*csv.Writer)(c).Flush()
	return (*csv.Writer)(c).Error()
}
