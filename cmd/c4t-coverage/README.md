% c4t-coverage 8

# NAME

c4t-coverage - makes a coverage testbed using a plan

# SYNOPSIS

c4t-coverage

```
[--config|-c]=[value]
[--verbose|-v]
[-d]=[value]
[-x]
```

**Usage**:

```
c4t-coverage [GLOBAL OPTIONS] command [COMMAND OPTIONS] [ARGUMENTS...]
```

# GLOBAL OPTIONS

**--config, -c**="": Path to config file for coverage (not the tester config file!) (default: coverage.toml)

**--verbose, -v**: enables verbose output

**-d**="": `directory` to which outputs will be written (default: coverage)

**-x**: if true, use 'dune exec' to run OCaml ACT binaries

