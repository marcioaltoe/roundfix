---
task: task_01
spec: 0065-loop-order-and-verification-honesty
status: pending
type: chore
complexity: medium
---

# Task 01: State the loop order once

## Overview

`docs/agents/autonomous-work.md` states the corrected order in one place and
the original order in another, so the contradiction is live and both are being
followed. The same order is carried to adopting repositories by
`internal/baseline/assets/modules/autonomous-work.json`, which is the third
place it can disagree.

The order itself was settled on 2026-08-05 and does not change here: ADR-0091
stands, so the QA gate is the graph's terminal Task and runs before any Pull
Request exists.

## Requirements

1. MUST state one order in every place it appears: implement the graph
   including its authored gate, archive, open the Pull Request, watch until
   Clean, merge.
2. MUST remove the contradicting restatement rather than annotate it.
3. MUST carry the same order in the Baseline module asset, so adopting
   repositories receive the corrected clause.
4. MUST record why the order is this one: a Spec whose acceptance observes its
   own Pull Request reaches `pass` through ADR-0080's environment-blocked rows
   with equivalent evidence, proven by Spec 0078's gate on 2026-08-05.
5. MUST change only `docs/agents/autonomous-work.md`, the Baseline module asset
   `internal/baseline/assets/modules/autonomous-work.json`, this Task file, and
   the ADR-0081 digest fallout under `DERIVED_DIGEST_PATHS`.
6. MUST run `make baseline-digests` after the module asset edit, then re-record
   the two characterization corpora that command does not reach:

   ```bash
   go test ./internal/baseline -count=1 \
     -run TestBaselinePlanCharacterization -update-baseline-plan-characterization
   go test ./internal/baseline -count=1 \
     -run TestCatalogDiagnosticCharacterization -update-catalog-diagnostics
   ```

7. MUST NOT change Go source, and MUST NOT add the divergence check — that is
   task_04, which depends on this Task precisely so the check never lands while
   the sources disagree.

## Subtasks

- [ ] Remove the contradicting restatement and state the order once.
- [ ] Carry the same order into the Baseline module asset.
- [ ] Run the regeneration chain and both re-records.

## Acceptance Criteria

- [ ] `docs/agents/autonomous-work.md` states the order exactly once.
- [ ] The Baseline module asset states the same order.
- [ ] The recorded reason names ADR-0091, ADR-0080, and the Spec 0078 evidence.
- [ ] No Go source changed.
- [ ] `make verify` exits 0.

## Context

- instruction: `docs/agents/autonomous-work.md`
- instruction: `docs/adr/0091-the-qa-gate-is-a-task-node-of-its-own-type.md`

## Verification

- `make verify` — expected: exit 0.
- `git diff --name-only HEAD | grep -E "\.go$" | grep -q . && exit 1 || exit 0`
  — expected: exit 0; no Go source changed.
- `git diff --name-only HEAD | grep -vE "^(docs/agents/autonomous-work\.md$|docs/specs/0065-loop-order-and-verification-honesty/task_01\.md$|internal/baseline/(assets/(modules|setups|source-baselines|formatter-fixtures|profiles)|testdata)/)" | grep -q . && exit 1 || exit 0`
  — expected: exit 0; only the bounded paths changed.

## References

- `_prd.md` → Core Features 1 and 2; Success Metric 1.
- `_techspec.md` → The order, decided; Build Order 1.
- `docs/workflow/authorizations/2026-08-04-queue-tail-tooling.md`.
- ADR-0091, ADR-0080, ADR-0081.
