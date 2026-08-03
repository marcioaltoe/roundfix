---
task: task_01
spec: 0064-spec-artifact-consistency-gate
status: pending
type: backend
complexity: high
---

# Task 01: Report a Constraint contradiction with both sides

## Overview

Deliver the Spec Consistency Check's report model together with its first
detector family, so the first settled Task already reports a real
contradiction rather than shipping an empty scaffold. The slice reads one Spec
folder's `Project Constraints` sections and the authorization records they
cite, and returns findings that each name a file and line on both sides. It is
verifiable on its own through fixture Spec folders; no command surface is
needed yet.

This family leads the Spec because it is the defect that cost four gate
executions closing Spec 0072 — a PRD citing a tooling authorization it was not
listed in.

## Requirements

1. MUST expose a read-only entry point that takes a Spec Root, a repository
   root, and a slug, and returns a result carrying findings and the detectors
   it skipped. It MUST open no network connection and write no file.
2. MUST model a finding as a stable diagnostic code, a severity of `error` or
   `gap`, a one-line summary naming both sides, the locations, and a concrete
   fix.
3. MUST give every `error` finding at least two locations — one per side of
   the contradiction — each with a repository-relative path and a 1-based line.
4. MUST implement `SC-CONSTRAINT-MISSING`, `SC-CONSTRAINT-UNREASONED`,
   `SC-CONSTRAINT-SOURCE`, `SC-TOOLING-UNAUTHORIZED`, and
   `SC-TOOLING-UNBOUNDED` over the PRD and every present TechSpec, as specified
   in the TechSpec's API Contracts table.
5. MUST report `SC-TOOLING-UNAUTHORIZED` when an applicable Tooling authority
   row cites an authorization record whose text does not name the Spec's slug
   or number, locating both the citing row and the record.
6. MUST skip a detector, and record the skip with the artifact that was
   missing, when an input artifact is absent — never fail for absence.
7. MUST render a result as text and as one `roundfix-speccheck/v1` JSON object.
8. SHOULD keep the package free of any dependency on the Daemon's Task Graph
   load path beyond reading it.

## Subtasks

- [ ] Add the package with the finding, location, severity, skip, and result
      models plus both renderers.
- [ ] Parse a `Project Constraints` section into rows carrying label,
      applicability, reason, cited source path, and line.
- [ ] Read the authorization records and answer whether a record names a given
      Spec.
- [ ] Implement the five constraint and tooling detectors.
- [ ] Add fixture Spec folders — one clean, one dirty per detector, one absent
      per skippable input — and table-driven tests over them.

## Acceptance Criteria

- [ ] A fixture Spec whose PRD omits a required Constraint row reports
      `SC-CONSTRAINT-MISSING` at severity `error`.
- [ ] A fixture Spec whose row states applicability with no reason reports
      `SC-CONSTRAINT-UNREASONED`.
- [ ] A fixture Spec whose row cites a source path that does not exist reports
      `SC-CONSTRAINT-SOURCE`, locating the citing line and naming the missing
      path.
- [ ] A fixture Spec whose Tooling authority row cites an authorization record
      that does not name it reports `SC-TOOLING-UNAUTHORIZED` with one location
      in the Spec and one in the record.
- [ ] A fixture Spec with an applicable Tooling authority row and no bounded
      files reports `SC-TOOLING-UNBOUNDED`.
- [ ] A fixture Spec with no TechSpec reports no finding from any TechSpec-only
      detector and lists that detector as skipped.
- [ ] Every `error` finding produced by the fixtures carries at least two
      locations, asserted by a test that walks all of them.
- [ ] A clean fixture Spec produces zero findings.

## Context

- interface: `internal/spec/spec.go`
- instruction: `docs/agents/agent-instructions.md`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/speccheck -count=1 -run 'Constraint|Tooling|Skip|Location' -v | grep -q -- "--- PASS"`
  — expected: exit 0; the constraint, tooling, skip, and location tests ran and
  passed.
- `go test ./internal/speccheck -count=1` — expected: exit 0.
- `go test ./internal/spec -count=1` — expected: exit 0; the Task Graph loader
  is unchanged.
- `if grep -rn "net/http\|os.Create\|os.WriteFile\|os.Remove" internal/speccheck --include="*.go" | grep -v "_test.go" | grep -q .; then exit 1; fi`
  — expected: exit 0; the checker opens no transport and writes no file outside
  its tests.

## References

- `_prd.md` → Core Features 3 and 7; Goals.
- `_techspec.md` → Interfaces; Data Models; API Contracts; Build Order 1 and 2.
- ADR-0093, ADR-0094.
