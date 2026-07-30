---
task: task_05
spec: 0053-qa-gate-reachability-and-verdict-semantics
status: completed
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

## Result

### Implementation

- The qa-gate Skill pair now uses typed blocked causes throughout planning,
  execution, coverage, frontmatter, and verdict settlement. It permits `pass`
  with environment-blocked rows only when each carries its cause and equivalent
  observed or supervised evidence, while any finding-blocked row prevents
  `pass`.
- The qa-gate Skill now treats the prompt's Pull Request fact as the
  reachability decision. A named Open Pull Request enables read-only approval,
  status, unresolved-thread, Merge-Ready, and review-artifact ancestry
  observation; an absent Pull Request marks those journeys
  environment-blocked without consulting the unpushed Run Worktree branch.
- QA Report creation now allows only the unsuffixed daily filename and numeric
  `-NN` same-day siblings. The report template and coverage contract carry
  `rows_blocked_environment` and `rows_blocked_finding`.
- The protected-tooling audit now traverses every Task and related chronology
  once, then reports all missing, late, folded, misordered, untraceable,
  unsanctioned-derived, and out-of-scope authorization-shape problems together
  before flow QA.
- The Roundfix Skill pair now documents all six reconciliation states,
  including proof-based `superseded`, and limits `--apply` mutation to freshly
  revalidated `safe` or `superseded` work. Branch Integrity Preflight excludes
  QA-report-only branches from automatic integration and directs the operator
  to `roundfix reconcile --apply`.
- The user guide carries the same reachability, verdict, naming,
  reconciliation, and integration-exclusion contracts. Only the `QA Report`,
  `Run Worktree Reconciliation`, `Branch Integrity Preflight`, and `Reconcile
  Command` glossary definitions changed.

### Focused checks

- `rtk proxy env GOCACHE=/Users/marcio/.roundfix/worktrees/roundfix-339f8dac/run_20260730T131049Z_a02854137f3dd85c/.gocache go test ./skills -run 'TestAuthorialSkillSync/(qa-gate|roundfix)$|TestProjectConstraintQAGate' -count=1`
  — passed, 5 tests. This proves both edited Skill pairs are byte-identical and
  the qa-gate Project Constraint contract retains its required audit anchors.
- `rtk cmp -s .agents/skills/qa-gate/SKILL.md skills/qa-gate/SKILL.md` and
  `rtk cmp -s .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md` —
  passed.
- `rtk proxy git -c core.fsmonitor=false diff --unified=0 -- CONTEXT.md` —
  showed exactly four changed glossary definitions: `QA Report`, `Run
  Worktree Reconciliation`, `Branch Integrity Preflight`, and `Reconcile
  Command`.
- `rtk proxy env GOCACHE=/Users/marcio/.roundfix/worktrees/roundfix-339f8dac/run_20260730T131049Z_a02854137f3dd85c/.gocache go test ./skills -run 'TestAuthorialSkillSync|TestProjectConstraintQAGate' -count=1`
  — the Project Constraint and edited-pair checks passed, while the generated
  setup snapshot subtests failed with the expected stale qa-gate and Roundfix
  `contentDigest` diagnostics naming `make baseline-digests`.
- `rtk git diff --check` — passed after the final implementation and Result
  edits.
- The four commands in `## Verification` were not run. In particular,
  `make baseline-digests` remains Daemon-owned because it is a declared
  Verification command; no derived pin was edited by hand.

### Acceptance evidence

- Blocked causes and verdict: the canonical qa-gate Skill names both row
  renderings and count keys, requires equivalent evidence for an
  environment-blocked `pass`, and states that a nonzero
  `rows_blocked_finding` prevents `pass`.
- Report template: both typed count keys are present in frontmatter and in the
  Coverage instructions, including the zero-count case.
- Naming: the Skill permits only `qa-report-YYYY-MM-DD.md` and numeric
  `qa-report-YYYY-MM-DD-NN.md` siblings, and rejects scope or build suffixes.
- Reconciliation: the Roundfix Skill defines `superseded`, its fresh
  `--apply` revalidation, and the QA-report-only automatic-integration
  exclusion with reconcile guidance.
- Glossary: the zero-context diff identifies exactly the four Task-authorized
  definitions and shows each shipped behavior in its semantic owner.
- Derived pins and full Skill check: pending Daemon Verification. The focused
  authorial test proves the canonical Skill hashes changed and identifies the
  sanctioned regeneration command; the worktree contains no hand-edited
  derived pin.
