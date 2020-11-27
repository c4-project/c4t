% c4t-perturb 8

# NAME

c4t-perturb - perturbs a test plan

# SYNOPSIS

c4t-perturb

```
[--corpus-size|-n]=[value]
[--full-ids|-I]
[--seed|-s]=[value]
[--verbose|-v]
[-C]=[value]
```

**Usage**:

```
c4t-perturb [GLOBAL OPTIONS] command [COMMAND OPTIONS] [ARGUMENTS...]
```

# GLOBAL OPTIONS

**--corpus-size, -n**="": `number` of corpus files to select for this test plan (default: 0)

**--full-ids, -I**: map compilers to their 'full' IDs on perturbance

**--seed, -s**="": `seed` to use for any randomised components of this test plan (default: -1)

**--verbose, -v**: enables verbose output

**-C**="": read tester config from this `file`

