// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

// Package csv handles outputting of analysis data as CSVs.
package csv

import (
	"encoding/csv"
	"fmt"
	"strconv"
	"time"

	"github.com/MattWindsor91/act-tester/internal/model/service/compiler"
	"github.com/MattWindsor91/act-tester/internal/model/status"
	"github.com/MattWindsor91/act-tester/internal/plan/analyser"
)

// CompilerWriter wraps a CSV writer and makes it output compiler analyses.
type CompilerWriter csv.Writer

// OnAnalysis observes an analysis by emitting a CSV with compiler information.
func (c *CompilerWriter) OnAnalysis(a analyser.Analysis) {
	c.writeHeader()
	for cname, can := range a.Compilers {
		c.writeCompiler(cname, can)
	}
	(*csv.Writer)(c).Flush()
}

var staticColumnHeaders = [...]string{
	"compilerID",
	"styleID",
	"archID",
	"opt",
	"mopt",
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
	return []string{"min" + name, "avg" + name, "max" + name}
}

func (c *CompilerWriter) writeCompiler(cname string, can analyser.Compiler) {
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

func (c *CompilerWriter) staticColumnsForCompiler(cname string, can analyser.Compiler) [len(staticColumnHeaders)]string {
	return [...]string{
		cname,
		can.Info.Style.String(),
		can.Info.Arch.String(),
		optName(can.Info),
		can.Info.SelectedMOpt,
	}
}

func counts(cs map[status.Status]int) []string {
	result := make([]string, status.Last+1)
	for i := status.Ok; i <= status.Last; i++ {
		result[i] = strconv.Itoa(cs[i])
	}
	return result
}

func timeset(ts *analyser.TimeSet) []string {
	return []string{
		duration(ts.Min),
		duration(ts.Mean()),
		duration(ts.Max),
	}
}

func duration(d time.Duration) string {
	return fmt.Sprint(d.Seconds())
}

func optName(i compiler.Compiler) string {
	if i.SelectedOpt == nil {
		return ""
	}
	return i.SelectedOpt.Name
}
