---
task: task_01
spec: 0334-the-template-tooling-row-passes-and-a-wave-collision-refuses-before-the-run
status: completed
type: backend
complexity: medium
---

# Task 01: The templates' no-mutation Tooling authority row passes the Spec Check when its record grants no path

## Overview

The Spec Check reads the Tooling authority row's own words before it treats
a cited authorization record as a mutation, and accepts "no protected tooling
mutation" when that record is an operative grant with `paths: []`. The
`write-prd` and `write-techspec` templates say how such a row ends, and a
documentation contract test composes the row from each template and runs the
checker on it. This answers the adopted Backlog Entry of 2026-10-27, "The
authoring templates and the tooling check disagree on no tooling change"
([references](references/2026-10-27-the-authoring-templates-and-the-tooling-check-disagree-on-no-tooling-change.md)).

## Requirements

1. MUST change `declaresProtectedToolingMutation` in
   `internal/speccheck/constraints.go` as `_techspec.md` → API Contract 1
   states: a row whose whitespace-normalized, lower-cased reason contains
   `no protected tooling mutation` declares no mutation when it cites no
   record, or when the record selected for the row resolves through
   `spec.ReadAuthorization` as `spec.AuthorizationGranted` for the asking
   Spec with an empty `Record.Paths`. A row citing a record that lists any
   path, and a row without those words, keep today's verdict.
   `SC-TOOLING-UNBOUNDED` keeps its code, summary and fix.
2. MUST add `TestANoMutationRowCitingAnEmptyGrantIsNotUnbounded` to
   `internal/speccheck/report_test.go`, using `writeToolingRowFixture` with
   the Spec-contained record `docs/specs/0114-tooling-row/_authorization.md`
   (`status: approved`, a `granted` date, the five delivery operations,
   `consuming: 0114-tooling-row`), with the three cases `_techspec.md` →
   Testing Approach names: the no-mutation row on `paths: []` gives no
   `SC-TOOLING-UNBOUNDED`; the same row on a `paths` list holding `Makefile`
   gives it; `applicable — protected tooling mutation proposed and
   authorized` on `paths: []` gives it. Keep every existing case of
   `TestToolingRowStatesApplicability` unchanged.
3. MUST change the no-mutation comment of the canonical
   `.agents/skills/write-prd/references/prd-template.md` and
   `.agents/skills/write-techspec/references/techspec-template.md`
   (`_techspec.md` → API Contract 2): keep
   `With no protected tooling mutation, record: ` followed by the reason in
   backticks exactly as today, and add that the row still cites the
   Spec-contained authorization record and, when its `paths` is `[]`, must
   `end the row with` followed by `bounded files: none` in backticks. Change
   no other template line.
4. MUST add `TestTheTemplateNoMutationToolingRowPassesTheSpecCheck` to the
   existing `//verify:always` file
   `internal/docscontract/command_documentation_test.go` (never a new file in
   that package). For each template it reads the canonical file from the
   repository root (`baselineDocumentationRepoRoot`), extracts the reason
   after `With no protected tooling mutation, record: ` and the
   `bounded files: …` value after `end the row with`, fails naming the
   template when either is missing, composes the full row (reason, Source,
   Spec-contained authorization record, bounded files) in a temporary
   repository with the other three Project Constraint rows, and requires
   `speccheck.CheckStage` at `StagePRD` for the PRD template and
   `StageTechSpec` for the TechSpec template to give no
   `SC-TOOLING-UNBOUNDED` on a record with `paths: []`, and to give it on a
   record whose `paths` lists `Makefile`.
5. MUST copy the canonical skills to `skills/` with `make skills-sync`, then
   raise both version fields of `.agents/skills/write-prd/SKILL.md`,
   `.agents/skills/write-techspec/SKILL.md` and their mirrors with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$'
   -record-skill-versions -count=1`, recording
   `skills/testdata/owned-skill-versions.json`. Never hand-edit a mirror or
   a version field. Spec 0333 may have raised these versions first; the
   record command raises above whatever is recorded.
6. MUST NOT touch any `### QA settlement` section, any Baseline module,
   `docs/agents/`, the Makefile, CI workflows, `go.mod` or `go.sum`.

