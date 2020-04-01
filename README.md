# `act-tester`

`act-tester` is the top-level testing framework
for the [ACT](https://github.com/MattWindsor91/act) project.

## Components

`act-tester` has the following components:

- `act-tester-plan`, which creates an initial test plan over some Litmus tests and compilers for a single machine;
- `act-tester-fuzz`, which runs `act-fuzz` over a test plan to create a more useful plan;
- `act-tester-lift`, which runs a harness maker over a test plan to _lift_ subjects compilable harnesses;
- `act-tester-mach`, which runs _local_ compilers over Litmus harnesses to produce results;
- `act-tester`, which combines the above into a looping test campaign over multiple machines.

It also contains the following utilities:

- `act-gccnt`, a wrapper around `gcc` that adds support for (controlled) misbehaviour, useful for testing testers'
  resilience to compiler crashes;
- `act-litmus`, a wrapper around `litmus7` that incorporates various workarounds useful for act-tester.

## Use

Note that `act-tester` is still pretty rough around the edges - please feel free to file issues about its user
experience and documentation.

- Install using the usual `go` tools: for example,
 `go get github.com/MattWindsor91/act-tester/cmd/...`.  All commands are in the `cmd` directory.
- Make sure that at least `act-tester-mach` is installed on any remote machine you wish to use for testing.
- Create a `tester.toml` file (see `tester-example.toml`).
  Note that this is _different_ from the `act.conf` used by ACT (and, in fact, supersedes it in several areas),
  but you'll need both.
- For now, the easiest thing to do is to `cd` into an ACT working directory, make an `act.conf` in there, and run
  a command such as `act-tester examples/c_litmus/memalloy/*.litmus`.

## Licence

As with the rest of ACT, `act-tester` uses the MIT licence.  (See `LICENSE` for details.)

This part of ACT _doesn't_ use any code from the
[Herdtools7](https://github.com/herd/herdtools7) project, but other parts do,
with all the licensing consequences that entails.

## Acknowledgements

`act-tester` is part of work carried out under the umbrella of the
[IRIS](https://interfacereasoning.com) project.
