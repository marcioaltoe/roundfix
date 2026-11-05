---
task: task_02
spec: 0357-a-numbered-file-collision-between-queued-specs-parks-before-the-run
status: completed
type: backend
complexity: medium
---

# Task 02: The queue parks a colliding item before its Run and a retry refuses while it collides

## Overview

Today Delivery Revalidation parks only on strict findings, and nothing in the
queue compares the numbered files of two items (`_techspec.md` → Current
behavior 1, 2 and 4). After this Task `internal/delivery` turns the Numbered
Files an item adds, those on the default branch and those of earlier unmerged
items into a report of collisions and next free numbers, the engine parks a
colliding item as `numbered-file-collision` before any Run, a retry of an item
that has not run refuses while it collides, and the park falls in the
conflict Park Class with its next action. This is the detector the
maintainer chose for point 3 of the report reported by fiscus on 2026-11-05
([references](references/2026-11-05-fiscus-delivery-queue-report.md)).

## Requirements

1. MUST, before changing code, run
   `go test -count=1 -timeout 20m -run '^(TestRevalidation.*|TestRetryOfAnItemWithoutARun.*|TestEveryBlockerHasAParkClass|TestParkClassesAndNextActions)$' ./internal/delivery`
   and record in the Result that it passes, and that no exported identifier
   of the package names a numbered file.
2. MUST create `internal/delivery/numbered_files.go` with
   `BlockerNumberedFileCollision`, `NumberedFileClaims`,
   `NumberedFileCollision`, `NumberedFileReport`,
   `FindNumberedFileCollisions`, `(NumberedFileReport).Detail` and
   `(NumberedFileReport).Blocker` exactly as `_techspec.md` → Interfaces and
   API Contracts 3 and 4 state, and add the `NumberedFiles NumberedFileReport`
   field to `Revalidation`. The comparison takes paths the caller already
   filtered with the declaration; it imports no configuration and runs no
   command.
3. MUST, in `(*Engine).advanceItem`, park an item whose revalidation carries
   collisions with `revalidation.NumberedFiles.Blocker()` after the
   `revalidation-failed` park and before the `running` stage (API Contract 5),
   and, in `(*Engine).Retry`, refuse an item with no Run whose revalidation
   carries collisions with
   `retry Delivery Queue item "<slug>": numbered file collision: <Detail>`
   after the findings refusal, before any carry-forward (API Contract 6).
4. MUST add to `ClassifyPark` the case `BlockerNumberedFileCollision` →
   `ParkClassConflict` with the next action, verbatim:
   `in <worktree>, renumber the item's files the blocker names from its next free number, in the Task declarations on the item branch and in any file the branch holds, or deliver the earlier item it names first and merge the default branch into the item branch; then run roundfix deliver retry <slug>`
   (API Contract 7).
5. MUST create `internal/delivery/numbered_files_test.go`, each top-level test
   calling `t.Parallel()` first, with:
   - `TestFindNumberedFileCollisions`, subtests: `ahead_item` (subject
     `drizzle/0052_b.sql`, earlier `alpha` with `drizzle/0052_a.sql`: one
     collision owned by `alpha`, next free `drizzle/0053`), `default_branch`
     (default `drizzle/0051_x.sql` and `drizzle/0052_main.sql`, subject
     `drizzle/0052_b.sql`: owner `""`, next free `drizzle/0053`), `same_path`
     (subject and `alpha` both `drizzle/0052_a.sql`: one collision),
     `by_value` (subject `drizzle/052_b.sql`, `alpha` `drizzle/0052_a.sql`:
     collision, next free `drizzle/053`), `other_directory` (subject
     `drizzle/meta/0052_snapshot.json`, `alpha` `drizzle/0052_a.sql`: none),
     `no_leading_digit` (subject `drizzle/meta/_journal.json`, default
     holding it: none), `next_free_above_every_claim` (default up to
     `drizzle/0051_x.sql`, `alpha` `drizzle/0052_a.sql`, `gamma`
     `drizzle/0053_c.sql`, subject `drizzle/0052_b.sql`: one collision with
     `alpha`, next free `drizzle/0054`), `several` (subject
     `drizzle/0052_b.sql` and `drizzle/meta/0052_snapshot.json`, `alpha`
     `drizzle/0052_a.sql`, default `drizzle/meta/0052_snapshot.json`: the
     exact `Blocker` of API Contract 4) and `empty` (no claims: no collision,
     `Detail` and `Blocker` both `""`);
   - `TestNumberedFileReportRendersTheBlocker`, asserting the example string
     of API Contract 4 byte for byte from a hand-built report.
