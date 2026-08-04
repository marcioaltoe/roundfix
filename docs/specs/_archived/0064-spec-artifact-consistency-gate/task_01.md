---
task: task_01
spec: 0064-spec-artifact-consistency-gate
status: completed
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

## Result

### Implementation

- Added the read-only `speccheck.Check` entry point with the stable finding,
  severity, location, skipped-detector, and result models from the TechSpec.
- Added heading-anchored, multiline parsing for all four required Project
  Constraint rows, including applicability, reason, operative source,
  authorization-record citation, bounded-file declaration, and 1-based line.
- Added the five Task 01 detectors. Every reported error carries the declaring
  artifact location and the cited, counterpart, missing, or governing-source
  location that forms its other side.
- Added compact `roundfix-speccheck/v1` JSON rendering and text rendering with
  every code, severity, summary, location, fix, and recorded skip.
- Added one clean fixture, one isolated dirty fixture per detector, and an
  absent-TechSpec fixture. The clean tooling path exercises a real cited
  authorization record that names the fixture Spec and records bounded files.

### Focused checks

- Red signal: `rtk env GOCACHE=/Users/marcio/.roundfix/worktrees/roundfix-339f8dac/run_20260803T233822Z_3ffcad0ced4ba246/.gocache go test -buildvcs=false ./internal/speccheck -run '^TestCheckConstraintAndTooling$' -count=1`
  exited 1 before implementation because `internal/speccheck` had no non-test
  Go files.
- The same filtered command exited 0 after implementation.
- `rtk env GOCACHE=/Users/marcio/.roundfix/worktrees/roundfix-339f8dac/run_20260803T233822Z_3ffcad0ced4ba246/.gocache go test -buildvcs=false -shuffle=on ./internal/speccheck -count=1`
  exited 0 after the final Go edit; shuffled execution covered every Task 01
  fixture and both renderers.
- `rtk env GOCACHE=/Users/marcio/.roundfix/worktrees/roundfix-339f8dac/run_20260803T233822Z_3ffcad0ced4ba246/.gocache go vet -buildvcs=false ./internal/speccheck`
  exited 0.
- The Task's declared `## Verification` commands were not run; they remain
  Daemon-owned.

### Acceptance criteria evidence

1. `TestCheckConstraintAndTooling/missing_required_PRD_row` asserts the
   `constraint-missing` fixture reports `SC-CONSTRAINT-MISSING` at `error`.
2. `TestCheckConstraintAndTooling/applicability_without_reason` asserts the
   `constraint-unreasoned` fixture reports `SC-CONSTRAINT-UNREASONED`.
3. `TestCheckConstraintSourceNamesMissingPath` asserts
   `SC-CONSTRAINT-SOURCE` names and locates
   `docs/agents/missing-guide.md` beside the citing PRD row.
4. `TestCheckToolingUnauthorizedLocatesSpecAndRecord` asserts
   `SC-TOOLING-UNAUTHORIZED` locates both the PRD row and
   `docs/workflow/authorizations/other-spec.md`.
5. `TestCheckConstraintAndTooling/applicable_tooling_has_no_bounded_files`
   asserts `SC-TOOLING-UNBOUNDED` for the applicable unbounded row.
6. `TestCheckSkipMissingTechSpec` asserts the absent-TechSpec fixture has zero
   findings and records `_techspec.md` as missing for all five detectors.
7. `TestCheckErrorLocations` walks every error from every dirty fixture and
   requires at least two repository-relative paths with 1-based lines.
8. `TestCheckCleanFixture` asserts the clean fixture produces zero findings;
   `TestRenderResultTextAndJSON` separately validates both public report
   formats and the exact `roundfix-speccheck/v1` schema identifier.

No follow-up outside Task 01's slice was discovered.
