% act-tester-perturb 8

# NAME

act-tester-perturb - perturbs a test plan

# SYNOPSIS

act-tester-perturb

```
[--corpus-size|-n]=[value]
[--seed|-s]=[value]
[-C]=[value]
```

**Usage**:

```
act-tester-perturb [GLOBAL OPTIONS] command [COMMAND OPTIONS] [ARGUMENTS...]
```

# GLOBAL OPTIONS

**--corpus-size, -n**="": `number` of corpus files to select for this test plan;
if positive, the planner will use all viable provided corpus files (default: 0)

**--seed, -s**="": `seed` to use for any randomised components of this test plan; -1 uses run time as seed (default: -1)

**-C**="": read tester config from this `file`

