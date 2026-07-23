---
task: task_14
spec: 0046-public-context-driven-baseline-command
status: pending
type: backend
complexity: high
---

# Task 14: Synchronize canonical Baseline assets through Go

## Overview

Port maintained setup-snapshot synchronization to
`roundfix baseline assets sync`. Maintainers can check or refresh the embedded
catalog source from explicit immutable provenance without reintroducing a
runtime dependency on the setup skill.

## Requirements

1. MUST preserve the maintained explicit source directory, check-only,
   provenance, immutable commit, portable path, snapshot digest, output, and
   refusal contracts.
2. MUST expose the TechSpec command with text/JSON results and stable exit
   categories under the Baseline command family.
3. MUST verify source Git identity, clean committed source bytes, catalog
   compatibility, and normalized tree identity before any refresh.
4. MUST make check mode read-only and use the recoverable transaction for a
   non-empty refresh.
5. MUST update only Go-owned canonical Baseline assets and must not restore an
   executable engine to the skill.
6. MUST reproduce every maintained exact, semantic, and designed-delta sync
   row in the compatibility matrix.

## Subtasks

- [ ] Port source provenance and snapshot comparison.
- [ ] Implement check-only drift reporting.
- [ ] Implement confirmed canonical asset refresh through the transaction.
- [ ] Expose stable text/JSON command output and help.
- [ ] Add provenance, drift, rollback, and compatibility tests.

## Acceptance Criteria

- [ ] Check mode reports current or drifted state with zero writes.
- [ ] Dirty, mutable, uncommitted, incompatible, or non-portable sources fail before refresh.
- [ ] A successful refresh produces the expected canonical tree digest.
- [ ] Failure during refresh restores the complete asset preimage.
- [ ] Runtime catalog loading remains independent of the installed skill.
- [ ] Maintained Python sync cases have Go parity evidence.

## Context

- instruction: `docs/adr/0072-baseline-go-cutover-preserves-python-contracts.md`
- interface: `.agents/skills/setup-context-driven/scripts/context_setup.py`
- interface: `.agents/skills/setup-context-driven/assets/setups`
- interface: `Makefile`

## Verification

- `rtk go test -count=1 ./internal/baseline ./internal/cli -run 'TestBaselineAssetsSync|TestAssetsSyncCheck|TestAssetsSyncProvenance|TestAssetsSyncRollback|TestAssetsSyncCompatibility'` — expected: check, refresh, proof, transaction, and parity cases pass.
- `rtk go run -buildvcs=false ./cmd/roundfix baseline assets sync --help` — expected: help matches the approved source, check, and format contract.
- `rtk make verify` — expected: the full repository gate passes.

## References

- `_prd.md` → Goals 4–5; Core Features 19–21; Success Metrics.
- `_techspec.md` → System Architecture: embedded assets; API Contracts: maintenance operations; Build Order 8.
- ADR-0072 → Go destination for maintained sync-setups behavior.
