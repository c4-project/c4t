// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package csvdump

import (
	"encoding/csv"
	"fmt"
	"io"
	"time"

	"github.com/c4-project/c4t/internal/plan/stage"

	"github.com/c4-project/c4t/internal/plan/analysis"
)

// StageWriter wraps a CSV writer and makes it output stage analyses.
type StageWriter csv.Writer

// NewStageWriter creates a new stage writer over w.
func NewStageWriter(w io.Writer) *StageWriter {
	return (*StageWriter)(csv.NewWriter(w))
}

// OnAnalysis observes an analysis by emitting a CSV with stage information.
func (s *StageWriter) OnAnalysis(a analysis.Analysis) {
	s.writeHeader()
	for _, rec := range a.Plan.Metadata.Stages {
		s.writeStage(rec)
	}
	(*csv.Writer)(s).Flush()
}

var columnHeaders = [...]string{
	"Stage",
	"CompletedAt",
	"Duration",
}

func (s *StageWriter) writeHeader() {
	s.write(columnHeaders[:])
}

func (s *StageWriter) writeStage(rec stage.Record) {
	s.write([]string{
		rec.Stage.String(),
		rec.CompletedOn.Format(time.RFC3339),
		fmt.Sprint(rec.Duration.Seconds()),
	})
}

func (s *StageWriter) write(record []string) {
	_ = (*csv.Writer)(s).Write(record)
}
