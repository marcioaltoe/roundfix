---
task: task_01
spec: 0143-a-repository-says-who-reviews-before-the-pull-request
status: pending
type: backend
complexity: medium
---

# Task 01: Carry the pre-Pull-Request reviewer in configuration

## Overview

Configuration gains the reviewer that runs before a Pull Request exists, with
the precedence every other key already uses and a refusal for anything outside
the supported set. Nothing reads the policy yet; this slice makes it declarable
and resolvable.

## Requirements

1. MUST accept `pre_pr_review.provider` with the values `codex`, `claude`,
   `coderabbit` and `none`.
2. MUST resolve Project Config over User Config over a built-in default of
   `codex`, and MUST treat an absent key as inherit, never as `none`.
3. MUST report, beside the resolved value, which layer supplied it: project,
   user or default.
4. MUST refuse an unsupported value while configuration loads, with an error
   naming the key, the offending value and the four supported values.
5. MUST leave `review_source` and every other configuration key with their
   current keys, defaults, validation and behavior.
6. MUST name the resolved type `PrePRReview`, so the Task's own check and later
   readers agree.

## Subtasks

- [ ] Add the overlay and the typed section.
- [ ] Layer it with the existing precedence and record the answering layer.
- [ ] Validate the value at load with the named error.
- [ ] Cover each supported value, each layer, both-silent and an unsupported
      value at the configuration unit seam.

## Acceptance Criteria

- [ ] Each supported value loads and resolves.
- [ ] Project beats user beats default, and the answering layer is reported.
- [ ] A configuration with no key resolves to `codex` from the default layer.
- [ ] An unsupported value fails the load with the named error.
- [ ] Existing configuration tests pass unedited.

## Context

- interface: `internal/config/config.go`

## Verification

- `grep -q "PrePRReview" internal/config/config.go` — expected: exit 0; configuration carries the policy. Before this Task it does not.
- `out="$(go test -count=1 -run "^TestPrePRReviewPolicyResolution$" ./internal/config 2>&1)" || { printf "%s\n" "$out"; exit 1; }; missing="$(printf "%s\n" "$out" | grep "no tests to run")"; test -z "$missing"` — expected: exit 0; before this Task the run reports no tests to run, so the command fails.

## References

`_prd.md` → Core Features 1-2 and 4; User Stories 1-3 and 5; Goals 1 and 3;
Success Metrics 1-3;
`_techspec.md` → Implementation Design: The key and its values, Resolution and
its source, Refusal at load; API Contract 1; Build Order 1; ADR-0002.
