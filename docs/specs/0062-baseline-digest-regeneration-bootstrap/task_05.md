---
task: task_05
spec: 0062-baseline-digest-regeneration-bootstrap
status: completed
type: backend
complexity: low
---

# Task 05: Say what the regenerator cannot supply

## Overview

Adding a Normative Clause to a Baseline module reports a missing Source
Baseline manifest row for every new clause. The message names the clause and
stops there, so the reader's natural next move is to run the regeneration
again — which cannot help, because the regenerator maintains digests and spans
for rows that already exist and never creates one. This Task makes the
limitation legible at the moment it bites.

## Requirements

1. MUST extend the missing-clause and missing-rule Source Baseline diagnostics
   so each states that the regenerator maintains manifest rows but never
   creates them.
2. MUST keep naming the specific clause or rule that is missing, so the message
   gains guidance without losing its subject.
3. MUST keep the existing diagnostic codes unchanged, so anything matching on
   code keeps working.
4. MUST NOT generate manifest rows, infer their spans, or otherwise change what
   the manifest is; this Task makes the limitation legible, not absent.
5. MUST leave every other diagnostic's text unchanged.

## Subtasks

- [ ] Extend the missing-clause diagnostic with the guidance sentence.
- [ ] Extend the missing-rule diagnostic the same way.
- [ ] Confirm the codes and subjects are unchanged.

## Acceptance Criteria

- [ ] A catalog missing a Source Baseline manifest row for a clause produces a
      diagnostic naming both the clause and the fact that the regenerator
      cannot add the row.
- [ ] The same holds for a missing rule.
- [ ] Both diagnostic codes are byte-identical to their previous values.
- [ ] No other diagnostic's text changed, proven by the characterization
      corpus, which is updated in this Task only for the two messages this Task
      changes.
- [ ] `git status --porcelain` shows no path outside `internal/baseline/` and
      this task file.

## Context

- interface: `internal/baseline/catalog_validate.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/baseline -run TestSourceBaselineManifestRowGuidance -count=1`
  — expected: exit 0; both diagnostics carry the guidance and still name their
  subject.
- `grep -q "catalog.sourceBaseline.required-clause.missing" internal/baseline/catalog_validate.go`
  — expected: exit 0; the code is unchanged.
- `grep -q "catalog.sourceBaseline.required-rule.missing" internal/baseline/catalog_validate.go`
  — expected: exit 0; the code is unchanged.
- `go test ./internal/baseline -run TestCatalogDiagnosticCharacterization -count=1`
  — expected: exit 0; the corpus reflects exactly the two intended message
  changes and nothing else.
- `go test ./internal/baseline -count=1` — expected: exit 0.

## References

- `_prd.md` → Core Features 3; Non-Goals (not generating manifest rows).
- `_techspec.md` → API Contracts; Build Order 5.

## Result

### Implementation

- Extended only the `required-clause.missing` and `required-rule.missing`
  diagnostic details with one shared sentence: the regenerator maintains
  manifest rows but never creates them, so the named row must be added first.
- Kept each existing diagnostic code and Source Baseline subject unchanged.
  The missing clause or rule identifier remains the first part of the detail.
- Added a package-level table test that constructs a newly required clause and
  rule without adding manifest rows, then asserts each exact code, subject, and
  guided detail. The fixture never creates or infers a manifest row.
- Re-recorded the characterization corpus. Its diff changes only the detail of
  the existing missing-rule records and the existing missing-clause record.

### Focused checks

- Red signal before the production edit:
  `GOCACHE=/private/tmp/roundfix-task05-gocache rtk proxy go test ./internal/baseline -run '^TestSourceBaselineManifestRowGuidance$/^missing (clause|rule)$' -count=1`
  exited 1 because both details contained only their missing identifiers.
- `GOCACHE=/private/tmp/roundfix-task05-gocache rtk go test ./internal/baseline -run '^TestCatalogDiagnosticCharacterization$' -update-catalog-diagnostics -count=1`
  exited 0 and re-recorded the intentional message changes.
- `GOCACHE=/private/tmp/roundfix-task05-gocache rtk go test ./internal/baseline -run '^(TestSourceBaselineManifestRowGuidance|TestCatalogDiagnosticCharacterization)$' -count=1`
  exited 0 with 5 passing test and subtest results after the final code edit.
- `rtk git diff --check` exited 0.
- The Task's declared `## Verification` commands were not run; they remain for
  the Daemon.

### Acceptance-criterion evidence

1. `TestSourceBaselineManifestRowGuidance/missing_clause` passed while
   asserting the exact missing clause identifier and the regenerator guidance.
2. `TestSourceBaselineManifestRowGuidance/missing_rule` passed with the same
   assertions for a required rule without a manifest row.
3. The table test asserts the byte-identical codes
   `catalog.sourceBaseline.required-clause.missing` and
   `catalog.sourceBaseline.required-rule.missing`; the production diff changes
   only the detail arguments at their existing call sites.
4. `TestCatalogDiagnosticCharacterization` passed, and the golden diff contains
   only three detail replacements under those two codes; no other diagnostic
   record changed.
5. `rtk git -c core.fsmonitor=false status --short` listed only this Task file,
   `internal/baseline/catalog_validate.go`, `internal/baseline/catalog_test.go`,
   and `internal/baseline/testdata/catalog.diagnostics.golden.json`.
