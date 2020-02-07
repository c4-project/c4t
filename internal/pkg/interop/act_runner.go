package interop

import (
	"io"
	"os/exec"
)

const (
	// BinActCompiler is the name of the ACT compiler services binary.
	BinActCompiler = "act-compiler"
)

// ActRunner stores information about how to run the core ACT binaries.
type ActRunner struct {
	// DuneExec toggles whether ACT should be run through dune.
	DuneExec bool

	// ConfFile is the path to the act.conf to use.
	// If missing, we use ACT's default.
	ConfFile string
}

// StandardArgs captures the ACT 'standard args'.
type StandardArgs struct {
	// Whether verbosity is enabled.
	Verbose bool
}

func (s StandardArgs) ToArgv() []string {
	var argv []string
	if s.Verbose {
		argv = append(argv, "-v")
	}
	return argv
}

// Run runs an ACT command in a blocking manner.
func (a *ActRunner) Run(cmd string, stdin io.Reader, stdout, stderr io.Writer, sargs StandardArgs, argv ...string) error {
	fullArgv := append(sargs.ToArgv(), argv...)

	c := a.Command(cmd, fullArgv...)
	c.Stdin = stdin
	c.Stdout = stdout
	c.Stderr = stderr

	return c.Run()
}

func (a *ActRunner) Command(cmd string, argv ...string) *exec.Cmd {
	if a.DuneExec {
		duneArgv := append([]string{"exec", cmd, "--"}, argv...)
		return exec.Command("dune", duneArgv...)
	}
	return exec.Command(cmd, argv...)
}
