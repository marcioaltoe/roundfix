---
task: task_03
spec: 0262-a-run-worktree-under-a-symlinked-root-is-cleaned
status: pending
type: backend
complexity: medium
---

# Task 03: Every consumer cleans under a symlinked root, and the reclaim is documented

## Overview

task_01 and task_02 change the worktree package's reconciliation and path
comparison. Their consumers are the Reconcile Command, the Delivery Queue
release of merged Spec Runs and the daemon's superseded-branch disposition
(`_techspec.md` → API Contract 3). This Task is the Integration Checkpoint:
it proves each consumer end to end under a symbolic-linked worktree location,
runs the consuming suites, and documents that `roundfix reconcile --apply`
reclaims the Run Worktrees an earlier version kept. No consumer's production
code changes.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-12,
   [a Run Worktree under a symlinked worktree root is never cleaned](references/2026-10-12-a-run-worktree-under-a-symlinked-root-is-never-cleaned.md),
   for its two observed symptoms: the Delivery Queue cleanup warning and the
   `unknown` reconcile result.
2. MUST add `internal/cli/reconcile_symlinked_root_test.go` with tests that
   call `t.Parallel()` first and record, through the Run Database, a terminal
   implement Run whose Run Worktree `worktree.Create` made under a worktree
   location that is a symbolic link to a real temporary directory, keeping the
   link spelling on the Run row as Roundfix records it:
   - `TestReconcileClassifiesARunUnderASymlinkedWorktreeRoot`: the JSON report
     of `roundfix reconcile <run-id> --format json` gives the clean merged Run
     the classification `safe`, and no result carries the refusal
     `recorded Run Worktree path is unsafe`.
   - `TestReconcileApplyReleasesARunUnderASymlinkedWorktreeRoot`:
     `roundfix reconcile --apply` exits 0 and the Run Worktree directory, its
     Git registration and the Run Branch are gone afterwards.
   - `TestReleaseMergedSpecRunsCleansUnderASymlinkedWorktreeRoot`:
     `releaseMergedSpecRuns` for the merged Spec returns nil, so the Delivery
     Queue prints no `warning: cleanup failed`, and the Run Worktree is gone.
3. MUST add `internal/daemon/reconcile_symlinked_root_test.go` with
   `TestClassifySupersededBranchUnderASymlinkedWorktreeRoot`, calling
   `t.Parallel()` first: a Run whose branch is contained in its target, under
   a linked location, is not refused with `worktree condition failed`.
4. MUST add to `docs/user-guide/commands/reconcile.md` the paragraph of
   `_techspec.md` → API Contract 4, containing the phrase
   `resolved within its Worktree Root` and the command
   `roundfix reconcile --apply`, and to `docs/user-guide/configuration.md`, in
   the section that lists the Run and Task Worktree paths, one sentence
   containing the phrase `resolved within its Worktree Root`.
5. MUST record a Coverage Review in `docs/references/skill-coverage.json` for
   the entries `command reconcile`, `guide docs/user-guide/commands/reconcile.md`,
   `guide docs/user-guide/configuration.md` and `config keys`, each naming
   `Spec 0262 task_03` and stating that a symbolic-linked `worktree.location`
   is now resolved within its Worktree Root and that no flag, key, output or
   exit code changes, and MUST re-record the Behavior Surface Record with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.
   No skill file changes (`_prd.md` → Skills).
6. MUST search the user guides, `CONTEXT.md`, the agent guides and the
   Roundfix skill references for wording that says a symbolic link in a
   recorded Run Worktree path is always unsafe, and change any it finds in
   this Task; the search at `94fbf85b` found none.
7. MUST run the new tests once with `-race` (ADR-0259) and MUST NOT change
   production code in `internal/cli`, `internal/daemon` or `internal/delivery`.

## Subtasks

- [ ] Add the three CLI tests and the daemon test on a linked location.
- [ ] Add the reconcile and configuration guide text.
- [ ] Record the Coverage Reviews and re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] `roundfix reconcile` classifies and `--apply` releases a Run under a linked location.
- [ ] The merged-Spec release and the superseded-branch disposition succeed there.
- [ ] The consuming suites pass, and the guides document the reclaim.

## Consumers

- Behavior Surface Record and Skill Coverage Map — consumers: `./internal/docscontract`

## Context

- instruction: `docs/adr/0268-a-recorded-worktree-path-is-resolved-within-its-worktree-root.md`
- instruction: `internal/cli/reconcile_test.go`
- instruction: `internal/cli/reconcile_merged_test.go`
- instruction: `internal/cli/deliver_release.go`
- instruction: `internal/daemon/reconcile.go`
- instruction: `internal/daemon/reconcile_test.go`
- creates: `internal/cli/reconcile_symlinked_root_test.go`
- creates: `internal/daemon/reconcile_symlinked_root_test.go`
- interface: `docs/user-guide/commands/reconcile.md`
- interface: `docs/user-guide/configuration.md`
- interface: `docs/references/skill-coverage.json`
- interface: `docs/references/behavior-surfaces.json`

## Verification

- `out="$(go test -count=1 -race -v -run '^(TestReconcileClassifiesARunUnderASymlinkedWorktreeRoot|TestReconcileApplyReleasesARunUnderASymlinkedWorktreeRoot|TestReleaseMergedSpecRunsCleansUnderASymlinkedWorktreeRoot)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestReconcileClassifiesARunUnderASymlinkedWorktreeRoot TestReconcileApplyReleasesARunUnderASymlinkedWorktreeRoot TestReleaseMergedSpecRunsCleansUnderASymlinkedWorktreeRoot; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing passing test %s\n' "$name" >&2; exit 1; }; done; go test -count=1 -run 'Reconcile|Release|ForceStop|Prune|Item|Stop' ./internal/cli || exit 1` — expected: exit 0. Before this Task the three tests do not exist, so the command fails; afterwards they pass and the reconcile, release, stop, prune and item tests of the consuming CLI suite pass.
- `out="$(go test -count=1 -race -v -run '^TestClassifySupersededBranchUnderASymlinkedWorktreeRoot$' ./internal/daemon 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestClassifySupersededBranchUnderASymlinkedWorktreeRoot ' || { printf 'missing passing daemon test\n' >&2; exit 1; }; go test -count=1 -run 'Reconcile|Supersed|Disposition|Worktree' ./internal/daemon || exit 1` — expected: exit 0. Before this Task the test does not exist, so the command fails.
- `tr -s '[:space:]' ' ' < docs/user-guide/commands/reconcile.md | grep -qF -- 'resolved within its Worktree Root' || { printf 'the reconcile guide does not document the Worktree Root\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/configuration.md | grep -qF -- 'resolved within its Worktree Root' || { printf 'the configuration guide does not name the Worktree Root\n' >&2; exit 1; }; test "$(grep -c 'Spec 0262 task_03' docs/references/skill-coverage.json)" -ge 4 || { printf 'the four Coverage Reviews are missing\n' >&2; exit 1; }; go test -count=1 -tags docscontract -run '^TestTheSkillCoverageMapIsCurrent$' ./internal/docscontract || exit 1` — expected: exit 0. Before this Task the guide text and the reviews are absent, so the command fails; afterwards the Skill Coverage Map is current.

## References

- `_prd.md` → Goals; User Story 1; User Story 2; User Story 4; Core Feature 4; Core Feature 5; Success Metric 1; Success Metric 3; Skills
- `_techspec.md` → API Contract 3; API Contract 4; Invariant 4; Testing Approach; Build Order 3
- ADR-0268; ADR-0053; ADR-0256; ADR-0259
