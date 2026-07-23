---
task: task_05
spec: 0046-public-context-driven-baseline-command
status: pending
type: backend
complexity: high
---

# Task 05: Plan root-instruction preservation

## Overview

Extend the read-only planner through greenfield, preservation, Source
Baseline, and Readoption decisions. The result names every root backup and
repository-rules disposition before any instruction carrier can change.

## Requirements

1. MUST offer explicit greenfield and preservation modes for unconfigured or
   incompatible repositories.
2. MUST plan immutable content-addressed root backups and preserve a safe alias
   target exactly once.
3. MUST import no existing rule in greenfield mode and require one validated
   disposition for every root rule in preservation mode.
4. MUST retain every maintained Source Baseline, Readoption, Upgrade Retention,
   Decision Document, and Repository-Specific Normative Rules contract.
5. MUST leave nested carriers unchanged and expose their conflicts only as
   warnings.
6. MUST emit complete Decision Document skeletons with the required schema
   fields and stable next actions.

## Subtasks

- [ ] Port Source Baseline, Readoption, and retention decision models.
- [ ] Resolve greenfield and preservation planning states.
- [ ] Derive content-addressed root backup identities.
- [ ] Render consolidated manual classification inputs and warnings.
- [ ] Add parity and refusal tests for every preservation state.

## Acceptance Criteria

- [ ] Greenfield plans backups but moves zero existing rules into repository rules.
- [ ] Preservation cannot become ready with an unclassified root rule.
- [ ] Backup paths use the full raw-content SHA-256 identity and reject collisions.
- [ ] Safe aliases back up one target; unsafe aliases remain blocking.
- [ ] Decision skeletons pass the strict runtime parser without manual schema repair.
- [ ] Exact maintained Readoption and retention fixtures match the Python corpus.

## Context

- instruction: `docs/adr/0058-upgrade-retention-requires-clause-level-accounting.md`
- instruction: `docs/adr/0070-baseline-audits-all-carriers-but-preserves-root-instructions.md`
- interface: `.agents/skills/setup-context-driven/scripts/context_baseline.py`
- interface: `docs/agents/repository-rules.md`

## Verification

- `rtk go test -count=1 ./internal/baseline ./internal/cli -run 'TestGreenfieldPlan|TestPreservationPlan|TestRootBackupIdentity|TestDecisionDocumentSkeleton|TestReadoptionCompatibility'` — expected: preservation modes, backups, strict decision inputs, and parity cases pass.
- `rtk go test -count=1 ./internal/baseline -run 'TestNestedCarrierWarning|TestPreservationRequiresEveryDisposition'` — expected: nested files stay out of mutation plans and incomplete classifications remain unresolved.
- `rtk make verify` — expected: the full repository gate passes.

## References

- `_prd.md` → User Stories 3 and 5; Core Features 5–9 and 19; User Experience.
- `_techspec.md` → Data Models: RepositorySnapshot and backups; Testing Approach; Build Order 3–4.
- ADR-0058 → clause-level retention.
- ADR-0070 → root-only automatic preservation.
