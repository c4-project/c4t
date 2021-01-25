% c4t-analyse 8

# NAME

c4t-analyse - analyses a plan file

# SYNOPSIS

c4t-analyse

```
[--csv-compilers]
[--csv-stages]
[--error-on-bad-status|-e]
[--filter-file]=[value]
[--num-workers|-j]=[value]
[--save-dir|-d]=[value]
[--show-compiler-logs|-L]
[--show-compilers|-C]
[--show-mutation|-M]
[--show-ok|-O]
[--show-plan-info|-P]
[--show-subjects|-S]
```

**Usage**:

```
c4t-analyse [GLOBAL OPTIONS] command [COMMAND OPTIONS] [ARGUMENTS...]
```

# GLOBAL OPTIONS

**--csv-compilers**: dump CSV of compilers and their run times

**--csv-stages**: dump CSV of stages and their run times

**--error-on-bad-status, -e**: report an error if plan contains subjects with bad statuses

**--filter-file**="": load compile result filters from this file

**--num-workers, -j**="": number of `workers` to run in parallel (default: 1)

**--save-dir, -d**="": if present, save failing corpora to this `directory`

**--show-compiler-logs, -L**: show breakdown of compiler logs (requires -show-compilers)

**--show-compilers, -C**: show breakdown of compilers and their run times

**--show-mutation, -M**: show results of any mutation testing involved in this plan

**--show-ok, -O**: show subjects that did not have compile or run issues

**--show-plan-info, -P**: show plan metadata and stage times

**--show-subjects, -S**: show subjects by status

