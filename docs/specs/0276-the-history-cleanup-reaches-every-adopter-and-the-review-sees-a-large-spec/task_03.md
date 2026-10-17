---
task: task_03
spec: 0276-the-history-cleanup-reaches-every-adopter-and-the-review-sees-a-large-spec
status: completed
type: backend
complexity: medium
---

# Task 03: Spec check reports a Task that edits an owned skill without its version record

## Overview

An owned skill changes only together with its version, yet a plan whose Task
edits `.agents/skills/<name>/` without the mirror, the version fields, the
version record and the record step passes `roundfix spec check --strict` and
fails only at settlement (`_techspec.md` → Current behavior). This Task adds
the gap `SC-SKILL-VERSION-UNTASKED`, pins it in the corpus, and documents it in
the `spec` guide, the Roundfix `spec` reference and the `write-tasks` rule. It
is verifiable on its own through detector tests on temporary repositories.

## Requirements

1. MUST answer the adopted Backlog Entry of 2026-10-16,
   [spec check does not flag an owned skill edit without its version record](references/2026-10-16-spec-check-does-not-flag-an-owned-skill-edit-without-its-version-record.md),
   by implementing API Contract 7 of `_techspec.md` in a new
   `internal/speccheck/skill_version.go`, called from `Check` in
   `internal/speccheck/constraints.go` next to `detectSurfaceRecord` with the
   loaded graph, and registered at `StageTasks` in `stagedDetectors` in
   `internal/speccheck/coherence.go` (ADR-0282, ADR-0189).
2. MUST read the owned skills from the `skills` keys of
   `skills/testdata/owned-skill-versions.json`; MUST NOT read a completed Task,
   the manifest's `qa:` Task or a Task of type `qa`; MUST NOT let another Task's
   declarations cover the Task that edits the skill; MUST skip the code, not
   fail, in a repository without the record and in a Spec without a Task
   Graph; and MUST treat an `instruction:` path as read-only. Each finding is a
   gap with no horizon.
3. MUST add the four tests named in Verification to the new
   `internal/speccheck/skill_version_test.go`, each calling `t.Parallel()`
   first: `TestSkillVersionGapNamesATaskWithoutItsRecord` (a Task declaring an
   owned skill reference without its mirror, without the `SKILL.md` pair,
   without the record and without the step, each with the summary, location
   and fix of API Contract 7, and one finding per skill),
   `TestSkillVersionGapAcceptsATaskThatRecordsItsVersion` (the complete
   declaration, a deletion with its mirror deletion, an upstream skill, a
   completed Task and a QA Task yield none),
   `TestSkillVersionGapSkipsWithoutTheRecord` (a repository without the record
   and a Spec without a graph list the skip), and
   `TestSkillVersionGapRunsAtTheTasksStageAndFailsStrict` (the PRD and
   TechSpec stages list the code as skipped; the full sweep reports it and
   `--strict` promotes it).
4. MUST add the code to `corpusFindingCodes` in
   `internal/docscontract/corpus_test.go`, to
   `internal/docscontract/testdata/corpus-golden.json` and to the
   characterization in `internal/spec/archive_layout_characterization_test.go`,
   at the active count the focused corpus test reports on the tree this Task
   starts from (expected 0, since this Spec's own pending task_04 declares its
   skill files and record), and append to the update text of both the
   identical sentence `SC-SKILL-VERSION-UNTASKED joins the characterized code
   set, with every earlier count unchanged.` No other count may move.
5. MUST add a row for the code to the code table of
   `docs/user-guide/commands/spec.md`, a bullet to
   `.agents/skills/roundfix/references/spec.md`, and to
   `.agents/skills/write-tasks/SKILL.md` a section "Owned skill versions"
   stating that a non-QA Task declaring a file of an owned skill also declares
   its mirror, both `SKILL.md` files and
   `skills/testdata/owned-skill-versions.json` and names the
   `-record-skill-versions` step in a requirement, that another Task cannot
   cover it, naming `SC-SKILL-VERSION-UNTASKED`. No text inside any
   `### QA settlement` section may change.
6. MUST run `make skills-sync`, raise and record the Roundfix and write-tasks
   skill versions with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
   and re-record the Behavior Surface Record with
   `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`,
   because this Task edits a user guide page. No record or mirror is
   hand-edited.

## Subtasks

- [ ] Write the detector and wire it into the full sweep and the stage registry.
- [ ] Add the four tests.
- [ ] Pin the code in the corpus golden and its characterization.
- [ ] Document the code in the guide, the Roundfix reference and write-tasks.
- [ ] Sync the mirrors, record the versions and the Behavior Surface Record.

## Acceptance Criteria

- [ ] A pending Task that edits an owned skill without its record yields one gap per skill at its Context line; `--strict` fails.
- [ ] The complete declaration yields none, and so do completed, QA and upstream-skill Tasks.
- [ ] A repository without the record and the PRD and TechSpec stages list the code as skipped.
- [ ] Every earlier corpus count is unchanged.

## Consumers

- The Spec Consistency Check finding set — consumers: `./internal/speccheck`, `./internal/cli`, `./internal/daemon`, `./internal/docscontract`

## Context

- instruction: `docs/adr/0282-the-history-cleanup-names-a-superseded-adr-in-its-successor-and-a-large-candidate-is-reviewed-in-segments.md`
- instruction: `docs/adr/0189-an-owned-skills-version-names-its-content-and-its-minimum-is-the-bundle.md`
- instruction: `internal/speccheck/surface_record.go`
- instruction: `internal/speccheck/surface_record_test.go`
- creates: `internal/speccheck/skill_version.go`
- creates: `internal/speccheck/skill_version_test.go`
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

