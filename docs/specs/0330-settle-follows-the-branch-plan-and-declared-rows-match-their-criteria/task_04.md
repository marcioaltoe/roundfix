---
task: task_04
spec: 0330-settle-follows-the-branch-plan-and-declared-rows-match-their-criteria
status: completed
type: backend
complexity: medium
---

# Task 04: Settle commits what a Task commit keeps

## Overview

Settle stages every change of its surface with `git add --all`, so it
commits an untracked executable file that the Daemon's Task commit refuses
unless the Spec's authorization lists it under `executables` (ADR-0340).
This Task makes settle refuse such a file before it writes anything, with
the authorization it already reads and the Daemon's own stage filter, keeps
a tracked file whose mode a Task changes as work output in both paths
(ADR-0157), and writes the rule into the Roundfix `settle` reference, the
`settle` guide and `CONTEXT.md`.

## Requirements

1. MUST answer the Backlog Entry
   [`references/2026-10-26-settle-and-a-tracked-mode-change-commit-an-executable-the-task-commit-refuses.md`](references/2026-10-26-settle-and-a-tracked-mode-change-commit-an-executable-the-task-commit-refuses.md)
   of 2026-10-26 by implementing API Contracts 7 and 11 of `_techspec.md`.
2. MUST add `UnauthorizedExecutables` to `internal/daemon/task_engine.go` as
   API Contract 7 states, built on `filterStageablePaths` and
   `authorizedExecutables`, with no change to the Task, QA Report or Batch
   commit.
3. MUST make `preflightSettle` in `internal/cli/settle.go` list the surface's
   untracked files with `git ls-files --others --exclude-standard -z` and
   refuse, before the branch comparison of task_03 and before any write, every
   path `UnauthorizedExecutables` returns, with exit `2` and the message of
   API Contract 7.
4. MUST add `internal/cli/settle_executables_test.go`, through CLI dispatch
   on the implement workspace fixture, with:
   - `TestSettleRefusesAnUnlistedExecutableInANewDirectory`: an untracked
     `newdir/helper.sh` with mode `0755` and no authorization; exit `2`
     naming `newdir/helper.sh (mode`, HEAD, status and Task file unchanged;
   - `TestSettleCommitsAListedExecutableWithItsMode`: a committed granted
     `_authorization.md` listing `scripts/launch.sh` under `executables`;
     settle exits `0` and the commit records `scripts/launch.sh` as
     `100755`;
   - `TestSettleCommitsATrackedFileMadeExecutable`: a tracked script whose
     mode the Task changed to `0755`; settle exits `0` and commits it as
     `100755`.
5. MUST add `TestTaskCycleCommitsATrackedFileMadeExecutable` to
   `internal/daemon/task_engine_executables_test.go`: a real-Git Task Cycle
   whose Agent makes a tracked file executable completes the Task and commits
   it as `100755`, as settle does.
6. MUST add API Contract 11's sentence, verbatim, to the `## Settle Command`
   section of `.agents/skills/roundfix/references/settle.md` and
   `docs/user-guide/commands/settle.md`, each with the refusal message; then
   run `make skills-sync` and
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
   which raises the `roundfix` version pair.
7. MUST revise **Authorized Executable** in `CONTEXT.md` through
   `domain-modeling` with "A settle commit keeps the same paths.", keeping
   its `_Avoid_` line.
8. MUST, as the last step after the last guide edit, re-record the Behavior
   Surface fingerprint with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.
   Hand-edit no mirror, version field or record.
9. MUST NOT change the Daemon's commit filters, the authorization record's
   fields, settle's staging of other files, any text inside a
   `### QA settlement` section, the Makefile or the CI workflows. Each new
   top-level test calls `t.Parallel()` first.

## Subtasks

- [ ] Export the Daemon's executable judgement and refuse unlisted executables in settle preflight.
- [ ] Add the three CLI tests and the Daemon test.
- [ ] Write the rule into the reference, guide and glossary; sync and record; re-record the Behavior Surface last.

## Acceptance Criteria

- [ ] Settle refuses an unlisted untracked executable, also inside a new directory, with nothing written.
- [ ] A listed executable and a tracked file made executable are committed with mode 100755 by settle, and the tracked one by the Task commit too.
- [ ] The reference, guide and glossary carry the rule, and the owned-skill and Behavior Surface records are current.

## Consumers

- `UnauthorizedExecutables` — consumers: `./internal/cli`, `./internal/daemon`

## Context

- instruction: `docs/adr/0340-an-untracked-executable-rides-a-task-commit-only-when-the-spec-authorization-names-it.md`
- instruction: `docs/adr/0157-tracked-source-is-work-output-whatever-its-mode.md`
- instruction: `.agents/skills/domain-modeling/SKILL.md`
- interface: `internal/daemon/task_engine.go`
- interface: `internal/daemon/task_engine_executables_test.go`
- interface: `internal/cli/settle.go`
- creates: `internal/cli/settle_executables_test.go`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `.agents/skills/roundfix/references/settle.md`
- interface: `skills/roundfix/references/settle.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/user-guide/commands/settle.md`
- interface: `CONTEXT.md`
- interface: `docs/references/behavior-surfaces.json`

