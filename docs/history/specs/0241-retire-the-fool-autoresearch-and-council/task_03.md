---
task: task_03
spec: 0241-retire-the-fool-autoresearch-and-council
status: completed
type: backend
complexity: medium
---

# Task 03: Baseline update lists the retired copies an adopter still holds

## Overview

Adopters keep their installed `council` and `the-fool` copies after the next
update, because Roundfix never deletes an installed skill tree. This Task
makes `roundfix baseline update` list each Retired Skill the repository still
holds under `Skills retired` (`skills.retired` in JSON) with the paths to
delete, read-only and without changing the update's state or exit code. It
gives the removal commands in the baseline user guide, describes the report in
the Roundfix Skill's baseline reference, and adds **Retired Skill** to
`CONTEXT.md`.

This is an authorized tooling Task. Its Governed Paths are the Roundfix
Skill's `SKILL.md` (both copies) and its canonical baseline reference, bounded
in `_authorization.md`.

## Requirements

1. MUST create `internal/baseline/retired_skills_installed.go` with
   `InstalledRetiredSkill` and `InstalledRetiredSkills` per the TechSpec's
   Interfaces and Invariants 2 and 3, reading only.
2. MUST add `Retired` to `baselineUpdateSkillsResult` in
   `internal/cli/baseline_update.go` as the Data Models section states, and
   fill it per Invariants 4 and 5 in the preview and after an applied plan's
   skills stage; `--no-skills` skips it. MUST print the text block of
   Invariant 6 only when the list is non-empty, so every existing transcript
   and golden without retired copies stays byte-identical.
3. MUST create `internal/baseline/retired_skills_installed_test.go` with
   `TestInstalledRetiredSkillsListsTreesLinksAndLockEntries`,
   `TestInstalledRetiredSkillsIsEmptyWithoutRetiredCopies` and
   `TestInstalledRetiredSkillsRefusesAMalformedLock` (Testing Approach 3).
4. MUST create `internal/cli/baseline_update_retired_skills_test.go` with
   `TestBaselineUpdateListsRetiredSkillsWithoutChangingItsState`,
   `TestBaselineUpdateWithoutSkillsOmitsRetiredSkills` and
   `TestBaselineUpdateAppliesAndListsRetiredSkills` (Testing Approach 4),
   on the adoption fixture the existing update tests use; the first asserts
   Surface Transcript 1's `Skills retired` lines verbatim, its state and its
   exit. No test reaches the network.
5. MUST add the task_03 Fixed text `### Retired skills` to
   `docs/user-guide/commands/baseline.md`, and the task_03 sentence to
   `.agents/skills/roundfix/references/baseline.md`.
6. MUST add the **Retired Skill** entry, as the task_03 Fixed text gives it,
   to `CONTEXT.md` after **Composed Setup Snapshot**, following the
   `domain-modeling` skill, and change no other glossary entry.
7. MUST run `make skills-sync` and re-record versions with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
   after the last skill edit, so the Roundfix Skill's version rises in both
   fields. MUST NOT edit a `### QA settlement` section or a version or digest
   by hand.
8. MUST NOT change Doctor, a Baseline asset, a lock or any other skill.

## Subtasks

- [ ] Add the read-only retired copy inspection and its tests.
- [ ] Report it from `baseline update` in text and JSON, with tests.
- [ ] Document the removal commands, the skill reference and the glossary term.
- [ ] Sync the Roundfix Skill and record its version.

## Acceptance Criteria

- [ ] A repository holding retired copies sees each one with its paths, and its update state and exit code are unchanged.
- [ ] `--no-skills`, or a repository without retired copies, prints no `Skills retired` and no `skills.retired`.
- [ ] A malformed lock is reported, not ignored.
- [ ] The guide gives the removal commands, `CONTEXT.md` defines **Retired Skill**, and the Roundfix Skill's mirrors equal their canonical files with the raised version recorded.

## Context

