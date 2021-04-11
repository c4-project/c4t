// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package litmus

import (
	"fmt"
	"io"

	"github.com/c4-project/c4t/internal/model/litmus"
)

// Fixset contains the various fix-ups that the litmus tool needs to do before and after the run to Litmus.
type Fixset struct {
	// InjectStdbool is true when the Litmus header needs includes to 'stdbool.h' enabling.
	InjectStdbool bool

	// UseAsCall is true when Litmus needs the '-ascall' flag enabling.
	UseAsCall bool

	// RemoveAtomicCasts is true when casts to (_Atomic int) in the outcome dumper must be removed.
	RemoveAtomicCasts bool
}

// Args converts the part of the fixset that relates to litmus7 arguments into an argument slice to send to litmus7.
func (f *Fixset) Args() []string {
	var args []string

	if f.UseAsCall {
		args = append(args, "-ascall", "true")
	}

	return args
}

// PopulateFromLitmus switches various fixes on as appropriate for the Litmus file l.
func (f *Fixset) PopulateFromLitmus(l *litmus.Litmus) {
	if l != nil && l.IsC() {
		f.PopulateFromStats(l.Stats)
	}
}

// PopulateFromStats switches various fixes on according to the statistics in s.
func (f *Fixset) PopulateFromStats(s *litmus.Statset) {
	if s == nil {
		return
	}

	// TODO(@MattWindsor91): this should only be turned on if atomic integers are present.
	// Even then, it should only appear when we're using `gcc`, but I'm unsure how to enforce that.
	f.RemoveAtomicCasts = true

	if 0 < s.LiteralBools {
		f.InjectStdbool = true
	}
	if 0 < s.Returns {
		f.UseAsCall = true
	}
}

// Dump dumps a human-readable description of the fixset to the given writer.
func (f *Fixset) Dump(w io.Writer) error {
	for _, c := range []struct {
		field string
		on    bool
	}{
		{"injecting stdbool", f.InjectStdbool},
		{"using -ascall", f.UseAsCall},
		{"removing _Atomic casts", f.RemoveAtomicCasts},
	} {
		if !c.on {
			continue
		}
		if _, err := fmt.Fprintln(w, c.field); err != nil {
			return err
		}
	}
	return nil
}

// NeedsPatch checks whether patching the litmus harness is necessary.
func (f *Fixset) NeedsPatch() bool {
	return f.RemoveAtomicCasts || f.InjectStdbool
}
