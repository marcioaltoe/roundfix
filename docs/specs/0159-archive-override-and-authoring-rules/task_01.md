---
task: task_01
spec: 0159-archive-override-and-authoring-rules
status: pending
type: backend
complexity: medium
---

# Task 01: The QA Archive Override

## Overview

The Baseline guidance and the archive skill describe archiving a Spec despite failed or missing QA on maintainer authority, stamped `qa_override: true`, but `roundfix archive` accepts no override input.

## Requirements

1. MUST accept `roundfix archive <slug> --qa-override --approval <source> --reason <text>`, and refuse `--approval` or `--reason` without `--qa-override` and `--qa-override` without both.
2. MUST still require every non-QA Task completed under an override.
3. MUST refuse an override when the newest QA Report already qualifies under the declared-acceptance policy.
4. MUST stamp `qa_override: true`, `qa_override_approval`, `qa_override_reason`, `qa_override_qa_outcome` (the observed verdict, `missing`, or the read error) and `qa_override_revision` (the archived `HEAD`), and report the archive as an override.
5. MUST move the QA Task file and every QA Report byte-identically, never changing a status or a verdict.
6. MUST leave archive without the flag exactly as today.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] An authorized override archives a Spec with failed QA and stamps every field, with the QA Task and report unchanged.
- [ ] An override is refused with a non-QA Task pending, with QA already qualifying, and without approval or reason.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/spec/archive.go`
- interface: `internal/cli/archive.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestArchiveQAOverrideStampsProvenance|TestArchiveQAOverrideStillRequiresNonQATasks|TestArchiveQAOverrideRefusedWhenQAQualifies|TestArchiveCommandQAOverrideRequiresApprovalAndReason)$" ./internal/spec ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestArchiveQAOverrideStampsProvenance TestArchiveQAOverrideStillRequiresNonQATasks TestArchiveQAOverrideRefusedWhenQAQualifies TestArchiveCommandQAOverrideRequiresApprovalAndReason; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task none of the four cases exists, so the command fails.

## References

- [_techspec.md](_techspec.md) — The override
