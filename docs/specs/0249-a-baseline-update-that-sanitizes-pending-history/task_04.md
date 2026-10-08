---
task: task_04
spec: 0249-a-baseline-update-that-sanitizes-pending-history
status: completed
type: backend
complexity: low
---

# Task 04: The upgrade notice names the Pending History

## Overview

After every release outcome that prints the recommendation notice,
`roundfix upgrade` appends one line about the Pending History of the working
directory's repository, or says to run `roundfix baseline update` when it runs
outside a repository (ADR-0254). The upgrade's stdout and exit code do not
change, and the line is computed in process from the filesystem. The Task is
verifiable on its own through the existing upgrade notice fixtures.

## Requirements

1. MUST add `printHistoryNotice(environment commandEnvironment, stderr io.Writer)`
   and call it at the end of `printUpgradeRecommendationNotice`. It runs after
   the recommendation lines, the installed executable's output or the
   `recommendations not checked` line, on both the in-process and the
   installed-executable path. It implements `_techspec.md` → API Contract 7:
   - It loads the configuration with `loadCommandConfig(environment, io.Discard)`.
     A load failure prints nothing more.
   - When `GitRoot` is empty it prints
     `roundfix: history: outside a repository; run roundfix baseline update in each adopted repository to plan its pending history`.
   - Otherwise it resolves the Spec Root (an external one prints nothing) and
     calls `historyInventory`.
   - With pending units it prints
     `roundfix: history: <u> unit(s) pending sanitize (<parts>); run roundfix baseline update to plan them`.
     `<parts>` is `<f> Legacy Archive Folder(s)` when `<f>` is positive,
     followed by each pending kind name in inventory order, joined by `, `.
   - With none pending it prints nothing.
   - On an inventory error it prints
     `roundfix: history not checked: <reason>`, with the reason on one line
     through task_02's `historyRefusalLine`.
2. MUST run no Git subprocess, network call or write in the notice
   (`_techspec.md` → Invariant 10). Help, usage errors and failed upgrades
   still print no notice line at all.
3. MUST add these tests to `internal/cli/upgrade_notice_test.go`, on the
   existing `prepareNoticeUpgrade` fixtures:
   - `TestUpgradeNoticeNamesPendingHistory`: with three Legacy Archive
     Folders and one unreduced retired Finding written into the workspace
     repository, the `current` and `installed` outcomes each print the exact
     line `roundfix: history: 4 unit(s) pending sanitize (3 Legacy Archive Folder(s), findings); run roundfix baseline update to plan them`
     as the last stderr line. Stdout and exit code equal the outcome's
     existing values.
   - `TestUpgradeNoticeOutsideARepositoryNamesBaselineUpdate`: a work
     directory with no `.git` in any parent prints the outside line as the
     last stderr line.
   - `TestUpgradeNoticeIsSilentWithoutPendingHistory`: the workspace with only
     already reduced history prints exactly what it printed before, that is
     `upgradeFixtureNotice()`.
4. MUST NOT change any existing assertion in the upgrade tests. Every existing
   notice test passes unchanged, because its workspace holds no history.

## Subtasks

- [ ] Add `printHistoryNotice` and call it after the recommendation notice.
- [ ] Write the three notice tests.

## Acceptance Criteria

- [ ] Inside a repository with Pending History, the upgrade notice names its count and kinds and points at `roundfix baseline update`.
- [ ] Outside a repository, the notice says to run the update in each adopted repository.
- [ ] Stdout, exit codes and the existing notices are unchanged.

## Context

- interface: `internal/cli/upgrade.go`
- interface: `internal/cli/upgrade_notice_test.go`
- instruction: `internal/cli/history.go`
- instruction: `docs/adr/0254-baseline-update-sanitizes-pending-history-in-the-change-it-plans.md`
- instruction: `docs/user-guide/commands/upgrade.md`

## Verification

