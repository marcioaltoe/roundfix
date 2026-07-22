---
task: task_04
spec: 0043-context-driven-baseline-coverage-and-skill-restoration
status: pending
type: backend
complexity: high
---

# Task 04: Authorize exact setup Change Plans

## Overview

Replace the split preview/apply projections with one concrete Change Plan that
owns both machine-readable operations and executable bytes. Maintainers can
authorize that exact plan by digest; stale state, ambiguous removal authority,
or an observed delta mismatch produces no unapproved mutation.

## Requirements

1. MUST derive definite public `plannedChanges` and apply mutations from the
   same concrete Change Plan after decisions resolve.
2. MUST report every create, refresh, remove, rename, and reference edit with
   path, managed identity, state, applicable condition, reason, and exact path
   preimage/postimage digests.
3. MUST calculate a deterministic plan digest over the selection, decision
   values, catalog, ordered operations, and path digests while excluding
   volatile timestamps.
4. MUST let read-only audit accept repeated decision inputs and expose an
   authorizable digest only for a fully resolved concrete plan.
5. MUST require `--confirm-plan` for a non-empty apply, preserve categorized
   exits `0/1/2/3`, and keep structured results on stdout and diagnostics on
   stderr.
6. MUST preserve every unmarked file outside proven manifest ownership unless
   an explicit adoption or removal decision appears in the presented plan.
7. MUST stage writes, verify the observed affected-path delta against the
   authorized plan, roll back on failure, and leave a second apply empty.

## Subtasks

- [ ] Consolidate public operations and executable file mutations under one
      Change Plan.
- [ ] Add exact operation metadata and canonical plan digest calculation.
- [ ] Add audit decision inputs and digest-bound apply confirmation.
- [ ] Add explicit preserve/remove decisions for ambiguous old inventory.
- [ ] Add rename and reference-edit planning without double-counting path
      deltas.
- [ ] Add postwrite delta proof and rollback across all affected paths.
- [ ] Reproduce the observed omitted-removal and false-positive-removal cases.

## Acceptance Criteria

- [ ] The resolved public plan's unique path/digest triples equal the complete
      observed before/after tree delta after apply.
- [ ] The Go-to-Rust transition previews every actual managed removal before
      confirmation.
- [ ] An unmarked conditional guide outside prior manifest ownership is
      preserved and never appears as a removal candidate.
- [ ] Missing confirmation and stale confirmation exit `3` with the current
      plan and no writes; malformed confirmation exits `2`.
- [ ] A matching digest applies only the listed mutations and returns
      structured success without mixing diagnostics into stdout.
- [ ] Missing or invalid old ownership markers require an explicit
      `preserve|remove` answer before any affected path can change.
- [ ] A forced postwrite mismatch or I/O failure restores all original bytes.
- [ ] Repeating confirmed apply against the resulting repository reports an
      empty plan and exits `0` without another confirmation.
- [ ] Canonical and embedded setup skill trees are synchronized after the
      slice.

## Context

- instruction: `.agents/skills/agentic-cli-design/SKILL.md`
- instruction: `docs/adr/0046-setup-owned-agent-instructions-are-declarative.md`
- instruction: `docs/adr/0047-setup-decisions-declare-their-effects.md`
- interface: `.agents/skills/setup-context-driven/scripts/context_setup.py`
- interface: `.agents/skills/setup-context-driven/tests/test_preview.py`
- interface: `.agents/skills/setup-context-driven/tests/test_apply.py`

## Verification

- `rtk python3 -B .agents/skills/setup-context-driven/tests/test_preview.py`
  — expected: resolved previews expose exact operations and digests; unresolved
  decisions remain conditional and deterministic.
- `rtk python3 -B .agents/skills/setup-context-driven/tests/test_apply.py`
  — expected: confirmation, stale-state rejection, removal authority, exact
  tree-delta parity, rollback, and idempotency cases pass.
- `rtk make verify` — expected: the full repository gate passes with the
  preview-first apply contract and stable CLI behavior.

## References

- `_prd.md` → Goal 4; User Story 4; Core Features 5–6, 9; User Experience;
  Success Metrics.
- `_techspec.md` → Interfaces: FileMutation and ChangePlan; Data Models:
  PlannedChange and removal authority; API Contracts: audit and apply; Build
  Order 4.
- ADR-0046 → explicit ownership and confirmation boundaries.
- ADR-0047 → one shared Decision Plan across preview, audit, and apply.
