---
task: task_14
spec: 0057-baseline-capability-evidence-and-retention
status: pending
type: backend
complexity: medium
---

# Task 14: Show the probe in the output a human reads

## Overview

The evaluated probe reaches machine output but not text. A maintainer whose
executable candidate was rejected sees only that local evidence is
insufficient, with a recommendation to install a tool that is already
installed — while the JSON beside it names the exact candidate and the reason
it failed. The Spec's promise is that remediation never requires reading
another format, and text is the format a human reads.

## Requirements

1. MUST render, in text output, the inspected candidate and the rejection
   reason for every unsatisfied executable capability, matching what machine
   output already carries.
2. MUST render, in text output, the inspected paths and their states for every
   unsatisfied declared-file capability.
3. MUST state the absence of any candidate distinctly from a candidate that was
   inspected and rejected.
4. MUST NOT recommend installing a tool whose candidate was found and rejected;
   the next action must address the reason, not the absence.
5. MUST keep text and machine output describing the same evidence, so the two
   cannot disagree about what was inspected.
6. MUST leave machine output unchanged.

## Subtasks

- [ ] Render the inspected candidate and reason in text.
- [ ] Render inspected paths and states for declared-file probes in text.
- [ ] Distinguish absence from inspected-and-rejected.
- [ ] Make the next action address the actual reason.

## Acceptance Criteria

- [ ] A rejected executable candidate is named in text output, with its
      rejection reason.
- [ ] A broken symlink candidate reports the broken-link reason in text, not a
      recommendation to install the tool.
- [ ] A capability with no candidate at all reports absence, distinctly.
- [ ] An unsatisfied declared-file capability lists its inspected paths and
      their states in text.
- [ ] Text and machine output name the same inspected subject for the same
      capability, asserted by comparing both for one repository.
- [ ] Machine output is byte-identical to before this Task for the same input.
- [ ] `git status --porcelain` shows no path outside `internal/baseline/`,
      `internal/cli/`, and this task file.

## Context

- interface: `internal/baseline/profile_alignment.go`
- interface: `internal/cli/cli.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/baseline ./internal/cli -run '^TestCapabilityTextRendersProbe$' -count=1 -v | grep -q -- "--- PASS: TestCapabilityTextRendersProbe"`
  — expected: exit 0; text names the inspected candidate and reason.
- `go test ./internal/baseline -run '^TestCapabilityTextAndJSONAgree$' -count=1 -v | grep -q -- "--- PASS: TestCapabilityTextAndJSONAgree"`
  — expected: exit 0; both formats name the same inspected subject.
- `go test ./internal/baseline -run '^TestBaselinePlanCharacterization$' -count=1 -v | grep -q -- "--- PASS: TestBaselinePlanCharacterization"`
  — expected: exit 0.
- `go test ./internal/baseline ./internal/cli -count=1` — expected: exit 0.

## References

- `_prd.md` → User Story 2; Core Features 5; User Experience (a blocking
  divergence reads as probe evidence, not prose).
- `qa/qa-report-2026-08-02-01.md` → F-001.
