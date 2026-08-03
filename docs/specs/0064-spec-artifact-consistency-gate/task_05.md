---
task: task_05
spec: 0064-spec-artifact-consistency-gate
status: pending
type: test
complexity: high
---

# Task 05: Replay the four findings that motivated this Spec

## Overview

Build this Spec's acceptance: four fixture Spec folders reproducing the
artifact shapes behind Spec 0058's QA-001 and QA-004 and Spec 0056's F-001 and
F-002, each asserting the exact code the check reports. Add the corpus sweep
that measures the check against every Spec in the repository, so the
false-positive rate is a number in a test rather than a claim in prose.

The fixtures are authored from each QA report's Expected and Actual, with the
report path recorded as provenance. They are not recovered from Git: the
pre-remediation artifacts lived on Run Branches that no longer exist and the
archived copies are post-fix. Saying so is part of the deliverable — a corpus
claiming a provenance it does not have is the defect class this Spec removes.

## Requirements

1. MUST add four fixture Spec folders, each recording in a `README` the QA
   report path it reproduces and the fact that the shape is authored from that
   report rather than recovered from Git.
2. MUST assert the 0058 QA-001 fixture — a PRD Core Feature the Coverage Map
   does not cover — reports `SC-COVERAGE-UNMAPPED` naming that feature.
3. MUST assert the 0058 QA-004 fixture — five emitted tokens, four documented —
   reports `SC-VOCABULARY-UNDOCUMENTED` naming the undocumented token.
4. MUST assert the 0056 F-001 fixture reports both halves: `SC-ADR-UNLISTED`
   for the ADR its TechSpec cites and its PRD row omits, and `SC-ADR-RELATED`
   for the ADR that cites two ADRs the Spec lists and is itself unlisted.
5. MUST assert the 0056 F-002 fixture — a PRD Core Feature the TechSpec
   narrows without a Coverage Map entry — reports `SC-COVERAGE-UNMAPPED`.
6. MUST add a sweep over every Spec in this repository that records the
   per-code finding counts as a checked-in golden characterization number.
   Archived Specs predate this contract and stay byte-identical, so their
   counts are measured and recorded, never asserted at zero. Bringing the
   **active** Specs to zero errors, and the test that holds them there,
   belong to the Task that performs that work.
7. MUST assert the sweep over the full corpus completes inside one second,
   measured in the test.
8. SHOULD make the golden number's update path obvious, so a deliberate
   detector change is a declared break rather than a mystery diff.

## Subtasks

- [ ] Author the four replay fixtures with their provenance READMEs.
- [ ] Assert the reported code and locations for each.
- [ ] Add the whole-corpus sweep recording golden per-code counts.
- [ ] Add the duration budget assertion.

## Acceptance Criteria

- [ ] Each of the four fixtures reports the code its QA report describes, and
      the test names the report path.
- [ ] Each fixture's `README` records the QA report path and states the shape
      is authored from the report, not recovered from Git.
- [ ] The whole-corpus sweep compares against a checked-in golden per-code
      count and fails when any count moves.
- [ ] The golden file records counts for archived and active Specs separately,
      so the Task that cleans the active Specs can move one without the other.
- [ ] The full-corpus sweep completes inside one second, asserted by the test.
- [ ] No archived Spec file is modified by this Task.

## Context

- instruction: `docs/agents/spec-routing.md`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/speccheck -count=1 -run 'Replay|Corpus|Budget' -v | grep -q -- "--- PASS"`
  — expected: exit 0; the replay, corpus, and budget tests ran and passed.
- `go test ./internal/speccheck -count=1` — expected: exit 0.
- `git diff --name-only HEAD -- docs/specs/_archived | grep -q . && exit 1 || exit 0`
  — expected: exit 0; no archived Spec file changed.

## References

- `_prd.md` → Success Metrics; Core Features 2, 4, and 6.
- `_techspec.md` → Testing Approach; Build Order 6.
- ADR-0093, ADR-0094.
