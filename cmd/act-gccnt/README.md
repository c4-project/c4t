# act-gccnt

`act-gccnt` (_gccn't_) is a wrapper over `gcc` that injects random failures.
These failures simulate things like compiler failure (returning nonzero exit code)
or divergence (not terminating).

## Usage

`act-gccnt [FLAGS] -o OUTFILE INFILE...`

### Flags

_gccn't_ takes only a very small subset of GCC flags.
It also takes various flags that control how it fails, prefixed with `nt` and
discussed below.

#### Trigger failure on optimisation level

These flags take an optimisation level (minus the leading `-O`), and can be repeated.

- `--nt-diverge-opt`: spin infinitely instead of compiling;
- `--nt-error-opt`: exit with code 1 instead of compiling.

#### Miscellaneous

- `--nt-dry-run`: print out a summary of what _gccn't_ _would_ do, but don't actually do anything.
