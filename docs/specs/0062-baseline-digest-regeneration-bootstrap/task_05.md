---
task: task_05
spec: 0062-baseline-digest-regeneration-bootstrap
status: pending
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
