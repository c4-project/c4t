// Copyright (c) 2020-2021 C4 Project
//
// This file is part of c4t.
// Licenced under the MIT licence; see `LICENSE`.

package stdflag

const (
	// FlagCompiler is a standard flag for compiler selection arguments.
	FlagCompiler = "c"
	flagConfFile = "C"
	// FlagOutDir is a short flag for specifying an output directory.
	FlagOutDir = "d"
	// FlagWorkerCount is a standard flag for arguments that set a worker count.
	FlagWorkerCount = "j"
	// FlagAltWorkerCount is a flag for arguments that set a secondary worker count.
	FlagAltWorkerCount = "J"
	flagSubjectCycles  = "k"
	// FlagMachine is a standard flag for machine selection arguments.
	FlagMachine = "m"
	// FlagNum is a standard flag for 'number of' arguments.
	FlagNum = "n"
	// FlagCompilerTimeout is a short flag for compiler timeout.
	FlagCompilerTimeout = "t"
	// FlagRunTimeout is a short flag for run timeout.
	FlagRunTimeout  = "T"
	flagC4fDuneExec = "x"

	// FlagCPUProfile is a standard flag for specifying a CPU profile output.
	FlagCPUProfile = "cpuprofile"

	flagCorpusSize = "corpus-size"
	// FlagCompilerTimeoutLong is a long flag for compiler timeout.
	FlagCompilerTimeoutLong = "compiler-timeout"
	// FlagRunTimeoutLong is a long flag for run timeout.
	FlagRunTimeoutLong = "run-timeout"
	// FlagWorkerCountLong is a long flag for arguments that set a worker count.
	FlagWorkerCountLong = "num-workers"
	// FlagCompilerWorkerCountLong is a long flag for arguments that set a compiler worker count.
	FlagCompilerWorkerCountLong = "num-compiler-workers"
	// FlagRunWorkerCountLong is a long flag for arguments that set a runner worker count.
	FlagRunWorkerCountLong = "num-run-workers"

	usageConfFile      = "read tester config from this `file`"
	usageCorpusSize    = "`number` of corpus files to select for this test plan"
	usageC4fDuneExec   = "if true, use 'dune exec' to run c4f binaries"
	usageOutDir        = "`directory` to which outputs will be written"
	usageSubjectCycles = "number of `cycles` to run for each subject in the corpus"
	usageCPUProfile    = "`file` into which we should dump pprof information"
)
