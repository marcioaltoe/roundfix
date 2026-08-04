---
task: task_05
spec: 0070-declared-unreachable-acceptance
status: pending
type: test
complexity: medium
---

# Task 05: Replay Spec 0058 and hold the corpus unchanged

## Overview

This Spec's first Success Metric is a replay: Spec 0058's QA report and
artifacts, with its release row declared, must archive without `qa_override`
and stamp the release action as unproven. That Spec is why this one exists —
its evidence was complete except for something no gate may ever perform, and
the only exit available spent the mechanism reserved for failed evidence.

The second half is the guard: every Spec that archives today still archives
identically. A widened gate is only safe if the widening is measured.

## Requirements

1. MUST replay Spec 0058's archived report and artifacts with its release row
   declared, and assert it archives without `qa_override`.
2. MUST assert the resulting archive record names the release action that
   remains unproven, rather than dropping it.
3. MUST assert a row declared unreachable that the environment could in fact
   reach is reported as wrongly declared rather than accepted.
4. MUST assert a blocked row with no matching declaration still blocks the
   archive.
5. MUST assert non-regression across the archived Spec corpus: every Spec that
   satisfied the archive precondition before this Spec still satisfies it, and
   no archived artifact is modified.
6. MUST make the replay fixture's provenance explicit — which report it
   reproduces, and that the declaration is added by this Spec rather than
   present in the original.

## Subtasks

- [ ] Build the Spec 0058 replay fixture with its provenance recorded.
- [ ] Assert the archive succeeds and the unproven action is stamped.
- [ ] Assert the wrongly-declared and unmatched-row refusals.
- [ ] Assert corpus non-regression.

## Acceptance Criteria

- [ ] The 0058 replay archives without `qa_override` and stamps the release
      action.
- [ ] The replay fixture records the report path it reproduces and states the
      declaration was added by this Spec.
- [ ] A wrongly declared row is reported, not accepted.
- [ ] An unmatched blocked row still blocks the archive.
- [ ] Every archived Spec still satisfies its archive precondition.
- [ ] No file under `docs/specs/_archived/` is modified.

## Context

- instruction: `docs/agents/spec-routing.md`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/spec -count=1 -run 'Replay|Corpus|Unreachable' -v | grep -q -- "--- PASS"`
  — expected: exit 0; the replay and corpus tests ran and passed.
- `go test ./internal/spec ./internal/cli -count=1` — expected: exit 0.
- `go test -parallel 16 ./... 2>&1 | grep -q "^FAIL" && exit 1 || exit 0`
  — expected: exit 0.
- `git diff --name-only HEAD -- docs/specs/_archived | grep -q . && exit 1 || exit 0`
  — expected: exit 0; no archived Spec file changed.

## References

- `_prd.md` → Success Metrics 1, 2, 3 and 4; Decisions.
- `_techspec.md` → Testing Approach; Build Order 5.
- ADR-0080.