## Verification

- `out="$(go test -count=1 -timeout 15m -v -run '^(TestSettleRefusesAnUnlistedExecutableInANewDirectory|TestSettleCommitsAListedExecutableWithItsMode|TestSettleCommitsATrackedFileMadeExecutable)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestSettleRefusesAnUnlistedExecutableInANewDirectory TestSettleCommitsAListedExecutableWithItsMode TestSettleCommitsATrackedFileMadeExecutable; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; out="$(go test -count=1 -timeout 15m -v -run '^TestTaskCycleCommitsATrackedFileMadeExecutable$' ./internal/daemon 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestTaskCycleCommitsATrackedFileMadeExecutable ' || { echo 'missing pass: TestTaskCycleCommitsATrackedFileMadeExecutable' >&2; exit 1; }; go test -count=1 -timeout 15m -run 'Settle|Executable' ./internal/cli ./internal/daemon` — expected: exit 0. Before this Task the four tests do not exist, so the first pass check fails. Afterwards settle refuses the unlisted executable, commits the listed and the tracked ones with their mode, the Task commit agrees, and every settle and executable test of both packages passes.
- `for file in .agents/skills/roundfix/references/settle.md skills/roundfix/references/settle.md docs/user-guide/commands/settle.md; do tr -s '[:space:]' ' ' < "$file" | grep -qF -- 'Settle refuses an untracked executable file that the Spec' || { printf '%s lacks the executables rule\n' "$file" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$file" | grep -qF -- 'a tracked file whose mode changes is committed by both.' || { printf '%s lacks the tracked-mode rule\n' "$file" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'A settle commit keeps the same paths.' || { echo 'CONTEXT.md lacks the revised **Authorized Executable**' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Authorized Executable**:' || { echo 'CONTEXT.md lacks **Authorized Executable**' >&2; exit 1; }; diff -r .agents/skills/roundfix skills/roundfix && go test -timeout 15m -count=1 -run '^(TestEveryOwnedSkillVersionIsRecorded|TestSettlementGuidanceIsOneTable)$' ./skills && go test -timeout 15m -count=1 -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task no reference or guide states the rule, so the first phrase check fails. Afterwards the skill and its mirror agree, the versions are recorded and the documentation contracts pass with the Behavior Surface record re-recorded last.

## References

- `_prd.md` → Goals 5 and 6; Core Features 5 and 6; Success Metrics 6 and 7; Glossary
- `_techspec.md` → Current behavior 8; API Contracts 7 and 11; Declared breaks 5 and 6; Invariant 4; Testing Approach; Build Order 4
- ADR-0340, ADR-0157


## Result

Implemented the Task 04 slice; status and settlement remain Daemon-owned.

- Acceptance criterion 1: added `daemon.UnauthorizedExecutables`, reusing
  `filterStageablePaths` and the granted authorization's executable paths.
  Settle enumerates untracked files with
  `git ls-files --others --exclude-standard -z` and refuses every executable
  drop with the prescribed exit-2 diagnostic before branch comparison or
  writes. The new CLI refusal test checks an executable inside a new
  directory and unchanged HEAD, full porcelain status and Task bytes.
- Acceptance criterion 2: CLI tests commit a granted, explicitly listed
  executable and a tracked script changed to executable, asserting Git mode
  `100755`. The real-Git Task Cycle test also commits a tracked script's
  mode change as `100755` and observes one completed Task and no failure.
  Daemon Task, QA Report and Batch commit filters are unchanged.
- Acceptance criterion 3: added API Contract 11 verbatim and its refusal
  message to the canonical settle reference and user guide. Revised
  **Authorized Executable** with the prescribed settle sentence, preserving
  its `_Avoid_` line. Synchronized skill mirrors and recorded both Roundfix
  version fields as `0.1.155`. Regenerated the Behavior Surface fingerprint
  after the final guide edit, using only the sanctioned recorders.

Focused evidence:

- Red check before the production change:
  `go test -count=1 -timeout 3m ./internal/cli -run '^TestSettleRefusesAnUnlistedExecutable'`
  exited 1: settle exited 0 and committed `newdir/helper.sh`, reproducing the
  recorded defect.
- After the production change:
  `go test -count=1 -timeout 4m ./internal/cli ./internal/daemon -run '^(TestSettle(RefusesAnUnlisted|CommitsAListed|CommitsATracked)|TestTaskCycleCommitsATracked)'`
  exited 0 for both packages, exercising all four new regressions.
- `make skills-sync` exited 0.
- `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
  exited 0 and recorded the version pair.
- As the final regeneration after guide edits:
  `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`
  exited 0.
- A read-only Python comparison confirmed every canonical Roundfix skill
  file equals its mirror and all three references contain the exact
  executable-rule sentence. `git diff --check` exited 0.

The initial changed-path baseline contained only this Task file. Final scope
is limited to this Task's declared interfaces, new CLI test and Result.
Declared Task Verification, repository Verification and incremental
Verification were not run. No Task status, Task Graph or other Task file was
edited; no commit, push or Pull Request was made. No follow-up identified.
