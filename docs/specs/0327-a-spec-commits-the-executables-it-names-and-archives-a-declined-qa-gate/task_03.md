---
task: task_03
spec: 0327-a-spec-commits-the-executables-it-names-and-archives-a-declined-qa-gate
status: pending
type: backend
complexity: medium
---

# Task 03: Settle names each worktree it removes

## Overview

After `roundfix settle` integrates a kept Task Worktree and the Run, it
removes them and exits `0` without saying so. Quaestor's shell was still in
the removed Run Worktree, so its next Git command exited `128` and the settle
looked like it had failed. This Task makes settle name each worktree it
removes on stderr, keeps its exit codes, and writes the rule into the settle
usage, the Roundfix `settle` and `spec-delivery` references, the `settle`
guide and `CONTEXT.md`.

## Requirements

1. MUST answer the Backlog Entry
   [`references/2026-10-26-a-spec-cannot-commit-its-executables-settle-exits-128-and-archive-refuses-a-declined-qa.md`](references/2026-10-26-a-spec-cannot-commit-its-executables-settle-exits-128-and-archive-refuses-a-declined-qa.md)
   of 2026-10-26, defect (a)'s settle exit, by implementing API Contracts 6
   and 9 of `_techspec.md`.
2. MUST print, in `runSettleCommand` of `internal/cli/settle.go`,
   `Task Worktree removed: <task worktree path>` on stderr after
   `integrateTaskWorktree` succeeds, and `Run Worktree removed: <run worktree path>`
   after `integrateSettledRun` returns no pending command; never for a
   conflict, a pending integration or an error. Stdout and every exit code
   keep their bytes (Invariant 4).
3. MUST change, in `internal/cli/settle_test.go`, the stderr expectation of
   `TestSettleTaskStatusRetargetsKeptRunWorktreeAndCleansUpAfterIntegration`,
   `TestRunSettleRetargetsKeptTaskWorktreeAndCleansUpAfterIntegration`,
   `TestSettleAcceptsCompletedTaskFromKeptTaskWorktreeAfterHookRefusal` and
   `TestSettleCommitsDeletedAndRenamedWorkFromTaskWorktree` to the exact
   stderr with the removal lines, in order (`_techspec.md` → Declared breaks
   4), and add:
   - `TestSettleNamesEachWorktreeItRemoves`: a kept Task Worktree settled
     from the user checkout exits `0`, and stderr is exactly
     `Settle surface: <task worktree>`, `Task Worktree removed: <task worktree>`
     and `Run Worktree removed: <run worktree>`, each path absent afterwards;
   - `TestSettleNamesNoRemovalWhenIntegrationIsPending`: a kept Run Worktree
     whose user branch gained a commit after the Run started, so the
     fast-forward is refused, exits `1` with `integration pending — ` on
     stdout, stderr exactly `Settle surface: <run worktree>`, and the Run
     Worktree still present.
4. MUST add API Contract 9's `settleUsage` sentence to `internal/cli/settle.go`.
5. MUST write API Contract 9's sentence, verbatim, with the two stderr lines
   and the advice to run settle from the user checkout, into
   `.agents/skills/roundfix/references/settle.md` (after the Run-level
   integration paragraph), `.agents/skills/roundfix/references/spec-delivery.md`
   (its recovery step) and `docs/user-guide/commands/settle.md`; then run
   `make skills-sync` and
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
   which raises the `roundfix` version pair.
