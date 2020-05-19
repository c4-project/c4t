% act-tester 8

# NAME

act-tester - makes documentation for act-tester commands

# SYNOPSIS

act-tester

```
[--cpuprofile]=[value]
[--machine-filter|-m]=[value]
[-A]=[value]
[-C]=[value]
[-k]=[value]
[-n]=[value]
[-x]
```

**Usage**:

```
act-tester [GLOBAL OPTIONS] command [COMMAND OPTIONS] [ARGUMENTS...]
```

# GLOBAL OPTIONS

**--cpuprofile**="": `file` into which we should dump pprof information

**--machine-filter, -m**="": A `glob` to use to filter incoming machines by ID.

**-A**="": read ACT config from this `file`

**-C**="": read ACT config from this `file`

**-k**="": number of `cycles` to run for each subject in the corpus (default: 10)

**-n**="": `number` of corpus files to select for this test plan;
if non-positive, the planner will use all viable provided corpus files (default: 0)

**-x**: read ACT config from this `file`
