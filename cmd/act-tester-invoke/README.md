% act-tester-invoke 8

# NAME

act-tester-invoke - runs the machine-dependent phase of an ACT test, potentially remotely

# SYNOPSIS

act-tester-invoke

```
[--compiler-timeout|-t]=[value]
[--num-compiler-workers|-j]=[value]
[--num-run-workers|-J]=[value]
[--run-timeout|-T]=[value]
[--skip-compiler]
[--skip-runner]
[-C]=[value]
[-d]=[value]
```

**Usage**:

```
act-tester-invoke [GLOBAL OPTIONS] command [COMMAND OPTIONS] [ARGUMENTS...]
```

# GLOBAL OPTIONS

**--compiler-timeout, -t**="": a `timeout` to apply to each compilation (default: 1m0s)

**--num-compiler-workers, -j**="": number of compiler `workers` to run in parallel (default: 1)

**--num-run-workers, -J**="": number of runner `workers` to run in parallel (not recommended except on manycore machines) (default: 1)

**--run-timeout, -T**="": a `timeout` to apply to each run (default: 1m0s)

**--skip-compiler**: if given, skip the compiler

**--skip-runner**: if given, skip the runner

**-C**="": read tester config from this `file`

**-d**="": `directory` to which outputs will be written (default: mach_results)

