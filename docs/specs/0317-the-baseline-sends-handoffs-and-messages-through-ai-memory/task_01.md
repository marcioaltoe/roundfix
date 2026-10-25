---
task: task_01
spec: 0317-the-baseline-sends-handoffs-and-messages-through-ai-memory
status: pending
type: backend
complexity: high
---

# Task 01: Every upstream setup follows the ai-memory commit and the core module dispatches its skills

## Overview

Upstream replaced `handoff` with six `ai-memory-*` skills at
`65d5f365cace27c4b468db1723213723b10c870b`. The catalog refuses a sync that
drops a skill a module still requires, so this Task first swaps `handoff`
for `ai-memory-handoff` and `ai-memory-messaging` in the core module, then
imports every upstream setup through the asset sync, the only writer of a
Setup Snapshot. The parity row, the derived files, this
repository's Setup Manifest, its changed upstream skills, its lock and its
recommended list follow through their sanctioned commands (`_prd.md` → Core
Feature 1; ADR-0330).

## Requirements

1. MUST answer the Backlog Entry
   [`references/2026-10-24-the-baseline-sends-handoffs-and-messages-through-ai-memory.md`](references/2026-10-24-the-baseline-sends-handoffs-and-messages-through-ai-memory.md)
   of 2026-10-24, its expected outcome 1.
2. MUST first check that the local checkout `~/dev/skills` holds commit
   `65d5f365cace27c4b468db1723213723b10c870b`; when it does not, MUST stop and
   report before writing anything. MUST NOT fetch from or otherwise reach a
   network remote.
3. MUST edit `internal/baseline/assets/modules/core.json` as text, keeping
   every other byte and the file's layout, exactly as `_techspec.md` → Data
   Model 1 states, and then run the Module Version Record command declared
   in `_authorization.md`.
4. MUST then run, verbatim and in order, the commands of `_techspec.md` → API
   Contract 1: the fetch into an empty repository, the binary build, the
   sync, the `--check` (exit `0`), the parity transform of
   `internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json`,
   `make baseline-digests`, the managed refresh twice (the second reports
   `File changes: 0`), the twelve-skill restore with its confirmed plan
   digest, the lock reconcile with its confirmed plan digest, the deletion of
   `.agents/skills/handoff/`, and the rewrite of `skills/recommended.txt`
   from the lock. Every command that runs the binary gets a fresh temporary
   `HOME`; the binary is built under the Run's own environment and rebuilt
   after the sync so the refresh, the restore and the reconcile read the new
   snapshots.
5. MUST hand-edit no snapshot, derived file, parity digest, skill file or
   `skills-lock.json` entry. A file these commands change that this Task does
   not declare is a stop-and-report condition.
6. MUST make the declared break of `_techspec.md` → Testing Approach (Build
   Order 1) in `internal/baseline/assets_sync_upstream_rename_test.go` and no
   other change to existing tests.
7. MUST create `internal/baseline/ai_memory_setup_test.go` with
   `TestEveryUpstreamSetupFollowsTheAiMemoryRevision` and
   `TestTheCoreModuleDispatchesTheAiMemorySkills` as `_techspec.md` → Testing
   Approach describes them. Each top-level test calls `t.Parallel()` first.
8. MUST NOT change any other Baseline module, profile source, template,
   coverage entry, Go production file, owned skill, user guide page,
   `.roundfixrc.yml`, the Makefile or `go.mod`/`go.sum`.
   `.agents/skills/golang-cli/assets/examples/go.mod` is upstream skill
   content the restore writes, not this repository's module.

## Subtasks

- [ ] Check the commit, swap the core module's skills and record the module.
- [ ] Sync and check, transform the parity row, regenerate the digests and refresh twice.
- [ ] Restore the twelve skills, reconcile the lock, delete the `handoff` tree and rewrite the recommended list.
- [ ] Make the declared test change and add the two tests.

## Acceptance Criteria

- [ ] Every upstream snapshot names `65d5f36` and holds the six ai-memory skills and no `handoff`.
- [ ] The core module requires and dispatches `ai-memory-handoff` and `ai-memory-messaging`; no module names `handoff`.
- [ ] This repository holds every required skill at the refreshed `go` snapshot, keeps no `handoff` copy or lock entry, and a second managed refresh reports no file change.

## Consumers

- the refreshed upstream snapshots and the core module's dispatch — consumers: `./internal/baseline`, `./skills`

## Context

