---
task: task_06
spec: 0037-terminal-outcome-integrity
status: pending
type: docs
complexity: medium
---

# Task 06: Align terminal-outcome operator guidance

## Overview

Publish the completed Force Stop, graceful Stop Request, terminal replay, and
cleanup contracts through supported operator documentation and command help.
This documentation slice leaves every protected Skill path unchanged so the
isolated tooling Task can publish the final Agent-facing wording.

## Requirements

1. MUST document that Force Stop reports Stopped only after owner exit proof.
2. MUST document the failed-closed result, retained Active Run lock, and
   actionable inspection or retry guidance.
3. MUST document graceful Stop Request observation during Review Source waits.
4. MUST document idempotent same-outcome replay and conflicting-outcome
   rejection without exposing internal storage details as user API.
5. MUST describe registered Agent Session cleanup and primary-before-secondary
   diagnostics.
6. MUST preserve canonical glossary terms and finding/ADR traceability.
7. MUST leave all protected tooling files unchanged.

## Subtasks

- [ ] Update Stop Command and operational usage guidance.
- [ ] Align command help with proof-before-completion behavior.
- [ ] Document graceful watch interruption and terminal replay.
- [ ] Document cleanup eligibility and diagnostic ordering.
- [ ] Resolve glossary, ADR, Spec, and finding links.
- [ ] Check supported examples and terminology.

## Acceptance Criteria

- [ ] A reader can predict successful and failed Force Stop state and lock
      behavior.
- [ ] Guidance states when a graceful watch stop is observed and what work will
      not run afterward.
- [ ] Registered-session cleanup and secondary-warning behavior are explicit.
- [ ] User docs do not promise terminal overwrite or warning-only owner
      reclamation.
- [ ] Every cited ADR, finding, and command reference resolves.
- [ ] No protected tooling path changes in this Task.

## Context

- instruction: `docs/agents/cli.md`
- instruction: `docs/agents/domain.md`
- instruction: `.agents/skills/tech-writer/SKILL.md`
- interface: `CONTEXT.md`
- interface: `README.md`
- interface: `docs/user-guide/commands.md`
- interface: `docs/user-guide/usage.md`
- interface: `docs/findings/2026-07-16-vortex-pr87-detached-watch-notification.md`

## Verification

- `rtk grep -n 'stop --force\\|Force Stop\\|Stop Request' docs/user-guide/commands.md docs/user-guide/usage.md`
  — expected: supported guidance covers force and graceful stop contracts.
- `rtk go test ./internal/cli -run 'Test.*Stop.*(Help|Usage|Report)' -count=1`
  — expected: command help and public reports match documented behavior.
- `rtk git diff --check`
  — expected: documentation contains no whitespace errors.

## References

- `_prd.md` → User Stories 1–5; User Experience; Non-Goals; Decisions.
- `_techspec.md` → API Contracts; Risks & Considerations; Build Order 6.
- `../../adr/0052-run-completion-is-compare-and-set.md` → public terminal
  integrity contract.
