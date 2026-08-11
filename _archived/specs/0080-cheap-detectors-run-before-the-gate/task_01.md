---
task: task_01
spec: 0080-cheap-detectors-run-before-the-gate
status: completed
type: backend
complexity: high
---

# Task 01: Detect the facts an agent turn resolves today

## Overview

The mechanical detectors, plus the typed result they produce and the
Daemon-owned materializer that turns it into report rows. Everything lands
inert: nothing calls it from the gate yet, so the repository is never left red
for the next Task.

Today those facts are resolved inside agent turns: each Task commit's changed
paths against the authorization, commit shape, the report's structure, and
evidence-path resolution, at roughly thirty minutes a round.
Every one of those is a Git read or a file read with a written answer.

## Requirements

1. MUST implement detectors for the facts the gate performs by hand: the
   authorization's bounded files against each Task commit's actual changed
   paths; consequent-fix commit shape, including a fix folded into its
   authorized commit or ordered before its cause; the report's structural
   contract — no pending row, a terminal per-row status, blocked strings
   matching their typed cause, the three counts equal to the actual row
   counts; and evidence-path resolution.
2. MUST return the typed `MechanicalResult` the TechSpec declares, carrying
   `Findings`, `Carried`, `Blocked`, `Skips`, and `Blocking`, with every
   Finding naming code, file, line, detail, fix, and the row it blocks when
   known.
3. MUST report every failure it finds in one pass, never only the first — a
   fail-fast stage spends the single Verification repair on the first of N
   defects, which this repository already measured and recorded.
4. MUST be citation-only per ADR-0093: each detector compares written
   declarations with observable repository facts, and none infers intent,
   judges prose, or evaluates whether a decision was correct.
5. MUST be presence-aware per ADR-0094: an absent input artifact produces a
   `MechanicalSkip` naming the detector and the missing artifact, never a
   failure and never silence.
6. MUST implement the Daemon-owned materializer that writes the typed result
   into report rows: Findings become finding sections, Carried rows retain
   their establishing report and head, Blocked rows use the
   `blocked (finding: …)` form, and each Skip records its detector and missing
   artifact.
7. MUST NOT compute a verdict, settle a Task, or edit an existing report. The
   materializer applies the semantics ADR-0080 already owns; it never
   redefines them.
8. MUST leave the gate calling nothing new: this Task adds capability only.

## Subtasks

- [ ] Implement the detectors with their carrier fixtures.
- [ ] Implement the typed result and the materializer.
- [ ] Prove presence-aware skips and corpus non-regression.

## Acceptance Criteria

- [ ] Each detector fails its red fixture with its own code and passes its
      green fixture.
- [ ] Each detector records a skip, not a failure, when its input is absent.
- [ ] A result carrying several failures reports all of them.
- [ ] The pre-existing fixture corpus reports no new diagnostics and the
      budget guard still passes.
- [ ] Nothing in the gate's runtime path calls the new code yet.

## Context

- interface: internal/speccheck/citations.go
- interface: internal/speccheck/report.go
- instruction: docs/workflow/authorizations/2026-08-06-proof-cost.md

## Verification

- `output="$(go test -count=1 ./internal/speccheck -run 'Mechanical|AuthPaths|ConsequentOrder|ReportShape|EvidencePath' -v 2>&1)"; st=$?; printf '%s\n' "$output" | grep -q -- '--- PASS' && [ "$st" -eq 0 ]`
  — expected: exit 0; the named detector tests exist, are selected, and pass —
  an empty selection cannot satisfy this.
- `output="$(go test -count=1 ./internal/speccheck -run '^TestCheckCorpusBudget$' -v 2>&1)"; st=$?; printf '%s\n' "$output" | grep -q -- '--- PASS' && [ "$st" -eq 0 ]`
  — expected: exit 0; the corpus budget guard still selects and passes.
  — expected: exit 0; both packages the change lands in stay green.
- `grep -rq 'MechanicalResult' internal/ && grep -rq 'MechanicalSkip' internal/`
  — expected: exit 0; the typed result and its skip record exist on the real
  surface.

These commands are deliberately absent: `go build -buildvcs=false ./...` and a
whole-package `go test` sweep both pass against a tree where no work has
happened, so each approves the Task before it starts. Compilation and
regression are the Run-level gate's job; the commands above name cases that
do not exist yet.

## References

- `_prd.md` → Core Features 1, 2, 3; User Story 3.
- `_techspec.md` → Implementation Design (Interfaces); Build Order 1.
- ADR-0096, ADR-0093, ADR-0094, ADR-0080.

## Result

Implemented the inert mechanical-stage capability in `internal/speccheck`.
`RunMechanicalStage` compares exact authorization paths with real
`git diff-tree` output, checks declared consequent-fix ancestry, validates QA
Report row statuses and typed blocked counts, and resolves report evidence
paths. It accumulates every diagnostic into the typed `MechanicalResult` and
records absent inputs as `MechanicalSkip` values. `WriteMechanicalResult`
writes finding sections, carried provenance, finding-blocked rows, and skips
to a supplied writer without choosing a report path or verdict.

Focused checks run after the final implementation edit:

- `rtk go test ./internal/speccheck -run '^(TestMechanicalAuthPaths|TestMechanicalConsequentOrder|TestMechanicalReportShape|TestMechanicalEvidencePath|TestMechanicalReportsAllFindings|TestMaterializeMechanicalResult|TestMechanicalCorpusNonRegression)$' -count=1`
  passed: 17 tests in one package.
- `rtk go test ./internal/speccheck -run '^TestRenderResultTextAndJSON$' -count=1`
  passed: the existing consistency-report rendering contract remains green.
- `rtk go vet ./internal/speccheck` passed with no diagnostics.
- `rtk grep -r -n 'RunMechanicalStage' internal/daemon internal/agent internal/cli`
  returned no matches (exit 1), proving the gate runtime does not call the new
  entry point in this slice.

Acceptance evidence:

- Each detector has file-backed green/red carriers or a real temporary Git
  carrier and returns its own stable code: `QA-AUTH-PATHS`,
  `QA-CONSEQUENT-ORDER`, `QA-REPORT-SHAPE`, or `QA-EVIDENCE-PATH`.
- The focused cases cover missing authorization, consequent-fix declaration,
  and QA Report inputs; each produces a detector-plus-artifact skip and no
  corresponding finding.
- `TestMechanicalReportsAllFindings` produces all four codes in one result,
  checks `Blocking`, verifies every required Finding field, and checks known
  row-to-finding links.
- `TestMechanicalCorpusNonRegression` walked the pre-existing fixture corpus
  through the same `assertMechanicalCorpusBudget` helper used by the authored
  `TestCheckCorpusBudget` guard and reported no new mechanical diagnostics.
- The runtime-path search above found no call. Task 03 remains the owner of
  Daemon wiring, Agent Session withholding, report lifecycle, and verdict
  settlement.

The commands under this Task's `## Verification` were not run; the Daemon owns
that complete selection and settlement evidence.
