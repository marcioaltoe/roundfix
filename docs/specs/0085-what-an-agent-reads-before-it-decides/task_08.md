---
task: task_08
spec: 0085-what-an-agent-reads-before-it-decides
status: pending
type: qa
complexity: high
---

# Task 08: Run the final QA gate

## Overview

The authored terminal gate. It walks every PRD goal against the built binary and
carries one acceptance no fixture substitutes for: running `roundfix baseline`
against this repository to observe the catalog edits reaching the two guides
through the real update path. This node also carries the glossary check.

## Requirements

1. MUST execute the `qa-gate` skill as this Spec's authored terminal gate and
   write its report to the Spec's QA directory with a machine-readable verdict.
2. MUST walk every PRD goal against the built binary rather than a fixture.
3. MUST include one row running `roundfix baseline` against this repository,
   observing the consultation clause reaching both rendered guides.
4. MUST include one row proving every retired artifact family resolves under the
   single archive root through the built CLI.
5. MUST include one row proving no ADR states its status only in its body, and
   that every retired record carries a forward pointer.
6. MUST audit both authorized tooling Tasks: the authorization landed before
   them, each touched only its bounded files, and derived pins were regenerated
   by the sanctioned command rather than hand-edited.
7. MUST record any row as blocked with its reason rather than dropping it when
   the environment makes it impossible.
8. MUST record the origin of the external evidence used, per ADR-0104.
9. MUST run the glossary check against `CONTEXT.md` and record whether this Spec
   coined, changed, or dropped a term.
10. MUST prove the authoritative gate on a genuinely cold cache — the one the
    Makefile exports, not the user-level one.
11. MUST report `rows_blocked_environment`, `rows_blocked_finding`, and
    `rows_blocked_declared` counts in the report frontmatter.

## Subtasks

- [ ] Run the gate and write the report.
- [ ] Run `roundfix baseline` against this repository.
- [ ] Prove the single archive root through the built CLI.
- [ ] Audit both tooling Tasks against their bounded files.
- [ ] Run the glossary check.

## Acceptance Criteria

- [ ] A QA report exists with a machine-readable verdict and all three
      blocked-row counts.
- [ ] Every PRD goal has a row.
- [ ] One row records the Baseline update reaching both guides.
- [ ] One row records every retired family under the single root.
- [ ] One row records the ADR corpus carrying frontmatter statuses and forward
      pointers.
- [ ] The tooling audit reports authorization order, bounded paths, and
      regeneration provenance together.

## Context

- instruction: `.agents/skills/qa-gate/SKILL.md`
- instruction: `docs/agents/domain.md`

## Verification

- `test -d docs/specs/0085-what-an-agent-reads-before-it-decides/qa` — expected: exits 0.
- `grep -lq 'rows_blocked_declared' docs/specs/0085-what-an-agent-reads-before-it-decides/qa/qa-report-*.md` — expected: exits 0.
- `grep -lqE 'verdict: (pass|fail|partial)' docs/specs/0085-what-an-agent-reads-before-it-decides/qa/qa-report-*.md` — expected: exits 0, proving a verdict was written rather than left pending.
- `grep -lq 'baseline' docs/specs/0085-what-an-agent-reads-before-it-decides/qa/qa-report-*.md` — expected: exits 0, proving the Baseline row was recorded.

Whole-package sweeps, `go build`, `go clean -testcache` and `make verify` are
deliberately absent: each passes against a tree where no work has happened, so
it approves the Task before it starts. Regression is the Run-level gate's job.

## References

- `_prd.md` → every Goal.
- `_techspec.md` → Testing Approach; Build Order 8.
- ADR-0081, ADR-0091, ADR-0104.
