% c4t-backend 8

# NAME

c4t-backend - runs backends standalone

# SYNOPSIS

c4t-backend

```
[--arch|-a]=[value]
[--backend-id|-n]=[value]
[--backend-style|-s]=[value]
[-C]=[value]
[-x]
```

**Usage**:

```
c4t-backend [GLOBAL OPTIONS] command [COMMAND OPTIONS] [ARGUMENTS...]
```

# GLOBAL OPTIONS

**--arch, -a**="": ID of `ARCH` to target for architecture-dependent backends

**--backend-id, -n**="": filter to backends whose names match `GLOB`

**--backend-style, -s**="": filter to backends whose styles match `GLOB`

**-C**="": read tester config from this `file`

**-x**: if true, use 'dune exec' to run OCaml ACT binaries

