---
task: task_07
spec: 0046-public-context-driven-baseline-command
status: pending
type: backend
complexity: high
---

# Task 07: Emit portable Baseline Plans

## Overview

Complete the non-interactive planning slice by turning equivalent repository
state and decisions into one portable, digest-bound JSON artifact. Humans and
automations can inspect the same canonical ledger through a concise file-level
projection.

## Requirements

1. MUST implement strict `roundfix/baseline-plan/v1` and
   `roundfix/baseline-result/v1` serializers with deterministic ordering,
   duplicate-key rejection, unknown-field rejection, and stable exit
   categories.
2. MUST include repository, catalog, profile, decision, retention, preimage,
   postimage, warning, Setup Manifest, and ordered managed-entry evidence.
3. MUST derive one `fileChanges` row per affected path from the canonical
   ledger and reject any projection mismatch.
4. MUST compute the Plan Digest from the approved canonical payload and exact
   postimages without volatile ACP metadata or absolute paths.
5. MUST allow `roundfix baseline plan` to emit text or a complete portable JSON
   artifact and never prompt.
6. MUST reproduce exact maintained planned bytes and digests and remain
   idempotent for equivalent snapshots and decisions.

## Subtasks

- [ ] Port Decision Plan, rendering, retention, and Setup Manifest generation.
- [ ] Define strict plan and result document codecs.
- [ ] Produce the canonical managed-entry ledger and derived file projection.
- [ ] Implement repository-relative postimages and Plan Digest calculation.
- [ ] Expose complete text/JSON planning results and stable failures.
- [ ] Add cross-clone, determinism, and compatibility tests.

## Acceptance Criteria

- [ ] Equivalent human-normalized and file-based inputs produce identical Plan Digests.
- [ ] JSON plans contain no absolute checkout path or hidden pending-state reference.
- [ ] Another clone with matching identity and bounded preimages accepts the same plan document.
- [ ] `fileChanges` has one row per path and always matches the canonical ledger.
- [ ] Missing decisions produce result-schema next actions and exit 3 without a partial plan.
- [ ] Exact maintained fixtures reproduce planned bytes, manifests, ledgers, and digests.
- [ ] Planning performs zero repository mutations.

## Context

- instruction: `docs/adr/0068-baseline-command-uses-one-confirmation-gated-workflow.md`
- instruction: `docs/adr/0071-baseline-plans-are-portable-and-preimage-bound.md`
- interface: `internal/releaseplan/plan.go`
- interface: `internal/cli/releaseplan_command.go`

## Verification

- `rtk go test -count=1 ./internal/baseline ./internal/cli -run 'TestPlanDocument|TestBaselinePlanCommand|TestFileChangesProjection|TestPlanDigest|TestCrossClonePlan|TestPlanDeterminism'` — expected: schemas, output, digest, portability, and no-write cases pass.
- `rtk go run -buildvcs=false ./cmd/roundfix baseline plan --help` — expected: help exposes the approved profile, decision, decision-file, repository, and format flags without interactive options.
- `rtk make verify` — expected: the full repository gate passes.

## References

- `_prd.md` → User Stories 2, 4, 7–8; Core Features 10, 12, 15–16, 18–20.
- `_techspec.md` → Data Models: PlanDocument and Result; API Contracts: Automation; Build Order 4 and 7.
- ADR-0068 → explicit automation planning stage.
- ADR-0071 → portable plan and bounded-preimage contract.
