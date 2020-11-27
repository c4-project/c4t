// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package quantity

import (
	"context"
	"log"
	"time"
)

// Timeout is a Duration with the semantics of non-positive values being an absence of timeout.
type Timeout time.Duration

// IsActive checks whether t represents a valid, active timeout.
func (t Timeout) IsActive() bool {
	return 0 < t
}

// Log dumps this timeout to the logger l, if it is active.
func (t Timeout) Log(l *log.Logger) {
	if t.IsActive() {
		l.Printf("timeout at %s", t)
	}
}

// OnContext lifts this timeout to a context with the given parent.
// If this timeout is active, the semantics are as context.WithTimeout; else, context.WithCancel.
func (t Timeout) OnContext(parent context.Context) (context.Context, context.CancelFunc) {
	if !t.IsActive() {
		return context.WithCancel(parent)
	}
	return context.WithTimeout(parent, time.Duration(t))
}

// String forwards the usual duration stringification for timeouts.
func (t Timeout) String() string {
	return time.Duration(t).String()
}

// MarshalText marshals a timeout by stringifying it.
func (t Timeout) MarshalText() (text []byte, err error) {
	return []byte(t.String()), nil
}

// UnmarshalText unmarshals a timeout by using ParseDuration.
func (t *Timeout) UnmarshalText(in []byte) error {
	d, err := time.ParseDuration(string(in))
	*t = Timeout(d)
	return err
}
