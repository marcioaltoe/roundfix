---
task: task_02
spec: 0064-spec-artifact-consistency-gate
status: pending
type: backend
complexity: high
---

# Task 02: Report ADR and coverage gaps from citations

## Overview

Deliver the detectors that read the accepted ADR corpus and the Spec's own
cross-references: an ADR the Spec cites but does not list, an accepted ADR that
cites an ADR the Spec lists and is itself unlisted, a PRD Core Feature or User
Story no Coverage Map entry covers, one no Task references, and a declared
reference path that does not resolve. Verifiable on its own through fixture
Spec folders paired with a fixture ADR corpus.

Detection is by citation only. The check never judges whether an ADR is
topically relevant, which is what keeps it deterministic and sub-second.

## Requirements

1. MUST build an ADR corpus from `status: accepted` ADR files, carrying each
   ADR's number, title, and the set of ADR numbers its body cites.
2. MUST report `SC-ADR-UNLISTED` at severity `error` when an ADR cited anywhere
   in the Spec's artifacts is absent from the Active ADR obligations row,
   locating the citation and the row.
3. MUST report `SC-ADR-RELATED` at severity `gap` for an accepted ADR that
   cites an ADR the Spec lists and is itself unlisted. It MUST use the
   depth-one closure only, and MUST NOT report an ADR the Spec already lists.
4. MUST report `SC-COVERAGE-UNMAPPED` when a PRD Core Feature or User Story is
   covered by no Coverage Map entry, and MUST accept collective entry forms
   such as `Core Features 1-5 →` as covering each member of the range.
5. MUST report `SC-COVERAGE-UNTASKED` when a PRD Core Feature or User Story is
   referenced by no Task file, reading the Task Graph through the existing
   loader.
6. MUST report `SC-REF-UNRESOLVED` for a Task Context path or a reference index
   entry that does not resolve on disk.
7. MUST skip each detector, recording the skip, when its input artifact is
   absent — a Spec with no TechSpec reports no coverage-map finding, and a Spec
   with no Task Graph reports no task-coverage finding.
8. SHOULD keep gap findings out of the non-zero exit path; severity is what the
   command layer reads.

## Subtasks

- [ ] Read and index the accepted ADR corpus and its citation graph.
- [ ] Extract cited ADR numbers from every Spec artifact and the Active ADR row.
- [ ] Implement the two ADR detectors, including the depth-one closure.
- [ ] Parse PRD Core Features and User Stories, and the TechSpec Coverage Map
      including collective ranges.
- [ ] Implement the two coverage detectors and the reference resolver.
- [ ] Add fixture Spec folders and a fixture ADR corpus, with tests covering
      each detector clean, dirty, and absent.

## Acceptance Criteria

- [ ] A fixture Spec whose TechSpec cites an ADR its PRD's Active ADR row omits
      reports `SC-ADR-UNLISTED` with a location in each file.
- [ ] A fixture corpus in which an accepted ADR cites two ADRs the Spec lists
      reports that ADR once as `SC-ADR-RELATED` at severity `gap`.
- [ ] An ADR the Spec already lists is never reported as `SC-ADR-RELATED`.
- [ ] A fixture Spec whose Coverage Map lists Goals and Stories but no Core
      Feature reports `SC-COVERAGE-UNMAPPED` once per uncovered Core Feature.
- [ ] A Coverage Map entry written as a range covers every member of that
      range and produces no finding.
- [ ] A fixture Spec with a Task Graph whose Tasks reference no Core Feature 4
      reports `SC-COVERAGE-UNTASKED` for it.
- [ ] A Task Context path that does not resolve reports `SC-REF-UNRESOLVED`
      naming the declaring line and the missing path.
- [ ] A fixture Spec with no Task Graph reports no `SC-COVERAGE-UNTASKED` and
      lists it as skipped.

## Context

- interface: `internal/spec/spec.go`
- interface: `internal/spec/task.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/speccheck -count=1 -run 'ADR|Coverage|Reference|Closure' -v | grep -q -- "--- PASS"`
  — expected: exit 0; the ADR, coverage, reference, and closure tests ran and
  passed.
- `go test ./internal/speccheck -count=1` — expected: exit 0.
- `go test ./internal/spec -count=1` — expected: exit 0.

## References

- `_prd.md` → Core Features 2, 4, and 5; Goals.
- `_techspec.md` → Data Models; API Contracts; Build Order 3.
- ADR-0093, ADR-0094.
