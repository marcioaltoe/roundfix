---
task: task_02
spec: 0020-run-browser
status: pending
type: backend
complexity: medium
---

# Task 02: runs list enrichment: columns, state and limit flags, notes

## Overview

Bring the deterministic text surface to the new contract: enriched columns,
an Active-by-default state filter with `--state`, a bounded default with
`--limit`, hidden-count notes, and the self-explaining attach wording for
unknown ids and non-interactive contexts. Demoable by running `runs list`
against seeded Runs.

## Requirements

1. MUST print one Run per line, newest first, with the stable column order:
   run id, state, kind, target, agent, absolute UTC start time, duration
   (`running <elapsed>` for Active Runs), and local branch — plus the
   repository as the final column with `--all`. Run ids are never truncated.
2. MUST default the state filter to Active and support
   `--state <active|terminal|all>`; the unreleased `--active` flag is
   removed.
3. MUST default the bound to the 20 newest matching Runs and support
   `--limit N` with `0` unbounded.
4. MUST print exactly one trailing stderr note when Runs are hidden by the
   state filter or the bound, naming the hidden count and the widening flag
   (shapes: `(N terminal Run(s) hidden; use --state all)`,
   `(N older Run(s) hidden; use --limit 0)`).
5. MUST keep the empty result a single stdout line with exit `0`, stdout
   report-only, and every existing exit-code contract.
6. MUST update `roundfix attach <unknown>` to fail with the error naming
   picker numbers as unstable, and the non-interactive no-run-id failures of
   `attach` (and `runs` without a subcommand) to name `runs list`.

## Subtasks

- [ ] Column rendering with the shared row formatter (absolute time, duration)
- [ ] `--state` flag replacing `--active`; default active
- [ ] `--limit` flag with default 20
- [ ] Hidden-count stderr notes
- [ ] attach unknown-id and non-interactive wording; `runs` non-TTY wording
- [ ] CLI tests: byte-pinned columns, each flag, notes, empty, wording

## Acceptance Criteria

- [ ] With seeded Runs of both kinds, `runs list` prints only Active Runs
      with the eight columns and pins the byte shape in a CLI test.
- [ ] `--state all` includes terminal Runs; `--state terminal` excludes
      Active ones; hiding produces exactly one matching stderr note.
- [ ] 25 seeded matching Runs print 20 lines plus the older-hidden note;
      `--limit 0` prints all 25 without a note.
- [ ] `attach 41` (unknown id) exits `2` with the picker-number error;
      no-run-id non-interactive attach and bare `runs` in non-TTY exit `2`
      naming `runs list`.

## Verification

- `rtk go test ./internal/cli/` — expected: all tests pass, including the new
  listing and wording tests.
- `rtk make verify` — expected: fmt-check, tests, skills check, and build all
  pass.

## References

`_prd.md` → User Stories 4-5; Core Features 4-6. `_techspec.md` → API
Contracts: runs list, attach; Build Order 2; Risks (pre-release flag
supersession).
