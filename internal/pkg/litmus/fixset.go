package litmus

import (
	"fmt"
	"io"

	"github.com/MattWindsor91/act-tester/internal/pkg/interop"
)

// Fixset contains the various fix-ups that the litmus tool needs to do before and after the run to Litmus.
type Fixset struct {
	// UseAsCall is true when Litmus needs the '-ascall' flag enabling.
	UseAsCall bool

	// InjectStdbool is true when the Litmus header needs includes to 'stdbool.h' enabling.
	InjectStdbool bool
}

// Args converts the part of the fixset that relates to litmus7 arguments into an argument slice to send to litmus7.
func (f *Fixset) Args() []string {
	var args []string

	if f.UseAsCall {
		args = append(args, "-ascall", "true")
	}

	return args
}

// PopulateFromStats switches various fixes on according to the statistics in s.
func (f *Fixset) PopulateFromStats(s *interop.Statset) {
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
		{"using -ascall", f.UseAsCall},
		{"injecting stdbool", f.InjectStdbool},
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
