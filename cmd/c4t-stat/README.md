% c4t-stat 8

# NAME

c4t-stat - inspects the statistics file

# SYNOPSIS

c4t-stat

```
[--csv-mutations]
[--input|-i]=[value]
[--mutations]=[value]
[--use-totals|-t]
[-C]=[value]
```

**Usage**:

```
c4t-stat [GLOBAL OPTIONS] command [COMMAND OPTIONS] [ARGUMENTS...]
```

# GLOBAL OPTIONS

**--csv-mutations**: dump CSV of mutation testing results

**--input, -i**="": read statistics from this `FILE`

**--mutations**="": show mutations matching `filter` ('all' or 'killed')

**--use-totals, -t**: use multi-session totals rather than per-session totals

**-C**="": read tester config from this `file`

