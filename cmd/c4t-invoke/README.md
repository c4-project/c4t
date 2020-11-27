% c4t-invoke 8

# NAME

c4t-invoke - runs the machine-dependent phase of an ACT test, potentially remotely

# SYNOPSIS

c4t-invoke

```
[--compiler-timeout|-t]=[value]
[--force|-f]
[--num-compiler-workers|-j]=[value]
[--num-run-workers|-J]=[value]
[--run-timeout|-T]=[value]
[--verbose|-v]
[-C]=[value]
[-d]=[value]
```

**Usage**:

```
c4t-invoke [GLOBAL OPTIONS] command [COMMAND OPTIONS] [ARGUMENTS...]
```

# GLOBAL OPTIONS

**--compiler-timeout, -t**="": a `timeout` to apply to each compilation (default: 0s)

**--force, -f**: allow invoke on plans that have already been invoked

**--num-compiler-workers, -j**="": number of compiler `workers` to run in parallel (default: 0)

**--num-run-workers, -J**="": number of runner `workers` to run in parallel (not recommended except on manycore machines) (default: 0)

**--run-timeout, -T**="": a `timeout` to apply to each run (default: 0s)

**--verbose, -v**: enables verbose output

**-C**="": read tester config from this `file`

**-d**="": `directory` to which outputs will be written (default: mach_results)

