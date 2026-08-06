---
task: task_07
spec: 0065-loop-order-and-verification-honesty
status: pending
type: chore
complexity: low
---

# Task 07: Correct the blocked-row count the loop clause cites

## Overview

Corrective Task from the QA gate's F-002 (`Trust-Damage`). The loop-order
rationale cites Spec 0078's gate as passing with **nine of eighteen** rows
blocked. Its archived report records **eleven**: nine of them blocked on the
absent Pull Request, plus one on host-load timing variance and one on a sandbox
that denied Agent state.

The conclusion the citation supports is true — ADR-0080's blocked-row typing
absorbs a Spec whose acceptance observes its own Pull Request, which is why
this Spec keeps the ADR-0091 order. What is false is the number offered as its
proof, and a rationale carrying a wrong measurement teaches the reader to
distrust the correct ones beside it.

## Requirements

1. MUST correct the count in every carrier to eleven of eighteen, naming that
   nine of those eleven were blocked on the absent Pull Request.
2. MUST correct it in all five carriers the gate found:
   `docs/agents/autonomous-work.md`, its Baseline formatter golden,
   `internal/baseline/assets/modules/autonomous-work.json`, this Spec's
   `_techspec.md`, and the Skills that restate it.
3. MUST verify the corrected number against the archived Spec 0078 report
   rather than against another carrier, so the correction cannot inherit the
   defect it repairs.
4. MUST keep the surrounding argument unchanged: the conclusion was never the
   error.
5. MUST run `make skills-sync` if any Skill text changed, then
   `make baseline-digests`, then re-record the two characterization corpora:

   ```bash
   go test ./internal/baseline -count=1 \
     -run TestBaselinePlanCharacterization -update-baseline-plan-characterization
   go test ./internal/baseline -count=1 \
     -run TestCatalogDiagnosticCharacterization -update-catalog-diagnostics
   ```

6. MUST leave `SC-LOOP-ORDER-DIVERGENT` passing: the three order statements
   must still agree after the edit.
7. MUST NOT change Go source.

## Subtasks

- [ ] Read the archived Spec 0078 report and record the exact counts.
- [ ] Correct every carrier and run the regeneration chain.

## Acceptance Criteria

- [ ] No carrier states `nine of eighteen`.
- [ ] Every carrier states eleven of eighteen, with nine attributed to the
      absent Pull Request.
- [ ] The number matches the archived Spec 0078 report, cited in the Result.
- [ ] `SC-LOOP-ORDER-DIVERGENT` passes.
- [ ] No Go source changed.

## Context

- instruction: `docs/agents/autonomous-work.md`
- interface: `docs/specs/_archived/0078-roundfix-asks-for-the-review/qa/qa-report-2026-08-05.md`

## Verification

- `grep -rn "nine of eighteen" docs/agents .agents skills internal/baseline/assets docs/specs/0065-loop-order-and-verification-honesty/_techspec.md | grep -q . && exit 1 || exit 0`
  — expected: exit 0; the wrong count is gone from every carrier.
- `grep -q "eleven of eighteen" docs/agents/autonomous-work.md` — expected:
  exit 0; the guide carries the corrected count.
- `go run -buildvcs=false ./cmd/roundfix spec check > /dev/null` — expected:
  exit 0; the order statements still agree.
- `make verify` — expected: exit 0.
- `git diff --name-only HEAD | grep -E "\.go$" | grep -q . && exit 1 || exit 0`
  — expected: exit 0; no Go source changed.

## References

- `qa/qa-report-2026-08-05.md` → F-002; row R10.
- `_techspec.md` → The order, decided.
- `docs/workflow/authorizations/2026-08-04-queue-tail-tooling.md`.
- ADR-0080, ADR-0081.
