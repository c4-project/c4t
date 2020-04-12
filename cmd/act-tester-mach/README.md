# act-tester-mach

`act-tester-mach` performs the parts of an ACT test run that are machine-dependent.
By default, it accepts a lifted, potentially-fuzzed plan and:

- compiles each lifted test harness to a binary using each planned compiler;
- runs each compiled binary;
- returns a plan with the results of said runs.

## Usage

`act-tester-mach [FLAGS] -i PLAN`

### Flags

- `-compiler-timeout DURATION` sets the compiler timeout.
- `-run-timeout DURATION` sets the run timeout.
- `-num-workers N` sets the number of parallel run workers.

- `-emit-json` causes the machine runner to emit progress information in JSON,
  in the format expected by `act-tester-rmach` and `act-tester`.
- `-skip-compiler` disables the compiling phase.
- `-skip-runner` disables the running phase.
