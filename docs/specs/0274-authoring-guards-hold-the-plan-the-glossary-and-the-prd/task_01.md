---
task: task_01
spec: 0274-authoring-guards-hold-the-plan-the-glossary-and-the-prd
status: pending
type: backend
complexity: medium
---

# Task 01: Spec check reports a guide Task that does not record its surface fingerprint

## Overview

A Task that edits a user guide page changes its fingerprint in the Behavior
Surface Record, and the record's check runs at that Task's own settlement.
Today `roundfix spec check --strict` is clean on a plan whose guide Task cannot
re-record it (`_techspec.md` → Current behavior). This Task adds the gap
`SC-SURFACE-RECORD-UNTASKED`, pins it in the corpus, and documents it in the
`spec` guide, the Roundfix `spec` reference and the `write-tasks` rule. It is
verifiable on its own through the detector tests on temporary repositories.

## Requirements

1. MUST answer the adopted Backlog Entry of 2026-10-14,
   [spec check does not flag a guide Task without the coverage record](references/2026-10-14-spec-check-does-not-flag-a-guide-task-without-the-coverage-record.md),
   by implementing API Contract 1 of `_techspec.md` in a new
   `internal/speccheck/surface_record.go`, called from `Check` in
   `internal/speccheck/constraints.go` next to `detectSkillsDeclaration` with
   the loaded graph, and registered at `StageTasks` in `stagedDetectors` in
   `internal/speccheck/coherence.go` (ADR-0280).
2. MUST report the guide Task even when another Task of the Spec declares the
   record and names the step, because the guide Task settles first; MUST NOT
   read a completed Task, the manifest's `qa:` Task or a Task of type `qa`;
   MUST NOT infer a command-help change from a declared source file; MUST
   skip the code, not fail, in a repository without
   `docs/references/behavior-surfaces.json` and in a Spec without a Task Graph.
   Each finding is a gap with no horizon (ADR-0280).
3. MUST add the four tests named in Verification to the new
   `internal/speccheck/surface_record_test.go`, each calling `t.Parallel()`
   first: `TestSurfaceRecordGapNamesAGuideTaskWithoutItsRecord` (the Spec 0268
   shape of `_techspec.md` → Current behavior, a Task declaring the record but
   not the step, and the record and step only in another Task, each with the
   summary and Context line of API Contract 1),
   `TestSurfaceRecordGapNamesLimitsThatForbidTheRecord` (a Limits line naming
   the record yields a gap at that line),
   `TestSurfaceRecordAcceptsAGuideTaskThatRecordsItself` (the fixed shape, a completed Task, a QA Task, a Task without a guide page and a
   repository without the record yield none, the last with its skip), and
   `TestSurfaceRecordRunsAtTheTasksStageAndFailsStrict` (the PRD and TechSpec
   stages list the code as skipped; the full sweep reports it and `--strict`
   promotes it).
4. MUST add the code to `corpusFindingCodes` in
   `internal/docscontract/corpus_test.go`, to
   `internal/docscontract/testdata/corpus-golden.json` and to the
   characterization in `internal/spec/archive_layout_characterization_test.go`,
   at the active count the focused corpus test reports on the tree this Task
   starts from (expected 0), and append to the update text of both the
   identical sentence `SC-SURFACE-RECORD-UNTASKED joins the characterized code
   set, with every earlier count unchanged.` No other count may move.
5. MUST add a row for the code to the code table of
   `docs/user-guide/commands/spec.md`, a bullet to
   `.agents/skills/roundfix/references/spec.md`, and to the "Skills Tasks for
   changed Behavior Surfaces" section of `.agents/skills/write-tasks/SKILL.md`
   the rule that a non-QA Task declaring a user guide page declares
   `docs/references/behavior-surfaces.json` and names the
   `-record-skill-coverage` step in a requirement, that another Task cannot
   cover it, and that the Spec's `_authorization.md` Limits never forbid that
   file, naming `SC-SURFACE-RECORD-UNTASKED`. No text inside any
   `### QA settlement` section may change.
