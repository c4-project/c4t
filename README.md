# c4t

_c4t_ ('C4 tester')
is the top-level testing framework
for the 'C Compiler Concurrency Checker' (C4) project.  It sits on top of
[c4f](https://github.com/MattWindsor91/c4f) and
[herdtools7](https://github.com/herd/herdtools7) and allows running of
multi-compiler, multi-machine testing campaigns.

_c4t_ is written in Go (making it fairly portable and cross-compilable, perhaps
unlike _c4f_), and is free software under the MIT licence.

## Components

**NOTE:** The _c4t_ binaries are currently named `act-tester*`; these will be
renamed to `c4t*` at a later date.

The main entry point into _c4t_ is `act-tester`, which sets up and executes
multiple parallel testing cycles for fully automated use.  Many parts of
_c4t_'s testing infrastructure, as well as other utilities, are also available
in separate binaries.

The main object of work in _c4t_ is a _test plan_, a bulky (and often gzipped)
JSON file that centralises information about a testing campaign for one machine.
Many of its components operate on one or more test plans, either implicitly or
explicitly.

At time of writing, `c4t` has the following components (due to be renamed from
`act-tester-*` to `c4t-*`:

### The main test cycle

- `act-tester-plan`, which reads the tester config and creates an initial test
   plan over some Litmus tests and compilers for a single machine;
- `act-tester-perturb`, which _perturbs_ a plan configuration, choosing a random
  sample of plans as well as randomised parameters to the compiler;
- `act-tester-fuzz`, which runs _c4f_ over a test plan to create a more useful
  plan;
- `act-tester-lift`, which runs backends like `litmus7` over a test plan to
  _lift_ subjects to harnesses;
- `act-tester-invoke` (on the machine running _c4t_) and `act-tester-mach` (on
   the target machine), which communicate with each other through SSH and
   perform the compilation and running phases of a test plan;

### Analysing test plans

- `act-tester-analyse`, which performs some basic analysis over a test plan and
  prints reports on failures, compiler warnings, etc.;
- `act-tester`, which combines the above into a looping test campaign over multiple machines.

### Utilities

- `act-tester-setc`, which overrides compiler parameters in an existing plan
  (useful for exploring particular optimisation levels);
- `act-tester-coverage`, which produces coverage testbeds (work in progress);
- `act-gccnt` (GCCn't), a wrapper over `gcc` that can inject compiler failures
  when certain parameters are triggered (useful for testing that the workflow
  handles such issues).

## Use

Note that _c4t_ is still pretty rough around the edges - please feel free to
file issues about its user experience and documentation.

- Install using the usual `go` tools: for example,
 `go get github.com/MattWindsor91/c4t/cmd/...`.  All commands are in the `cmd` directory.
- Make sure that the [c4f](https://github.com/MattWindsor91/c4f) tools are
  in `PATH` on the test-running machine (eg run `make install`)
- Make sure that at least `act-tester-mach` is installed on any remote machine you wish to use for testing.
- Create a `tester.toml` file in
  [UserConfigDir](https://golang.org/pkg/os/#UserConfigDir)`/act`
  (see `tester-example.toml`).
- The easiest way to check if _c4t_ is working is
  `act-tester path/to/c4f/examples/c_litmus/memalloy`.

## Acknowledgements

_c4t_ is part of work carried out under the umbrella of the
[IRIS](https://interfacereasoning.com) project.
