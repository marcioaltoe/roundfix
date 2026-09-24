---
task: task_02
spec: 0160-a-review-that-reaches-a-verdict
status: pending
type: backend
complexity: low
---

# Task 02: Run the claude provider

## Overview

`claude` is a valid pre-PR review policy that `roundfix review` refuses, although the read-only session it needs is the one Codex already uses.

## Requirements

1. MUST run a review with `pre_pr_review.provider: claude` through the `review` profile on a read-only session, with the same evidence, fallback and blocking rules as `codex`.
2. MUST keep refusing `coderabbit`, with a reason naming the missing local review surface.
3. MUST keep refusing a `review` profile whose runtime does not match the provider.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] A `claude` policy yields a record with provider `claude`.
- [ ] `coderabbit` is refused naming the missing surface.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/cli/review.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestReviewRunsTheClaudeProvider|TestReviewRefusesCodeRabbitNamingTheMissingSurface)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestReviewRunsTheClaudeProvider TestReviewRefusesCodeRabbitNamingTheMissingSurface; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task the named cases do not exist, so the command fails.

## References

- [_techspec.md](_techspec.md) — The claude provider
