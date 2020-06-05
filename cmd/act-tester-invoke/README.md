% act-tester-rmach 8

# NAME

act-tester-rmach - runs the machine-dependent phase of an ACT test, potentially remotely

# SYNOPSIS

act-tester-rmach

```
[--compiler-timeout|-t]=[value]
[--emit-json|-J]
[--num-workers|-j]=[value]
[--run-timeout|-T]=[value]
[--skip-compiler]
[--skip-runner]
[-C]=[value]
[-d]=[value]
```

**Usage**:

```
act-tester-rmach [GLOBAL OPTIONS] command [COMMAND OPTIONS] [ARGUMENTS...]
```

# GLOBAL OPTIONS

**--compiler-timeout, -t**="": a `timeout` to apply to each compilation (default: 1m0s)

**--emit-json, -J**: emit progress reports in JSON form on stderr

**--num-workers, -j**="": number of `workers` to run in parallel (default: 1)

**--run-timeout, -T**="": a `timeout` to apply to each run (default: 1m0s)

**--skip-compiler**: if given, skip the compiler

**--skip-runner**: if given, skip the runner

**-C**="": read ACT config from this `file`

**-d**="": `directory` to which outputs will be written (default: mach_results)

