% c4t-plan 8

# NAME

c4t-plan - runs the planning phase of an ACT test standalone

# SYNOPSIS

c4t-plan

```
[--filter-compilers|-c]=[value]
[--filter-machines|-m]=[value]
[--num-workers|-j]=[value]
[--verbose|-v]
[-C]=[value]
[-d]=[value]
[-x]
```

**Usage**:

```
c4t-plan [GLOBAL OPTIONS] command [COMMAND OPTIONS] [ARGUMENTS...]
```

# GLOBAL OPTIONS

**--filter-compilers, -c**="": `glob` to use to filter compilers to enable

**--filter-machines, -m**="": `glob` to use to filter machines to plan

**--num-workers, -j**="": number of `workers` to run in parallel (default: 1)

**--verbose, -v**: enables verbose output

**-C**="": read tester config from this `file`

**-d**="": `directory` to which outputs will be written

**-x**: if true, use 'dune exec' to run OCaml ACT binaries

