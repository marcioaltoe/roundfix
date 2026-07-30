---
task: task_05
spec: 0053-qa-gate-reachability-and-verdict-semantics
status: pending
type: docs
complexity: medium
---

# Task 05: Align the Skill pairs, the guides, and the authorized digests

## Overview

The verdict rule lives in the qa-gate Skill, not in Roundfix (ADR-0080), so the
behavior the earlier Tasks made possible only becomes usable once the Skill says
how to reach it. Update the authorized Skill pairs, the guides, and the
`CONTEXT.md` glossary entries whose shipped behavior these Tasks changed, then
let the sanctioned command own the derived digest fallout.

## Requirements

1. MUST state the verdict rule in the qa-gate Skill pair: a row the environment
   made unreachable is recorded with its cause and counted in
   `rows_blocked_environment`, and it does not by itself prevent `pass`; a row a
   finding blocks counts in `rows_blocked_finding` and does prevent it.
2. MUST document the read-only observation journeys and how the Pull Request
   fact decides whether a Pull Request journey is runnable or
   environment-blocked.
3. MUST state the report naming contract as numeric same-day suffixes only, and
   include the typed counts in the Skill's report template.
4. MUST teach the roundfix Skill pair the `superseded` reconcile vocabulary and
   the QA-report-only exclusion from automatic integration.
5. MUST fold in the one-pass authorization-audit reporting Spec 0054 asks of the
   gate, since this Spec owns the only authorized qa-gate mutation.
6. MUST update exactly these `CONTEXT.md` glossary entries: `Run Worktree
   Reconciliation` gains `superseded`; `Reconcile Command` widens to revalidated
   `superseded` work; `Branch Integrity Preflight` gains the QA-report-only
   exclusion; `QA Report` mentions the typed blocked-cause counts.
7. MUST confine protected-tooling edits to the exact paths the TechSpec's
   Tooling authority section authorizes, and obtain every derived pin from
   `make baseline-digests` rather than by hand.
8. MUST NOT edit the glossary entries ahead of the behavior landing in
   task_01 through task_04.

## Subtasks

- [ ] Update the qa-gate Skill pair: verdict rule, journeys, naming, template,
      authorization audit.
- [ ] Update the roundfix Skill pair: reconcile vocabulary and integration
      exclusion.
- [ ] Update the user guide and the four `CONTEXT.md` glossary entries.
- [ ] Run `make baseline-digests` and commit its output as the authorized
      fallout.

## Acceptance Criteria

- [ ] The qa-gate Skill states both blocked causes, their counts, and which one
      prevents `pass`.
- [ ] The Skill's report template carries both count keys.
- [ ] The Skill instructs numeric same-day suffixes only.
- [ ] The roundfix Skill describes `superseded` and the integration exclusion.
- [ ] The four glossary entries describe the behavior these Tasks shipped.
- [ ] `roundfix skills check` passes and every changed derived pin came from
      `make baseline-digests`.

## Context

- interface: `.agents/skills/qa-gate/SKILL.md`
- interface: `skills/qa-gate/SKILL.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `CONTEXT.md`
- interface: `docs/user-guide/context-driven-development.md`

## Verification

- `make skills-sync-check` — expected: the embedded bundle matches the canonical
  skills.
- `roundfix skills check` — expected: pass.
- `make baseline-digests` — expected: `ok: true`, with only authorized paths
  changed.
- `make verify` — expected: exit 0.

## References

`_prd.md` → Goal 1 Stories 1–2, Goal 2 Story 5, Goal 4 Story 6;
`_techspec.md` → Build Order 5, Decisions, Tooling authority; ADR-0080.
