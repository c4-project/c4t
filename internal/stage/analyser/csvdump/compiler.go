// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package csvdump

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/c4-project/c4t/internal/model/service/compiler"
	"github.com/c4-project/c4t/internal/plan/analysis"
	"github.com/c4-project/c4t/internal/subject/status"
)

// CompilerWriter wraps a CSV writer and makes it output compiler analyses.
type CompilerWriter csv.Writer

// NewCompilerWriter creates a new compiler writer over w.
func NewCompilerWriter(w io.Writer) *CompilerWriter {
	return (*CompilerWriter)(csv.NewWriter(w))
}

// OnAnalysis observes an analysis by emitting a CSV with compiler information.
func (c *CompilerWriter) OnAnalysis(a analysis.Analysis) {
	c.writeHeader()
	for cname, can := range a.Compilers {
		c.writeCompiler(cname, can)
	}
	(*csv.Writer)(c).Flush()
}

var staticColumnHeaders = [...]string{
	"CompilerID",
	"StyleID",
	"ArchID",
	"Opt",
	"MOpt",
}

func (c *CompilerWriter) writeHeader() {
	rec := staticColumnHeaders[:]
	rec = append(rec, timesetHeader("Compile")...)
	rec = append(rec, timesetHeader("Run")...)
	for i := status.Ok; i <= status.Last; i++ {
		rec = append(rec, i.String())
	}
	c.write(rec)
}

func timesetHeader(name string) []string {
	return []string{"Min" + name, "Avg" + name, "Max" + name}
}

func (c *CompilerWriter) writeCompiler(cname string, can analysis.Compiler) {
	scs := c.staticColumnsForCompiler(cname, can)
	rec := scs[:]
	rec = append(rec, timeset(can.Time)...)
	rec = append(rec, timeset(can.RunTime)...)
	rec = append(rec, counts(can.Counts)...)
	c.write(rec)
}

func (c *CompilerWriter) write(record []string) {
	_ = (*csv.Writer)(c).Write(record)
}

func (c *CompilerWriter) staticColumnsForCompiler(cname string, can analysis.Compiler) [len(staticColumnHeaders)]string {
	return [...]string{
		cname,
		can.Info.Style.String(),
		can.Info.Arch.String(),
		optName(can.Info),
		can.Info.SelectedMOpt,
	}
}

func counts(cs map[status.Status]int) []string {
	result := make([]string, status.Last)
	for i := status.Ok; i <= status.Last; i++ {
		result[i-1] = strconv.Itoa(cs[i])
	}
	return result
}

func timeset(ts *analysis.TimeSet) []string {
	return []string{
		duration(ts.Min),
		duration(ts.Mean()),
		duration(ts.Max),
	}
}

func duration(d time.Duration) string {
	return fmt.Sprint(d.Seconds())
}

func optName(i compiler.Instance) string {
	if i.SelectedOpt == nil {
		return ""
	}
	return i.SelectedOpt.Name
}