6. MUST add to `internal/delivery/revalidate_test.go`, through the existing
   `recordingItemRevalidator`, `newRevalidationTestEngine`,
   `newRetryRevalidationTestEngine` and `seedNoRunRetryItem` seams:
   - `TestRevalidationParksANumberedFileCollisionBeforeTheRun`: the item parks
     with the report's exact `Blocker`, keeps its branch, worktree and
     provisioned flag, and records no `run` event;
   - `TestStrictFindingsParkBeforeANumberedFileCollision`: findings and a
     collision together park as `revalidation-failed: <codes>`;
   - `TestRetryOfAnItemWithoutARunRefusesWhileANumberedFileCollides`: the
     exact refusal, the item byte-identical after it, and no carry-forward.
7. MUST add to `TestParkClassesAndNextActions` the row
   `{BlockerNumberedFileCollision, ParkClassConflict}` and an exact check of
   its next action for slug `example` and worktree `/worktrees/example`.
8. MUST run sabotages 3 to 8 of `_techspec.md` → Testing Approach on a
   temporary copy and record which named test fails.
9. MUST keep every existing test of `./internal/delivery` passing unchanged
   apart from the row requirement 7 adds.
10. MUST NOT change `internal/cli`, `internal/config`, the Run Database
    schema, any other blocker, Park Class or next action, any guide, skill or
    reference, the Makefile, CI workflows, `go.mod` or `go.sum`.

## Subtasks

- [ ] Characterize Delivery Revalidation and the Park Classes.
- [ ] Add the pure comparison and the report.
- [ ] Park in `advanceItem`, refuse in `Retry`, classify in `ClassifyPark`.
- [ ] Add the tests and run the sabotages.

## Acceptance Criteria

- [ ] A collision is found by directory and number value against the default branch and every earlier claim, and the next free number is above all of them.
- [ ] A colliding item parks before any Run in the conflict Park Class with its exact next action, findings still park first, and a retry of an item that has not run refuses while it collides.

## Context

- creates: `internal/delivery/numbered_files.go`
- creates: `internal/delivery/numbered_files_test.go`
- interface: `internal/delivery/engine.go`
- interface: `internal/delivery/park_class.go`
- interface: `internal/delivery/park_class_test.go`
- interface: `internal/delivery/revalidate_test.go`

## Consumers

- `Revalidation.NumberedFiles`, `ClassifyPark` and `BlockerNumberedFileCollision` — consumers: `./internal/delivery`, `./internal/cli`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestFindNumberedFileCollisions|TestNumberedFileReportRendersTheBlocker|TestRevalidationParksANumberedFileCollisionBeforeTheRun|TestStrictFindingsParkBeforeANumberedFileCollision|TestRetryOfAnItemWithoutARunRefusesWhileANumberedFileCollides|TestRevalidation.*|TestRetryOfAnItemWithoutARun.*|TestEveryBlockerHasAParkClass|TestParkClassesAndNextActions|TestParkClassRequiresAnExactBlockerOrColonPrefix|TestExistingGenericParkAnswersAreUnchanged)$' ./internal/delivery 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestFindNumberedFileCollisions TestNumberedFileReportRendersTheBlocker TestRevalidationParksANumberedFileCollisionBeforeTheRun TestStrictFindingsParkBeforeANumberedFileCollision TestRetryOfAnItemWithoutARunRefusesWhileANumberedFileCollides TestFindNumberedFileCollisions/by_value TestFindNumberedFileCollisions/next_free_above_every_claim TestFindNumberedFileCollisions/other_directory TestParkClassesAndNextActions/numbered-file-collision; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 20m -run '^(TestDeliverStatus.*|TestRevalidate.*|TestFollowUp.*)$' ./internal/cli` — expected: exit 0. Before this Task none of the five new tests exists, so the first pass line is missing; afterwards the comparison, the park, the retry refusal and the Park Class hold, and the command suites that read Park Classes and revalidation still pass.

## References

- `_prd.md` → Goal 2; Goal 3; User Story 2; User Story 3; Core Feature 3; Core Feature 4; Core Feature 5; Core Feature 6; Success Metric 2; Success Metric 3
- `_techspec.md` → Current behavior 1; Current behavior 2; Current behavior 4; Interfaces; API Contract 3; API Contract 4; API Contract 5; API Contract 6; API Contract 7; Invariants 2; Invariants 3; Invariants 5; Testing Approach; Build Order 2

## Result

Skills loaded: implement-task, coding-guidelines, golang-error-handling, golang-safety, golang-structs-interfaces, golang-testing, testing-boss, evidence-gate, tech-writer, systematic-debugging.

Implemented the pure comparison and report, the pre-Run park, the no-Run retry
refusal, and the conflict Park Class with its exact next action. The comparison
reads only supplied paths; it imports no configuration and runs no command.

Pre-change characterization:

- Before changing code, ran the required
  `go test -count=1 -timeout 20m -run '^(TestRevalidation.*|TestRetryOfAnItemWithoutARun.*|TestEveryBlockerHasAParkClass|TestParkClassesAndNextActions)$' ./internal/delivery`:
  passed (25 tests including subtests).
