---
task: task_02
spec: 0057-baseline-capability-evidence-and-retention
status: pending
type: backend
complexity: medium
---

# Task 02: Resolve symlinked executables without running them

## Overview

Executable capability discovery inspects each PATH candidate with a link-level
stat and then requires a regular file, so a symlink never qualifies. Every tool
installed through Homebrew or Docker Desktop — the norm on maintainer machines —
reports missing while being present and working. This Task resolves a bounded
link chain and judges the target, still without executing anything, and reports
the candidate it inspected instead of an empty result.

## Requirements

1. MUST resolve a symlinked PATH candidate through a bounded chain of links to
   a regular executable, and treat that as discovered evidence.
2. MUST NOT execute any candidate or target; discovery stays offline and
   side-effect free per ADR-0087.
3. MUST bound the chain so a cycle terminates, and report a cycle, a broken
   link, and a non-executable target as three distinct reasons.
4. MUST report the inspected candidate whenever one existed, and report the
   absence of any candidate distinctly from a candidate that failed.
5. MUST keep a directly-present regular executable resolving exactly as it does
   today.
6. MUST leave every other capability evidence kind unchanged.

## Subtasks

- [ ] Add bounded link resolution to the candidate probe.
- [ ] Distinguish cycle, broken link, and non-executable target.
- [ ] Report the inspected candidate and the absence case.
- [ ] Confirm no candidate is executed.

## Acceptance Criteria

- [ ] A one-hop symlink to a regular executable resolves as discovered.
- [ ] A multi-hop chain within the bound resolves as discovered.
- [ ] A link cycle terminates and reports the cycle reason.
- [ ] A broken link reports the broken-link reason, not absence.
- [ ] A link to a non-executable target reports the non-executable reason.
- [ ] A directly-present regular executable still resolves as discovered.
- [ ] No candidate is executed, proven by a probe target that would record its
      own invocation.
- [ ] The characterization corpus from task 01 shows only the intended
      capability status changes.
- [ ] `git status --porcelain` shows no path outside `internal/baseline/` and
      this task file.

## Context

- interface: `internal/baseline/profile_alignment.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/baseline -run TestExecutableCandidateResolution -count=1` —
  expected: exit 0; all seven probe cases assert their own reason.
- `go test ./internal/baseline -run TestExecutableCandidateNeverExecutes -count=1`
  — expected: exit 0; a target that would record its invocation records nothing.
- `go test ./internal/baseline -run TestBaselinePlanCharacterization -count=1` —
  expected: exit 0.
- `go test ./internal/baseline -count=1` — expected: exit 0.

## References

- `_prd.md` → User Story 3; Core Features 4.
- `_techspec.md` → Implementation Design: Interfaces; Build Order 2.
- ADR-0087.
