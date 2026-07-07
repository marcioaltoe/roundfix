---
task: task_02
spec: 0017-run-discovery
status: completed
type: backend
complexity: medium
---

# Task 02: runs list command with stable columns and empty-result contract

## Overview

Expose the listing query as `roundfix runs list`: a new `runs` command
namespace whose first subcommand prints this repository's Runs as stable
plain-text columns for humans and agents. Demoable end to end: seed Runs, run
the command, read the report.

## Requirements

1. MUST add a `runs` command with a `list` subcommand following the existing
   subcommand dispatch shape; unknown subcommands and unexpected arguments
   fail with exit `2` and a usage pointer.
2. MUST print one Run per line, newest first, in the stable column order:
   run id, state, kind, target. The target is the Open Pull Request number for
   review-kind Runs and the Spec slug for spec Runs. Active Runs are visibly
   marked; run ids are never truncated.
3. MUST scope to the current repository by default, support `--all` to list
   every repository (each line then also names the repository), and support
   `--active` to filter to Active Runs; the flags compose.
4. MUST print a single clear line and exit `0` when no Runs match.
5. MUST exit `2` with an actionable error when run outside a git repository
   without `--all`, naming `--all` as the alternative.
6. MUST keep stdout report-only; diagnostics go to stderr. Help and top-level
   usage text name the new command truthfully.

## Subtasks

- [x] `runs` dispatch and `list` flag parsing with usage text
- [x] Line formatting: columns, active marker, target derivation, `--all`
      repository column
- [x] Repository resolution and the outside-a-repository failure
- [x] Empty-result line and exit code
- [x] CLI tests: format pinning, ordering, filters, empty, usage errors

## Acceptance Criteria

- [x] With seeded Runs of both kinds, `runs list` prints them newest first with
      id, state, kind, and the correct per-kind target, and the byte shape is
      pinned by a CLI test.
- [x] `runs list --active` prints only Active Runs; `runs list --all` includes
      other repositories' Runs and names each repository.
- [x] With no Runs, the command prints one line and exits `0`.
- [x] `runs bogus` and `runs list <extra>` exit `2` with a usage pointer.
- [x] `roundfix --help` lists the `runs` command.

## Verification

- `rtk go test ./internal/cli/` — expected: all tests pass, including the new
  `runs list` tests.
- `rtk go run ./cmd/roundfix runs list` — expected: this repository's Runs or
  the single empty-result line; exit `0`.
- `rtk make verify` — expected: fmt-check, tests, skills check, and build all
  pass.

## References

`_prd.md` → User Stories 1-2; Core Features 1-4. `_techspec.md` → API
Contracts: runs list; Build Order 2; Risks (column stability).

## Result

Implemented `roundfix runs list` through the new `runs` command namespace. The
command scopes to the current Git repository by default, supports `--all` and
`--active`, prints stable plain-text report rows with full Run ids, marks Active
Runs with `*`, derives targets as `pr:<number>` or `spec:<slug>`, and keeps
diagnostics on stderr.

Acceptance evidence:

- `TestRunRunsListPrintsStableColumnsNewestFirst` seeds review and spec Runs
  and pins the exact byte output order and columns.
- `TestRunRunsListActiveAndAllFlagsCompose` covers the active-only filter,
  repository scoping, `--all`, and composed `--all --active` behavior.
- `TestRunRunsListEmptyResultExitsZero` covers the single empty-result line and
  exit `0`.
- `TestRunRunsListUsageErrors` covers `runs bogus` and `runs list <extra>`
  exiting `2` with usage pointers.
- `TestRunHelp` and `TestRunCommandHelp/runs` cover top-level help and `runs`
  usage text.
- `TestRunRunsListOutsideRepositoryRequiresAll` covers the outside-repository
  exit `2` error naming `--all`.

Verification:

- `rtk go test ./internal/cli/` passed: 303 CLI tests.
- `rtk go run ./cmd/roundfix runs list` passed and printed `No Runs found.`
  with exit `0`.
- `rtk make verify` passed: `rtk go test ./...` reported 832 tests across 18
  packages, `roundfix skills check` passed, and `go build` completed.
