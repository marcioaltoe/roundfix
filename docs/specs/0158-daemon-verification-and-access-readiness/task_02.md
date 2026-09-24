---
task: task_02
spec: 0158-daemon-verification-and-access-readiness
status: pending
type: backend
complexity: high
---

# Task 02: Let a named Task repair a red repository gate

## Overview

When a Task's Verification carries the configured repository command, the Daemon runs it before the Task and blocks the Task on a red result, which also blocks a Task written to repair that red gate.

## Requirements

1. MUST accept `precondition_repairs` in `_authorization.md` frontmatter as a list of Task identifiers of the consuming Spec.
2. MUST refuse at Run planning an identifier that is not a Task in the Graph, and a named Task whose Verification does not carry the configured repository command verbatim, naming the Task and the reason.
3. MUST let a named Task start when the repository precondition is red, publishing the known-red entry as a Daemon Verification event.
4. MUST settle a named Task completed only when every Verification command, the repository command included, passes.
5. MUST keep every unnamed Task blocked with "repository not green on entry".
6. MUST read the list only from the authorization resolved at Run start, never from the Run Worktree or a Task file.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] With the repository command red, a named Task reaches its Agent turn and settles completed once its commands pass.
- [ ] With the repository command red, an unnamed Task settles failed with "repository not green on entry".
- [ ] A named Task lacking the verbatim command is refused before the Run starts.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/authorization/authorization.go`
- interface: `internal/daemon/task_engine.go`
- interface: `internal/cli/implement.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestAuthorizationParsesPreconditionRepairs|TestNamedTaskRepairsKnownRedPrecondition|TestUnnamedTaskStaysBlockedByRedPrecondition|TestPreconditionRepairWithoutConfiguredCommandIsRefused)$" ./internal/authorization ./internal/daemon ./internal/cli 2>&1)" || { printf "%s\n" "$out"; exit 1; }; for name in TestAuthorizationParsesPreconditionRepairs TestNamedTaskRepairsKnownRedPrecondition TestUnnamedTaskStaysBlockedByRedPrecondition TestPreconditionRepairWithoutConfiguredCommandIsRefused; do printf "%s\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task none of the four cases exists, so the command fails.

## References

- [_techspec.md](_techspec.md) — Precondition repair
