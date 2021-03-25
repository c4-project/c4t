// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package stdflag

const (
	// Short flags are registered here where possible, to make sure there are
	// no duplicates.

	// FlagCompiler is a standard flag for compiler selection arguments.
	FlagCompiler = "c"
	// FlagOutDir is a short flag for specifying an output directory.
	FlagOutDir = "d"
	// FlagWorkerCount is a standard flag for arguments that set a worker count.
	FlagWorkerCount   = "j"
	flagSubjectCycles = "k"
	// FlagMachine is a standard flag for machine selection arguments.
	FlagMachine = "m"
	// FlagNum is a standard flag for 'number of' arguments.
	FlagNum = "n"
	// FlagCompilerTimeout is a short flag for compiler timeout.
	FlagCompilerTimeout = "t"

	flagC4fDuneExec = "x"

	// FlagConfigFile is a standard flag for overriding the config file.
	FlagConfigFile = "C"
	// FlagAltWorkerCount is a flag for arguments that set a secondary worker count.
	FlagAltWorkerCount = "J"
	// FlagRunTimeout is a short flag for run timeout.
	FlagRunTimeout = "T"

	// FlagCPUProfile is a standard flag for specifying a CPU profile output.
	FlagCPUProfile = "cpuprofile"

	flagCorpusSize = "corpus-size"

	usageConfFile      = "read tester config from this `file`"
	usageCorpusSize    = "`number` of corpus files to select for this test plan"
	usageC4fDuneExec   = "if true, use 'dune exec' to run c4f binaries"
	usageOutDir        = "`directory` to which outputs will be written"
	usageSubjectFuzzes = "number of `times` to fuzz each subject in the corpus"
	usageCPUProfile    = "`file` into which we should dump pprof information"
)
