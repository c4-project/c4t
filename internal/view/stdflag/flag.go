// Copyright (c) 2020 Matt Windsor and contributors
//
// This file is part of act-tester.
// Licenced under the MIT licence; see `LICENSE`.

package stdflag

const (
	// FlagRunTimeout is a short flag for run timeout.
	FlagRunTimeout  = "T"
	flagActConfFile = "A"
	flagConfFile    = "C"
	// FlagUseJSON is a short flag for enabling JSON output where available.
	FlagUseJSON = "J"
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
	flagActDuneExec     = "x"

	// FlagCPUProfile is a standard flag for specifying a CPU profile output.
	FlagCPUProfile = "cpuprofile"

	FlagSkipCompiler = "skip-compiler"
	FlagSkipRunner   = "skip-runner"

	// FlagCompilerTimeoutLong is a long flag for compiler timeout.
	FlagCompilerTimeoutLong = "compiler-timeout"
	// FlagRunTimeoutLong is a long flag for run timeout.
	FlagRunTimeoutLong = "run-timeout"
	// FlagUseJSONLong is a long flag for JSON emission.
	FlagUseJSONLong = "emit-json"
	// FlagWorkerCountLong is a long flag for arguments that set a worker count.
	FlagWorkerCountLong = "num-workers"

	usageActConfFile = "read ACT config from this `file`"
	usageCorpusSize  = "`number` of corpus files to select for this test plan;\n" +
		"if non-positive, the planner will use all viable provided corpus files"
	usageActDuneExec   = "if true, use 'dune exec' to run OCaml ACT binaries"
	usageOutDir        = "`directory` to which outputs will be written"
	usageSubjectCycles = "number of `cycles` to run for each subject in the corpus"
	usageCPUProfile    = "`file` into which we should dump pprof information"
)
