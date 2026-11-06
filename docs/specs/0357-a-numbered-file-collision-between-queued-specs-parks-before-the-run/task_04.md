---
task: task_04
spec: 0357-a-numbered-file-collision-between-queued-specs-parks-before-the-run
status: completed
type: backend
complexity: high
---

# Task 04: Delivery Revalidation reads the numbered files of the item, the default branch and earlier items

## Overview

task_02 parks on a report the command workflow does not fill yet
(`_techspec.md` → Current behavior 3 and 6). After this Task
`(*commandDeliveryWorkflow).Revalidate` reads `delivery.numbered_files` at the
default-branch commit and, only when it declares patterns, collects the
Numbered Files the item adds, those of the default branch and those of every
earlier unmerged queue item, from Git, the Spec folders and the Run Database,
and returns the report, so a real queue parks before the Run. It answers the
case reported by fiscus on 2026-11-05, five queued Specs that all started at
migration `0052`
([references](references/2026-11-05-fiscus-delivery-queue-report.md)).

## Requirements

1. MUST, before changing code, run
   `go test -count=1 -timeout 20m -run '^(TestRevalidate.*|TestDeliverStatus.*)$' ./internal/cli`
   and record in the Result that it passes, and that the
   `newDeliveryRevalidationGitFixture` workflow carries no Run Database.
2. MUST add `numberedFileReport` to the command workflow in a new file
   `internal/cli/deliver_numbered_files.go`, following the seven steps of
   `_techspec.md` → API Contract 8, and call it from `Revalidate` after the
   premise comparison, setting `NumberedFiles` and wrapping an error as
   `check numbered files: %w`. Without a declaration it runs no Git command
   after reading the configuration and never touches `workflow.store`
   (Invariant 1); it writes no file and no ref other than the default
   branch's remote-tracking ref (Invariant 4).
3. MUST create `internal/cli/deliver_numbered_files_test.go`, each top-level
   test calling `t.Parallel()` first, over temporary repositories with a bare
   `origin` on disk (as `newDeliveryRevalidationGitFixture` builds them), a
   default branch whose `.roundfixrc.yml` declares
   `delivery:\n  numbered_files: [db/*.sql]\n` and holds `db/0001_init.sql`,
   and active Specs `alpha` and `beta` (and `gamma` where named), each with one
   pending backend Task whose Context declares `creates:` of its migration,
   queued in a temporary Run Database from `store.Open` in the order alpha,
   beta, gamma:
   - `TestRevalidateReportsANumberedFileCollision`, revalidating `beta` in a
     linked worktree (`git worktree add`) of the queue repository, itself a
     clone of `origin`, as the queue creates item worktrees, so the earlier
     items' branches resolve there, with subtests:
     - `off`: the same repositories without the key and a workflow whose store
       is nil give an empty report and no error;
     - `ahead_branch`: `alpha` parked with a recorded item branch that adds
       `db/0002_alpha.sql`, while `alpha`'s plan on the default branch
       declares nothing; `beta` declares `db/0002_beta.sql`: the report's
       `Blocker` is exactly
       `numbered-file-collision: db/0002_beta.sql collides with db/0002_alpha.sql (alpha); next free number: db/0003`;
     - `ahead_plan`: `alpha` queued with no branch, its plan declaring
       `db/0002_alpha.sql`: the same `Blocker`;
     - `default_branch`: `beta`'s worktree made before `origin` gains
       `db/0002_main.sql`, no earlier item claiming anything: the `Blocker` is
       `numbered-file-collision: db/0002_beta.sql collides with db/0002_main.sql (default branch); next free number: db/0003`;
     - `ignored`: `alpha` `merged` with a recorded branch adding
       `db/0002_alpha.sql` that the default branch does not hold, `gamma`
       after `beta` declaring `db/0002_gamma.sql`, a completed Task of `beta`
       declaring `db/0002_done.sql`, and `beta` also declaring
       `db/meta/_journal.json`: an empty report;
     - `renumbered`: the `ahead_branch` setup with `beta` declaring
       `db/0003_beta.sql`: an empty report.
   - `TestDeliverStatusNamesTheNumberedFileCollisionPark`: a queue whose item
     `beta` is parked with the `ahead_branch` blocker and worktree
     `/worktrees/beta`; `runCLI` of `deliver status` exits `0` and its stdout
     contains the lines
     `Park: beta conflict: in /worktrees/beta, renumber the item's files the blocker names from its next free number, in the Task declarations on the item branch and in any file the branch holds, or deliver the earlier item it names first and merge the default branch into the item branch; then run roundfix deliver retry beta`
     and `Answer: ` followed by the same next action; it asserts those lines,
     not the whole output, because other Specs change the status lines.
