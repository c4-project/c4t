% act-tester-lift 8

# NAME

act-tester-lift - runs the harness-lifter phase of an ACT test

# SYNOPSIS

act-tester-lift

```
[-A]=[value]
[-d]=[value]
[-x]
```

**Usage**:

```
act-tester-lift [GLOBAL OPTIONS] command [COMMAND OPTIONS] [ARGUMENTS...]
```

# GLOBAL OPTIONS

**-A**="": read ACT config from this `file`

**-d**="": `directory` to which outputs will be written (default: lift_results)

**-x**: if true, use 'dune exec' to run OCaml ACT binaries

