---
task: task_04
spec: 0158-daemon-verification-and-access-readiness
status: completed
type: docs
complexity: low
---

# Task 04: Record the decisions

## Overview

Two decisions change what the Daemon allows and must outlive this Spec: independence is declared, never inferred, and the entry to a red repository gate is granted only by the frozen authorization record.

## Requirements

1. MUST add ADR-0159 recording that Verification commands run past a failure only when the Task declares them independent, and why inference from command text was rejected.
2. MUST add ADR-0160 recording that only Tasks named in the authorization resolved at Run start may enter a red repository gate, and that they settle only on the same gate passing.
3. MUST document in `docs/user-guide/configuration.md` that requested full access is proven during profile readiness and refused before a Run when it cannot be honoured.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] Both ADRs exist with accepted status.
- [ ] The configuration guide states the readiness refusal.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- creates: `docs/adr/0159-verification-runs-past-a-failure-only-when-declared-independent.md`
- creates: `docs/adr/0160-only-the-frozen-authorization-opens-a-red-repository-gate.md`
- interface: `docs/user-guide/configuration.md`

## Verification

- `grep -q "^status: accepted" docs/adr/0159-verification-runs-past-a-failure-only-when-declared-independent.md && grep -q "^status: accepted" docs/adr/0160-only-the-frozen-authorization-opens-a-red-repository-gate.md && grep -q "precondition_repairs" docs/adr/0160-only-the-frozen-authorization-opens-a-red-repository-gate.md && grep -qi "readiness" docs/user-guide/configuration.md && grep -q "agent_full_access" docs/user-guide/configuration.md && grep -qi "refuse" docs/user-guide/configuration.md` — expected: exit 0; before this Task neither ADR exists, so the command fails.

## References

- [_techspec.md](_techspec.md) — Decisions

## Result

Recorded the two accepted decisions required by this Task. ADR-0159 makes
independence an explicit `verification: independent` declaration and rejects
inferring it from command text. ADR-0160 makes the Run-start frozen
`precondition_repairs` authorization the only entry to a red repository gate
and requires the same repository command to pass before settlement. The
configuration guide now documents that requested `agent_full_access` is proven
during profile readiness and refused before Run creation when it cannot be
honoured.

Focused-check evidence:

- Pre-change: both ADR paths were absent, and the configuration guide had no
  sentence covering refusal of an unhonoured requested full-access mode before
  a Run.
- `rtk git diff --check` passed.
- `rtk rg -n '^status: accepted$' docs/adr/0159-verification-runs-past-a-failure-only-when-declared-independent.md docs/adr/0160-only-the-frozen-authorization-opens-a-red-repository-gate.md` found accepted status in both ADRs.
- `rtk rg -n -i 'access policy readiness|agent_full_access|readiness refuses before a Run|requested full-access mode' docs/user-guide/configuration.md` found the new readiness refusal guidance and the existing configuration key.
- Referenced ADR and TechSpec paths were checked as present.

The Daemon-owned `## Verification` command was not run in this Agent turn.
