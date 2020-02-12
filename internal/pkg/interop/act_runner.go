package interop

import (
	"os/exec"
)

// ActRunner stores information about how to run the core ACT binaries.
type ActRunner struct {
	// DuneExec toggles whether ACT should be run through dune.
	DuneExec bool

	// ConfFile is the path to the act.conf to use.
	// If missing, we use ACT's default.
	ConfFile string
}

// StandardArgs captures the ACT 'standard arguments', less those covered by ActRunner itself.
type StandardArgs struct {
	// Whether verbosity is enabled.
	Verbose bool
}

// ToArgv converts s to an argument vector fragment.
func (s StandardArgs) ToArgv() []string {
	var argv []string
	if s.Verbose {
		argv = append(argv, "-v")
	}
	return argv
}

// Command constructs a Cmd for running the ACT command cmd with subcommand sub and arguments argv.
func (a ActRunner) Command(cmd, sub string, sargs StandardArgs, argv ...string) *exec.Cmd {
	fargv := a.actArgv(sub, sargs, argv)
	dcmd, dargv := liftDuneExec(a.DuneExec, cmd, fargv)
	return exec.Command(dcmd, dargv...)
}

func (a ActRunner) actArgv(sub string, sargs StandardArgs, argv []string) []string {
	sargv := sargs.ToArgv()

	// Reserving room for the subcommand, and optionally '-config FOO'.
	fargv := make([]string, 1, 3+len(sargv)+len(argv))
	fargv[0] = sub
	fargv = append(fargv, sargs.ToArgv()...)

	if a.ConfFile != "" {
		fargv = append(fargv, "-config", a.ConfFile)
	}

	return append(fargv, argv...)
}

func liftDuneExec(duneExec bool, cmd string, argv []string) (string, []string) {
	if duneExec {
		return "dune", append([]string{"exec", cmd, "--"}, argv...)
	}
	return cmd, argv
}
