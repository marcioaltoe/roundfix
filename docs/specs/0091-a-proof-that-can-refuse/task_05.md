---
task: task_05
spec: 0091-a-proof-that-can-refuse
status: pending
type: qa
complexity: high
---

# Task 05: Run the final QA gate

## Overview

The authored terminal gate. It walks every PRD goal against the built binary and
carries the acceptance no fixture substitutes for: a live refusal on all three
runtimes, proved against the adapters as they actually behave rather than
against a recorded payload.

## Requirements

1. MUST execute the `qa-gate` skill as this Spec's authored terminal gate and
   write its report to the Spec's QA directory with a machine-readable verdict.
2. MUST walk every PRD goal against the built binary rather than a fixture.
3. MUST include one row per runtime — `codex`, `claude`, `opencode` — in which a
   configured selection naming a model that runtime does not offer is refused by
   `roundfix profiles validate`, recording the exit status and the advertised set
   named in the message.
4. MUST include one row proving the previously-passing tuple
   `claude` / `opus-9-does-not-exist` / `high` is now refused, quoting the
   before-state recorded in this Spec's PRD.
5. MUST include one row proving every selection in this repository's own
   `.roundfixrc.yml` still proves, with unchanged encodings.
6. MUST include one row proving preflight sent no prompt, so the change kept
   proof token-free.
7. MUST record that row as blocked with its reason rather than dropping it when
   a runtime is unavailable.
8. MUST record the origin of the external evidence used, per ADR-0104.
9. MUST run the glossary check against `CONTEXT.md` and record whether this Spec
   coined, changed, or dropped a term. `Runtime Catalogue` is a coined term.
10. MUST prove the authoritative gate on a genuinely cold cache — the one the
    Makefile exports, not the user-level one.
11. MUST report `rows_blocked_environment`, `rows_blocked_finding`, and
    `rows_blocked_declared` counts in the report frontmatter.

## Subtasks

- [ ] Run the gate and write the report.
- [ ] Refuse an unoffered model on each of the three runtimes, live.
- [ ] Prove this repository's own profile still validates unchanged.
- [ ] Run the glossary check.

## Acceptance Criteria

- [ ] A QA report exists with a machine-readable verdict and all three
      blocked-row counts.
- [ ] Every PRD goal has a row.
- [ ] Three rows record a live refusal, one per runtime, or a recorded reason
      each was blocked.
- [ ] One row records the previously-passing tuple now refused.
- [ ] One row records this repository's five distinct tuples still proving with
      their existing encodings.
- [ ] The report records the glossary check result for `Runtime Catalogue`.

## Context

- instruction: `.agents/skills/qa-gate/SKILL.md`
- instruction: `docs/agents/domain.md`

## Verification

- `test -d docs/specs/0091-a-proof-that-can-refuse/qa` — expected: exits 0.
- `grep -lq 'rows_blocked_declared' docs/specs/0091-a-proof-that-can-refuse/qa/qa-report-*.md` — expected: exits 0.
- `grep -lq 'Runtime Catalogue' docs/specs/0091-a-proof-that-can-refuse/qa/qa-report-*.md` — expected: exits 0.
- `grep -lqE 'verdict: (pass|fail|partial)' docs/specs/0091-a-proof-that-can-refuse/qa/qa-report-*.md` — expected: exits 0, proving a verdict was written rather than left pending.
- `GOCACHE="$PWD/.gocache" go clean -testcache` — expected: exits 0.
- `make verify` — expected: exits 0 on a genuinely cold cache, with no Go package line reporting `(cached)`.

## References

- `_prd.md` → every Goal.
- `_techspec.md` → Testing Approach; Build Order 5.
- ADR-0091, ADR-0104, ADR-0112.
