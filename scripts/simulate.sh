#!/usr/bin/env bash
# Copyright (c) 2020 Matt Windsor and contributors
#
# This file is part of act-tester.
# Licenced under the MIT licence; see `LICENSE`.

# Simulates a list of C litmus tests without fuzzing.

set -euo pipefail

OUTPUT_DIR="$(mktemp -d)" || exit 2
trap 'rm -rf ${OUTPUT_DIR}' EXIT

act-tester-plan -m localhost "$@" |
act-tester-lift -d "${OUTPUT_DIR}/lift" |
act-tester-mach -d "${OUTPUT_DIR}/compile"
