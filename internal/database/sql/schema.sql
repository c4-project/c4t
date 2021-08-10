--  Copyright (c) 2020-2021 C4 Project
--
--  This file is part of c4t.
--  Licenced under the MIT licence; see `LICENSE`.

-- Schema for c4t analysis.

-- A machine.
CREATE TABLE machine
(
    -- Database ID of the machine.
    machine_id INTEGER PRIMARY KEY,
    -- c4t ID of the machine.
    name       TEXT UNIQUE NOT NULL
);

-- An unconfigured compiler.
--
-- Because a compiler's configuration can change over time, there can be multiple different compilers with the same
-- name.  A compiler record will generally be reused if possible, however.
CREATE TABLE compiler
(
    -- Database ID of the compiler.
    compiler_id INTEGER PRIMARY KEY,
    -- Database ID of the parent machine.
    machine_id  INTEGER NOT NULL REFERENCES machine,
    -- c4t ID of the compiler.
    name        TEXT    NOT NULL,
    -- Resolved architecture ID.
    arch        TEXT    NOT NULL,
    UNIQUE (machine_id, name, arch)
);

-- A group of plans coming from one instance of the director.
CREATE TABLE experiment
(
    -- Database ID of the experiment.
    experiment_id INTEGER PRIMARY KEY,
    -- Start time of the experiment.
    start_time    TIMESTAMP
);

-- A testing plan.
CREATE TABLE plan
(
    -- Database ID of this plan.
    plan_id       INTEGER PRIMARY KEY,
    -- ID of any experiment to which this plan belongs.
    experiment_id INTEGER REFERENCES experiment,
    -- ID of the machine on which this plan is executed.
    machine_id    INTEGER NOT NULL REFERENCES machine
);

-- A mutation set-up.
CREATE TABLE mutation
(
    -- Database ID of this mutation.
    mutation_id INTEGER NOT NULL,
    -- Compiler ID to which this mutation is attached.
    compiler_id INTEGER NOT NULL REFERENCES compiler,
    -- The mutation operator, if explicitly named.
    operator TEXT,
    -- The variant of the mutation within the operator, if explicitly numbered.
    variant INTEGER,
    -- The index of the mutation.
    number INTEGER NOT NULL,
    -- The operator/variant pair must be unique.
    UNIQUE(operator, variant),
    -- Indices must be unique per compiler id.
    UNIQUE(compiler_id, number)
);

-- Stores information about a compiler instance.
CREATE TABLE compiler_instance
(
    -- Database ID of the plan containing this configuration.
    plan_id     INTEGER NOT NULL REFERENCES plan,
    -- Database ID of the compiler being configured.
    compiler_id INTEGER NOT NULL REFERENCES compiler,
    -- The selected machine optimisation level (may be null).
    mopt TEXT,
    -- The selected optimisation level (may be null).
    opt TEXT,
    -- The mutation selected for this instance (may be null).
    mutation_id INTEGER REFERENCES mutation
);