---
task: task_06
spec: 0153-a-reviewer-the-workflow-runs
status: pending
type: backend
complexity: medium
---

# Task 06: A record that cannot claim the wrong reviewer

## Overview

Pre-PR review found two defects in the reauthored command. Both were confirmed
in source, and both are about the record telling the truth.

**The policy and the profile can disagree.** The selection loop passes
`profiles.review.preferred` and its fallbacks to `runtimeForProfileSelection`
without ever comparing them to `policy.Provider`. That value appears exactly
once, writing `codex` into the record. So a `review` profile naming `claude` or
`opencode` would execute that runtime while the persisted evidence claims the
configured provider ran.

**Overlapping invocations share a session.** `reviewSessionRef(headCommit,
gitRoot, index)` is deterministic, so two reviews of the same head reuse one
ACPX named session: their prompts can share conversation state, and the first
to finish can close the session the second is still using.

## Requirements

1. MUST refuse, as a named configuration error, a `review` profile selection
   whose runtime is not the provider the policy selected.
2. MUST make that refusal before the readiness pass, so an unselected provider
   is never probed.
3. MUST give the session name a component unique to the invocation, so two
   concurrent reviews of the same head cannot share a session.
4. MUST keep every behavior Tasks 01 through 04 delivered, including the diff
   in the prompt and each blocking signal.

## Subtasks

- [ ] Compare each selection's runtime to the selected provider and refuse a
      mismatch by name.
- [ ] Add an invocation-unique component to the session reference.
- [ ] Add a test for each.

## Acceptance Criteria

- [ ] A `review` profile whose preferred runtime is not the selected provider
      refuses with both named, and no readiness probe runs.
- [ ] A mismatching fallback refuses the same way.
- [ ] Two session references built for the same head in one process differ.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/cli/review.go`

## Verification

- `out="$(go test -count=1 -v -run "^TestReviewCommandRefusesProviderProfileMismatch" ./internal/cli 2>&1)" || { printf "%s\n" "$out"; exit 1; }; printf "%s\n" "$out" | grep -q -- "--- PASS: TestReviewCommandRefusesProviderProfileMismatch"` — expected: exit 0; before this Task the case does not exist, so the command fails.
- `out="$(go test -count=1 -v -run "^TestReviewSessionRefIsUniquePerInvocation" ./internal/cli 2>&1)" || { printf "%s\n" "$out"; exit 1; }; printf "%s\n" "$out" | grep -q -- "--- PASS: TestReviewSessionRefIsUniquePerInvocation"` — expected: exit 0; before this Task the case does not exist, so the command fails.

## References

- [_techspec.md](_techspec.md) — What blocks
