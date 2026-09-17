---
task: task_01
spec: 0140-a-spec-traces-the-promises-it-makes
status: pending
type: backend
complexity: medium
---

# Task 01: Read a declared promise and refuse an undeclared one

## Overview

The Spec Consistency Check gains the reader that classifies a promise section,
and the two coined codes that report a section carrying no declaration. This
slice is complete on its own: a PRD with no Success Metrics section, and a
TechSpec whose API Contracts are prose, each produce their finding at the stage
that writes them, while every existing detector keeps its behavior.

## Requirements

1. MUST classify the PRD's Success Metrics section and the TechSpec's API
   Contracts section, read at heading level two or three, as one of three
   states: declared units, explicit none, or no declaration, exactly as the
   TechSpec's "Declaring a promise" defines them.
2. MUST treat a section that mixes `None.` with numbered items as declared
   units, never as an explicit none.
3. MUST report `SC-METRIC-UNDECLARED` at the PRD stage and
   `SC-CONTRACT-UNDECLARED` at the TechSpec stage, each with severity `gap`,
   each naming the artifact it read and the section it could not settle.
4. MUST leave every existing detector's code, stage, severity and message
   unchanged.
5. MUST NOT emit either code when the section declares numbered items or an
   explicit none with its reason.

## Subtasks

- [ ] Add the two coverage unit kinds to the checker's unit vocabulary.
- [ ] Parse each promise section into its three states.
- [ ] Coin the two codes and place each at its stage.
- [ ] Cover the states with table-driven tests at the checker's fixture seam,
      including the mixed `None.` case and both silent cases.

## Acceptance Criteria

- [ ] A fixture PRD with no Success Metrics section reports
      `SC-METRIC-UNDECLARED` as a `gap` at the PRD stage.
- [ ] A fixture TechSpec whose API Contracts section is prose reports
      `SC-CONTRACT-UNDECLARED` as a `gap` at the TechSpec stage.
- [ ] A fixture declaring `None.` with a reason reports neither code.
- [ ] A fixture mixing `None.` with numbered items is read as declared units and
      reports neither code.
- [ ] The stage table still binds every pre-existing code to the stage it had.

## Context

- interface: `internal/speccheck/citations.go`
- interface: `internal/speccheck/coherence.go`

## Verification

- `grep -q "SC-METRIC-UNDECLARED" internal/speccheck/citations.go` — expected: exit 0; the code does not exist before this Task.
- `grep -q "SC-CONTRACT-UNDECLARED" internal/speccheck/citations.go` — expected: exit 0.
- `grep -q "CodeMetricUndeclared" internal/speccheck/coherence.go` — expected: exit 0; the stage table names the PRD-stage code.
- `grep -q "CodeContractUndeclared" internal/speccheck/coherence.go` — expected: exit 0; the stage table names the TechSpec-stage code.
- `grep -rq "func TestPromiseSectionDeclaration" internal/speccheck` — expected: exit 0; the declaration states are covered by a named test.
- `out="$(go test -count=1 -run "^TestPromiseSectionDeclaration$" ./internal/speccheck 2>&1)" || { printf "%s\n" "$out"; exit 1; }; missing="$(printf "%s\n" "$out" | grep "no tests to run")"; test -z "$missing"` — expected: exit 0; before this Task the run reports no tests to run, so the command fails.

## References

`_prd.md` → Core Features 1, 3, 4, 5; User Stories 1, 3; Success Metric 3;
`_techspec.md` → Implementation Design: Declaring a promise, Severity and why it
splits, Codes and their stages; API Contracts 1-2; Build Order 1;
ADR-0093; ADR-0117; ADR-0156.
