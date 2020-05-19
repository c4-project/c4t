% act-tester-plan 8

# NAME

act-tester-plan - runs the planning phase of an ACT test standalone

# SYNOPSIS

act-tester-plan

```
[--corpus-size|-n]=[value]
[--num-workers|-j]=[value]
[--seed|-s]=[value]
[-A]=[value]
[-C]=[value]
[-m]=[value]
[-x]
```

**Usage**:

```
act-tester-plan [GLOBAL OPTIONS] command [COMMAND OPTIONS] [ARGUMENTS...]
```

# GLOBAL OPTIONS

**--corpus-size, -n**="": `number` of corpus files to select for this test plan;
if positive, the planner will use all viable provided corpus files (default: 0)

**--num-workers, -j**="": number of `workers` to run in parallel (default: 1)

**--seed, -s**="": `seed` to use for any randomised components of this test plan; -1 uses run time as seed (default: -1)

**-A**="": read ACT config from this `file`

**-C**="": read ACT config from this `file`

**-m**="": ID of machine to use for this test plan

**-x**: if true, use 'dune exec' to run OCaml ACT binaries

