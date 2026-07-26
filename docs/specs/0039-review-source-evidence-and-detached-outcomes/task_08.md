---
task: task_08
spec: 0039-review-source-evidence-and-detached-outcomes
status: pending
type: docs
complexity: medium
---

# Task 08: Align review-evidence docs and glossary

## Overview

Publish Review Source Evidence, Review Skipped, bounded retry, issue knowledge,
Detached monitoring, notification receipts, and artifact-only inheritance
through canonical vocabulary and supported user guidance. This slice leaves
every protected Skill path unchanged.

## Requirements

1. MUST define Review Source Evidence and Review Skipped in canonical
   vocabulary.
2. MUST document accepted current-head check, status, and CodeRabbit approval
   evidence plus stale/unresolved refusal.
3. MUST document typed bounded retry and the absence of log-text inference or
   new retry configuration.
4. MUST document wait phases, deadlines, retry episodes, and unknown Review
   Issue counts.
5. MUST document Detached Supervisor monitoring and notification receipt
   semantics.
6. MUST document the exact artifact-only inheritance proof and its refusals.
7. MUST preserve ADR/finding traceability and leave protected tooling
   unchanged.

## Subtasks

- [ ] Add canonical Evidence and Review Skipped vocabulary.
- [ ] Update Watch and Detached command guidance.
- [ ] Document retry, waits, and unknown issue knowledge.
- [ ] Document notification context and receipt events.
- [ ] Document artifact-only inheritance boundaries.
- [ ] Resolve ADR, finding, Spec, and command links.

## Acceptance Criteria

- [ ] A reader can distinguish pending, reviewed, verified, skipped, and failed
      Evidence on the expected head.
- [ ] Review Skipped cannot be mistaken for Clean, Clean Unverified, or a
      zero-issue Round.
- [ ] Retry guidance names positive typed conditions and existing bounds only.
- [ ] Detached examples contain the stable Supervisor outcome command.
- [ ] Notification delivery is described as best-effort and separately
      receipted.
- [ ] Artifact inheritance guidance refuses every non-Daemon or mixed-path
      descendant.
- [ ] Every link resolves, canonical terms are used, and no protected tooling
      path changes.

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

- `rtk grep -n 'Review Source Evidence\\|Review Skipped\\|events .*--follow.*outcome\\|Review Issues: unknown' CONTEXT.md docs/user-guide/commands.md docs/user-guide/usage.md`
  — expected: canonical Evidence, outcome, monitor, and unknown-count guidance
  is present.
- `rtk go test ./internal/cli -run 'Test.*(Watch.*Help|Detached.*Monitor|DocumentationContract)' -count=1`
  — expected: supported command examples match implemented behavior.
- `rtk git diff --check`
  — expected: documentation contains no whitespace errors.

## References

- `_prd.md` → User Stories 1–8; User Experience; Non-Goals; Decisions.
- `_techspec.md` → API Contracts; Risks & Considerations; Build Order 8.
- `../../adr/0054-review-source-evidence-determines-review-outcomes.md` →
  canonical review outcome behavior.
