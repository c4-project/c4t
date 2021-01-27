// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package directorobs_test

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/c4-project/c4t/internal/director"

	"github.com/c4-project/c4t/internal/model/service/compiler"
	"github.com/c4-project/c4t/internal/model/service/compiler/optlevel"

	"github.com/stretchr/testify/require"

	"github.com/c4-project/c4t/internal/ux/directorobs"

	"github.com/c4-project/c4t/internal/stage/analyser/saver"

	"github.com/c4-project/c4t/internal/helper/iohelp"
	"github.com/c4-project/c4t/internal/model/id"
)

// ExampleLogger_OnPrepare is a runnable example indirectly exercising Logger.OnPrepare.
func ExampleLogger_OnPrepare() {
	l, _ := directorobs.NewLogger(iohelp.NopWriteCloser{Writer: os.Stdout}, 0)
	r, _ := directorobs.NewForwardObserver(0, l)

	director.OnPrepare(director.PrepareInstancesMessage(5), r)

	// Output:
	// running on 5 instances
}

// ExampleLogger_OnCycle is a runnable example indirectly exercising Logger.OnCycle.
func ExampleLogger_OnCycle() {
	l, _ := directorobs.NewLogger(iohelp.NopWriteCloser{Writer: os.Stdout}, 0)
	r, _ := directorobs.NewForwardObserver(0, l)
	i, _ := r.Instance(id.FromString("localhost"))

	c := director.Cycle{
		Instance:  0,
		MachineID: id.FromString("localhost"),
		Iter:      10,
		Start:     time.Time{},
	}

	go func() {
		// These messages will arrive through l.OnCycle.
		director.OnCycle(director.CycleStartMessage(c), i)
		director.OnCycle(director.CycleErrorMessage(c, errors.New("the front fell off")), i)
		// Important, else the logger will keep waiting for the instance to provide observations.
		i.OnInstance(director.InstanceClosedMessage())
	}()
	_ = r.Run(context.Background())

	// Output:
	// * localhost starts cycle 10 *
	// * localhost ERROR: the front fell off *
	// [instance 0 has closed]
}

// ExampleLogger_OnCycleSave is a runnable example indirectly exercising Logger.OnCycleSave.
func ExampleLogger_OnCycleSave() {
	l, _ := directorobs.NewLogger(iohelp.NopWriteCloser{Writer: os.Stdout}, 0)
	r, _ := directorobs.NewForwardObserver(0, l)
	i, _ := r.Instance(id.FromString("localhost"))

	go func() {
		// These messages will arrive through l.OnCycleSave.
		saver.OnArchiveStart("subj", "subj.tar.gz", 2, i)
		saver.OnArchiveFileAdded("subj", "a.out", 0, i)
		saver.OnArchiveFileMissing("subj", "compile.log", 1, i)
		saver.OnArchiveFinish("subj", i)
		// Important, else the logger will keep waiting for the instance to provide observations.
		i.OnInstance(director.InstanceClosedMessage())
	}()
	_ = r.Run(context.Background())

	// Output:
	// saving (cycle [0:  #0 (Jan  1 00:00:00)]) subj to subj.tar.gz
	// when saving (cycle [0:  #0 (Jan  1 00:00:00)]) subj: missing file compile.log
	// [instance 0 has closed]
}

// ExampleLogger_OnCycleCompiler is a runnable example indirectly exercising Logger.OnCycleCompiler.
func ExampleLogger_OnCycleCompiler() {
	l, _ := directorobs.NewLogger(iohelp.NopWriteCloser{Writer: os.Stdout}, 0)
	r, _ := directorobs.NewForwardObserver(0, l)
	i, _ := r.Instance(id.FromString("localhost"))

	go func() {
		// These messages will arrive through l.OnCycleCompiler.
		compiler.OnCompilerConfigStart(3, i)
		compiler.OnCompilerConfigStep(0,
			compiler.Named{
				ID: id.FromString("gcc.4"),
				Instance: compiler.Instance{
					SelectedMOpt: "arch=native",
					SelectedOpt:  &optlevel.Named{Name: "3", Level: optlevel.Level{}},
					Compiler:     compiler.Compiler{Style: id.CStyleGCC, Arch: id.ArchArm7},
				},
			}, i)
		compiler.OnCompilerConfigStep(1,
			compiler.Named{
				ID: id.FromString("gcc.9"),
				Instance: compiler.Instance{
					SelectedMOpt: "arch=skylake",
					SelectedOpt:  &optlevel.Named{Name: "2", Level: optlevel.Level{}},
					Compiler:     compiler.Compiler{Style: id.CStyleGCC, Arch: id.ArchArm8},
				},
			}, i)
		compiler.OnCompilerConfigStep(2,
			compiler.Named{
				ID: id.FromString("msvc"),
				Instance: compiler.Instance{
					Compiler: compiler.Compiler{Style: id.FromString("msvc"), Arch: id.ArchX8664},
				},
			}, i)
		compiler.OnCompilerConfigEnd(i)
		// Important, else the logger will keep waiting for the instance to provide observations.
		i.OnInstance(director.InstanceClosedMessage())
	}()
	_ = r.Run(context.Background())

	// Output:
	// [0:  #0 (Jan  1 00:00:00)] compilers 3:
	// - gcc.4: gcc@arm.7 opt "3" march "arch=native"
	// - gcc.9: gcc@arm.8 opt "2" march "arch=skylake"
	// - msvc: msvc@x86.64
	// [instance 0 has closed]
}

// TestLogger_Run_empty tests that running a logger with no attached instances works out.
func TestLogger_Run_empty(t *testing.T) {
	t.Parallel()

	l, err := directorobs.NewLogger(iohelp.DiscardCloser(), 0)
	require.NoError(t, err, "logger should construct without errors")
	r, err := directorobs.NewForwardObserver(0, l)
	require.NoError(t, err, "forwarder should construct without errors")
	err = r.Run(context.Background())
	require.NoError(t, err, "no channels = no error")
}

// TestLogger_Run_noMessages tests that running a logger with no messages, works out.
func TestLogger_Run_noMessages(t *testing.T) {
	t.Parallel()

	l, err := directorobs.NewLogger(iohelp.DiscardCloser(), 0)
	require.NoError(t, err, "logger should construct without errors")
	r, err := directorobs.NewForwardObserver(0, l)
	require.NoError(t, err, "forwarder should construct without errors")
	i, err := r.Instance(id.FromString("foo"))
	require.NoError(t, err, "instance should construct without errors")
	go func() {
		i.OnInstance(director.InstanceClosedMessage())
	}()
	err = r.Run(context.Background())
	require.NoError(t, err, "should have stopped running with no error")
}
