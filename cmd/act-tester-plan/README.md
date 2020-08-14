% act-tester-plan 8

# NAME

act-tester-plan - runs the planning phase of an ACT test standalone

# SYNOPSIS

act-tester-plan

```
[--filter-compilers|-c]=[value]
[--filter-machines|-m]=[value]
[--num-workers|-j]=[value]
[-A]=[value]
[-C]=[value]
[-d]=[value]
[-x]
```

**Usage**:

```
act-tester-plan [GLOBAL OPTIONS] command [COMMAND OPTIONS] [ARGUMENTS...]
```

# GLOBAL OPTIONS

**--filter-compilers, -c**="": `glob` to use to filter compilers to enable

**--filter-machines, -m**="": `glob` to use to filter machines to plan

**--num-workers, -j**="": number of `workers` to run in parallel (default: 1)

**-A**="": read ACT config from this `file`

**-C**="": read tester config from this `file`

**-d**="": `directory` to which outputs will be written

**-x**: if true, use 'dune exec' to run OCaml ACT binaries