- instruction: `docs/adr/0246-a-retired-skill-leaves-the-baseline-and-the-adopter-removes-its-copy.md`
- instruction: `internal/cli/baseline_update_test.go`
- instruction: `internal/cli/baseline_update_repository_profile_test.go`
- creates: `internal/baseline/retired_skills_installed.go`
- creates: `internal/baseline/retired_skills_installed_test.go`
- creates: `internal/cli/baseline_update_retired_skills_test.go`
- interface: `internal/cli/baseline_update.go`
- interface: `docs/user-guide/commands/baseline.md`
- interface: `CONTEXT.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `.agents/skills/roundfix/references/baseline.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/roundfix/references/baseline.md`
- interface: `skills/testdata/owned-skill-versions.json`

## Verification

- `out="$(go test -count=1 -v -run '^(TestInstalledRetiredSkillsListsTreesLinksAndLockEntries|TestInstalledRetiredSkillsIsEmptyWithoutRetiredCopies|TestInstalledRetiredSkillsRefusesAMalformedLock)$' ./internal/baseline 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestInstalledRetiredSkillsListsTreesLinksAndLockEntries TestInstalledRetiredSkillsIsEmptyWithoutRetiredCopies TestInstalledRetiredSkillsRefusesAMalformedLock; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done && out="$(go test -count=1 -v -run '^(TestBaselineUpdateListsRetiredSkillsWithoutChangingItsState|TestBaselineUpdateWithoutSkillsOmitsRetiredSkills|TestBaselineUpdateAppliesAndListsRetiredSkills)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestBaselineUpdateListsRetiredSkillsWithoutChangingItsState TestBaselineUpdateWithoutSkillsOmitsRetiredSkills TestBaselineUpdateAppliesAndListsRetiredSkills; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the six tests do not exist, so their pass lines are missing and the command fails.
- `for phrase in '### Retired skills' 'Skills retired' 'git rm -r -q --ignore-unmatch .agents/skills/council .agents/skills/the-fool'; do tr -s '[:space:]' ' ' < docs/user-guide/commands/baseline.md | grep -qF -- "$phrase" || { printf 'missing phrase in baseline guide: %s\n' "$phrase" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Retired Skill**: A skill the Baseline no longer requires' || { printf 'missing glossary term: Retired Skill\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/baseline.md | grep -qF -- 'Skills retired' || { printf 'missing phrase in skill reference\n' >&2; exit 1; }; cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md && cmp .agents/skills/roundfix/references/baseline.md skills/roundfix/references/baseline.md && out="$(go test -count=1 -v -run '^(TestEveryOwnedSkillVersionIsRecorded)$' ./skills 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- "--- PASS: TestEveryOwnedSkillVersionIsRecorded" || { printf 'missing pass: TestEveryOwnedSkillVersionIsRecorded\n' >&2; exit 1; }` — expected: exit 0; before this Task neither the guide, the glossary nor the skill reference has the new text, so the command fails; after it the mirrors equal their canonical files and the raised version is recorded.

## References

- `_prd.md` → Goal 4; User Story 2; Core Feature 5; Success Metrics 5-6; Release note; Declared breaks
- `_techspec.md` → Interfaces; Data Models; Fixed texts (task_03); API Contract 2; Surface Transcript 1; Testing Approach 3-4; Vocabulary Contract; Build Order 3
- ADR-0246, ADR-0191, ADR-0189, ADR-0233

## Result

Implemented the task_03 slice for Daemon Verification; status and the authored
Verification commands remain Daemon-owned. The starting worktree had only
this Task file's Daemon-written `in_progress` status change.

- Retired copies: `InstalledRetiredSkills` reads the existing validated lock
  and uses `Lstat` for both skill locations. It combines trees, dangling
  links and lock-only entries in lexical skill order, with stable relative
  paths, and excludes `autoresearch`. Preview and applied update output add
  the optional report without changing the successful update state, message,
  category, next action or exit code. The CLI tests assert the Surface
  Transcript 1 lines verbatim, structured entries and unchanged retired
  trees/lock bytes. Additional checks cover lock-only deletion text and all
  three deletion targets.
- Omission: tests cover current repositories without retired copies and
  `--no-skills` in both formats, in preview and with `--yes`. They assert no
  retired text block or JSON field and preserve repository bytes.
- Malformed lock: inspection refuses malformed JSON and an invalid skills
  object, naming `skills-lock.json`. The CLI reports an execution failure
  with next action `repair skills-lock.json and rerun roundfix baseline update`.
  Additional tests cover a non-regular lock and a skill-path `Lstat` error.
- Documentation and versions: added the TechSpec's fixed guide section,
  Roundfix baseline-reference sentence, and only the **Retired Skill** entry
  after **Composed Setup Snapshot** in `CONTEXT.md`. `make skills-sync` copied
  the reference. The record generator raised both Roundfix version fields
  from 0.1.44 to 0.1.45 in both copies and appended the generated digest.
  Python comparisons confirmed both mirrors equal their canonical files,
  the guide matches the fixed text and no other glossary entry changed.

Focused-check evidence:

- Before implementation, `GOCACHE=/private/tmp/roundfix-task03-cache go test
  ./internal/baseline -run InstalledRetiredSkills -count=1` exited 1 because
  `InstalledRetiredSkill` and `InstalledRetiredSkills` did not exist.
- After the final code/test edits, `GOCACHE=/private/tmp/roundfix-task03-cache
  go test ./internal/baseline ./internal/cli -run
  'InstalledRetiredSkills|BaselineUpdate.*RetiredSkill' -count=1` exited 0
  (baseline 1.141s; CLI 6.930s).
- `make skills-sync` exited 0. The required `go test ./skills -run
  '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`, with the
  task-scoped cache, first hit the protected canonical file's sandbox write
  restriction, then exited 0 with authorized escalation (0.164s). No version
  or digest was edited by hand.
- `GOCACHE=/private/tmp/roundfix-task03-cache make baseline-digests` exited 0
  with `changed:false`; it changed no Baseline asset or derived pin.
- `GOCACHE=/private/tmp/roundfix-task03-cache go test ./internal/cli -run
  BaselineUpdate -count=1` exited 0 (33.032s) on a stable worktree. An earlier
  overlapping run passed its assertions but exited 1 because suiteguard
  observed the concurrently authorized skill-version generator writes;
  stopping those mutations resolved the guard failure without changing tests.
- `git -c core.fsmonitor=false diff --check` exited 0. Changed-file inspection
  shows only this Task's declared source, tests, docs, Roundfix Skill copies,
  generated owned-version record and this Result. No Doctor, lock, other skill,
  other Task or Task Graph change was made.

The Task's declared Verification commands and the repository-wide gate were
not run in this child turn. No commit, push or Pull Request was made. No
follow-up work was identified.

## Carry-forward provenance

- Source Run: `run_20261006T214048Z_131051a4f3c19cbe`
- Source commit: `f807e63f7f26f15cc312e5d66f74a92bd5895035`
