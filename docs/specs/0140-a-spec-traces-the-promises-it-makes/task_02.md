---
task: task_02
spec: 0140-a-spec-traces-the-promises-it-makes
status: pending
type: backend
complexity: medium
---

# Task 02: Trace a declared promise to its Task

## Overview

A declared promise now carries the obligations a user story already carries: a
Success Metric appears in the TechSpec Coverage Map, and both kinds are named in
some Task's References. This slice makes the two existing coverage detectors
read the new unit kinds, and makes a contract finding point at the TechSpec that
declared it rather than at the PRD.

## Requirements

1. MUST require each declared Success Metric in the TechSpec Coverage Map, and
   report a missing one through the existing unmapped code with its existing
   severity.
2. MUST require each declared Success Metric and each declared API Contract to
   be named in some Task's References, and report a missing one through the
   existing untasked code with its existing severity.
3. MUST recognize a reference written as the words a reader writes followed by a
   number sequence, so that `Success Metrics 1-2` and `API Contract 3` both
   resolve.
4. MUST report a contract finding against the artifact that declared it, so that
   a reader is sent to the TechSpec and not to the PRD.
5. MUST keep the Coverage Map requirement skipped for a Spec with no TechSpec,
   while its declared metrics still require a Task reference.
6. MUST NOT require a declared API Contract in the Coverage Map.

## Subtasks

- [ ] Feed the two new unit kinds into the Coverage Map and References detectors.
- [ ] Carry each unit's declaring artifact through to its finding location.
- [ ] Extend the reference vocabulary to both new phrases, naming the two
      patterns `metricRefPattern` and `contractRefPattern` beside the existing
      story and feature patterns.
- [ ] Cover mapped, unmapped, named and unnamed promises with tests, including a
      Spec with no TechSpec.

## Acceptance Criteria

- [ ] A fixture declaring two metrics and mapping one reports the unmapped code
      for the second, at its existing severity.
- [ ] A fixture declaring a metric and a contract that no Task names reports the
      untasked code for both, and the contract finding names the TechSpec.
- [ ] A fixture whose Task References name `Success Metrics 1-2` and
      `API Contract 1` reports neither code.
- [ ] A fixture with no TechSpec reports no Coverage Map finding for its metric
      and still reports the untasked code when no Task names it.
- [ ] A declared API Contract absent from the Coverage Map reports nothing.

## Context

- interface: `internal/speccheck/citations.go`

## Verification

- `grep -q "metricRefPattern" internal/speccheck/citations.go` — expected: exit 0; a Task reference to a Success Metric resolves through its own pattern. Task 01 introduces the unit kinds but no reference pattern, so this fails before this Task.
- `grep -q "contractRefPattern" internal/speccheck/citations.go` — expected: exit 0; a Task reference to an API Contract resolves through its own pattern.
- `grep -rq "func TestPromiseCoverageReachesItsTask" internal/speccheck` — expected: exit 0; coverage of both kinds is covered by a named test.
- `out="$(go test -count=1 -run "^TestPromiseCoverageReachesItsTask$" ./internal/speccheck 2>&1)" || { printf "%s\n" "$out"; exit 1; }; missing="$(printf "%s\n" "$out" | grep "no tests to run")"; test -z "$missing"` — expected: exit 0; before this Task the run reports no tests to run, so the command fails.

## References

`_prd.md` → Core Feature 2; User Stories 2, 4; Success Metric 3;
`_techspec.md` → Implementation Design: Codes and their stages, Interfaces;
API Contract 3; Coverage Map; Build Order 2; ADR-0093; ADR-0156.
