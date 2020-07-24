// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package csv

import (
	"encoding/csv"
	"fmt"
	"time"

	"github.com/MattWindsor91/act-tester/internal/plan/stage"

	"github.com/MattWindsor91/act-tester/internal/plan/analyser"
)

// StageWriter wraps a CSV writer and makes it output stage analyses.
type StageWriter csv.Writer

// OnAnalysis observes an analysis by emitting a CSV with stage information.
func (s *StageWriter) OnAnalysis(a analyser.Analysis) {
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
