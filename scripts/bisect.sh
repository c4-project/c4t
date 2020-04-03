#!/usr/bin/env bash

set -euo pipefail

SCRIPTDIR="${SCRIPTDIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"}"
readonly SCRIPTDIR

if [[ $# -ne 2 ]]; then
    echo "usage: $0 TRACE ORIGINAL" >&2
    exit 2
fi

trace="$1"
original="$2"

# Try to remove some entropy from the bisection process
export ACT_SEED=0

# TODO(@MattWindsor91): this is working around a failure in the litmusifying process when a litmus file has a name
# of the form foo.bar.baz (the output C file becomes foo.c).
OUTPUT_DIR="$(mktemp -d)" || exit 2
trap 'rm -rf ${OUTPUT_DIR}' EXIT

act-fuzz replay -trace "${trace}" "${original}" -o "${OUTPUT_DIR}/in.litmus"

! "${SCRIPTDIR}/simulate.sh" "${OUTPUT_DIR}/in.litmus" | grep "flagged"