6. MUST run `make skills-sync`, raise and record the Roundfix and write-tasks
   skill versions with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
   and re-record the Behavior Surface Record with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`,
   because this Task edits a user guide page. No record or mirror is
   hand-edited.

## Subtasks

- [ ] Write the detector and wire it into the full sweep and the stage registry.
- [ ] Add the four tests.
- [ ] Pin the code in the corpus golden and its characterization.
- [ ] Document the code in the guide, the Roundfix reference and write-tasks.
- [ ] Sync the mirrors, record the versions and the Behavior Surface Record.

## Acceptance Criteria

- [ ] A pending guide Task without the record or the step yields one gap at its Context line; `--strict` fails.
- [ ] A Limits line naming the record yields a gap at that line.
- [ ] A guide Task that declares the record and names the step yields none, and so do completed and QA Tasks.
- [ ] A repository without the record and the PRD and TechSpec stages list the code as skipped.
- [ ] Every earlier corpus count is unchanged.

## Consumers

- The Spec Consistency Check finding set — consumers: `./internal/speccheck`, `./internal/cli`, `./internal/daemon`, `./internal/docscontract`

## Context

- instruction: `docs/adr/0280-a-guide-task-records-its-own-surface-fingerprint-a-prd-records-who-decided-and-the-glossary-is-held-to-its-own-rule.md`
- instruction: `docs/adr/0256-a-release-waits-for-the-skills-that-describe-a-changed-surface.md`
- creates: `internal/speccheck/surface_record.go`
- creates: `internal/speccheck/surface_record_test.go`
- interface: `internal/speccheck/constraints.go`
- interface: `internal/speccheck/coherence.go`
- interface: `internal/docscontract/corpus_test.go`
- interface: `internal/docscontract/testdata/corpus-golden.json`
- interface: `internal/spec/archive_layout_characterization_test.go`
- interface: `docs/user-guide/commands/spec.md`
- interface: `.agents/skills/roundfix/references/spec.md`
- interface: `skills/roundfix/references/spec.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `.agents/skills/write-tasks/SKILL.md`
- interface: `skills/write-tasks/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`

## Verification

- `out="$(go test -count=1 -v -run '^(TestSurfaceRecordGapNamesAGuideTaskWithoutItsRecord|TestSurfaceRecordGapNamesLimitsThatForbidTheRecord|TestSurfaceRecordAcceptsAGuideTaskThatRecordsItself|TestSurfaceRecordRunsAtTheTasksStageAndFailsStrict)$' ./internal/speccheck 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestSurfaceRecordGapNamesAGuideTaskWithoutItsRecord TestSurfaceRecordGapNamesLimitsThatForbidTheRecord TestSurfaceRecordAcceptsAGuideTaskThatRecordsItself TestSurfaceRecordRunsAtTheTasksStageAndFailsStrict; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; grep -qF -- '"SC-SURFACE-RECORD-UNTASKED":' internal/docscontract/testdata/corpus-golden.json || { printf 'the corpus golden lacks SC-SURFACE-RECORD-UNTASKED\n' >&2; exit 1; }; for f in docs/user-guide/commands/spec.md .agents/skills/roundfix/references/spec.md .agents/skills/write-tasks/SKILL.md; do grep -qF -- 'SC-SURFACE-RECORD-UNTASKED' "$f" || { printf 'missing SC-SURFACE-RECORD-UNTASKED in %s\n' "$f" >&2; exit 1; }; done; make skills-sync-check || exit 1; go test -count=1 -run '^TestEveryOwnedSkillVersionIsRecorded$' ./skills || exit 1; go test -count=1 ./internal/speccheck || exit 1; go test -count=1 -run '^TestArchiveLayoutCharacterizationPinsCorpusGoldenAfterSpec0095$' ./internal/spec || exit 1; go test -count=1 -run '^(TestRunSpecCheckCleanText|TestSpecCheckRunVerification|TestRunSpecCheckGapStrictPromotion)$' ./internal/cli || exit 1; go test -count=1 -run 'SpecConsistency' ./internal/daemon || exit 1; go test -count=1 -tags docscontract -run '^(TestCheckCorpusGolden|TestCheckActiveCorpusHasNoErrors|TestTheSkillCoverageMapIsCurrent)$' ./internal/docscontract` — expected: exit 0. Before this Task the four detector tests do not exist, so the command fails at its first pass check; afterwards the detector, corpus, spec check, settlement, skill version and Behavior Surface Record tests pass.

## References

- `_prd.md` → Goals; User Story 1; Core Feature 1; Core Feature 5; Success Metric 1; Success Metric 5
- `_techspec.md` → API Contract 1; API Contract 5; Invariants 1, 2 and 6; Current behavior; Testing Approach; Build Order 1
- ADR-0280; ADR-0256; ADR-0182; ADR-0117
