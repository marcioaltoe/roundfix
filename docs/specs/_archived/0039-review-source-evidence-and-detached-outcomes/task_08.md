---
task: task_08
spec: 0039-review-source-evidence-and-detached-outcomes
status: completed
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

- [x] Add canonical Evidence and Review Skipped vocabulary.
- [x] Update Watch and Detached command guidance.
- [x] Document retry, waits, and unknown issue knowledge.
- [x] Document notification context and receipt events.
- [x] Document artifact-only inheritance boundaries.
- [x] Resolve ADR, finding, Spec, and command links.

## Acceptance Criteria

- [x] A reader can distinguish pending, reviewed, verified, skipped, and failed
      Evidence on the expected head.
- [x] Review Skipped cannot be mistaken for Clean, Clean Unverified, or a
      zero-issue Round.
- [x] Retry guidance names positive typed conditions and existing bounds only.
- [x] Detached examples contain the stable Supervisor outcome command.
- [x] Notification delivery is described as best-effort and separately
      receipted.
- [x] Artifact inheritance guidance refuses every non-Daemon or mixed-path
      descendant.
- [x] Every link resolves, canonical terms are used, and no protected tooling
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

- `rtk grep -n 'Review Source Evidence\|Review Skipped\|events .*--follow.*outcome\|Review Issues: unknown' CONTEXT.md docs/user-guide/commands.md docs/user-guide/usage.md`
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

## Result

Published the canonical Review Source Evidence states and Review Skipped
boundary, then aligned the README, command reference, operational guide, and
detached-watch finding with the implemented watch and Detached Run contracts.
The finding keeps its original observations and now carries a dated routing
addendum; its `done` status records Spec routing rather than a QA or release
verdict.

Verification:

- `rtk grep -n 'Review Source Evidence\|Review Skipped\|events .*--follow.*outcome\|Review Issues: unknown' CONTEXT.md docs/user-guide/commands.md docs/user-guide/usage.md`
  — passed; the output found canonical Evidence and Review Skipped terms, the
  stable Supervisor outcome command, and the unknown-count report in the
  expected guides.
- `rtk go test ./internal/cli -run 'Test.*(Watch.*Help|Detached.*Monitor|DocumentationContract)' -count=1`
  — passed: 24 tests in one package.
- `rtk git diff --check` — passed with no whitespace errors.
- A local target-existence check passed for the Task 08 glossary, ADR-0054,
  Specs 0037–0039, and detached-watch finding links. The repaired Spec 0038
  link now follows its archived location.
- `rtk git -c core.fsmonitor=false diff --name-only -- .agents/skills skills/roundfix`
  — passed with no output; protected tooling paths are unchanged.

Verification feedback repair:

- Attempt 1 exposed a verification-contract escaping error: the stored basic
  regular expression used two backslashes before each alternation, so grep
  searched for literal separator text and returned no matches.
- The Task command now uses one backslash for each basic-regular-expression
  alternation. The repaired command exited `0` and found every required term;
  `rtk git diff --check` also passed afterward.

Acceptance evidence:

- Expected-head Evidence is distinguished as `pending`, `reviewing`,
  `reviewed`, `verified`, `skipped`, or `failed`; current-head CodeRabbit
  checks, statuses, and approvals require zero unresolved CodeRabbit threads
  before `verified`, while stale signals remain non-verifying.
- Review Skipped has its own reason and next-action report, exit `3`, and
  explicit refusals to mean Clean, Clean Unverified, or a zero-issue Round.
- Retry guidance lists only positively typed timeout, DNS, reset, HTTP `429`,
  and GitHub `5xx` conditions, uses the existing poll interval, Review Source
  timeout, and Run Budget, and rejects progress, event, or Console Log text as
  a retry signal.
- Both Detached examples print
  `roundfix events <run-id> --follow --filter outcome`; wait guidance names
  `WaitingForReview`, `WaitingForReviewCheck`, their deadlines, and retry
  episodes.
- Notification guidance keeps delivery best-effort and outcome-independent,
  records separate `sent`, `skipped`, or `failed` receipts, and lists the
  additive `notify.command` context.
- Artifact inheritance requires the exact current Daemon-created commit,
  verified sole parent, non-empty review-root-only diff, refreshed Evidence,
  and no unresolved CodeRabbit threads. Non-Daemon, user-authored, mixed-path,
  out-of-root, wrong-parent, stale, and unresolved descendants fall back to
  ordinary current-head polling.

Follow-ups: none for this Task slice.
