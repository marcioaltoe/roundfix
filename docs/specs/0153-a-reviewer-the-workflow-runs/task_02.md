---
task: task_02
spec: 0153-a-reviewer-the-workflow-runs
status: pending
type: backend
complexity: high
---

# Task 02: The command, the Codex path and every refusal

## Overview

The operation itself. `roundfix review [--base <ref>]` resolves the policy, runs
a read-only Codex reviewer session over the candidate through `agent.Runner`,
and reports through the exit status.

The three exits are the point. A reviewer that found problems is not the same
event as a reviewer that could not run, and collapsing them makes a broken
reviewer look like a strict one.

## Requirements

1. MUST exit 0 when the reviewer ran and returned no findings, and when the
   policy is `none` and a configured omission was recorded.
2. MUST exit 1 when the reviewer ran and returned findings, with the findings in
   the record.
3. MUST exit 2 when Preflight Validation fails, and when the selected mode is
   blocked by a runtime failure, a timeout, or output it cannot read.
4. MUST perform no reviewer call and no readiness probe when the policy is
   `none`.
5. MUST refuse `claude` and `coderabbit` by naming the provider, and MUST NOT
   record them as omitted.
6. MUST NOT fall back to another provider or to `none` on any failure.
7. MUST run the reviewer through `agent.Runner` rather than starting a process
   of its own.
8. MUST refuse unknown flags, matching the surrounding commands.

## Subtasks

- [ ] Add the command and its preflight.
- [ ] Run the Codex reviewer session through the runner and map its outcome.
- [ ] Register the command on the public surface.
- [ ] Add tests for each exit, each refusal and the no-call case.

## Acceptance Criteria

- [ ] With a stubbed runtime the Codex path writes a record and exits 0 or 1
      according to findings.
- [ ] `none` exits 0, writes an omitted record, and the stub records no `Run`
      and no `Probe`.
- [ ] A runtime failure, a timeout and unreadable output each exit 2 with the
      reason named, and each records blocked.
- [ ] `claude` and `coderabbit` exit 2 naming the provider and write no omitted
      record.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/agent/agent.go`
- interface: `internal/cli/doctor.go`

## Verification

- `out="$(go test -count=1 -v -run "^TestReviewCommand" ./internal/cli 2>&1)" || { printf "%s\n" "$out"; exit 1; }; printf "%s\n" "$out" | grep -q -- "--- PASS: TestReviewCommandBlocksWithoutSelectingNone"` — expected: exit 0; before this Task the case does not exist, so the command fails.
- `go run -buildvcs=false ./cmd/roundfix --help 2>&1 | grep -q "roundfix review"` — expected: exit 0; the command appears on the public surface. Before this Task the usage has no review line, so the command fails.

## References

- [_techspec.md](_techspec.md) — The command
