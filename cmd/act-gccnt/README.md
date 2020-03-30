# act-gccnt

`act-gccnt` (as in _gccn't_) is a wrapper over `gcc` that injects random failures.
These failures simulate things like compiler failure (returning nonzero exit code)
or divergence (not terminating).

## Usage

`act-gccnt [FLAGS] -o OUTFILE INFILE...`

### Flags

`act-gccnt` takes only a very small subset of GCC flags.

It also takes various flags that control how it fails:

- `-gccnt error`: exit with code 1 instead of compiling;
- `-gccnt timeout`: spin infinitely instead of compiling.