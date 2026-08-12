---
task: task_06
spec: 0051-doctor-readiness-contract-reconciliation
status: completed
type: chore
complexity: medium
---

# Task 06: Reconcile the derived Baseline skill snapshot

## Overview

Repair the deterministic QA failure caused by Task 05's canonical Roundfix
Skill edit. Refresh the TypeScript/Bun setup snapshot's repository-owned skill
digest and its complete derived catalog/parity chain. This is a bounded
generated-artifact reconciliation, not a Baseline behavior change.

## Requirements

1. MUST derive the Roundfix `contentDigest` from the current canonical
   `.agents/skills/roundfix` tree rather than copying the stale QA value.
2. MUST update the enclosing TypeScript/Bun setup digest and regenerate every
   directly affected catalog and parity authority.
3. MUST limit implementation mutations to
   `internal/baseline/assets/setups/typescript-bun.json`,
   `internal/baseline/testdata/catalog.digest`,
   `internal/baseline/testdata/catalog.normalized.json`,
   `internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json`,
   `internal/baseline/testdata/parity-corpus/v1/manifest.json`, and this Task
   file.
4. MUST NOT change Baseline implementation, schemas, source setup authority,
   any other Baseline asset, either Roundfix Skill copy, an upstream-managed
   skill, `skills-lock.json`, or `skills/recommended.txt`.
5. MUST preserve the failed QA report as evidence and leave final flow
   validation to a fresh QA run.

## Subtasks

- [x] Reproduce the authorial snapshot failure from the failed QA evidence.
- [x] Compute the canonical Roundfix Skill digest and refresh the
      TypeScript/Bun setup snapshot.
- [x] Regenerate the normalized catalog/digest and parity fixture/manifest.
- [x] Prove the complete derived chain and exact changed-file allowlist.

## Acceptance Criteria

- [x] The TypeScript/Bun setup snapshot records the current canonical Roundfix
      Skill digest.
- [x] The enclosing setup digest, normalized catalog/digest, parity fixture,
      and parity manifest are internally consistent.
- [x] The focused authorial snapshot and Baseline compatibility tests pass.
- [x] No runtime code, source setup, lock authority, upstream-managed skill,
      or out-of-allowlist artifact changes.
- [x] The full repository verification gate passes.

## Context

- instruction: `docs/agents/agent-instructions.md`
- instruction: `docs/agents/go.md`
- instruction: `.agents/skills/no-workarounds/SKILL.md`
- instruction: `.agents/skills/systematic-debugging/SKILL.md`
- interface: `skills/baseline_skill_contract_test.go`
- interface: `internal/baseline/catalog_test.go`
- interface: `internal/baseline/compatibility_corpus_test.go`
- interface: `internal/baseline/assets/setups/typescript-bun.json`
- interface: `docs/specs/0051-doctor-readiness-contract-reconciliation/qa/qa-report-2026-07-26.md`

## Verification

- `rtk go test ./skills -run 'TestAuthorialSkillSync/typescript-bun.json' -count=1` — expected: the canonical Roundfix Skill digest matches the setup snapshot.
- `rtk go test ./internal/baseline -run 'Test(CatalogCompatibility|AssetsSyncCompatibilityMatchesMaintainedPythonContract|BaselineCompatibilityCorpus)' -count=1` — expected: catalog and parity authorities accept the regenerated chain.
- `rtk make verify` — expected: the repository gate passes.
- `rtk git -c core.fsmonitor=false status --porcelain | rtk awk '{path=substr($0,4); if (path != "internal/baseline/assets/setups/typescript-bun.json" && path != "internal/baseline/testdata/catalog.digest" && path != "internal/baseline/testdata/catalog.normalized.json" && path != "internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json" && path != "internal/baseline/testdata/parity-corpus/v1/manifest.json" && path != "docs/specs/0051-doctor-readiness-contract-reconciliation/task_06.md") {print; bad=1}} END {exit bad}'` — expected: no out-of-allowlist path.

## References

- `_prd.md` → Core Features 7; Success Metrics.
- `_techspec.md` → Documentation and skill synchronization; Build Order 6.
- `qa/qa-report-2026-07-26.md` → failed static-gate evidence.

## Result

Reconciled the TypeScript/Bun Baseline snapshot and every directly derived
catalog and parity identity from the current canonical Roundfix Skill bytes.

### Changes

- Derived the canonical Roundfix Skill digest as
  `1e4ea59e0d6e553e42c6c63e16ad95165a78be8bbf31b8e0cd8b56ce13cc9146`
  and recorded it in the TypeScript/Bun setup snapshot.
- Recomputed the enclosing setup digest as
  `48592a566d3be7589f00e6895a5d2edb4bd54b59d3e92d171b333558e936c5d3`
  and refreshed the normalized catalog plus its domain-separated digest.
- Recomputed the asset-sync fixture's TypeScript/Bun digest and refreshed the
  parity manifest identity for the changed fixture.
- Corrected the Task-local compatibility regex after the declared names
  selected no tests, and disabled Git fsmonitor in the Task-local allowlist
  command after the worktree daemon emitted a diagnostic.

### Verification

- Pre-change
  `rtk go test ./skills -run 'TestAuthorialSkillSync/typescript-bun.json' -count=1`
  reproduced the expected `d5269b...f59f` versus `1e4ea5...9146` mismatch.
- The same authorial command passed after reconciliation: 2 test results in
  one package.
- The original Baseline regex exited zero with `No tests found`; it received
  no compatibility credit. The corrected
  `rtk go test ./internal/baseline -run 'Test(CatalogCompatibility|AssetsSyncCompatibilityMatchesMaintainedPythonContract|BaselineCompatibilityCorpus)' -count=1`
  passed all three catalog, asset-sync, and parity owners.
- Independent digest recomputation matched the recorded setup digest,
  `sha256:b4328a2c299dda6a312b47098babc2e92f68c910fe3c610e8e9522fe7628071a`
  catalog digest, asset-sync setup digest
  `886c338c717ff084675db5563d60ade18b54b37dad7f30c10a817e4cf464fbbe`,
  and the fixture's 83,662-byte
  `1166f01b5f2f5c0f0fe361813e065fd9aa415d0633101fb6348903071644d24f`
  manifest identity.
- `rtk make verify` passed 2,420 tests across 23 packages, the focused
  synchronization guards, shipped Skill validation, and the CLI build.
- `rtk git -c core.fsmonitor=false status --porcelain` and both unstaged and
  staged diff inspection found exactly the five authorized generated artifacts
  plus this Task file; the allowlist filter and `rtk git diff --check` exited
  zero.

### Acceptance evidence

- The passing authorial test proves the snapshot's Roundfix `contentDigest`
  was derived from the current `.agents/skills/roundfix` tree.
- Matching recorded and independently recomputed identities prove the setup,
  normalized catalog/digest, parity fixture, and manifest form one consistent
  derived chain.
- The three focused compatibility owners and the complete repository gate
  passed after the last generated-artifact edit.
- Git postflight proves no runtime code, source setup, lock authority,
  upstream-managed skill, failed QA report, or out-of-allowlist artifact
  changed.

### Follow-ups

None.