## Subtasks

- [ ] Characterize: the no-mutation row citing a `paths: []` record reports `SC-TOOLING-UNBOUNDED` today.
- [ ] Change the row reading and add the Spec Check test.
- [ ] Add the templates' ending, sync the mirrors and record the versions.
- [ ] Add the documentation contract test.

## Acceptance Criteria

- [ ] The templates' no-mutation row with a `paths: []` record passes at the PRD and TechSpec stages.
- [ ] A row whose record lists a Governed Path, or that declares a mutation, without bounded files is still `SC-TOOLING-UNBOUNDED`.
- [ ] Both templates name the `bounded files: none` ending, the mirrors equal their canonical copies, and the owned-skill versions are recorded.

## Context

- instruction: `docs/adr/0347-a-row-citing-an-empty-grant-declares-no-mutation-and-a-completed-task-never-collides.md`
- instruction: `docs/adr/0131-the-tooling-row-states-applicability-not-mutation.md`
- instruction: `docs/adr/0179-an-explicitly-empty-paths-list-grants-operations-only.md`
- interface: `internal/speccheck/constraints.go`
- interface: `internal/speccheck/report_test.go`
- interface: `internal/docscontract/command_documentation_test.go`
- interface: `.agents/skills/write-prd/references/prd-template.md`
- interface: `skills/write-prd/references/prd-template.md`
- interface: `.agents/skills/write-prd/SKILL.md`
- interface: `skills/write-prd/SKILL.md`
- interface: `.agents/skills/write-techspec/references/techspec-template.md`
- interface: `skills/write-techspec/references/techspec-template.md`
- interface: `.agents/skills/write-techspec/SKILL.md`
- interface: `skills/write-techspec/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`

## Consumers

- The Tooling authority row reading — consumers: `./internal/speccheck`, `./internal/docscontract`
- The owned-skill version record and the authoring template contracts — consumers: `./skills`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^TestANoMutationRowCitingAnEmptyGrantIsNotUnbounded$' ./internal/speccheck 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestANoMutationRowCitingAnEmptyGrantIsNotUnbounded ' || { echo 'missing pass: TestANoMutationRowCitingAnEmptyGrantIsNotUnbounded' >&2; exit 1; }; go test -count=1 -timeout 20m -run 'Tooling|Constraint' ./internal/speccheck` — expected: exit 0. Before this Task the test does not exist, so the pass line is missing; afterwards it passes with the existing tooling and constraint tests.
- `out="$(go test -count=1 -timeout 20m -tags docscontract -v -run '^TestTheTemplateNoMutationToolingRowPassesTheSpecCheck$' ./internal/docscontract 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestTheTemplateNoMutationToolingRowPassesTheSpecCheck ' || { echo 'missing pass: TestTheTemplateNoMutationToolingRowPassesTheSpecCheck' >&2; exit 1; }` — expected: exit 0. Before this Task the test does not exist; afterwards the row composed from each template passes and the Makefile variant is refused.
- `for f in .agents/skills/write-prd/references/prd-template.md skills/write-prd/references/prd-template.md .agents/skills/write-techspec/references/techspec-template.md skills/write-techspec/references/techspec-template.md; do tr -s '[:space:]' ' ' < "$f" | grep -q 'end the row with .bounded files: none.' || { printf 'missing row ending in %s\n' "$f" >&2; exit 1; }; done; cmp .agents/skills/write-prd/SKILL.md skills/write-prd/SKILL.md && cmp .agents/skills/write-techspec/SKILL.md skills/write-techspec/SKILL.md && go test -count=1 -timeout 20m ./skills` — expected: exit 0. Before this Task no template names the ending; afterwards all four copies do, the `SKILL.md` mirrors equal their canonical copies, and the owned-skill version and authoring template contracts pass.

## References

