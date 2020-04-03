#!/bin/bash
# Copyright (c) 2020 Matt Windsor and contributors
#
# This file is part of act-tester.
# Licenced under the MIT licence; see `LICENSE`.

# Expects to be run in an ACT working directory.
# Runs all of the various single-shot binaries of act-tester in a pipeline, outputting to ./tester_test.

# TODO(@MattWindsor91): make this parametric and more robust.

set -euo pipefail

rm -rf tester_test || true
mkdir tester_test

act-tester-plan -m localhost examples/c_litmus/memalloy/*.litmus |
	tee tester_test/plan.toml |
act-tester-fuzz -d tester_test/fuzz -n 10 |
	tee tester_test/plan.fuzzed.toml |
act-tester-lift -d tester_test/lift |
	tee tester_test/plan.lifted.toml |
act-tester-mach -d tester_test/compile > tester_test/plan.out.toml
