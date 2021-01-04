% c4t 8

# NAME

c4t - runs compiler tests

# SYNOPSIS

c4t

```
[--corpus-size|-n]=[value]
[--cpuprofile]=[value]
[--machine-filter|-m]=[value]
[--no-dashboard|-D]
[-C]=[value]
[-k]=[value]
[-x]
```

**Usage**:

```
c4t [GLOBAL OPTIONS] command [COMMAND OPTIONS] [ARGUMENTS...]
```

# GLOBAL OPTIONS

**--corpus-size, -n**="": `number` of corpus files to select for this test plan (default: 0)

**--cpuprofile**="": `file` into which we should dump pprof information

**--machine-filter, -m**="": a `glob` to use to filter incoming machines by ID

**--no-dashboard, -D**: turns off the dashboard

**-C**="": read tester config from this `file`

**-k**="": number of `cycles` to run for each subject in the corpus (default: 10)

**-x**: if true, use 'dune exec' to run c4f binaries