6. MUST revise, through `domain-modeling`, **Run Worktree** ("removed after a
   Clean integrated spec outcome or after settle integrates it") and
   **Settle Command** ("names each worktree it removes on stderr") in
   `CONTEXT.md`, keeping their `_Avoid_` lines.
7. MUST, as the last step after the last guide edit, re-record the Behavior
   Surface fingerprint with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.
   Hand-edit no mirror, version field or record.
8. MUST NOT change settle's staging, surface resolution or integration, any
   other command, any text inside a `### QA settlement` section, the Makefile
   or the CI workflows.

## Subtasks

- [ ] Print the removal lines and update the four pinned settle tests.
- [ ] Add the two new settle tests and the usage sentence.
- [ ] Write the rule into the references, the guide and the glossary; sync and record; re-record the Behavior Surface last.

## Acceptance Criteria

- [ ] Settle names each worktree it removed on stderr and keeps its stdout and exit codes.
- [ ] A pending integration names no removal and keeps the Run Worktree.
- [ ] The usage, references, guide and glossary carry the rule, and the owned-skill and Behavior Surface records are current.

## Consumers

- Settle's stderr — consumers: `./internal/cli`

## Context

- instruction: `.agents/skills/domain-modeling/SKILL.md`
- interface: `internal/cli/settle.go`
- interface: `internal/cli/settle_test.go`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `.agents/skills/roundfix/references/settle.md`
- interface: `.agents/skills/roundfix/references/spec-delivery.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/roundfix/references/settle.md`
- interface: `skills/roundfix/references/spec-delivery.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/user-guide/commands/settle.md`
- interface: `CONTEXT.md`
- interface: `docs/references/behavior-surfaces.json`

## Verification

- `out="$(go test -count=1 -timeout 15m -v -run '^(TestSettleNamesEachWorktreeItRemoves|TestSettleNamesNoRemovalWhenIntegrationIsPending|TestSettleTaskStatusRetargetsKeptRunWorktreeAndCleansUpAfterIntegration|TestRunSettleRetargetsKeptTaskWorktreeAndCleansUpAfterIntegration|TestSettleAcceptsCompletedTaskFromKeptTaskWorktreeAfterHookRefusal|TestSettleCommitsDeletedAndRenamedWorkFromTaskWorktree)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestSettleNamesEachWorktreeItRemoves TestSettleNamesNoRemovalWhenIntegrationIsPending TestSettleTaskStatusRetargetsKeptRunWorktreeAndCleansUpAfterIntegration TestRunSettleRetargetsKeptTaskWorktreeAndCleansUpAfterIntegration TestSettleAcceptsCompletedTaskFromKeptTaskWorktreeAfterHookRefusal TestSettleCommitsDeletedAndRenamedWorkFromTaskWorktree; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < internal/cli/settle.go | grep -qF -- 'settle removes the kept Task Worktree and Run Worktree it integrated and names each on stderr' || { echo 'settle usage does not name the removal' >&2; exit 1; }` — expected: exit 0. Before this Task the two new tests do not exist, so the first pass check fails. Afterwards both removal cases and the pending case pass, the four updated tests pass, and the usage names the removal.
- `for file in .agents/skills/roundfix/references/settle.md skills/roundfix/references/settle.md .agents/skills/roundfix/references/spec-delivery.md skills/roundfix/references/spec-delivery.md docs/user-guide/commands/settle.md; do tr -s '[:space:]' ' ' < "$file" | grep -qF -- 'because Git exits 128 in a removed directory' || { printf '%s lacks the settle rule\n' "$file" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'removed after a Clean integrated spec outcome or after settle integrates it' || { echo 'Run Worktree entry not revised' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'names each worktree it removes on stderr' || { echo 'Settle Command entry not revised' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Run Worktree**:' || { echo 'CONTEXT.md lacks **Run Worktree**' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Settle Command**:' || { echo 'CONTEXT.md lacks **Settle Command**' >&2; exit 1; }; diff -r .agents/skills/roundfix skills/roundfix && go test -timeout 15m -count=1 -run '^(TestEveryOwnedSkillVersionIsRecorded|TestSettlementGuidanceIsOneTable)$' ./skills && go test -timeout 15m -count=1 -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task no reference or guide states the settle rule, so the first phrase check fails. Afterwards the skill and its mirror agree, the owned-skill versions are recorded, the `### QA settlement` table is unchanged, and the documentation contracts pass with the Behavior Surface record re-recorded last.

## References

- `_prd.md` → Goals 3 and 5; Core Features 3 and 5; Success Metrics 6 and 7; Glossary
- `_techspec.md` → Current behavior 5, 6 and 8; API Contracts 6 and 9; Declared breaks 4 and 5; Invariant 4; Testing Approach; Build Order 3
