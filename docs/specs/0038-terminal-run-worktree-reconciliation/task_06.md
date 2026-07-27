---
task: task_06
spec: 0038-terminal-run-worktree-reconciliation
status: pending
type: docs
complexity: medium
---

# Task 06: Align reconciliation docs and glossary

## Overview

Publish the Reconcile Command, five reconciliation states, dry-run/apply flow,
and Runs List discovery contract through supported documentation and canonical
vocabulary. This slice leaves every protected Skill path unchanged.

## Requirements

1. MUST define the Reconcile Command and Run Worktree Reconciliation states in
   canonical glossary vocabulary.
2. MUST document single-Run and repository-wide dry-run examples.
3. MUST document `--apply` as the only mutation switch and the absence of a
   force bypass.
4. MUST explain positive cleanliness and ancestry proof plus conservative
   preservation of dirty, unintegrated, and unknown work.
5. MUST document Integration Pending promotion and unchanged other outcomes.
6. MUST document Runs List stderr discovery without promising classification.
7. MUST preserve ADR/finding traceability and leave protected tooling
   unchanged.

## Subtasks

- [ ] Add canonical reconciliation vocabulary.
- [ ] Update command and usage documentation.
- [ ] Add dry-run, JSON, and explicit apply examples.
- [ ] Explain state, proof, and preservation semantics.
- [ ] Document Runs List retained-worktree guidance.
- [ ] Resolve ADR, Spec, finding, and command links.

## Acceptance Criteria

- [ ] A reader can predict all five states from the documented evidence rules.
- [ ] Examples distinguish requested stdout from diagnostic stderr.
- [ ] No documentation suggests age, outcome, missing path, or force is
      sufficient for deletion.
- [ ] Integration Pending and other terminal outcome behavior is explicit.
- [ ] Runs List guidance points to Reconcile Command for classification.
- [ ] Every cited link resolves and terminology matches `CONTEXT.md`.
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

- `rtk grep -n 'roundfix reconcile\\|Reconcile Command\\|safe\\|unintegrated\\|released' CONTEXT.md docs/user-guide/commands.md docs/user-guide/usage.md`
  — expected: canonical command and state guidance is present.
- `rtk go test ./internal/cli -run 'Test.*(Reconcile.*Help|DocumentationContract|RunsList.*Retained)' -count=1`
  — expected: command and Runs List public contracts match documentation.
- `rtk git diff --check`
  — expected: documentation contains no whitespace errors.

## References

- `_prd.md` → User Stories 1–5; User Experience; Non-Goals; Decisions.
- `_techspec.md` → API Contracts; Risks & Considerations; Build Order 6.
- `../../adr/0053-terminal-run-worktree-reconciliation-is-proof-based.md` →
  canonical reconciliation behavior.
