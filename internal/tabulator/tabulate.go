// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

// Package tabulate provides a generic interface for tabulating data (to CSV, to tabwriter, etc).
package tabulator

import (
	"fmt"
	"strconv"
)

// Tabulator is the interface of things that can create tables.
type Tabulator interface {
	// Header adds a header to the current table, if one has not yet been set.
	// Tabulators should silently ignore attempts to install a duplicate header; this means that sub-table functions
	// can call Header even if the parent table calls Header.
	Header(labels ...string)

	RowTabulator

	// Flush commits the table to the underlying writer, returning any error that occurred during tabulation.
	Flush() error
}

// RowTabulator is the interface of things that can tabulate a row.
// Each method should return the current or parent tabulator for method chaining.
type RowTabulator interface {
	// Cell adds some representation of the given cell value to the table's current row.
	Cell(value interface{}) RowTabulator

	// EndRow ends the current row.
	EndRow() Tabulator
}

// LineWriter can write table lines to some underlying flushable writer.
type LineWriter interface {
	// Write writes a single row with the given cells.
	Write(cells []string) error

	// Flush commits any lines written to the writer.
	Flush() error
}

// LineTabulator uses a LineWriter to tabulate a table.
type LineTabulator struct {
	w      LineWriter
	ncells int
	err    error
	row    []string
}

// NewLineTabulator constructs a new tabulator using a LineWriter w.
func NewLineTabulator(w LineWriter) *LineTabulator {
	return &LineTabulator{w: w}
}

func (t *LineTabulator) Header(labels ...string) {
	if 0 < t.ncells || t.err != nil {
		return
	}
	t.ncells = len(labels)
	t.resetRow()
	for _, l := range labels {
		t.Cell(l)
	}
	_ = t.EndRow()
}

func (t *LineTabulator) Cell(value interface{}) RowTabulator {
	return t.stringCell(stringify(value))
}

func stringify(value interface{}) string {
	switch v := value.(type) {
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	case string:
		return v
	case fmt.Stringer:
		return v.String()
	default:
		return "?"
	}
}

func (t *LineTabulator) stringCell(value string) RowTabulator {
	if t.err == nil {
		t.row = append(t.row, value)
	}
	return t
}

func (t *LineTabulator) EndRow() Tabulator {
	if t.err != nil {
		return t
	}
	t.err = t.w.Write(t.row)
	t.resetRow()

	return t
}

func (t *LineTabulator) resetRow() {
	t.row = make([]string, 0, t.ncells)
}

func (t LineTabulator) Flush() error {
	if t.err != nil {
		return t.err
	}
	return t.w.Flush()
}