4. MUST run sabotages 9 to 12 of `_techspec.md` → Testing Approach on a
   temporary copy and record which named test fails.
5. MUST keep every existing `TestRevalidate*`, `TestDeliverStatus*` and
   `TestDeliver*` test of `./internal/cli` passing unchanged.
6. MUST NOT change `internal/delivery`, `internal/config`, `CreateItemBranch`,
   `ResolveConflict`, the Delivery Plan, the Run Database schema, any guide,
   skill or reference, the Makefile, CI workflows, `go.mod` or `go.sum`.

## Subtasks

- [ ] Characterize Delivery Revalidation without a declaration.
- [ ] Read the declaration at the default commit; return early when empty.
- [ ] Collect the subject's and earlier items' claims from Git, the plans and the queue.
- [ ] Add the command tests and run the sabotages.

## Acceptance Criteria

- [ ] With a declaration, a real revalidation reports each collision with the earlier item's branch or plan, or with the default branch, with the exact blocker; later and merged items, completed Tasks, unnumbered names and a renumbered plan report nothing.
- [ ] Without a declaration nothing new runs and the store is never read.
- [ ] `deliver status` prints the collision park and its next action.

## Context

- interface: `internal/cli/deliver_revalidate.go`
- creates: `internal/cli/deliver_numbered_files.go`
- creates: `internal/cli/deliver_numbered_files_test.go`
- instruction: `internal/cli/deliver_revalidate_test.go`
- instruction: `internal/cli/deliver_park_status_test.go`

## Consumers

