// Copyright (c) 2021 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package dash

import (
	"fmt"
	"strconv"

	"github.com/c4-project/c4t/internal/director"
	"github.com/mum4k/termdash/cell"

	"github.com/mum4k/termdash/widgets/text"
)

// maxLogLines is the maximum number of lines that can be written to the log before it resets.
const maxLogLines = 1000

// syslog is the dashboard's system log.
type syslog struct {
	log *text.Text
	// nlines is the number of lines written, so far, to the log.
	nlines uint
}

// newSysLog constructs a new system log.
func newSysLog() (*syslog, error) {
	log, err := text.New(text.RollContent(), text.WrapAtWords())
	if err != nil {
		return nil, err
	}
	return &syslog{log: log, nlines: 0}, nil
}

// reportPrepare reports a director preparation message.
func (s *syslog) reportPrepare(m director.PrepareMessage) {
	if m.Kind != director.PrepareInstances {
		return
	}
	s.write(fmt.Sprintln("Instances:", strconv.Itoa(m.NumInstances)))
}

// reportCycleError logs a cycle error to the system log.
func (s *syslog) reportCycleError(cycle director.Cycle, err error) {
	s.write(fmt.Sprintf("ERROR on %s:\n%s\n", cycle, err), text.WriteCellOpts(cell.FgColor(cell.ColorMaroon)))
}

// write writes text to syslog, using options opts.
func (s *syslog) write(text string, opts ...text.WriteOption) {
	s.nlines += countNewlines(text)
	if maxLogLines <= s.nlines {
		s.nlines -= maxLogLines
		s.log.Reset()
		_ = s.log.Write("[log reset]")
	}

	if err := s.log.Write(text, opts...); err != nil {
		s.logError(err)
	}
}

// logError logs an error to syslog, giving up if it can't.
func (s *syslog) logError(err error) {
	// not s.Write, else we'd get a possibly infinite loop
	_ = s.log.Write(err.Error(), text.WriteCellOpts(cell.FgColor(cell.ColorRed)))
}
