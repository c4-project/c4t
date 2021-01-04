% c4t-fuzz 8

# NAME

c4t-fuzz - runs the batch-fuzzer phase of a C4 test

# SYNOPSIS

c4t-fuzz

```
[--corpus-size|-n]=[value]
[--verbose|-v]
[-C]=[value]
[-d]=[value]
[-k]=[value]
[-x]
```

**Usage**:

```
c4t-fuzz [GLOBAL OPTIONS] command [COMMAND OPTIONS] [ARGUMENTS...]
```

# GLOBAL OPTIONS

**--corpus-size, -n**="": `number` of corpus files to select for this test plan (default: 0)

**--verbose, -v**: enables verbose output

**-C**="": read tester config from this `file`

**-d**="": `directory` to which outputs will be written (default: fuzz_results)

**-k**="": number of `cycles` to run for each subject in the corpus (default: 10)

**-x**: if true, use 'dune exec' to run c4f binaries