- `out="$(go test -count=1 -v -run '^(TestUpgradeNoticeNamesPendingHistory|TestUpgradeNoticeOutsideARepositoryNamesBaselineUpdate|TestUpgradeNoticeIsSilentWithoutPendingHistory)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for t in TestUpgradeNoticeNamesPendingHistory TestUpgradeNoticeOutsideARepositoryNamesBaselineUpdate TestUpgradeNoticeIsSilentWithoutPendingHistory; do printf '%s\n' "$out" | grep -q -- "--- PASS: $t " || { printf 'missing passing test %s\n%s\n' "$t" "$out" >&2; exit 1; }; done; go test -count=1 -run '^TestUpgrade' ./internal/cli` — expected: exit 0. Before this Task the three tests do not exist, so the first check fails. After it, they pass and every existing upgrade test passes unchanged.

## References

- `_prd.md` → Core Feature 5; Goals; User Story 4; Success Metric 5
- `_techspec.md` → API Contract 7; Invariant 10; Build Order 4
- ADR-0254

## Result

Implemented the filesystem-only Pending History notice at the end of
`printUpgradeRecommendationNotice`, after either recommendation path and its
failure diagnostic. Configuration load failures, unresolved or external Spec
Roots, and inventories with no pending units add no line. Inventory failures,
including errors retained on individual units, use `historyRefusalLine`.

Acceptance evidence:

- Inside a repository: `TestUpgradeNoticeNamesPendingHistory` checks both
  `current` and `installed` outcomes with three Legacy Archive Folders and an
  unreduced retired Finding. It asserts the exact four-unit notice as the
  last stderr line, the existing stdout, and exit 0.
- Outside a repository:
  `TestUpgradeNoticeOutsideARepositoryNamesBaselineUpdate` asserts the exact
  outside line after the recommendation notice, unchanged stdout, and exit 0.
- Existing output: `TestUpgradeNoticeIsSilentWithoutPendingHistory` exercises
  reduced Findings and Backlog Entries and asserts `upgradeFixtureNotice()`
  exactly. Existing notice assertions remain byte-identical. Focused checks
  also exercise every release outcome, the installed child, recommendation
  failures, and help, usage errors and failed upgrades.
- Additional checks: `TestUpgradeNoticeHistoryInventoryFailures` covers both
  archive enumeration errors and per-kind symlink refusals, asserting one
  diagnostic line and unchanged stdout/exit; the external Spec Root test
  confirms silence even with pending repository history.
- Inspection of the new notice and its configuration, root resolution and
  inventory callees confirms filesystem reads only: no Git subprocess,
  network request or write. Tests use the existing `.git` directory fixture,
  which is not an initialized Git repository.

Focused checks:

- Before implementation,
  `GOCACHE=/private/tmp/roundfix-task04-gocache rtk proxy go test ./internal/cli -run 'TestUpgradeNotice' -count=1`
  exited 1: pending-history and outside-repository tests reproduced the
  missing lines; the reduced-history test already passed.
- After the final Go edits,
  `GOCACHE=/private/tmp/roundfix-task04-gocache rtk proxy go test ./internal/cli -run 'TestUpgrade(Notice|Writes|Asks|Keeps)' -count=1`
  exited 0 (`ok roundfix/internal/cli`).
- `rtk proxy git -c core.fsmonitor=false diff --check` exited 0.
- `GOCACHE=/private/tmp/roundfix-task04-gocache rtk make verify-incremental`
  exited 0 with sandbox escalation: formatting, vet, package tests (including
  the complete CLI suite), skill synchronization/checks and build passed.
  The initial sandboxed attempt was interrupted by denied network access to
  `cafe.github.com` and provided no passing result. Its two orphaned CLI test
  fixtures were stopped. The successful rerun held the worktree steady;
  no implementation edits followed it.

Declared Task Verification is reserved for the Daemon. Status remains
`in_progress`; no commit, push, Pull Request, Task Graph edit or other Task
edit was performed. The only pre-existing worktree change was the
Daemon-owned status in this file. No follow-up implementation was identified.