- `out="$(go test -count=1 -v -run '^(TestSkillVersionGapNamesATaskWithoutItsRecord|TestSkillVersionGapAcceptsATaskThatRecordsItsVersion|TestSkillVersionGapSkipsWithoutTheRecord|TestSkillVersionGapRunsAtTheTasksStageAndFailsStrict)$' ./internal/speccheck 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestSkillVersionGapNamesATaskWithoutItsRecord TestSkillVersionGapAcceptsATaskThatRecordsItsVersion TestSkillVersionGapSkipsWithoutTheRecord TestSkillVersionGapRunsAtTheTasksStageAndFailsStrict; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; grep -qF -- '"SC-SKILL-VERSION-UNTASKED":' internal/docscontract/testdata/corpus-golden.json || { printf 'the corpus golden lacks SC-SKILL-VERSION-UNTASKED\n' >&2; exit 1; }; for f in docs/user-guide/commands/spec.md .agents/skills/roundfix/references/spec.md .agents/skills/write-tasks/SKILL.md; do grep -qF -- 'SC-SKILL-VERSION-UNTASKED' "$f" || { printf 'missing SC-SKILL-VERSION-UNTASKED in %s\n' "$f" >&2; exit 1; }; done; make skills-sync-check || exit 1; go test -count=1 ./skills || exit 1; go test -count=1 ./internal/speccheck || exit 1; go test -count=1 -run '^TestArchiveLayoutCharacterizationPinsCorpusGoldenAfterSpec0095$' ./internal/spec || exit 1; go test -count=1 -run '^(TestRunSpecCheckCleanText|TestSpecCheckRunVerification|TestRunSpecCheckGapStrictPromotion)$' ./internal/cli || exit 1; go test -count=1 -run 'SpecConsistency' ./internal/daemon || exit 1; go test -count=1 -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the four detector tests do not exist, so the command fails at its first pass check; afterwards the detector, corpus, spec check, settlement, skill version and Behavior Surface Record tests pass.

## References

- `_prd.md` → Goals; User Story 5; Core Feature 6; Core Feature 9; Success Metric 4
- `_techspec.md` → API Contract 7; API Contract 10; Invariant 9; Current behavior; Testing Approach; Build Order 3
- ADR-0282; ADR-0189; ADR-0280


## Result

Implemented the Task 03 slice for Daemon Verification. Status remains
Daemon-owned; the declared Verification and repository gates were not run.

### Implementation and acceptance evidence

- Missing owned-skill work: `detectSkillVersion` reads ownership from the
  version record's `skills` keys and evaluates each eligible Task separately.
  It emits one gap per skill at its first mutating Context declaration, with
  API Contract 7's summary and ordered missing-declaration fix. Mirrors must
  match deletions; other mirror, SKILL.md and record declarations must be
  writable. A requirement must name `-record-skill-versions`; Verification
  text and another Task's declarations cannot supply the step.
  `TestSkillVersionGapNamesATaskWithoutItsRecord` exercises each missing
  component, all missing components, another Task, instruction-only and
  deleted records, wrong mirror kinds, and multiple paths and skills. An
  instruction citing the same path first does not move the finding's line.
- Accepted declarations and exclusions:
  `TestSkillVersionGapAcceptsATaskThatRecordsItsVersion` exercises complete
  interface and creates declarations, a reference deletion with its mirror
  deletion, multiple skills, an upstream skill, a completed Task, a QA Task,
  and an instruction-only skill reference. They yield no skill-version gap.
- Missing inputs and stage behavior:
  `TestSkillVersionGapSkipsWithoutTheRecord` exercises an absent record, a
  non-file record and an absent graph, asserting their skip paths; malformed
  and unreadable records return errors. The detector is next to
  `detectSurfaceRecord` in `Check` and registered at `StageTasks`.
  `TestSkillVersionGapRunsAtTheTasksStageAndFailsStrict` asserts PRD and
  TechSpec stage skips, findings in both full and Tasks sweeps, and strict
  promotion to error. The detector has no contract horizon.
- Corpus counts: the focused corpus test passed before the detector change.
  After registering the code, its diagnostic measured
  `SC-SKILL-VERSION-UNTASKED: 0` and every earlier count unchanged. The golden
  and archive-layout characterization now pin that count and append the
  required identical update sentence. The final focused corpus and
  characterization checks passed.
- Documentation: added the guide table row, Roundfix reference bullet and
  write-tasks "Owned skill versions" section. Generated mirrors and version
  records carry Roundfix 0.1.88 and write-tasks 0.0.17. Re-recorded the
  Behavior Surface Record for the changed guide. Baseline regeneration
  reported no derived changes.

### Focused checks and generation

- Before implementation,
  `go test ./internal/speccheck -run '^TestSkillVersionGapNamesATaskWithoutItsRecord$/interface/record$' -count=1`
  exited 1: the expected finding was absent.
- `go test ./internal/speccheck -run 'TestSkillVersionGap' -count=1` exited 0
  with all four new detector tests and their subtests passing.
- `go test -tags docscontract ./internal/docscontract -run '^(TestCheckCorpusGolden|TestStageScopeDefaultSweepIsUnchanged)$' -count=1`
  exited 0.
- `go test ./internal/spec -run 'ArchiveLayoutCharacterization' -count=1`
  exited 0.
- `make skills-sync` exited 0.
- `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
  exited 0 and generated both skills' version changes and record entries.
- `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`
  exited 0 and regenerated the Behavior Surface Record.
- `make baseline-digests` exited 0 with `changed: false`.
- `git diff --check` exited 0. A focused scope inspection confirmed all
  17 changed paths are declared by this Task, its authored fields are
  preserved, both skill mirrors match, and QA settlement text is unchanged.

The Daemon owns the declared Verification, settlement and any Task commit.
