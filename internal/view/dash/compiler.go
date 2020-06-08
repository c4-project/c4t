// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package dash

import (
	"fmt"

	"github.com/MattWindsor91/act-tester/internal/model/service/compiler"
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/widgets/text"
)

// OnCompilerPlanSet prepares for receiving compiler plans by clearing out any existing compilers shown on the dash.
func (o *Observer) OnCompilerPlanStart(_ int) {
	o.compilers.Reset()
}

// OnCompilerPlan outputs compiler information to this instance's compiler log.
func (o *Observer) OnCompilerPlan(c compiler.Named) {
	opts := text.WriteCellOpts(cell.FgColor(optColour(c.SelectedOpt)))
	if err := o.compilers.Write(fmt.Sprintf("%s: ", c.ID), opts); err != nil {
		o.logError(err)
	}
	if err := o.compilers.Write(c.String() + "\n"); err != nil {
		o.logError(err)
	}
}

// OnCompilerPlanFinish does nothing.
func (o *Observer) OnCompilerPlanFinish() {
}
