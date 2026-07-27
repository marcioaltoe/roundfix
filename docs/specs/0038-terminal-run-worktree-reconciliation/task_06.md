---
task: task_06
spec: 0038-terminal-run-worktree-reconciliation
status: completed
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

- [x] Add canonical reconciliation vocabulary.
- [x] Update command and usage documentation.
- [x] Add dry-run, JSON, and explicit apply examples.
- [x] Explain state, proof, and preservation semantics.
- [x] Document Runs List retained-worktree guidance.
- [x] Resolve ADR, Spec, finding, and command links.

## Acceptance Criteria

- [x] A reader can predict all five states from the documented evidence rules.
- [x] Examples distinguish requested stdout from diagnostic stderr.
- [x] No documentation suggests age, outcome, missing path, or force is
      sufficient for deletion.
- [x] Integration Pending and other terminal outcome behavior is explicit.
- [x] Runs List guidance points to Reconcile Command for classification.
- [x] Every cited link resolves and terminology matches `CONTEXT.md`.
- [x] No protected tooling path changes in this Task.

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

- `rtk grep -n 'roundfix reconcile\|Reconcile Command\|safe\|unintegrated\|released' CONTEXT.md docs/user-guide/commands.md docs/user-guide/usage.md`
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

## Result

Published the Reconcile Command as a dry-run-first, proof-based workflow with
single-Run, repository-wide, JSON, and explicit apply examples. The canonical
vocabulary and user guides now define all five reconciliation states, separate
requested stdout from diagnostic stderr, explain conservative preservation,
and route Runs List discovery to the Reconcile Command without treating the
listing as a classification surface.

Verification:

- `rtk grep -n 'roundfix reconcile\|Reconcile Command\|safe\|unintegrated\|released' CONTEXT.md docs/user-guide/commands.md docs/user-guide/usage.md`
  passed and found the canonical command, all-state guidance, dry-run/apply
  examples, and retained-worktree note across the three supported documents.
- `rtk go test ./internal/cli -run 'Test.*(Reconcile.*Help|DocumentationContract|RunsList.*Retained)' -count=1`
  passed with 30 tests in one package.
- `rtk git diff --check` passed.
- Link and anchor inspection found `CONTEXT.md`, ADR-0053, Spec 0038, the
  detached-watch finding, the Reconcile Command heading, the Stop Command
  heading, and finding 4 at their cited paths and anchors.
- `rtk git -c core.fsmonitor=false diff --name-only` listed only `CONTEXT.md`,
  this Task file, `docs/user-guide/commands.md`, and
  `docs/user-guide/usage.md`.

Acceptance evidence:

- The state table and glossary define `safe`, `unintegrated`, `dirty`,
  `unknown`, and `released` from cleanliness, ref resolution, ancestry, and
  surface-presence evidence.
- Redirection examples send Reconcile reports and Runs List rows to stdout
  files while sending Runs List discovery guidance to a separate stderr file.
- Both guides state that age, terminal outcome, one missing path, and a force
  assertion are insufficient; only freshly revalidated `safe` evidence can
  release work.
- The apply contract records a safe Integration Pending Run as Clean and keeps
  every other terminal outcome unchanged.
- Runs List keeps its stdout rows unchanged, reports the retained count on
  stderr without classification, and points to `roundfix reconcile`.
- The link inspection and glossary comparison resolved every new citation and
  kept canonical capitalization and terminology.
- The changed-path inspection proves that neither protected Roundfix Skill path
  changed in this Task.

Follow-ups: none.

### Verification feedback attempt 1

The Daemon exposed an escaping defect in the first Verification command: its
task-file pattern contained two literal backslashes before each basic regular
expression alternation, so the exact extracted command matched none of the
documented terms. The Verification contract now stores one backslash before
each alternation. The corrected command passed locally against all three
documented surfaces; the Daemon remains authoritative for the full rerun.
