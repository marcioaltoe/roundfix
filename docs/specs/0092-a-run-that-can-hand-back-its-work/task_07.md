---
task: task_07
spec: 0092-a-run-that-can-hand-back-its-work
status: pending
type: qa
complexity: high
---

# Task 07: Run the final QA gate

## Overview

The authored terminal gate. It walks every PRD goal against the built binary and
carries the acceptance no fixture substitutes for: a real Run whose selected
runtime cannot serve it, observed activating its configured Fallback Chain. This
node also carries the glossary check.

## Requirements

1. MUST execute the `qa-gate` skill as this Spec's authored terminal gate and
   write its report to the Spec's QA directory with a machine-readable verdict.
2. MUST walk every PRD goal against the built binary rather than a fixture.
3. MUST include one row in which a real Run selects a runtime that cannot serve
   it and the configured Fallback Selection is attempted, recording the Run
   identifier, the selection failure, and the fallback tuple that ran. The
   2026-08-08 instance was a quota exhaustion that cannot be summoned on demand;
   the row MUST state which unavailability it used instead.
4. MUST include one row in which a Batch fails after resolving at least one
   Review Issue, recording that the resolved issue kept its outcome and that the
   Run still ended `Unresolved`.
5. MUST include one row proving no path reports `Clean` while an unresolved
   Review Issue remains.
6. MUST include one row in which a superseded Run Branch is discarded through
   Roundfix, recording the branch record written before removal, and one in
   which an unreachable commit refuses the discard.
7. MUST include one row in which a stopped Run's settled Task is carried forward,
   and one in which a moved input refuses it.
8. MUST record any of the above as blocked with its reason rather than dropping
   it when the runtime or environment makes it impossible.
9. MUST record the origin of the external evidence used, per ADR-0104, and
   separate it from this repository's own artifacts.
10. MUST run the glossary check against `CONTEXT.md` and record whether this
    Spec coined, changed, or dropped a term. `Selection Failure` and
    `Branch Disposition` are coined terms.
11. MUST prove the authoritative gate on a genuinely cold cache — the one the
    Makefile exports, not the user-level one.
12. MUST report `rows_blocked_environment`, `rows_blocked_finding`, and
    `rows_blocked_declared` counts in the report frontmatter.

## Subtasks

- [ ] Run the gate and write the report.
- [ ] Prove a fallback activates on a runtime that cannot serve.
- [ ] Prove a failed Batch keeps its outcomes and the Run stays Unresolved.
- [ ] Prove both dispositions and both refusals.
- [ ] Run the glossary check.

## Acceptance Criteria

- [ ] A QA report exists with a machine-readable verdict and all three
      blocked-row counts.
- [ ] Every PRD goal has a row.
- [ ] One row records a live fallback activation, naming the unavailability used,
      or a recorded reason it was blocked.
- [ ] One row records a failed Batch preserving a resolved issue with the Run
      still `Unresolved`.
- [ ] Two rows record the branch disposition and its refusal.
- [ ] Two rows record the carry-forward and its refusal.
- [ ] The report records the glossary check for both coined terms.

## Context

- instruction: `.agents/skills/qa-gate/SKILL.md`
- instruction: `docs/agents/domain.md`

## Verification

- `test -d docs/specs/0092-a-run-that-can-hand-back-its-work/qa` — expected: exits 0.
- `grep -lq 'rows_blocked_declared' docs/specs/0092-a-run-that-can-hand-back-its-work/qa/qa-report-*.md` — expected: exits 0.
- `grep -lq 'Selection Failure' docs/specs/0092-a-run-that-can-hand-back-its-work/qa/qa-report-*.md` — expected: exits 0.
- `grep -lq 'Branch Disposition' docs/specs/0092-a-run-that-can-hand-back-its-work/qa/qa-report-*.md` — expected: exits 0.
- `grep -lqE 'verdict: (pass|fail|partial)' docs/specs/0092-a-run-that-can-hand-back-its-work/qa/qa-report-*.md` — expected: exits 0, proving a verdict was written rather than left pending.
- `GOCACHE="$PWD/.gocache" go clean -testcache` — expected: exits 0.
- `make verify` — expected: exits 0 on a genuinely cold cache, with no Go package line reporting `(cached)`.

## References

- `_prd.md` → every Goal.
- `_techspec.md` → Testing Approach; Build Order 7.
- ADR-0010, ADR-0050, ADR-0091, ADR-0104, ADR-0113, ADR-0114, ADR-0115.
