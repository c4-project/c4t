% c4t-mach 8

# NAME

c4t-mach - runs the machine-dependent phase of a C4 test

# SYNOPSIS

c4t-mach

```
[--compiler-timeout|-t]=[value]
[--num-compiler-workers|-j]=[value]
[--num-run-workers|-J]=[value]
[--run-timeout|-T]=[value]
[-d]=[value]
```

**Usage**:

```
c4t-mach [GLOBAL OPTIONS] command [COMMAND OPTIONS] [ARGUMENTS...]
```

# GLOBAL OPTIONS

**--compiler-timeout, -t**="": a `timeout` to apply to each compilation (default: 0s)

**--num-compiler-workers, -j**="": number of compiler `workers` to run in parallel (default: 0)

**--num-run-workers, -J**="": number of runner `workers` to run in parallel (not recommended except on manycore machines) (default: 0)

**--run-timeout, -T**="": a `timeout` to apply to each run (default: 0s)

**-d**="": `directory` to which outputs will be written (default: mach_results)

