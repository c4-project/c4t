// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package timing

import (
	"fmt"
	"time"
)

// Span is a pair of times, representing the start and end of a process.
type Span struct {
	// Start is the start time of the timespan.
	Start time.Time `json:"start,omitempty"`
	// End is the end time of the timespan.
	End time.Time `json:"end,omitempty"`
}

// SpanFromInstant constructs a Span representing the instant in time t.
func SpanFromInstant(t time.Time) Span {
	return Span{
		Start: t,
		End:   t,
	}
}

// SpanFromDuration constructs a Span from the start time start and duration d.
func SpanFromDuration(start time.Time, d time.Duration) Span {
	return Span{
		Start: start,
		End:   start.Add(d),
	}
}

// SpanSince constructs a Span representing the period of time between start and now.
func SpanSince(start time.Time) Span {
	return Span{
		Start: start,
		End:   time.Now(),
	}
}

// Duration gets the duration of the timespan.
// Duration is zero if either of the ends of the timespan are zero.
func (t *Span) Duration() time.Duration {
	if t.IsUndefined() {
		return 0
	}
	return t.End.Sub(t.Start)
}

// IsUndefined gets whether this span is ill-defined.
// This happens if either time is zero, or the end is before the start.
func (t *Span) IsUndefined() bool {
	return t.Start.IsZero() || t.End.IsZero() || t.End.Before(t.Start)
}

// IsInstant gets whether this
func (t *Span) IsInstant() bool {
	return t.Start.Equal(t.End)
}

// String retrieves a string representation of this timespan.
func (t Span) String() string {
	switch {
	case t.IsUndefined():
		return "(undefined)"
	case t.IsInstant():
		return t.Start.Format(time.RFC3339)
	default:
		return fmt.Sprintf("%s (from %s to %s)", t.Duration(), t.Start.Format(time.RFC3339), t.End.Format(time.RFC3339))
	}
}

// StartNoLaterThan moves this timespan's start to become start, if start is nonzero and before this timespan's current start.
func (t *Span) StartNoLaterThan(start time.Time) {
	if t.Start.IsZero() || (!start.IsZero() && start.Before(t.Start)) {
		t.Start = start
	}
}

// EndNoEarlierThan moves this timespan's end to become end, if end is nonzero and after this timespan's current end.
func (t *Span) EndNoEarlierThan(end time.Time) {
	if t.End.IsZero() || (!end.IsZero() && end.After(t.End)) {
		t.End = end
	}
}

// Union sets this timespan's start to other's start if it is earlier, and its end to other's end if it is later.
func (t *Span) Union(other Span) {
	t.StartNoLaterThan(other.Start)
	t.EndNoEarlierThan(other.End)
}

// MockDate is an arbitrary date useful for timespan tests.
var MockDate = time.Date(1997, time.May, 1, 21, 0, 0, 0, time.UTC)
