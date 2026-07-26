---
task: task_06
spec: 0051-doctor-readiness-contract-reconciliation
status: pending
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

- [ ] Reproduce the authorial snapshot failure from the failed QA evidence.
- [ ] Compute the canonical Roundfix Skill digest and refresh the
      TypeScript/Bun setup snapshot.
- [ ] Regenerate the normalized catalog/digest and parity fixture/manifest.
- [ ] Prove the complete derived chain and exact changed-file allowlist.

## Acceptance Criteria

- [ ] The TypeScript/Bun setup snapshot records the current canonical Roundfix
      Skill digest.
- [ ] The enclosing setup digest, normalized catalog/digest, parity fixture,
      and parity manifest are internally consistent.
- [ ] The focused authorial snapshot and Baseline compatibility tests pass.
- [ ] No runtime code, source setup, lock authority, upstream-managed skill,
      or out-of-allowlist artifact changes.
- [ ] The full repository verification gate passes.

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
- `rtk go test ./internal/baseline -run 'TestEmbeddedCatalogCompatibility|TestCompatibilityCorpus' -count=1` — expected: catalog and parity authorities accept the regenerated chain.
- `rtk make verify` — expected: the repository gate passes.
- `rtk git status --porcelain | rtk awk '{path=substr($0,4); if (path != "internal/baseline/assets/setups/typescript-bun.json" && path != "internal/baseline/testdata/catalog.digest" && path != "internal/baseline/testdata/catalog.normalized.json" && path != "internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json" && path != "internal/baseline/testdata/parity-corpus/v1/manifest.json" && path != "docs/specs/0051-doctor-readiness-contract-reconciliation/task_06.md") {print; bad=1}} END {exit bad}'` — expected: no out-of-allowlist path.

## References

- `_prd.md` → Core Features 7; Success Metrics.
- `_techspec.md` → Documentation and skill synchronization; Build Order 6.
- `qa/qa-report-2026-07-26.md` → failed static-gate evidence.

## Result

Pending.
