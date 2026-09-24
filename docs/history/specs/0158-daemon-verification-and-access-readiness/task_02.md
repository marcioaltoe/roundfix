---
task: task_02
spec: 0158-daemon-verification-and-access-readiness
status: completed
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

## Result

Implemented authorization-scoped repository precondition repair. The
authorization reader preserves an exact, duplicate-free
`precondition_repairs` Task list. Implement planning validates that every
identifier belongs to the committed Task Graph and that each named Task
carries the configured repository command verbatim before any Run record is
created.

The Task engine uses only the authorization resolution frozen at Run start. A
named Task may continue after the repository command publishes its known-red
Daemon Verification event, but its normal post-Agent Verification still runs
every Task command and is the only path to completed settlement. An unnamed
Task retains the existing `repository not green on entry` failure.

Focused-check evidence:

- Pre-change: `rtk env GOCACHE=/tmp/roundfix-task02-gocache go test -count=1 -run '^TestAuthorizationParsesPreconditionRepairs$' ./internal/authorization` failed to compile because `AuthorizationRecord.PreconditionRepairs` did not exist.
- `rtk env GOCACHE=/tmp/roundfix-task02-gocache go test -count=1 -run '^TestAuthorizationParsesPreconditionRepairs$' ./internal/authorization` passed, proving the ordered Task identifiers survive authorization parsing.
- `rtk env GOCACHE=/tmp/roundfix-task02-gocache go test -count=1 -run '^(TestNamedTaskRepairsKnownRedPrecondition|TestUnnamedTaskStaysBlockedByRedPrecondition|TestTaskCycleRepositoryGatePreconditionFailureStartsNoAgentSession)$' ./internal/daemon` passed. The named Task received one Agent turn after a red entry, retained the known-red Verification event after the Worktree authorization was changed, reran both Task commands, and settled completed; both unnamed entry points settled failed without Agent work and retained `repository not green on entry`.
- `rtk env GOCACHE=/tmp/roundfix-task02-gocache go test -count=1 -run '^(TestPreconditionRepairWithoutConfiguredCommandIsRefused|TestUnknownPreconditionRepairTaskIsRefused)$' ./internal/cli` passed. A whitespace-only near-match for `make verify` and an identifier absent from the Graph were each refused with the Task and reason named, and neither case created a Run record.
- `rtk git diff --check` passed.

An additional `go test -count=1 ./internal/authorization ./internal/spec`
attempt passed `internal/authorization` but reached the repository's
`TestCoverageEquivalence` guard in `internal/spec`, which rejects focused
package execution when its observed repository-wide corpus differs from the
recorded corpus. No coverage expectation or generated baseline was changed.

The Daemon-owned `## Verification` command was not run in this Agent turn.