- The pre-change search
  `rg -n '^(type|const|func) [A-Z].*(Numbered|numbered)|NumberedFile|numbered-file' internal/delivery`
  returned no matches: no exported identifier named a numbered file. Adding
  the tests before implementation then produced the expected build failure
  for the absent `NumberedFileClaims`, `NumberedFileReport` and
  `NumberedFileCollision` types.
- Initial Git status contained only the Daemon's `pending` → `in_progress`
  change in this Task file. That status remains untouched.

Acceptance evidence:

- Directory and numeric-value comparison against the default branch and all
  earlier claims: `TestFindNumberedFileCollisions` covers all nine required
  cases, including different padding, identical paths, ignored journal names,
  separate directories, and the next free number above every external claim.
  Additional cases cover the root directory, the 18-digit limit, ASCII digits,
  subject deduplication, ordering by rendered owner, padding from colliding
  subject paths, and unchanged input slices. `TestNumberedFileReportRendersTheBlocker`
  checks the hand-built example's Detail and Blocker byte for byte.
- `TestRevalidationParksANumberedFileCollisionBeforeTheRun` checks the exact
  blocker, retained branch/worktree/provisioned flag, no Run ID or `run`
  event, and preserved premise and owner warnings. Its first focused run
  exposed a missing fake-workspace recording callback; adding the existing
  `recordWorkspace` fixture setup resolved it without changing production
  persistence or weakening assertions.
- `TestStrictFindingsParkBeforeANumberedFileCollision` checks the exact
  findings blocker when both findings and collisions are present.
  `TestRetryOfAnItemWithoutARunRefusesWhileANumberedFileCollides` checks the
  exact refusal, equality of every persisted item field before and after,
  and absence of carry-forward.
- `TestParkClassesAndNextActions/numbered-file-collision` checks the conflict
  class and exact next action for both the bare blocker and a colon suffix;
  the existing Pending Question assertion checks the same answer.
- After the fixture correction,
  `go test -count=1 -timeout 5m ./internal/delivery` exited 0
  (`ok roundfix/internal/delivery 2.601s`). This package check exercised the
  new tests and all existing delivery tests. Existing tests were left
  unchanged apart from the required Park Class row and its exact-action check.

Sabotages 3–8 ran in six separate temporary source checkouts, each built
from `git archive HEAD` plus this Task's six delivery files. Every mutant
compiled and exited 1 on the expected behavioral assertion. Commands used
`go test -count=1 -timeout 5m -v -run <selector> ./internal/delivery`:

| Sabotage | Selector | Observed failure |
| --- | --- | --- |
| 3: compare digit strings | `^TestFindNumberedFileCollisions$/^by_value$` | `TestFindNumberedFileCollisions/by_value`: missing collision for `052` versus `0052` |
| 4: default-only next free | `^TestFindNumberedFileCollisions$/^next_free_above_every_claim$` | `TestFindNumberedFileCollisions/next_free_above_every_claim`: `drizzle/0052` instead of `drizzle/0054` |
| 5: ignore directory | `^TestFindNumberedFileCollisions$/^other_directory$` | `TestFindNumberedFileCollisions/other_directory`: unexpected cross-directory collision |
| 6: remove engine collision park | `^TestRevalidationParksANumberedFileCollisionBeforeTheRun$` | `TestRevalidationParksANumberedFileCollisionBeforeTheRun`: item merged and events included `run` |
| 7: remove retry collision refusal | `^TestRetryOfAnItemWithoutARunRefusesWhileANumberedFileCollides$` | `TestRetryOfAnItemWithoutARunRefusesWhileANumberedFileCollides`: nil error instead of refusal |
| 8: remove classification case | `^(TestEveryBlockerHasAParkClass\|TestParkClassesAndNextActions)$/^numbered-file-collision$` | `TestEveryBlockerHasAParkClass` and `TestParkClassesAndNextActions/numbered-file-collision`: unclassified blocker and generic next action |

Temporary transcripts: `/private/tmp/roundfix-0357-task02-7twr02su/sabotage-{3,4,5,6,7,8}.log`;
command/results index: `/private/tmp/roundfix-0357-task02-7twr02su/results.json`;
runner: `/private/tmp/roundfix-0357-task02-sabotages.py`.

Scope: only the six declared delivery files and this Result were authored.
The command workflow that supplies claims remains task_04's slice, and the
park documentation remains task_03's slice. Declared Task Verification,
repository Verification and CLI checks were not run; the Daemon owns them.
No Task status, manifest, other Task, configuration, guide, skill, schema,
module, Makefile or workflow was edited. No commit, push or Pull Request.

## Carry-forward provenance

- Source Run: `run_20261105T211401Z_cae788fdc7c4d6e1`
- Source commit: `7ba96664b8388136d31e5ec75d8aa34f3e7abea4`
