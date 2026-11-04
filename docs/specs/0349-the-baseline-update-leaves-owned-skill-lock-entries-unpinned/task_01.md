---
task: task_01
spec: 0349-the-baseline-update-leaves-owned-skill-lock-entries-unpinned
status: completed
type: backend
complexity: medium
---

# Task 01: The bundle install removes the ref of each owned lock entry it rewrites and reports it

## Overview

`RecordInstalledOwnedSkills` rewrites only the `computedHash` of an owned
skill's `skills-lock.json` entry and keeps its `ref`. The skills CLI then
reinstalls the older catalog copy that `ref` names. This Task removes the
`ref` member of every entry the writer rewrites, keeps every other byte,
returns which entries got a hash and which lost a `ref`, and makes
`roundfix skills install` and `roundfix baseline update` report both. It
answers the marcioaltoe-skills report of 2026-11-04 (ai-memory message
`01a1a5f7`, adopted as
[references/2026-11-04-the-baseline-update-keeps-a-stale-ref-in-owned-skill-lock-entries.md](references/2026-11-04-the-baseline-update-keeps-a-stale-ref-in-owned-skill-lock-entries.md)).

## Requirements

1. MUST characterize first: add
   `TestRecordInstalledOwnedSkillsRemovesTheRefOfAnInstalledOwnedSkill` in
   `internal/baseline/owned_skills_lock_unpin_test.go` with the four rows of
   `_techspec.md` → Testing Approach. Run its `skills tool format` row against
   the unchanged writer and record in the Task Result that every `ref` stayed
   and only the `roundfix` hash changed. Then set the expectations of API
   Contract 1.
2. MUST change `RecordInstalledOwnedSkills` in
   `internal/baseline/owned_skills_lock.go` to return
   `OwnedSkillsLockUpdate` with `Hashed`, `Unpinned` and `Names()`
   (`_techspec.md` → Interfaces). It removes the `ref` member of each selected
   entry with the span API Contract 1 defines, whatever the member's value,
   and inserts a missing `computedHash` with a leading comma only when another
   member remains. Entry selection, the hash rule, the temporary-file write,
   the kept file mode and "no edit, no write" stay as they are.
3. MUST report both lists (`_techspec.md` → API Contract 2):
   - `runSkillsInstallCommand` in `internal/cli/cli.go` prints the `Hashed`
     names as `Updated skills-lock.json computedHash: <name>`, then the
     `Unpinned` names as `Removed skills-lock.json ref: <name>`;
   - `runBaselineUpdateSkillsStageWith` in `internal/cli/baseline_update.go`
     sets `LockUpdated` to `Names()` and the new
     `LockUnpinned []string` field (`json:"lockUnpinned,omitempty"`, after
     `lockUpdated`) to `Unpinned`;
   - `writeBaselineUpdateResult` prints one `- lock ref removed <name>` line
     per `LockUnpinned` name after the `- lock updated` lines;
   - `recordLock` in `baselineUpdateSkillsDependencies` takes the new return
     type, and the stubs in `internal/cli/baseline_update_test.go` (two) and
     `internal/cli/doctor_trailing_skills_test.go` (one) follow it.
4. MUST add `TestSkillsInstallRemovesTheRefOfOwnedSkillLockEntries`, which
   asserts the whole Surface Transcript 1 of `_techspec.md` (stdout line by
   line, empty stderr, exit 0) and the whole lock file, with the external
   entry byte-identical. Also add
   `TestBaselineUpdateSkillsStageReportsTheUnpinnedLockEntries`, both in
   `internal/cli/owned_skills_lock_test.go`, as `_techspec.md` → Testing
   Approach describes. Every new top-level test calls `t.Parallel()` first.
