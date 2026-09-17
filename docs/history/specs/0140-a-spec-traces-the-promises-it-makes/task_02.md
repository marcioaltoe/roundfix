---
task: task_02
spec: 0140-a-spec-traces-the-promises-it-makes
status: completed
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

## Result

Implemented the promise coverage extension in the shared citation detectors.
Success Metrics now join User Stories and Core Features in Coverage Map
checking, Success Metrics and API Contracts join Task References checking, and
each coverage unit carries the artifact that declared it. The reference reader
recognizes `Success Metric(s)` and `API Contract(s)` followed by the existing
number-sequence grammar. API Contracts are explicitly excluded from Coverage
Map enforcement, and the authoring-stage Coverage Map check applies the same
Success Metric rule.

Focused-check evidence:

- Before the implementation, `GOCACHE=/private/tmp/roundfix-go-cache rtk go
  test -count=1 ./internal/speccheck` reported 379 passing and four failing
  tests. The failures were the new promise coverage cases for the missing
  metric mapping, unnamed promises, and the no-TechSpec metric path.
- After the implementation, `GOCACHE=/private/tmp/roundfix-go-cache rtk go
  test -count=1 ./internal/speccheck` exited 0 with 383 passing tests.
- `rtk git diff --check` exited 0.

Acceptance evidence:

- `TestPromiseCoverageReachesItsTask/unmapped_metric_keeps_the_coverage_error`
  declares two metrics, maps one, and asserts one `SC-COVERAGE-UNMAPPED` error
  naming Success Metric 2.
- `TestPromiseCoverageReachesItsTask/unnamed_promises_retain_their_declaring_artifacts`
  asserts two `SC-COVERAGE-UNTASKED` errors at the existing severity, with the
  metric located in the PRD and the contract summary and location naming the
  TechSpec.
- `TestPromiseCoverageReachesItsTask/written_metric_range_and_contract_reference_resolve`
  uses `Success Metrics 1-2` and `API Contract 1` and asserts that neither
  coverage code is reported. Its Coverage Map omits the API Contract, which
  also proves that contracts are not required there.
- `TestPromiseCoverageReachesItsTask/metric_without_TechSpec_skips_mapping_but_still_needs_a_Task`
  asserts no Coverage Map finding and one untasked Success Metric finding when
  the TechSpec is absent.

The Daemon-owned Verification commands were not run in this turn.

## Carry-forward provenance

- Source Run: `run_20260917T164016Z_ac67737203a1aa3d`
- Source commit: `474609e696eaa641f67a700f8b45b876fb3f8a8e`