- instruction: `docs/adr/0330-the-secondbrain-is-read-only-and-ai-memory-carries-handoffs-and-messages.md`
- instruction: `internal/baseline/assets_sync.go`
- instruction: `internal/baseline/skills_snapshot_profile_test.go`
- instruction: `internal/baseline/repository_skills_snapshot_test.go`
- creates: `internal/baseline/ai_memory_setup_test.go`
- creates: `.agents/skills/ai-memory-durable-pages/SKILL.md`
- creates: `.agents/skills/ai-memory-durable-pages/LICENSE`
- creates: `.agents/skills/ai-memory-handoff/SKILL.md`
- creates: `.agents/skills/ai-memory-handoff/LICENSE`
- creates: `.agents/skills/ai-memory-learning-maintenance/SKILL.md`
- creates: `.agents/skills/ai-memory-learning-maintenance/LICENSE`
- creates: `.agents/skills/ai-memory-messaging/SKILL.md`
- creates: `.agents/skills/ai-memory-messaging/LICENSE`
- creates: `.agents/skills/ai-memory-retrieval/SKILL.md`
- creates: `.agents/skills/ai-memory-retrieval/LICENSE`
- creates: `.agents/skills/ai-memory-routing-install/SKILL.md`
- creates: `.agents/skills/ai-memory-routing-install/LICENSE`
- creates: `.agents/skills/golang-cli/assets/examples/go.mod`
- deletes: `.agents/skills/handoff/SKILL.md`
- deletes: `.agents/skills/handoff/agents/openai.yaml`
- deletes: `.agents/skills/context7-cli/references/skills.md`
- interface: `internal/baseline/assets/modules/core.json`
- interface: `internal/baseline/module-versions.json`
- interface: `internal/baseline/assets/setups/go.json`
- interface: `internal/baseline/assets/setups/rust.json`
- interface: `internal/baseline/assets/setups/typescript.json`
- interface: `internal/baseline/assets/setups/go-cli-typescript-bun.json`
- interface: `internal/baseline/assets/profiles/standard-typescript-monorepo.json`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/skill-dispatch.md`
- interface: `internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json`
- interface: `internal/baseline/testdata/parity-corpus/v1/manifest.json`
- interface: `internal/baseline/testdata/catalog.diagnostics.golden.json`
- interface: `internal/baseline/testdata/catalog.digest`
- interface: `internal/baseline/testdata/catalog.normalized.json`
- interface: `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`
- interface: `docs/agents/skill-dispatch.md`
- interface: `docs/agents/setup-context.json`
- interface: `skills-lock.json`
- interface: `skills/recommended.txt`
- interface: `.agents/skills/context7-cli/SKILL.md`
- interface: `.agents/skills/context7-cli/references/docs.md`
- interface: `.agents/skills/context7-cli/references/setup.md`
- interface: `.agents/skills/grilling/SKILL.md`
- interface: `.agents/skills/knowledge-workspace/SKILL.md`
- interface: `.agents/skills/systematic-debugging/SKILL.md`
- interface: `.agents/skills/typesafe-ai/SKILL.md`
- interface: `internal/baseline/assets_sync_upstream_rename_test.go`

## Verification

- `for name in TestEveryUpstreamSetupFollowsTheAiMemoryRevision TestTheCoreModuleDispatchesTheAiMemorySkills; do go test -list "^${name}\$" ./internal/baseline | grep -qx "$name" || { printf 'missing test %s\n' "$name" >&2; exit 1; }; done; go test -count=1 ./internal/baseline && go test -count=1 ./skills` — expected: exit 0. Before this Task neither test exists, so the command fails at its first check. Afterwards both exist and both packages pass: the snapshots agree on every external skill, this repository holds its required skills at the refreshed `go` snapshot, the parity row and the derived-file characterizations are current, the rename fixture follows `ai-memory-handoff`, and the recommended list equals the lock.
- `for f in go rust typescript; do grep -qF -- '"ref": "65d5f365cace27c4b468db1723213723b10c870b"' "internal/baseline/assets/setups/$f.json" || { printf '%s not at the ai-memory commit\n' "$f" >&2; exit 1; }; if grep -qF -- 'b3c45a45f1bccd3b33aaecaaa22947d942f2fc02' "internal/baseline/assets/setups/$f.json"; then printf '%s keeps the old commit\n' "$f" >&2; exit 1; fi; done; test -f .agents/skills/ai-memory-messaging/SKILL.md || { echo 'ai-memory-messaging not restored' >&2; exit 1; }; if test -e .agents/skills/handoff; then echo 'handoff tree kept' >&2; exit 1; fi; if grep -qF -- '"handoff":' skills-lock.json; then echo 'handoff lock entry kept' >&2; exit 1; fi; grep -qx 'ai-memory-handoff' skills/recommended.txt || { echo 'recommended list not updated' >&2; exit 1; }; go test -count=1 -run '^(TestEmbeddedSetupSnapshotsAgreeOnEveryExternalSkill|TestThisRepositoryHoldsItsRequiredSkillsAtTheSetupSnapshot|TestAssetSyncFollowsASkillRenamedUpstream)$' ./internal/baseline` — expected: exit 0. Before this Task the snapshots name `b3c45a45`, so the command fails at its first check. Afterwards every upstream snapshot names only the new commit, the ai-memory skills are restored, `handoff` has no tree and no lock entry, the recommended list names the ai-memory skills, and the agreement, required-skill and rename tests pass.

## References

- `_prd.md` → Goals 1-2; Core Feature 1; Success Metrics 1 and 5; User Story 1
- `_techspec.md` → Current behavior 1-3 and 8; Data Model 1; API Contracts 1-2; Invariants 1 and 6; Testing Approach; Build Order 1
- ADR-0330
