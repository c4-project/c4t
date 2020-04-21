# `app` directory

This directory contains the CLI glue used by the various `cmd` stubs.

## Why not put this code in the `cmd` stubs?

Having the app bodies separately means we can do things like extract documentation for the commands
from their `cli.App`s.