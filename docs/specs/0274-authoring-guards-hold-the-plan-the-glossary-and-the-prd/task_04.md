---
task: task_04
spec: 0274-authoring-guards-hold-the-plan-the-glossary-and-the-prd
status: completed
type: backend
complexity: high
---

# Task 04: A PRD carries a Clarification Record and spec check reports one without it

## Overview

PRDs of unattended cycles are written without an interview by construction,
and nothing in them says who settled their questions (`_techspec.md` →
Current behavior). This Task adds the Clarification Record to the `write-prd`
skill, its template and a new guide, to the minimal PRD of `write-techspec`,
and the gap `SC-DECISIONS-UNRECORDED` that the Spec Consistency Check reports
from the PRD stage for PRDs committed after that guide. It documents the code
and adds the glossary term. It is verifiable on its own through the detector
tests on temporary repositories.

## Requirements

1. MUST answer the maintainer's note of 2026-10-14 recorded in `_prd.md` →
   Decisions ("a skill `grilling` não é usada corretamente na criação de
   PRDs") by implementing API Contract 3 of `_techspec.md` in a new
   `internal/speccheck/clarification_record.go`, called from `Check` in
   `internal/speccheck/constraints.go` and from the authoring-stage sweep in
   `internal/speccheck/coherence.go`, and registered there at `StagePRD`
   (ADR-0280).
2. MUST read only the shape of the record and never judge whether a person
   was present (ADR-0280); MUST skip, not fail, a PRD committed before the
   oldest commit that added `ClarificationGuidePath`, a repository without that
   guide, and a missing PRD; MUST hold an uncommitted PRD.
3. MUST add the four tests named in Verification to the new
   `internal/speccheck/clarification_record_test.go`, each calling
   `t.Parallel()` first: `TestClarificationRecordGapNamesAPRDWithoutTheRecord`
   (no section, a section with neither shape, and an invalid date, each with
   the summary, location and fix of API Contract 3),
   `TestClarificationRecordAcceptsADatedAnswerOrTheAutonomousRoute` (each shape
   alone and both together yield none),
   `TestClarificationRecordSkipsAPRDBeforeItsHorizon` (a PRD committed before
   the guide's adding commit and a repository without the guide are skipped
   with the horizon text; an uncommitted PRD is held), and
   `TestClarificationRecordRunsFromThePRDStageAndFailsStrict`.
4. MUST add the code to `corpusFindingCodes`, the corpus golden and its
   characterization at the active count the focused corpus test reports
   (expected 0, since every active PRD predates the horizon), appending
   `SC-DECISIONS-UNRECORDED joins the characterized code set, with every
   earlier count unchanged.` to both update texts. No other count may move.
5. MUST create `.agents/skills/write-prd/references/clarification-record.md`
   stating the record's two line shapes, that a dated line names a question a
   person answered, that the autonomous-route line names the source the
   author decided from when no person was present, that grilling stays the
   interactive entry when a person drives the authoring, the horizon, and the
   code. MUST update `.agents/skills/write-prd/SKILL.md` (the Clarify and
   Record decisions steps and the reference list) and the `## Decisions`
   guidance of `.agents/skills/write-prd/references/prd-template.md` to write
   the record, with the phrase
   `Autonomous route: no person present; decisions taken from`, and add the
   Clarification Record to the minimal PRD contract of
   `.agents/skills/write-techspec/SKILL.md`.
6. MUST add a row for the code to `docs/user-guide/commands/spec.md` and a
   bullet to `.agents/skills/roundfix/references/spec.md`, both naming the
   horizon. No text inside any `### QA settlement` section may change.
7. MUST, through `domain-modeling`, add **Clarification Record** to
   `CONTEXT.md` near **Glossary Declaration**, linked to ADR-0280: the
   `## Decisions` section of a PRD that says who settled its questions, one
   dated line for each question a person answered, or the autonomous-route
   line naming the source of its decisions. The entry MUST pass the glossary
   contract of task_03.
8. MUST run `make skills-sync`, raise and record the Roundfix, write-prd and
   write-techspec skill versions with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
   and re-record the Behavior Surface Record with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`,
   because this Task edits a user guide page. No record or mirror is
   hand-edited.

## Subtasks

- [ ] Write the detector, its horizon and its wiring.
- [ ] Add the four tests and pin the code in the corpus.
- [ ] Write the guide, the skill and template changes, and the minimal PRD rule.
- [ ] Document the code and add the glossary entry.
- [ ] Sync the mirrors, record the versions and the Behavior Surface Record.

## Acceptance Criteria

- [ ] A PRD after the horizon without the record yields one gap at `--stage prd`; `--strict` fails.
- [ ] A dated answer or the autonomous-route line satisfies it.
- [ ] A PRD before the horizon is skipped with the horizon text.
- [ ] The skills, guide, reference and glossary describe the record, and every earlier corpus count is unchanged.

## Consumers

- The Spec Consistency Check finding set — consumers: `./internal/speccheck`, `./internal/cli`, `./internal/daemon`, `./internal/docscontract`

## Context

- instruction: `docs/adr/0280-a-guide-task-records-its-own-surface-fingerprint-a-prd-records-who-decided-and-the-glossary-is-held-to-its-own-rule.md`
- instruction: `docs/adr/0244-a-spec-declares-the-domain-terms-it-introduces-and-the-check-holds-them-to-the-glossary.md`
- instruction: `.agents/skills/domain-modeling/SKILL.md`
- instruction: `internal/speccheck/glossary.go`
- creates: `internal/speccheck/clarification_record.go`
- creates: `internal/speccheck/clarification_record_test.go`
- creates: `.agents/skills/write-prd/references/clarification-record.md`
- creates: `skills/write-prd/references/clarification-record.md`
- interface: `internal/speccheck/constraints.go`
- interface: `internal/speccheck/coherence.go`
- interface: `internal/docscontract/corpus_test.go`
- interface: `internal/docscontract/testdata/corpus-golden.json`
- interface: `internal/spec/archive_layout_characterization_test.go`
- interface: `.agents/skills/write-prd/SKILL.md`
- interface: `skills/write-prd/SKILL.md`
- interface: `.agents/skills/write-prd/references/prd-template.md`
- interface: `skills/write-prd/references/prd-template.md`
- interface: `.agents/skills/write-techspec/SKILL.md`
- interface: `skills/write-techspec/SKILL.md`
- interface: `docs/user-guide/commands/spec.md`
- interface: `.agents/skills/roundfix/references/spec.md`
- interface: `skills/roundfix/references/spec.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`
- interface: `CONTEXT.md`

## Verification

- `out="$(go test -count=1 -v -run '^(TestClarificationRecordGapNamesAPRDWithoutTheRecord|TestClarificationRecordAcceptsADatedAnswerOrTheAutonomousRoute|TestClarificationRecordSkipsAPRDBeforeItsHorizon|TestClarificationRecordRunsFromThePRDStageAndFailsStrict)$' ./internal/speccheck 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestClarificationRecordGapNamesAPRDWithoutTheRecord TestClarificationRecordAcceptsADatedAnswerOrTheAutonomousRoute TestClarificationRecordSkipsAPRDBeforeItsHorizon TestClarificationRecordRunsFromThePRDStageAndFailsStrict; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; grep -qF -- '"SC-DECISIONS-UNRECORDED":' internal/docscontract/testdata/corpus-golden.json || { printf 'the corpus golden lacks SC-DECISIONS-UNRECORDED\n' >&2; exit 1; }; for f in docs/user-guide/commands/spec.md .agents/skills/roundfix/references/spec.md .agents/skills/write-prd/SKILL.md .agents/skills/write-prd/references/clarification-record.md; do grep -qF -- 'SC-DECISIONS-UNRECORDED' "$f" || { printf 'missing SC-DECISIONS-UNRECORDED in %s\n' "$f" >&2; exit 1; }; done; for f in .agents/skills/write-prd/references/prd-template.md .agents/skills/write-prd/references/clarification-record.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'Autonomous route: no person present; decisions taken from' || { printf 'missing the autonomous-route line in %s\n' "$f" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < .agents/skills/write-techspec/SKILL.md | grep -qF -- 'Clarification Record' || { printf 'write-techspec lacks the Clarification Record\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Clarification Record**:' || { printf 'CONTEXT.md does not define **Clarification Record**\n' >&2; exit 1; }; make skills-sync-check || exit 1; go test -count=1 -run '^TestEveryOwnedSkillVersionIsRecorded$' ./skills || exit 1; go test -count=1 ./internal/speccheck || exit 1; go test -count=1 -run '^TestArchiveLayoutCharacterizationPinsCorpusGoldenAfterSpec0095$' ./internal/spec || exit 1; go test -count=1 -run '^(TestRunSpecCheckCleanText|TestSpecCheckRunVerification|TestRunSpecCheckGapStrictPromotion)$' ./internal/cli || exit 1; go test -count=1 -run 'SpecConsistency' ./internal/daemon || exit 1; go test -count=1 -tags docscontract -run '^(TestCheckCorpusGolden|TestCheckActiveCorpusHasNoErrors|TestTheSkillCoverageMapIsCurrent|TestGlossaryPhrasesAreDefinedTerms|TestGlossaryLinksNameAcceptedADRs|TestGlossaryAvoidLinesArePlainLists|TestGlossaryStatesEachRuleOnce)$' ./internal/docscontract` — expected: exit 0. Before this Task the four detector tests do not exist, so the command fails at its first pass check; afterwards the detector, corpus, spec check, settlement, skill version, Behavior Surface Record and glossary contract tests pass.

## References

- `_prd.md` → Goals; User Story 4; Core Feature 4; Core Feature 5; Success Metric 4; Success Metric 5; Glossary; Decisions
- `_techspec.md` → API Contract 3; API Contract 5; Invariants 3 and 6; Current behavior; Data Models; Testing Approach; Build Order 4
- ADR-0280; ADR-0244; ADR-0117

## Result

Implemented the Task 04 slice for Daemon Verification; status and settlement
remain Daemon-owned.

### Implementation and acceptance evidence

- A PRD at the guide horizon without either record shape now yields exactly
  one `SC-DECISIONS-UNRECORDED` gap. Both `Check` and the authoring-stage sweep
  call the detector before their missing-PRD return, and the code is registered
  at `StagePRD`. `TestClarificationRecordGapNamesAPRDWithoutTheRecord` checks
  the exact summary, heading-or-line-1 location, severity and fix for absent
  sections, missing shapes, invalid calendar dates and text outside Decisions.
  `TestClarificationRecordRunsFromThePRDStageAndFailsStrict` exercises PRD,
  TechSpec, Tasks and full stages and the shared `PromoteGaps` strict policy.
  The CLI exit checks remain part of Daemon Verification.
- A valid dated answer, the autonomous-route line, or both together satisfy
  the check without judging whether a person was present.
  `TestClarificationRecordAcceptsADatedAnswerOrTheAutonomousRoute` exercises
  each shape through every stage, including a valid leap-day answer.
- The detector uses the oldest guide-adding commit and Git ancestry. Earlier
  PRDs and repositories without the guide skip with the specified horizon
  text; missing PRDs skip with their path; an uncommitted PRD is held.
  `TestClarificationRecordSkipsAPRDBeforeItsHorizon` exercises those cases,
  the inclusive same-commit boundary, and a guide deleted and later re-added.
  Every new top-level test calls `t.Parallel()` first.
- The new clarification guide, write-prd Clarify and Record decisions steps,
  PRD template, write-techspec minimal PRD contract, spec user guide and
  Roundfix reference describe the record and its horizon. The new glossary
  entry sits beside Glossary Declaration and links accepted ADR-0280. The
  focused glossary contracts pass. QA settlement sections retain their bytes.
- The focused corpus measurement reported `SC-DECISIONS-UNRECORDED` at 0,
  with all earlier counts unchanged. The code joins `corpusFindingCodes`, the
  corpus golden and its archive-layout characterization; both update texts
  append the required sentence. A subsequent focused corpus check passes.
- Canonical skills were synchronized through `make skills-sync`; the owned
  version recorder raised Roundfix to `0.1.85`, write-prd to `0.0.12`, and
  write-techspec to `0.0.13` in both copies and recorded their digests. The
  coverage recorder changed only the spec user guide's fingerprint. Neither
  mirrors nor version/coverage records were edited by hand.

### Focused commands and outcomes

- `rtk go test ./internal/speccheck -run ClarificationRecord -count=1`:
  initially failed to compile because the detector constants did not exist;
  after implementation, the missing-PRD case exposed a full-sweep early return.
  Moving the detector before that return made all four tests pass.
- `rtk go test ./internal/speccheck -run 'ClarificationRecord|Glossary|CheckStage' -count=1`:
  exit 0, 44 tests/subtests passed, including the additional horizon boundaries.
- `rtk proxy go test -tags docscontract ./internal/docscontract -run '^TestCheckCorpusGolden$' -count=1`:
  measured the new code at 0 before updating the golden; the intentional code-set
  mismatch was the only change in the reported counts.
- `rtk go test -tags docscontract ./internal/docscontract -run 'TestCheckCorpusGolden|TestGlossary' -count=1`:
  exit 0, 43 tests/subtests passed after updating the pins and glossary.
- `make skills-sync`: exit 0.
- `rtk proxy go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`:
  exit 0; generated the three skill version increments and their records.
- `rtk proxy go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`:
  exit 0; generated the spec guide fingerprint.
- `make baseline-digests`: automatic approval review rejected workspace
  regeneration because it could rewrite derived files outside this Task's
  scope. It was not retried in this workspace.
- `rtk proxy go test -tags repocontract ./skills -run '^TestOwnedSkillEditLeavesDerivedArtifactsByteIdentical$' -count=1`:
  exit 0; the existing contract regenerates only in a disposable repository
  copy and confirms Baseline, characterization and archived artifacts retain
  their bytes after an owned skill edit. No Baseline artifact changed here.
- `rtk git diff --check`: exit 0. A read-only Python postflight also confirmed
  all 24 changed paths are declared by task_04, canonical skill bytes match
  mirrors, QA settlement sections are unchanged, earlier corpus counts are
  identical, only the spec guide fingerprint moved, and status is still
  `in_progress`.

### Handoff limits

The authored Verification command, repository Verification and incremental
Verification were not run. The Daemon owns those gates and the Task verdict.
No other Task file or Task Graph manifest was edited; no commit, push or pull
request was created. No implementation follow-up was identified.

## Carry-forward provenance

- Source Run: `run_20261017T101202Z_03bf9e36c3f11f8c`
- Source commit: `ab913cb454d59fc92a27c27e517ed1bb83bf776d`
