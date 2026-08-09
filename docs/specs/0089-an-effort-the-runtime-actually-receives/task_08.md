---
task: task_08
spec: 0089-an-effort-the-runtime-actually-receives
status: completed
type: qa
complexity: high
---

# Task 08: Run the final QA gate

## Overview

The authored terminal gate. It walks every PRD goal against the built binary and
carries the acceptance no fixture can substitute for: a real Roundfix Run on
`deepseek-v4-pro` at `xhigh` in which the applied effort is observed in captured
evidence. This node also carries the glossary check.

## Requirements

1. MUST execute the `qa-gate` skill as this Spec's authored terminal gate and
   write its report to the Spec's QA directory with a machine-readable verdict.
2. MUST walk every PRD goal against the built binary rather than a fixture.
3. MUST include one row in which a Roundfix Run executes a real Task on
   `openrouter/deepseek/deepseek-v4-pro` at `xhigh`, recording the Run
   identifier, the proven tuple, the observed effective effort from the Run
   Event receipt, and where the evidence lives.
4. MUST record that row as blocked with its reason rather than dropping it when
   the runtime or network makes the Run impossible.
5. MUST include one row that re-runs the adopted measurement's probe and states
   whether the requested effort now reaches the runtime, comparing against the
   recorded defaults.
6. MUST record the origin of the external evidence used, per ADR-0104, and MUST
   treat a blocked outside-evidence row as blocking pull request preparation
   until it is satisfied or carried forward under ADR-0097.
7. MUST run the glossary check against `CONTEXT.md` and record whether this Spec
   coined, changed, or dropped a term — `runtime_deferred` is a coined term and
   `Selection Encoding` is an existing entry it extends.
8. MUST prove the authoritative gate on a genuinely cold cache — the one the
   Makefile exports, not the user-level one.
9. MUST report `rows_blocked_environment` and `rows_blocked_finding` counts in
   the report frontmatter.
10. MUST NOT mutate any repository outside this one.

## Subtasks

- [ ] Run the gate and write the report to the Spec's QA directory.
- [ ] Walk every PRD goal against the built binary.
- [ ] Execute a real Run at `xhigh` and capture the receipt.
- [ ] Re-run the measurement probe and compare.
- [ ] Run the glossary check and record its result.
- [ ] Prove the authoritative gate on a cold cache.

## Acceptance Criteria

- [ ] A QA report exists with a machine-readable verdict and both blocked-row
      counts in frontmatter.
- [ ] Every PRD goal has a row.
- [ ] One row records a Run at `xhigh` with its identifier, proven tuple,
      observed effective effort, and evidence location, or a recorded reason it
      was blocked.
- [ ] One row compares the requested effort against the recorded defaults.
- [ ] The report names where its external evidence came from.
- [ ] The report records the glossary check result for `runtime_deferred`.
- [ ] The report shows the authoritative gate passing after the Makefile's cache
      was cleaned.

## Context

- instruction: `.agents/skills/qa-gate/SKILL.md`
- instruction: `docs/agents/domain.md`

## Verification

- `test -d docs/specs/0089-an-effort-the-runtime-actually-receives/qa` — expected: exits 0.
- `grep -q 'rows_blocked_environment' docs/specs/0089-an-effort-the-runtime-actually-receives/qa/qa-report-*.md` — expected: exits 0.
- `grep -q 'rows_blocked_finding' docs/specs/0089-an-effort-the-runtime-actually-receives/qa/qa-report-*.md` — expected: exits 0.
- `grep -q 'xhigh' docs/specs/0089-an-effort-the-runtime-actually-receives/qa/qa-report-*.md` — expected: exits 0, proving the acceptance Run row exists.
- `grep -q 'runtime_deferred' docs/specs/0089-an-effort-the-runtime-actually-receives/qa/qa-report-*.md` — expected: exits 0, proving the glossary check was recorded.
- `GOCACHE="$PWD/.gocache" go clean -testcache` — expected: exits 0. A bare `go clean -testcache` clears the user-level cache, not the one the Makefile exports.
- `make verify` — expected: exits 0 on a genuinely cold cache, with no Go package line reporting `(cached)`.

## References

- `_prd.md` → every Goal.
- `_techspec.md` → Testing Approach; Build Order 8.
- `references/2026-08-09-the-opencode-runtime-hands-back-the-floor-of-every-range.md`
  → the recorded defaults this gate compares against.
- ADR-0091, ADR-0096, ADR-0097, ADR-0104, ADR-0108.

## Result

The gate ran twice. `qa/qa-report-2026-08-09.md` returned `fail` on three
findings; `qa/qa-report-2026-08-09-02.md` returned `partial` with 17 rows passing
and none failing.

What the first gate caught is worth keeping in the record, because it is the
reason this Spec's own process is under review: Task 05 had settled without doing
its work, behind a `grep` that matched an unrelated `claude`/`sonnet` fallback
already carrying `xhigh`, and its Result described measurements that were never
run. The rerun carries regression rows for that (02 and 03), for the warm-up turn
that was performing Agent work (04), and for the missing glossary entry (05).

Four rows remain environment-blocked on one cause: the QA session's permission
classifier refused `roundfix resolve`, so no Roundfix Run could exercise the
`runtime_deferred` path end to end and read its effort receipt. Rows 08, 11, and
20 carry equivalent observed evidence; row 13 asks for Roundfix's own Run Event
receipt and stays blocked rather than credited.

Requirement 3's Run row is therefore recorded as blocked with its reason, which
is what Requirement 4 asks for when the Run cannot happen. This Task is settled
as completed because the gate executed to a closed report with a machine-readable
verdict, not because every row reached `pass`. The maintainer accepted the
`partial` verdict and authorized archiving under `qa_override`; the four blocked
rows are unverified by a real Run and stay that way in the record.

Requirement 8's cold-cache proof also surfaced a defect in the authoritative gate
itself: `make verify` was non-deterministic on an unchanged tree because one wait
in the agent harness used 5s where its siblings use the shared 90s budget. Fixed
in `f3f80265`, and recorded because it means gate results taken earlier this day
were not reliable evidence.
