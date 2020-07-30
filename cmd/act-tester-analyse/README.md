% act-tester-analyse 8

# NAME

act-tester-analyse - analyses a plan file

# SYNOPSIS

act-tester-analyse

```
[--csv-compilers]
[--csv-stages]
[--num-workers|-j]=[value]
[--save-dir|-d]=[value]
[--show-compilers|-C]
[--show-ok|-O]
[--show-plan-info|-P]
[--show-subjects|-S]
```

**Usage**:

```
act-tester-analyse [GLOBAL OPTIONS] command [COMMAND OPTIONS] [ARGUMENTS...]
```

# GLOBAL OPTIONS

**--csv-compilers**: dump CSV of compilers and their run times

**--csv-stages**: dump CSV of stages and their run times

**--num-workers, -j**="": number of `workers` to run in parallel (default: 1)

**--save-dir, -d**="": if present, save failing corpora to this `directory`

**--show-compilers, -C**: show breakdown of compilers and their run times

**--show-ok, -O**: show subjects that did not have compile or run issues

**--show-plan-info, -P**: show plan metadata and stage times

**--show-subjects, -S**: show subjects by status

