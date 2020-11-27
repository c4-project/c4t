% c4t-gccnt 8

# NAME

c4t-gccnt - wraps gcc with various optional failure modes

# SYNOPSIS

c4t-gccnt

```
[--O0]
[--O1]
[--O2]
[--O3]
[--Ofast]
[--Og]
[--Os]
[--Oz]
[--march]=[value]
[--mcpu]=[value]
[--nt-bin]=[value]
[--nt-diverge-opt]=[value]
[--nt-dryrun]
[--nt-error-opt]=[value]
[--pthread]
[--std]=[value]
[-O]
[-o]=[value]
```

**Usage**:

```
c4t-gccnt [GLOBAL OPTIONS] command [COMMAND OPTIONS] [ARGUMENTS...]
```

# GLOBAL OPTIONS

**--O0**: optimisation level '0'

**--O1**: optimisation level '1'

**--O2**: optimisation level '2'

**--O3**: optimisation level '3'

**--Ofast**: optimisation level 'fast'

**--Og**: optimisation level 'g'

**--Os**: optimisation level 's'

**--Oz**: optimisation level 'z'

**--march**="": architecture optimisation to pass through to gcc

**--mcpu**="": cpu optimisation to pass through to gcc

**--nt-bin**="": the 'real' compiler `command` to run (default: gcc)

**--nt-diverge-opt**="": o-levels (minus the '-O') on which gccn't should diverge

**--nt-dryrun**: print the outcome of running gccn't instead of doing it

**--nt-error-opt**="": o-levels (minus the '-O') on which gccn't should exit with an error

**--pthread**: passes through pthread to gcc

**--std**="": standard to pass through to gcc

**-O**: optimisation level ''

**-o**="": output file (default: a.out)

