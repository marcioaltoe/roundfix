---
task: task_14
spec: 0046-public-context-driven-baseline-command
status: completed
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

- [x] Port source provenance and snapshot comparison.
- [x] Implement check-only drift reporting.
- [x] Implement confirmed canonical asset refresh through the transaction.
- [x] Expose stable text/JSON command output and help.
- [x] Add provenance, drift, rollback, and compatibility tests.

## Acceptance Criteria

- [x] Check mode reports current or drifted state with zero writes.
- [x] Dirty, mutable, uncommitted, incompatible, or non-portable sources fail before refresh.
- [x] A successful refresh produces the expected canonical tree digest.
- [x] Failure during refresh restores the complete asset preimage.
- [x] Runtime catalog loading remains independent of the installed skill.
- [x] Maintained Python sync cases have Go parity evidence.

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

## Result

Implemented `roundfix baseline assets sync` as the Go-owned maintainer
operation for canonical Setup Snapshot synchronization. The command accepts an
explicit canonical setups directory, proves its clean immutable Git identity,
normalizes complete safe skill trees, validates the resulting catalog in
memory, and reports current or drifted state without writing in `--check`
mode.

A non-empty refresh updates only
`internal/baseline/assets/setups/*.json` through the Git-private recoverable
Baseline transaction. The runtime catalog remains compiled through `go:embed`;
the command does not read or restore an executable engine from the installed
setup skill.

### Acceptance evidence

1. `TestAssetsSyncCheckIsReadOnlyAndReportsDrift` captures the complete
   canonical asset tree before and after check mode, reports all three drifted
   snapshots, and proves byte-identical post-state.
2. `TestAssetsSyncProvenanceAndPreMutationRefusals` covers dirty or untracked
   source bytes, non-GitHub identity, mutable refs, non-portable paths, and
   incompatible empty setup input. Every refusal preserves the complete asset
   preimage.
3. `TestBaselineAssetsSyncRefreshProducesCanonicalTreeAndIsIdempotent`
   refreshes all setup snapshots, loads the resulting catalog, verifies each
   source commit and lowercase snapshot digest, preserves the Roundfix-owned
   `setup-context-driven` content digest, and proves a second refresh is
   byte-idempotent.
4. `TestAssetsSyncRollbackRestoresCompleteAssetPreimage` injects a failure
   after the first replacement and proves the recoverable transaction restores
   every original asset byte.
5. `TestAssetsSyncCompatibilityMatchesMaintainedPythonContract` reconstructs
   the exact maintained `asset-sync` parity fixture with the current immutable
   commit and compares every generated Setup Snapshot byte-for-byte. It also
   loads the embedded catalog independently of filesystem skill state.
6. `TestBaselineAssetsSyncCommand` covers public dispatch, help, required
   source input, check mode, text and JSON rendering, stdout/stderr separation,
   and stable exit categories.

### Verification

- `rtk go test -count=1 ./internal/baseline ./internal/cli -run 'TestBaselineAssetsSync|TestAssetsSyncCheck|TestAssetsSyncProvenance|TestAssetsSyncRollback|TestAssetsSyncCompatibility'`
  passed with 14 tests across 2 packages.
- `rtk go run -buildvcs=false ./cmd/roundfix baseline assets sync --help`
  passed and printed the source, check-only, format, provenance, transaction,
  and exit-category contract.
- `rtk make verify` passed: 1,990 Go tests, both 256-test maintained Python
  suites, Baseline asset validation, Roundfix skill checks, and the binary
  build.
- `rtk git diff --check` passed.
