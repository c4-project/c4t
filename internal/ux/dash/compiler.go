// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package dash

import (
	"fmt"

	"github.com/mum4k/termdash/container/grid"

	"github.com/mum4k/termdash/container"

	"github.com/MattWindsor91/act-tester/internal/observing"

	"github.com/MattWindsor91/act-tester/internal/model/service/compiler"
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/widgets/text"
)

const (
	headerCompilers = "Compilers"

	compilersContainerIDSuffix = "Compilers"
)

// compilersContainerID calculates the container ID of the compiler section of the machine at location i.
// This is used to rename the container to account for the number of compilers.
func (o *Instance) compilersContainerID() string {
	return o.id + compilersContainerIDSuffix
}

type compilerObserver struct {
	// text contains a readout of the currently planned compilers for this instance.
	text *text.Text

	// id is the ID of the compiler container, used to update the container's name.
	id string
	// parent is the parent container of the compiler, used to update its container's name.
	parent *container.Container
}

func newCompilerObserver(parent *container.Container, id string) (*compilerObserver, error) {
	var err error
	o := compilerObserver{parent: parent, id: id}
	if o.text, err = text.New(text.WrapAtWords()); err != nil {
		return nil, err
	}

	return &o, nil
}

// OnCompilerConfig forwards a build observation.
func (o *Instance) OnCompilerConfig(m compiler.Message) {
	if o.compilers == nil {
		return
	}
	if err := o.compilers.onCompilerConfig(m); err != nil {
		o.logError(err)
	}
}

func (o *compilerObserver) onCompilerConfig(m compiler.Message) error {
	switch m.Kind {
	case observing.BatchStart:
		return o.onPlanStart(m.Num)
	case observing.BatchStep:
		return o.onPlan(*m.Configuration)
	default:
		return nil
	}
}

// onPlanStart prepares for receiving compiler plans by clearing out any existing compilers shown on the dash.
func (o *compilerObserver) onPlanStart(ncompilers int) error {
	o.text.Reset()
	return o.updateName(ncompilers)
}

func (o *compilerObserver) updateName(ncompilers int) error {
	return o.parent.Update(o.id, container.BorderTitle(fmt.Sprintf("%s (%d)", headerCompilers, ncompilers)))
}

// onPlan outputs compiler information to this instance's compiler log.
func (o *compilerObserver) onPlan(c compiler.Named) error {
	opts := text.WriteCellOpts(cell.FgColor(optColour(c.SelectedOpt)))
	if err := o.text.Write(fmt.Sprintf("%s: ", c.ID), opts); err != nil {
		return err
	}
	return o.text.Write(c.String() + "\n")
}

func (o *compilerObserver) grid() []grid.Element {
	return []grid.Element{grid.Widget(o.text)}
}
