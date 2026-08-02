---
task: task_04
spec: 0057-baseline-capability-evidence-and-retention
status: pending
type: backend
complexity: high
---

# Task 04: Render the probe and group by requirement strength

## Overview

A maintainer facing a blocking divergence today gets prose and has to open
catalog assets or Go source to learn what was actually checked. The evidence
now travels on the divergence; this Task renders it. It also groups divergences
by requirement strength, so an advisory stops reading like a blocker.

## Requirements

1. MUST render, for every unsatisfied divergence, the probe that was evaluated:
   a declared-file probe lists each inspected path with its state and the
   expected content; an executable probe names the inspected candidate or
   states that none existed.
2. MUST state, for a required stack capability, the selected technology, both
   resolutions, and whether removing it cascades to any decision.
3. MUST group divergences by requirement strength — blocking, advisory,
   informational — in both rendered and machine output.
4. MUST state on every advisory that it does not block readiness or apply,
   before any optional next action it suggests.
5. MUST render from the evidence carried on the divergence, never by
   re-deriving it, so the rendering and the verdict share one source.
6. MUST leave verdicts, blocking status, and readiness unchanged; this Task
   changes presentation only.

## Subtasks

- [ ] Render declared-file probes as inspected paths with states.
- [ ] Render executable probes as the inspected candidate or its absence.
- [ ] Render stack-capability resolutions and decision cascade.
- [ ] Group by requirement strength in both outputs.
- [ ] Put the non-blocking statement ahead of every advisory next action.

## Acceptance Criteria

- [ ] A blocking declared-file divergence renders every inspected path with its
      state and the expected content.
- [ ] A blocking executable divergence names the inspected candidate, or states
      that no candidate existed.
- [ ] A required stack capability renders both resolutions and whether removal
      cascades to a decision.
- [ ] Divergences appear grouped as blocking, advisory, informational in
      rendered output, and carry the same grouping in machine output.
- [ ] Every advisory states it does not block readiness or apply before any
      next action.
- [ ] No verdict, blocking flag, or readiness value changed, proven by the
      characterization corpus.
- [ ] `git status --porcelain` shows no path outside `internal/baseline/` and
      this task file.

## Context

- interface: `internal/baseline/profile_alignment.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/baseline -run '^TestDivergenceRendersProbe$' -count=1 -v | grep -q -- "--- PASS: TestDivergenceRendersProbe"` —
  expected: exit 0; each probe kind renders its inspected subject.
- `go test ./internal/baseline -run '^TestDivergenceGroupsByRequirement$' -count=1 -v | grep -q -- "--- PASS: TestDivergenceGroupsByRequirement"` —
  expected: exit 0; grouping and the advisory statement hold in both outputs.
- `go test ./internal/baseline -run '^TestBaselinePlanCharacterization$' -count=1 -v | grep -q -- "--- PASS: TestBaselinePlanCharacterization"` —
  expected: exit 0 with goldens re-recorded only for the intended rendering
  change; no verdict or readiness value differs.
- `go test ./internal/baseline -count=1` — expected: exit 0.

## References

- `_prd.md` → User Story 2; Core Features 5 and 8; User Experience.
- `_techspec.md` → API Contracts; Build Order 4.

