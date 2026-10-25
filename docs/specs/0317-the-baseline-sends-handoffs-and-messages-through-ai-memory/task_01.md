---
task: task_01
spec: 0317-the-baseline-sends-handoffs-and-messages-through-ai-memory
status: completed
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


## Result

Implemented the upstream setup migration and the core skill swap for this
Task's slice. Task status is unchanged; declared Verification and repository
Verification remain Daemon-owned.

### Implementation and acceptance evidence

1. Every upstream snapshot follows the ai-memory revision:
   - Before any write, the local `~/dev/skills` commit check
     `git -C /Users/marcio/dev/skills cat-file -e 65d5f365cace27c4b468db1723213723b10c870b^{commit}`
     exited 0. The starting catalog still required/dispatched `handoff` and
     lacked the six ai-memory entries.
   - Fetched that commit from the local checkout into an empty temporary
     repository, checked it out detached and added its provenance remote
     without contacting that remote. Built the binary under the Run's
     environment; each binary invocation used a fresh temporary `HOME`.
   - Asset sync exited 0 and reported exactly four
     `skills.setup-snapshot.updated` findings; the immediately following
     `--check` exited 0 with `setup-context-driven audit: ok`.
   - Applied the exact API Contract 1 parity transform through a temporary
     copy, excluding composed snapshots; `make baseline-digests` exited 0
     and reported `baseline-digests: regenerated`. No generated digest,
     snapshot or golden was edited by hand.
   - Added `TestEveryUpstreamSetupFollowsTheAiMemoryRevision`: every setup,
     including the composed setup, contains the six ai-memory skills and
     excludes `handoff`, `brand-guidelines` and `gws-*`; every GitHub setup
     and skill pins the declared full revision. Its focused run passed.
2. Core requires and dispatches the ai-memory skills:
   - Applied only the Data Model 1 text substitutions to `core.json`,
     preserving the remaining layout and bytes. Ran
     `go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1`;
     exit 0. The recorder advanced core from version 21 to 22 and wrote the
     module record.
   - Added `TestTheCoreModuleDispatchesTheAiMemorySkills`: no module
     requires or dispatches `handoff`; core requires both new skills with
     their exact trigger IDs and the MCP-availability suffix. Its focused
     run passed. Both new top-level tests call `t.Parallel()` first.
   - The sole existing-test edit moves the upstream rename fixture from
     `handoff` → `handoff-next` to `ai-memory-handoff` →
     `ai-memory-handoff-next`, including its module, directory, list and
     messages. The focused rename/refusal cases passed.
3. Repository skills and idempotent refresh:
   - Rebuilt the binary after sync/regeneration and ran the prescribed
     managed refresh twice. Both exited 0; the first changed only
     `docs/agents/setup-context.json` and `docs/agents/skill-dispatch.md`,
     and the second reported `File changes: 0` and `Idempotence: verified`.
     Both retained the two reported nested instruction carriers unchanged.
   - Reviewed the twelve-skill restoration preview (expected exit 3), then
     confirmed plan
     `c4ccee9bc42e0d93b0e56c7c26652c91cf060b1f6d2a15098b5e3434f9365a20`;
     exit 0, `restore.completed`. All skill bytes and lock entries came
     from the sanctioned restore.
   - Reviewed the reconciliation preview (expected exit 3), which listed
     only `remove-lock-entry skills-lock.json [handoff]`, then confirmed
     plan `76a7e63c2c96477b10e5b8268d4b618c1f42dd390be2aca1086530e128b2d876`;
     exit 0, `reconcile.completed`. Ran the prescribed
     `git rm -r .agents/skills/handoff` and derived `skills/recommended.txt`
     with `jq -r '.skills | keys[]' skills-lock.json`.
   - Focused snapshot-agreement, installed-required-skill and recommended
     list checks passed. A read-only postflight confirmed no retired tree
     or lock entry and exact sorted recommended-list equality.

### Focused checks

- `go test ./internal/baseline -run '^(TestEveryUpstreamSetupFollowsTheAiMemoryRevision|TestTheCoreModuleDispatchesTheAiMemorySkills|TestAssetSync.*Rename.*|TestAssetSyncStillRefusesAnInvalidCatalogWithoutDrift|TestEmbeddedSetupSnapshotsAgreeOnEveryExternalSkill|TestThisRepositoryHoldsItsRequiredSkillsAtTheSetupSnapshot)$' -count=1`
  — exit 0 (`ok roundfix/internal/baseline`).
- `go test ./skills -run '^TestRecommended(SkillsMatchLock|ListsExternallyManagedSkills)$' -count=1`
  — exit 0 (`ok roundfix/skills`).
- `git diff --check` — exit 0. Changed-path audit included staged,
  unstaged and untracked files: 47 paths, zero outside this Task's exact
  declared interfaces/creates/deletes plus this Task file.

### Execution limits and follow-ups

The configured sandbox refused process creation because its writable root
contains a symlink, even with the resolved checkout as the working directory.
Authorized commands ran through approved escalated execution. No network
remote was fetched, no real Roundfix home or Secondbrain was read or written,
and no ai-memory or provider call was made.

No follow-up implementation was added. No other Task file or Task Graph
manifest was edited. No commit, push or pull request was created. The two
retired handoff files are staged deletions because the Spec explicitly
requires `git rm`; other edits remain unstaged/untracked for the Daemon.

## Carry-forward provenance

- Source Run: `run_20261025T133232Z_0a2b2319e6c91868`
- Source commit: `e0f0bb0a49de0ffffe5156367a22373aacf21a9a`
