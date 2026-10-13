---
task: task_01
spec: 0258-authoring-that-proves-what-it-claims
status: pending
type: backend
complexity: medium
---

# Task 01: Spec check reports a TechSpec without its absence evidence, ADR alternatives or cross-read

## Overview

The Spec Consistency Check gains the TechSpec-stage detector of
`_techspec.md` → Interfaces and Data Models: `SC-ABSENCE-UNPROVEN`,
`SC-ADR-ALTERNATIVES-UNRELATED` and `SC-COVERAGE-UNREAD`, each a gap on every
active Spec with a TechSpec. The six fixture TechSpecs the prototype measured
gain the three proofs so the existing tests keep their expectations, and the
corpus pins gain the three codes at 0. It is verifiable on its own through the
new detector tests and the consuming suites.

## Requirements

1. MUST answer gaps 2, 3 and 4 of the adopted Backlog Entry of 2026-10-12,
   [Spec authoring does not prove what it claims](references/2026-10-12-spec-authoring-does-not-prove-what-it-claims.md),
   by implementing API Contracts 1, 2 and 3 of `_techspec.md` with the grammar
   of `_techspec.md` → Data Models, in a new file beside the existing
   detectors.
2. MUST run the detector in the full sweep where the PRD and TechSpec are
   read, and at the TechSpec stage, and MUST register the three codes at the
   TechSpec stage so earlier stages list them as skipped. A Spec without a
   TechSpec MUST list the three codes as skipped (ADR-0094).
3. MUST hold every active Spec, with no commit horizon (ADR-0265), and MUST
   make each finding a gap whose fix names the section or line and the
   write-techspec rule: "Absence evidence", "Rejected alternatives" or
   "Coverage cross-read".
4. MUST NOT execute any command an Absence Evidence item names, and MUST NOT
   read content inside fenced blocks as a declaration.
5. MUST add the three proofs to the TechSpec of the `clean`,
   `vocabulary-missing`, `citation-dirty`, `coverage-unmapped`,
   `coverage-untasked` and `reference-unresolved` fixtures, at the end of each
   file or inside its Coverage Map, without changing the findings those
   fixtures exist to produce.
6. MUST add the three codes at 0 to `corpusFindingCodes`, to the corpus golden
   and to the characterization that pins it, and append one identical sentence
   naming them to the update text of both.
7. Every new top-level test MUST call `t.Parallel()` first, and each detector
   test MUST prove the gap with the section or line missing and no finding of
   that code with it present.

## Subtasks

- [ ] Write the detector and its five tests.
- [ ] Wire it into the full sweep and the TechSpec stage.
- [ ] Complete the six fixture TechSpecs.
- [ ] Add the codes to the corpus pins.

## Acceptance Criteria

- [ ] A TechSpec missing Absence Evidence, an ADR line or the Deviations line
      yields the matching gap; with them present it yields none.
- [ ] An Absence Evidence item without a command or a result is a gap.
- [ ] A deviation bullet that names no PRD unit is a gap.
- [ ] `--strict` turns each gap into a failing result, and a Spec without a
      TechSpec lists the three codes as skipped.
- [ ] Every fix names its section or line and its skill rule.
- [ ] The existing speccheck, spec check command and settlement tests pass
      unchanged.

## Consumers

- The Spec Consistency Check finding set — consumers: `./internal/speccheck`, `./internal/cli`, `./internal/daemon`, `./internal/docscontract`

## Context

- instruction: `docs/adr/0265-a-spec-proves-what-it-claims-and-every-active-spec-is-held-to-it.md`
- instruction: `.agents/skills/roundfix/references/spec.md`
- creates: `internal/speccheck/authoring_proof.go`
- creates: `internal/speccheck/authoring_proof_test.go`
- interface: `internal/speccheck/citations.go`
- interface: `internal/speccheck/coherence.go`
- interface: `internal/speccheck/testdata/repo/docs/specs/clean/_techspec.md`
- interface: `internal/speccheck/testdata/repo/docs/specs/vocabulary-missing/_techspec.md`
- interface: `internal/speccheck/testdata/repo/docs/specs/citation-dirty/_techspec.md`
- interface: `internal/speccheck/testdata/repo/docs/specs/coverage-unmapped/_techspec.md`
- interface: `internal/speccheck/testdata/repo/docs/specs/coverage-untasked/_techspec.md`
- interface: `internal/speccheck/testdata/repo/docs/specs/reference-unresolved/_techspec.md`
- interface: `internal/docscontract/corpus_test.go`
- interface: `internal/docscontract/testdata/corpus-golden.json`
- interface: `internal/spec/archive_layout_characterization_test.go`

## Verification

- `out="$(go test -count=1 -v -run '^(TestAbsenceEvidenceGapFollowsTheSection|TestADRAlternativesGapNamesEachListedADR|TestCoverageCrossReadGapFollowsTheDeviationsLine|TestTechSpecProofRunsAtTheTechSpecStageAndInTheFullSweep|TestTechSpecProofFixesNameTheSectionAndTheSkillRule)$' ./internal/speccheck 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestAbsenceEvidenceGapFollowsTheSection TestADRAlternativesGapNamesEachListedADR TestCoverageCrossReadGapFollowsTheDeviationsLine TestTechSpecProofRunsAtTheTechSpecStageAndInTheFullSweep TestTechSpecProofFixesNameTheSectionAndTheSkillRule; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; for code in SC-ABSENCE-UNPROVEN SC-ADR-ALTERNATIVES-UNRELATED SC-COVERAGE-UNREAD; do grep -qF -- "\"$code\": 0" internal/docscontract/testdata/corpus-golden.json || { printf 'the corpus golden lacks %s\n' "$code" >&2; exit 1; }; done; go test -count=1 ./internal/speccheck || exit 1; go test -count=1 -run '^TestArchiveLayoutCharacterizationPinsCorpusGoldenAfterSpec0095$' ./internal/spec || exit 1; go test -count=1 -run '^(TestRunSpecCheckCleanText|TestSpecCheckRunVerification|TestSpecCheckClassifiesTheGateBoundary|TestRunSpecAuditPreservesSpecCheckBehavior)$' ./internal/cli || exit 1; go test -count=1 -run 'SpecConsistency' ./internal/daemon || exit 1; go test -count=1 -tags docscontract -run '^(TestCheckCorpusGolden|TestCheckActiveCorpusHasNoErrors)$' ./internal/docscontract` — expected: exit 0; before this Task the five detector tests do not exist, so the command fails at its first pass check.

## References

- [_prd.md](_prd.md) — Goals 1-3 and 5; User Stories 1-3 and 5; Core Features 1-3 and 5; Success Metric 1; Success Metric 4
- [_techspec.md](_techspec.md) — Interfaces; Data Models; API Contract 1; API Contract 2; API Contract 3; API Contract 5; Invariants 1-7; Testing Approach; Build Order 1
- ADR-0265; ADR-0093; ADR-0094; ADR-0156; ADR-0259
