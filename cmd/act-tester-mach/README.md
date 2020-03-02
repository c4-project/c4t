# act-tester-mach

`act-tester-mach` performs the parts of an ACT test run that are machine-dependent.
By default, it accepts a lifted, potentially-fuzzed plan and:

- compiles each lifted test harness to a binary using each planned compiler;
- runs each compiled binary;
- returns a plan with the results of said runs.

## Usage

`act-tester-mach [-c] [-r] -i PLAN`

- `-c` disables the compiling phase.
- `-r` disables the running phase.