- `_prd.md` → Goals 1-3; Core Features 1-2; Success Metric 1; Success Metric 2; Success Metric 3
- `_techspec.md` → Current behavior 1-3; API Contract 1; API Contract 2; Invariants 1-2; Testing Approach; Build Order 1
- ADR-0347; ADR-0131; ADR-0179


## Result

Skills loaded: implement-task, coding-guidelines, golang-testing, testing-boss, systematic-debugging, no-workarounds, tech-writer, skill-creator, evidence-gate.

The row's normalized reason now decides whether it declares a mutation before
its record citation does. A selected record exempts the no-mutation row only
when `spec.ReadAuthorization` returns `spec.AuthorizationGranted` for the
asking Spec with empty `Record.Paths`. Record selection, other authorization
diagnostics, and the unbounded finding's code, summary and fix retain their
existing behavior.

Acceptance evidence:

- The PRD and TechSpec template rows with `paths: []` produce no
  `SC-TOOLING-UNBOUNDED` in the new documentation contract. The test reads
  each canonical template, extracts its reason and ending, composes the row
  with its Source and Spec-contained record in a temporary repository, and
  calls `CheckStage` at the corresponding authoring stage.
- Both template rows still produce exactly one unbounded finding when the
  record lists `Makefile`. The new three-case Spec Check regression also
  proves that a declared mutation with an empty grant remains unbounded.
  Every case of `TestToolingRowStatesApplicability` is unchanged and passed
  in the focused check.
- Each canonical template has exactly one added sentence, inside its
  existing no-mutation comment, retaining the original reason and naming
  `bounded files: none`. `make skills-sync` regenerated the mirrors. The
  sanctioned version recorder raised both version fields in each skill and
  its mirror from `0.0.14` to `0.0.15` and appended their digests to
  `skills/testdata/owned-skill-versions.json`. The focused synchronization,
  version, template-label, authorization-journey and constraint-ownership
  contracts passed.

Focused checks and characterization:

- Before the checker edit,
  `go test ./internal/speccheck -run '^TestANoMutationRowCitingAnEmptyGrantIsNotUnbounded$/no_mutation_with_an_empty_grant$' -count=1`
  exited 1 with `SC-TOOLING-UNBOUNDED`, proving the reported defect.
- Before the template edit, a focused PRD-template subtest exited 1 naming
  `prd-template.md` and its missing bounded-files ending.
- `go test -count=1 -tags docscontract ./internal/speccheck ./internal/docscontract -run '^(TestANoMutationRowCitingAnEmptyGrantIsNotUnbounded|TestTheTemplateNoMutationToolingRowPassesTheSpecCheck|TestToolingRowStatesApplicability)$' -v`
  exited 0: all three new checker cases, all four template/grant cases and
  all four existing applicability cases passed.
- `make skills-sync` exited 0.
- `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions -count=1`
  exited 0; this is the requirement's source/record generator.
- A focused `./skills` selection of `TestEveryOwnedSkillVersionIsRecorded`,
  `TestAuthorialSkillSync`, `TestAuthoringTemplatesUseTheBoundedFilesLabel`,
  `TestToolingAuthorizationJourney` and `TestAuthoringConstraintOwnership`
  exited 0 after regeneration.
- Two temporary `go test -overlay` probes exercised the new documentation
  contract without modifying repository sources: restoring the old checker
  exited 1 with the unexpected unbounded finding for both empty-grant
  stages; supplying the pre-edit templates exited 1 with a missing-ending
  diagnostic naming each template. Both failures matched the intended
  regression, and the probes removed their temporary files.
- Final scope inspection and `git diff --check` found only this Task's
  declared paths changed. The initial working-tree change was the Daemon's
  `task_01.md` status transition; it remains `in_progress`.

Declared Verification and repository-wide verification were not run; the
Daemon owns those checks and Task settlement. No Task 02 implementation,
Task Graph, other Task, Baseline artifact, `docs/agents/` guide, Makefile,
workflow, Go module, QA settlement text or commit was changed.

## Carry-forward provenance

- Source Run: `run_20261031T173520Z_ca95397d60f28c1a`
- Source commit: `fdcee763fe7f15c39fc61d6d6f30efb19d036925`
