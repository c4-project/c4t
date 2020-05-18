% act-tester-fuzz 8

# NAME

act-tester-fuzz - runs the batch-fuzzer phase of an ACT test

# SYNOPSIS

act-tester-fuzz

```
[-A]=[value]
[-d]=[value]
[-i]=[value]
[-k]=[value]
[-n]=[value]
[-x]
```

**Usage**:

```
act-tester-fuzz [GLOBAL OPTIONS] command [COMMAND OPTIONS] [ARGUMENTS...]
```

# GLOBAL OPTIONS

**-A**="": read ACT config from this `file`

**-d**="": `directory` to which outputs will be written (default: fuzz_results)

**-i**="": read from this plan `file` instead of stdin

**-k**="": number of `cycles` to run for each subject in the corpus (default: 10)

**-n**="": `number` of corpus files to select for this test plan;
if non-positive, the planner will use all viable provided corpus files (default: 0)

**-x**: read ACT config from this `file`

