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
[--nt-diverge-mutant-period]=[value]
[--nt-diverge-opt]=[value]
[--nt-dryrun]
[--nt-error-mutant-period]=[value]
[--nt-error-opt]=[value]
[--nt-hit-mutant-period]=[value]
[--nt-mutant]=[value]
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

**--march**="": architecture optimisation `spec` to pass through to gcc

**--mcpu**="": cpu optimisation `spec` to pass through to gcc

**--nt-bin**="": the 'real' compiler `command` to run (default: gcc)

**--nt-diverge-mutant-period**="": diverge when the mutant number is a multiple of this `period` (default: 0)

**--nt-diverge-opt**="": optimisation `levels` (minus the '-O') on which gccn't should diverge

**--nt-dryrun**: print the outcome of running gccn't instead of doing it

**--nt-error-mutant-period**="": error when the mutant number is a multiple of this `period` (default: 0)

**--nt-error-opt**="": optimisation `levels` (minus the '-O') on which gccn't should exit with an error

**--nt-hit-mutant-period**="": report a hit when the mutant number is a multiple of this `period` (default: 0)

**--nt-mutant**="": the mutant `number` to use if simulating mutation testing (default: 0)

**--pthread**: passes through pthread to gcc

**--std**="": `standard` to pass through to gcc

**-O**: optimisation level ''

**-o**="": output file (default: a.out)

