// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package dash

import (
	"time"

	"github.com/mum4k/termdash/container"

	"github.com/c4-project/c4t/internal/director"
	"github.com/mum4k/termdash/widgets/text"
)

// timeKeeper logs the start and end time of an experiment.
type timeKeeper struct {
	start, end *text.Text
}

func newTimeKeeper() (*timeKeeper, error) {
	tk := &timeKeeper{}
	if err := tk.initStart(); err != nil {
		return nil, err
	}
	if err := tk.initEnd(); err != nil {
		return nil, err
	}
	return tk, nil
}

func (tk *timeKeeper) initStart() error {
	var err error
	if tk.start, err = text.New(text.DisableScrolling()); err != nil {
		return err
	}
	return tk.start.Write("-- no start time --")
}

func (tk *timeKeeper) makePane() container.Option {
	return container.SplitHorizontal(
		container.Top(container.PlaceWidget(tk.start)),
		container.Bottom(container.PlaceWidget(tk.end)),
	)
}

func (tk *timeKeeper) initEnd() error {
	var err error
	if tk.end, err = text.New(text.DisableScrolling()); err != nil {
		return err
	}
	return tk.end.Write("-- no end time --")
}

func (tk *timeKeeper) OnPrepare(p director.PrepareMessage) {
	switch p.Kind {
	case director.PrepareStart:
		logTime(tk.start, "Start: ", p.Time)
	case director.PrepareTimeout:
		logTime(tk.end, "End:   ", p.Time)
	}
}

func logTime(t *text.Text, header string, tm time.Time) {
	_ = t.Write(header, text.WriteReplace())
	if err := t.Write(tm.Format(time.Stamp)); err != nil {
		_ = t.Write("ERROR")
	}
}