5. MUST apply the declared breaks of `_techspec.md` → Declared breaks:
   - `TestRecordInstalledOwnedSkillsRewritesOnlyTheHashOfAnInstalledOwnedSkill`
     and `TestSkillsInstallRewritesTheOwnedSkillLockHash` keep their names,
     and their fixtures write `sourceType` where they wrote `ref`;
   - the other tests in `internal/baseline/owned_skills_lock_test.go` read
     the result through `Names()`;
   - no test is renamed or removed.
6. MUST update `docs/user-guide/commands/baseline.md` and
   `docs/user-guide/commands/skills.md` (`_techspec.md` → API Contract 5):
   - replace the sentence that says every other field is kept (or preserved)
     with one that says the install writes the installed folder's
     `computedHash`, removes the entry's `ref`, and keeps every other field;
   - name `skills.lockUnpinned` and the `- lock ref removed <name>` line in
     the baseline guide;
   - name `Removed skills-lock.json ref: <name>` in the skills guide.
7. MUST re-record the Behavior Surface Record with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`
   as the Task's last step, after the last edit of
   `internal/cli/baseline_update.go` and of both guides, and never hand-edit
   it.
8. MUST NOT touch an external lock entry, `source`, `sourceType` or
   `skillPath` of an owned entry, `internal/cli/doctor.go`, any other guide,
   any skill, the root `skills-lock.json`, the Makefile, CI workflows,
   `go.mod` or `go.sum`.

## Subtasks

- [ ] Characterize today's writer with the unpin table.
- [ ] Remove the `ref` member and return the two lists.
- [ ] Report them in `skills install` and `baseline update`, and adapt the stubs and fixtures.
- [ ] Update the baseline and skills command guides, then re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] Each installed owned entry loses its `ref` and gets the installed hash; every other byte, the external entries and the entry of an owned skill not installed stay.
- [ ] `roundfix skills install` prints Surface Transcript 1, and `baseline update` reports `- lock ref removed` and `skills.lockUnpinned`.
- [ ] A second run writes and reports nothing.

## Context

- interface: `internal/baseline/owned_skills_lock.go`
- interface: `internal/baseline/owned_skills_lock_test.go`
- creates: `internal/baseline/owned_skills_lock_unpin_test.go`
- interface: `internal/cli/baseline_update.go`
- interface: `internal/cli/cli.go`
- interface: `internal/cli/owned_skills_lock_test.go`
- interface: `internal/cli/baseline_update_test.go`
- interface: `internal/cli/doctor_trailing_skills_test.go`
- interface: `docs/user-guide/commands/baseline.md`
- interface: `docs/user-guide/commands/skills.md`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- `baseline.RecordInstalledOwnedSkills` and `baseline.OwnedSkillsLockUpdate` — consumers: `./internal/baseline`, `./internal/cli`
- The `skills install` lines and the `baseline update` field `skills.lockUnpinned` and its text lines — consumers: `./internal/cli`
- The baseline and skills command guides and the Behavior Surface Record — consumers: `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestRecordInstalledOwnedSkillsRemovesTheRefOfAnInstalledOwnedSkill|TestSkillsInstallRemovesTheRefOfOwnedSkillLockEntries|TestBaselineUpdateSkillsStageReportsTheUnpinnedLockEntries)$' ./internal/baseline ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestRecordInstalledOwnedSkillsRemovesTheRefOfAnInstalledOwnedSkill TestSkillsInstallRemovesTheRefOfOwnedSkillLockEntries TestBaselineUpdateSkillsStageReportsTheUnpinnedLockEntries; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 40m ./internal/baseline ./internal/cli` — expected: exit 0. Before this Task none of the three tests exists, so the first pass line is missing; afterwards they pass with both whole packages.
- `tr -s '[:space:]' ' ' < docs/user-guide/commands/baseline.md | grep -qF -- '- lock ref removed <name>' || { echo 'missing rule in the baseline guide' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/skills.md | grep -qF -- 'Removed skills-lock.json ref: <name>' || { echo 'missing rule in the skills guide' >&2; exit 1; }; go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task neither guide names the removed `ref`; afterwards both do and the documentation contracts, the Behavior Surface fingerprint among them, pass.

## References

- `_prd.md` → Goals; Core Features 1-2 and 4; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 5
- `_techspec.md` → Current behavior 1-3 and 5; Interfaces; Data Models; API Contracts 1-2 and 5; Surface Transcript 1; Invariants 1-4; Declared breaks; Testing Approach; Build Order 1

## Result

Skills loaded: implement-task, coding-guidelines, agentic-cli-design, golang-cli, golang-structs-interfaces, golang-error-handling, golang-testing, testing-boss, systematic-debugging, no-workarounds, tech-writer, evidence-gate.

Characterization before changing the writer:
`go test -count=1 -v ./internal/baseline -run '^TestRecordInstalledOwnedSkillsRemovesTheRefOfAnInstalledOwnedSkill$/skills_tool_format$'`
exited 1 against the new whole-file unpin expectation. The actual document
kept every `ref` (`roundfix`, external `testing-boss`, and uninstalled
`write-prd`); only the `roundfix` hash changed from zero to the installed
folder digest. The writer was unchanged for this check.

Implementation and acceptance evidence:

- Installed owned entries now receive their installed hash and lose `ref`
  through the specified member spans. `OwnedSkillsLockUpdate.Hashed` and
  `.Unpinned` are sorted, and `Names()` returns their sorted union. Whole-file
  assertions cover middle, first, last/current-hash and sole/null `ref`
  members, duplicate input names, unchanged external and uninstalled-owned
  entries, and preserved mode. Existing hash-only tests retain their names
  and use `sourceType` fixtures; other writer tests use `Names()`.
- The install regression asserts Surface Transcript 1 in full (stdout,
  empty stderr, exit 0) and the complete resulting lock. The baseline stage
  regression asserts both result lists, complete lock bytes with the external
  `ref` kept, the contiguous text report, and adjacent JSON `lockUpdated` /
  `lockUnpinned` fields. All three dependency stubs use the new return type.
- Every new regression checks a second invocation: no lock-update report,
  unchanged file identity, modification time and bytes. The baseline JSON
  omits both empty lock fields. All new top-level tests begin with
  `t.Parallel()`.
- Both command guides describe the installed `computedHash`, removed `ref`
  and preserved other fields. The baseline guide names
  `skills.lockUnpinned` and `- lock ref removed <name>`; the skills guide
  names `Removed skills-lock.json ref: <name>`. A focused Python inspection
  of whitespace-normalized guide text confirmed those phrases.

Focused checks (each exited 0 after implementation):

- `go test -count=1 -timeout 3m ./internal/baseline -run '^TestRecordInstalledOwnedSkills'`
- `go test -count=1 -timeout 5m ./internal/cli -run '^Test(SkillsInstall(RewritesTheOwnedSkillLockHash|RemovesTheRefOfOwnedSkillLockEntries)|BaselineUpdateSkillsStage(ReportsThe(Rewritten|Unpinned)LockEntries|FailsOnAnUnsafeLock))$'`
- `go test -count=1 -timeout 3m ./internal/cli -run '^TestBaselineUpdateSkillsStage(UsesPreviewThenConfirmation|DegradesUnreachableSourcePerSkill)$'`
- `go test -count=1 -timeout 3m ./internal/cli -run '^TestBaselineUpdateKeepsASetupSnapshotDifference$'`
- `git diff --check`

Final implementation step: re-recorded the Behavior Surface Record with
`go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`
(exit 0), after the last edit of `internal/cli/baseline_update.go` and both
command guides. No generated record was edited by hand; only this Result
was appended afterward.

The initial changed path was this Task file's Daemon-owned `in_progress`
status. That status and the graph are unchanged by the Agent. No declared
Verification command, repository Verification, commit, push or PR operation
was run. Doctor, skill guidance and glossary changes remain with the later
Tasks; no follow-up outside the authored graph was discovered.
