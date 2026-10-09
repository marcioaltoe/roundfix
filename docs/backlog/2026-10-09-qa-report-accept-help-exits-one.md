---
type: fix
status: open
created: 2026-10-09
spec: null
---

# `roundfix qa-report accept --help` exits 1

## Problem

`roundfix qa-report accept --help` exits 1 instead of printing the
subcommand's help with exit 0. Spec 0251's author found it while
fingerprinting command help, and had to use the parent command's help for
that surface. Reproduced on 2026-10-09 with v0.64.0.

## Direction

Make `--help` on `qa-report accept` print its usage and exit 0, like every
other subcommand. Add the case to the repository's help-exit test. Then point
the Behavior Surface Record at the subcommand's own help.

## Sources

Spec 0251's authoring report; reproduced by the operator.
