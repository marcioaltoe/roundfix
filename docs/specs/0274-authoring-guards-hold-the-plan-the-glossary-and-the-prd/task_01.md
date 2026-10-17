---
task: task_01
spec: 0274-authoring-guards-hold-the-plan-the-glossary-and-the-prd
status: completed
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


## Result

Implemented the Task 01 slice for Daemon Verification. Status and the authored
Verification remain Daemon-owned; no commit, push or Pull Request was made.
The starting worktree contained only the Daemon's `pending` to `in_progress`
change in this Task file, which was preserved.

### Implementation

- Added `detectSurfaceRecord` and `CodeSurfaceRecordUntasked`, called beside
  `detectSkillsDeclaration` in the full sweep and registered at `StageTasks`.
  Each eligible guide Task must itself declare the record under `interface:`
  or `creates:` and name `-record-skill-coverage` in a requirement. Findings
  are gaps with no history horizon, located at the guide's Context line.
- Completed Tasks, the manifest QA Task, Tasks of type `qa`, and Tasks without
  a changed Markdown guide do not enter the guide check. Declared CLI source
  files do not imply help changes. Missing/non-regular records and absent
  Task Graphs produce explicit skips. Authorization Limits lines naming the
  record produce separate gaps at those lines. Both fixes name the
  write-tasks rule "Skills Tasks for changed Behavior Surfaces".
- Added the four required parallel tests using the existing disposable Spec
  fixtures. They cover the three-guide Spec 0268 shape, record without step,
  step without record, another Task owning both, invalid record declaration
  kinds, a step only in Verification, fixed declarations, completed and QA
  Tasks, absent/non-regular records, and stage/strict behavior.
- Added the code to the corpus registry, golden and archive-layout
  characterization at the measured active count of zero. Appended the exact
  required update sentence to both pins; all 43 earlier counts stay zero.
- Documented the code in the spec guide and Roundfix reference, and the
  per-guide-Task record requirement and authorization Limits rule in
  write-tasks. The skill generators raised Roundfix from `0.1.82` to
  `0.1.83` and write-tasks from `0.0.15` to `0.0.16`, recorded both versions,
  synchronized the mirrors, and re-recorded the spec guide's fingerprint.
  No mirror, version record or Behavior Surface Record was hand-edited.

### Focused checks and acceptance evidence

1. **Pending guide Task without its record or step yields one Context gap;
   strict promotion fails:**
   `rtk proxy go test ./internal/speccheck -run '^TestSurfaceRecord' -count=1 -v`
   exited 0, with all four required top-level tests passing. The gap test
   asserts the exact API Contract summary, one location at the guide Context
   line, gap severity and the fix's rule name for each declaration kind.
   The stage/strict test asserts promotion to error, the mechanism used by
   `--strict`. Another Task's record and step do not satisfy the guide Task.
2. **Limits naming the record yields a gap at that line:**
   `TestSurfaceRecordGapNamesLimitsThatForbidTheRecord` passes in that focused
   run, asserting line 9 and proving mentions before and after Limits do not
   create findings. The Spec 0268 case asserts two gaps total: the guide Task
   and its authorization Limits. Before wiring the detector,
   `rtk proxy go test ./internal/speccheck -run '^TestSurfaceRecordGapNamesLimitsThatForbidTheRecord$' -count=1`
   exited 1 because it found zero gaps instead of one, reproducing the defect.
3. **Fixed, completed and QA Tasks yield none:**
   `TestSurfaceRecordAcceptsAGuideTaskThatRecordsItself` passes in the focused
   run for both allowed record declaration kinds, multiple guides, completed
   and QA Tasks, instruction-only/non-Markdown references and a declared CLI
   source without a guide. Initial QA fixtures failed graph validation because
   manifest `qa:` and type `qa` must agree; the fixture was corrected to that
   existing contract, without changing production graph validation.
4. **Missing record and early stages skip:**
   The acceptance test asserts the record-path skip for missing and directory
   records; the stage/strict test asserts PRD and TechSpec stage skips, full
   and Tasks-stage findings, and the `_tasks.md` skip for an absent graph.
5. **Every earlier corpus count is unchanged:**
   `rtk proxy go test -tags docscontract ./internal/docscontract -run '^TestCheckCorpusGolden$' -count=1 -v`
   first printed the measured counts, with the new code at zero and every
   existing code at zero; its mismatch was the intentionally unpinned new key.
   After updating the two pins, the same focused command exited 0. A final
   read-only comparison against `HEAD` confirmed the 43 earlier counts are
   identical and both update strings carry exactly the requested sentence.

Required source regeneration also exited 0:

- `rtk proxy make skills-sync`
- `rtk proxy go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
- `rtk proxy go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`
- `rtk proxy make baseline-digests` — repository-required regeneration after
  owned skill edits; reported `changed:false`, with derived artifacts already
  matching their sources.

Final inspection found both changed skill mirrors byte-identical to their
canonical directories, every `### QA settlement` section unchanged, and only
the 17 paths declared by this Task changed or new. `rtk proxy git diff --check`
exited 0. No other Task file or Task Graph manifest changed. The full declared
Verification and repository Verification were not run; the Daemon runs those
after this handoff. No follow-up outside Task 01's slice was identified.

## Carry-forward provenance

- Source Run: `run_20261017T091437Z_5c530d78d16d2f6a`
- Source commit: `718944ac9bf45b4d763a15d34030020754d43838`
