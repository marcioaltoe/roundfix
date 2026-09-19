---
task: task_02
spec: 0148-a-profile-that-declares-both-tiers
status: pending
type: backend
complexity: medium
---

# Task 02: Publish both tiers in the generated guidance

## Overview

An Agent reads the generated instructions, not the Profile. This slice renders
the incremental command where the template already renders the complete gate,
and keeps saying the contract is unmet when no incremental command is declared.

## Requirements

1. MUST render the declared incremental command in the guide template, beside
   the complete gate it renders today.
2. MUST render, for a Profile that declares no incremental command, the unmet
   two-tier statement the clauses already use, rather than an empty value.
3. MUST leave every other rendered value and its wording unchanged.
4. MUST NOT hand-edit any generated guide; rendering changes come from the
   template.

## Subtasks

- [ ] Render the incremental command in the template.
- [ ] Render the unmet statement when it is absent.
- [ ] Cover both cases at the generation seam.

## Acceptance Criteria

- [ ] A Profile declaring both publishes both.
- [ ] A Profile declaring none publishes the unmet statement.
- [ ] Every other rendered value is unchanged.

## Context

- interface: `internal/baseline/assets/templates/guides/agent-instructions.md`

## Verification

- `grep -q "verification.incremental" internal/baseline/assets/templates/guides/agent-instructions.md` — expected: exit 0; the template publishes the tier. Before this Task it does not.
- `out="$(go test -count=1 -run "^TestGeneratedGuidancePublishesBothTiers$" ./internal/baseline 2>&1)" || { printf "%s\n" "$out"; exit 1; }; missing="$(printf "%s\n" "$out" | grep "no tests to run")"; test -z "$missing"` — expected: exit 0; before this Task the run reports no tests to run, so the command fails.

## References

`_prd.md` → Core Features 3-4; User Stories 1 and 4; Goals 3-4;
Success Metrics 1-2;
`_techspec.md` → Implementation Design: Rendering; API Contracts 1-2;
Build Order 2; `_authorization.md`.
