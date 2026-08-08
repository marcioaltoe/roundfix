---
task: task_08
spec: 0088-a-third-runtime-that-can-run
status: pending
type: qa
complexity: high
---

# Task 08: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It runs after every implementation leaf
settles and walks the PRD's goals against the built binary. One row is the
acceptance the maintainer chose and no fixture can substitute for: a real Roundfix
Run executing a real Task on an `opencode-go` model. This node also carries the
glossary check, so a term this Spec coined, changed, or dropped is noticed while
the work is still open.

## Requirements

1. MUST execute the `qa-gate` skill as this Spec's authored terminal gate and
   write its report to the Spec's QA directory with a machine-readable verdict.
2. MUST walk every PRD goal against the built binary rather than against a test
   fixture.
3. MUST include one row in which a Roundfix Run executes a real Task on an
   `opencode-go` model, recording the Run identifier, the proven Agent Selection
   tuple, the Task outcome, and where the captured evidence lives.
4. MUST record that row as blocked with its reason, rather than dropping it, when
   the subscription, the network, or the runtime makes the Run impossible.
5. MUST include one row that re-runs the adopted measurement's failing sequence —
   a configured optional-category profile that fails — and states whether the
   Doctor Command now reports it, comparing against the measurement's recorded
   `profiles: ok (5 distinct tuples; 10 category references)`.
6. MUST record the origin of the external evidence used, per ADR-0104.
7. MUST run the glossary check against `CONTEXT.md` and record whether this Spec
   coined, changed, or dropped a term, and what was done about it.
8. MUST prove the authoritative gate on a clean test cache, because a stale cache
   reported exit 0 over a failing test on 2026-08-08.
9. MUST report `rows_blocked_environment` and `rows_blocked_finding` counts in the
   report frontmatter.
10. MUST NOT mutate any repository outside this one.

## Subtasks

- [ ] Run the gate and write the report to the Spec's QA directory.
- [ ] Walk every PRD goal against the built binary.
- [ ] Execute a real Run on an `opencode-go` model and capture its evidence.
- [ ] Re-run the measurement's failing sequence and compare the outcome.
- [ ] Run the glossary check and record its result.
- [ ] Prove the authoritative gate on a clean test cache.

## Acceptance Criteria

- [ ] A QA report exists in the Spec's QA directory with a machine-readable
      verdict and both blocked-row counts in frontmatter.
- [ ] Every PRD goal has a row.
- [ ] One row records a Roundfix Run on an `opencode-go` model with its Run
      identifier, proven tuple, Task outcome, and evidence location, or a recorded
      reason it was blocked.
- [ ] One row compares the Doctor Command's current output against the
      measurement's recorded output for a failing configured profile.
- [ ] The report names where its external evidence came from.
- [ ] The report records the glossary check result.
- [ ] The report shows the authoritative gate passing after the test cache was
      cleaned.

## Context

- instruction: `.agents/skills/qa-gate/SKILL.md`
- instruction: `docs/agents/domain.md`
- instruction: `docs/agents/spec-routing.md`

## Verification

- `test -d docs/specs/0088-a-third-runtime-that-can-run/qa` — expected: exits 0.
- `grep -q 'rows_blocked_environment' docs/specs/0088-a-third-runtime-that-can-run/qa/qa-report-*.md` — expected: exits 0.
- `grep -q 'rows_blocked_finding' docs/specs/0088-a-third-runtime-that-can-run/qa/qa-report-*.md` — expected: exits 0.
- `grep -q 'opencode-go' docs/specs/0088-a-third-runtime-that-can-run/qa/qa-report-*.md` — expected: exits 0, proving the acceptance Run row exists.
- `go clean -testcache` — expected: exits 0.
- `make verify` — expected: exits 0 on a clean cache, proving the authoritative gate passes on the finished Spec.

## References

- `_prd.md` → every Goal; Core Features.
- `_techspec.md` → Testing Approach; Build Order 8.
- `references/2026-08-08-what-the-opencode-adapter-answers-before-its-first-prompt.md`
  → the pre-Spec outcomes this gate compares against.
- ADR-0091, ADR-0096, ADR-0097, ADR-0104, ADR-0105, ADR-0106, ADR-0107.
