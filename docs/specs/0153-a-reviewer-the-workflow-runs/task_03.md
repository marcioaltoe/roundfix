---
task: task_03
spec: 0153-a-reviewer-the-workflow-runs
status: pending
type: backend
complexity: high
---

# Task 03: The command, its exits and everything that blocks

## Overview

`roundfix review [--base <ref>]` resolves the policy, runs the reviewer over the
diff Task 02 hands it, and reports through the exit status.

The blocking list is the substance. A first attempt treated a clean-looking
string as clean regardless of how the session ended, and review found two ways
that lies: a timeout followed by a fallback returning clean, and an adapter that
exited non-zero while its parsed output said `No findings`.

## Requirements

1. MUST exit 0 only on an explicit clean answer, and when the policy is `none`
   and a configured omission was recorded.
2. MUST exit 1 when the reviewer returned findings, with the findings recorded.
3. MUST exit 2 and record blocked for each of: a runtime failure, a timeout, a
   non-empty transport anomaly, empty agent output, output it cannot classify,
   and a provider this slice does not execute.
4. MUST classify the agent's message rather than the raw protocol stream.
5. MUST activate a configured fallback only for a selection that failed to start
   before the prompt was sent, and never after a failure of the review itself.
6. MUST perform no reviewer call and no readiness probe when the policy is
   `none`.
7. MUST refuse `claude` and `coderabbit` by naming the provider, and MUST NOT
   record them as omitted.
8. MUST prove the artifact directory writeable before spending the reviewer
   call.
9. MUST refuse unknown flags, matching the surrounding commands.

## Subtasks

- [ ] Add the command, its preflight and its policy resolution.
- [ ] Map every terminal signal to its outcome.
- [ ] Register the command on the public surface.
- [ ] Add a test per exit and per blocking signal.

## Acceptance Criteria

- [ ] Each of the six blocking signals exits 2, records blocked, and names its
      reason.
- [ ] A transport anomaly blocks even when the agent output reads clean.
- [ ] A timeout activates no fallback.
- [ ] `none` exits 0 with no `Run` and no `Probe` on the stub.
- [ ] `claude` and `coderabbit` exit 2 naming the provider.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/cli/doctor.go`
- interface: `internal/agent/acpx_runner.go`

## Verification

- `out="$(go test -count=1 -v -run "^TestReviewCommand" ./internal/cli 2>&1)" || { printf "%s\n" "$out"; exit 1; }; printf "%s\n" "$out" | grep -q -- "--- PASS: TestReviewCommandBlocksOnTransportAnomaly"` — expected: exit 0; before this Task the case does not exist, so the command fails.
- `out="$(go test -count=1 -v -run "^TestReviewCommand" ./internal/cli 2>&1)" || { printf "%s\n" "$out"; exit 1; }; printf "%s\n" "$out" | grep -q -- "--- PASS: TestReviewCommandKeepsATimeoutBlocked"` — expected: exit 0; before this Task the case does not exist, so the command fails.
- `go run -buildvcs=false ./cmd/roundfix --help 2>&1 | grep -q "roundfix review"` — expected: exit 0; the command appears on the public surface. Before this Task the usage has no review line, so the command fails.

## References

- [_techspec.md](_techspec.md) — What blocks