- `Revalidation.NumberedFiles` as the command workflow fills it — consumers: `./internal/cli`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestRevalidateReportsANumberedFileCollision|TestDeliverStatusNamesTheNumberedFileCollisionPark|TestRevalidate.*|TestDeliverStatus.*)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestRevalidateReportsANumberedFileCollision TestDeliverStatusNamesTheNumberedFileCollisionPark TestRevalidateReportsANumberedFileCollision/off TestRevalidateReportsANumberedFileCollision/ahead_branch TestRevalidateReportsANumberedFileCollision/ahead_plan TestRevalidateReportsANumberedFileCollision/default_branch TestRevalidateReportsANumberedFileCollision/ignored TestRevalidateReportsANumberedFileCollision/renumbered; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done` — expected: exit 0. Before this Task neither new test exists, so the first pass line is missing; afterwards a real revalidation reports each collision with the exact blocker, ignores later and merged items, completed Tasks and unnumbered names, does nothing without a declaration, and `deliver status` prints the park.

## References

- `_prd.md` → Goal 2; User Story 2; User Story 4; User Story 5; Core Feature 3; Core Feature 4; Success Metric 2; Success Metric 4
- `_techspec.md` → Current behavior 3; Current behavior 6; API Contract 8; Invariants 1; Invariants 2; Invariants 4; Testing Approach; Build Order 4


## Result

Skills loaded: implement-task, coding-guidelines, golang-cli, golang-error-handling, golang-structs-interfaces, golang-testing, testing-boss, tech-writer, evidence-gate.

Implemented the seven-step numbered-file collector in the command workflow and
called it after the premise comparison. Errors carry `check numbered files:`.
The collector reads declarations from the default commit, refreshes only its
remote-tracking ref, combines committed additions with pending non-QA Task
claims absent at the merge base, and compares only earlier unmerged items.
Earlier plans come from their recorded worktree when its manifest exists,
otherwise the subject's Specs Root; an unreadable plan retains branch claims.
Fetch disables tag following and FETCH_HEAD writes. No queue row, local branch,
Task declaration or worktree file is changed by the collector.

### Characterization and focused checks

- Before changing code, `go test -count=1 -timeout 20m -run '^(TestRevalidate.*|TestDeliverStatus.*)$' ./internal/cli` passed (`ok`, 2.935s).
  Inspection confirmed that `newDeliveryRevalidationGitFixture.workflow`
  initializes `loaded` and `git` only and carries no Run Database (`store` is nil).
- Tests were added before the collector. `go test -count=1 -timeout 5m -run '^TestRevalidateReportsANumberedFileCollision$' ./internal/cli` failed in
  `ahead_branch`, `ahead_plan` and `default_branch`: each returned an empty
  blocker instead of the expected collision.
- With the collector, `go test -count=1 -timeout 5m -v -run '^(TestRevalidateReportsANumberedFileCollision|TestDeliverStatusNamesTheNumberedFileCollisionPark)$' ./internal/cli`
  passed, including every required subtest and the added boundary cases.
- After the final code and test edits, `go test -count=1 -timeout 20m -run '^(TestRevalidate|TestDeliver)' ./internal/cli` passed (`ok`, 17.004s).
  This focused package run covers the existing `TestRevalidate*`,
  `TestDeliverStatus*` and `TestDeliver*` tests without changing them.
- `git diff --check` passed.

### Acceptance evidence

- Declared collisions: real bare origins, queue clones, linked item worktrees
  and temporary Run Databases exercise `ahead_branch`, `ahead_plan` and
  `default_branch`, asserting the exact blockers. `subject_branch` checks an
  addition held only by HEAD; `recorded_plan`, `missing_plan` and
  `missing_branch` check plan selection and fallback. `ahead_branch` also
  commits an empty item configuration to prove it cannot disable the default
  branch's declaration.
- Ignored claims: `ignored` excludes the merged alpha branch, later gamma,
  completed beta Task and unnumbered journal; `renumbered` reports nothing.
  `ignored/completed` separately places a colliding file on the refreshed
  default branch and changes beta's pending claim to 0003, so counting its
  completed 0002 declaration is independently observable. `base_file`
  excludes a declaration already present at the merge base. `absent_item`
  proves that a subject missing from the queue has no earlier items.
- Disabled behavior: `off` uses no configuration key and a nil workflow store;
  it returns an empty report. The recording Git adapter verifies that the
  last command reads the default configuration, with no subsequent Git call.
  `cleared_after_fetch` proves the refreshed empty declaration also returns
  before touching a nil store. Each initial revalidation leaves the queue,
  local branch/tag refs and worktree unchanged.
- Status: `TestDeliverStatusNamesTheNumberedFileCollisionPark` exits 0,
  produces no stderr, and asserts both exact `Park: beta conflict: ...` and
  `Answer: ...` lines for `/worktrees/beta` and the alpha collision blocker.

### Required sabotages

Each mutation ran separately on an initialized, staged temporary copy of HEAD
plus this Task's three Go files; the Run checkout was not mutated. Logs are at
`/private/tmp/roundfix-0357-sabotages-ahs_8j1u/sabotage-<number>.log`.
The command for 9–11 was `go test -count=1 -timeout 5m -v -run '^TestRevalidateReportsANumberedFileCollision$/ignored$' ./internal/cli`;
12 used the same command with `/off$`.

| Sabotage | Observed named failure (each command exited 1) |
| --- | --- |
| 9: compare later items too | `TestRevalidateReportsANumberedFileCollision/ignored`: beta collides with `db/0002_gamma.sql (gamma)`. |
| 10: compare merged items too | `TestRevalidateReportsANumberedFileCollision/ignored`: beta collides with `db/0002_alpha.sql (alpha)`. |
| 11: count completed Task declarations | `TestRevalidateReportsANumberedFileCollision/ignored/completed`: `db/0002_done.sql` collides with `db/0002_main.sql (default branch)`. |
| 12: read the store before checking the declaration | `TestRevalidateReportsANumberedFileCollision/off`: nil-pointer panic in the Run Database read. |

Task status remains Daemon-owned. The authored Verification command, repository
Verification, commits, pushes and pull requests were not run. No other Task,
Task Graph, guide, skill, reference, config/delivery package, tooling or module
file was changed. No follow-up implementation was included.
