---
task: task_02
spec: 0350-the-vendored-catalog-skills-follow-the-catalog-through-a-pinned-bump
status: pending
type: backend
complexity: high
---

# Task 02: The vendored catalog skills and every upstream setup follow catalog commit c8b2e85

## Overview

Run the Catalog Skill Bump once at `c8b2e8511a7ff89dbe9e87a56cc3c29a6acbcb84`
from the operator's clone, so the 33 catalog entries, their folders and the
five upstream Setup Snapshots name that commit and `systematic-debugging`
loses the five files the catalog removed. Make the setup tests read the
Catalog Revision from the lock and refuse the trees earlier fixes replaced,
and add the documentation contract that holds every vendored folder to the
tree its snapshots pin (`_techspec.md` → API Contracts 4 and 5).

## Requirements

1. MUST answer the Backlog Entry
   [`references/2026-11-04-the-vendored-catalog-skills-trail-the-catalog-and-name-branch-refs.md`](references/2026-11-04-the-vendored-catalog-skills-trail-the-catalog-and-name-branch-refs.md),
   its Expected section.
2. MUST first check, read-only, that the clone
   `/Users/marcio/.roundfix-operator/scratch/spec-0350/skills` holds the
   commit (`git -C /Users/marcio/.roundfix-operator/scratch/spec-0350/skills cat-file -e c8b2e8511a7ff89dbe9e87a56cc3c29a6acbcb84^{commit}`);
   when it does not, MUST stop and report before writing anything. MUST NOT
   check out, fetch into or otherwise write that clone, MUST NOT read or
   write `/Users/marcio/dev/skills`, and MUST NOT reach a network remote.
3. MUST run, once, `GIT_ALLOW_PROTOCOL=file make skills-third-party-bump SKILLS_CATALOG=/Users/marcio/.roundfix-operator/scratch/spec-0350/skills SKILLS_REF=c8b2e8511a7ff89dbe9e87a56cc3c29a6acbcb84`
   and require exit `0` and the final line `skills-third-party-bump: 33
   catalog skills and every upstream Setup Snapshot pinned to
   c8b2e8511a7ff89dbe9e87a56cc3c29a6acbcb84`; then run it a second time and
   require that `git status --porcelain` is the same after it as before it.
   A refusal, a restore preview that plans anything but the 33 lock updates
   and the five `systematic-debugging` removals, or a changed file this Task
   does not declare is a stop-and-report condition. When the restore no
   longer writes each entry's `ref` as the snapshot commit (Spec 0349 owns
   that code), stop and report.
4. MUST hand-edit no snapshot, vendored folder, lock entry, derived file or
   Setup Manifest value.
5. MUST change the tests as `_techspec.md` → Interfaces and API Contract 5
   state: replace the constant `upstreamSetupRevision` with the function that
   reads the one full commit every catalog entry of `skills-lock.json` names
   and fails on a branch, a short or mixed ref or no entry; update every call
   site in `upstream_setup_revision_test.go`, `ai_memory_setup_test.go` and
   `python_setup_test.go`; replace `onioncryTreeDigest`,
   `uiUxProMaxTreeDigest` and `knowledgeWorkspaceTreeDigest` with the three
   replaced-tree constants of Interfaces, each test keeping its name and its
   revision check and asserting inequality with a message naming the Spec
   that replaced the tree; and update each file's `Invariant:` header line to
   say so.
6. MUST create `internal/baseline/vendored_catalog_skills_docscontract_test.go`
   under `//go:build docscontract` and `//verify:always`, with the header
   comment lines of `internal/cli/live_documents_docscontract_test.go`, and
   `TestEveryVendoredCatalogSkillIsTheTreeItsCatalogCommitPublishes`, which
   calls `t.Parallel()` first and asserts, for every `skills-lock.json` entry
   in sorted order, that it is not a `roundskills.Names()` skill, that an
   embedded Setup Snapshot pins a `treeDigest` for it and that
   `portableRestoreDigest` of its `.agents/skills` folder (read with
   `inspectRestoreTarget`) equals that digest, naming the skill and
   `make skills-third-party-bump` on a difference; and that every `github`
   snapshot names `upstreamSetupRevision(t)`.
7. MUST make `TestRepositoryContractTestsDeclareTheirRelevance` accept a
   `docscontract` file of `internal/baseline` as `always`, as it accepts
   those of `internal/cli`, and count three such files.

## Subtasks

- [ ] Check the clone holds the commit, run the bump, and run it again with no change.
- [ ] Make the revision tests read the lock and the tree tests refuse the replaced trees.
- [ ] Add the vendored skill contract and admit its file in the relevance test.

## Acceptance Criteria

- [ ] The 33 catalog entries and every upstream snapshot name `c8b2e85…`; `systematic-debugging` has none of the five removed files.
- [ ] A second bump changes nothing.
- [ ] No test names `86df860…` or a tree the snapshots carry; the six setup revision and tree tests and the new contract pass.
- [ ] A hand edit of a vendored folder with a re-hashed lock entry fails the new contract.

## Context

- instruction: `Makefile`
- instruction: `docs/adr/0360-the-vendored-catalog-skills-follow-the-catalog-through-one-pinned-commit.md`
- instruction: `docs/adr/0224-this-repository-holds-its-upstream-skills-to-the-snapshot-and-the-lock-not-to-a-pinned-digest.md`
- instruction: `internal/baseline/skills_restore.go`
- instruction: `internal/baseline/skills_trailing.go`
- instruction: `internal/baseline/repository_skills_snapshot_test.go`
- instruction: `internal/cli/live_documents_docscontract_test.go`
- interface: `skills-lock.json`
- interface: `internal/baseline/assets/setups/go-cli-typescript-bun.json`
- interface: `internal/baseline/assets/setups/go.json`
- interface: `internal/baseline/assets/setups/python.json`
- interface: `internal/baseline/assets/setups/rust.json`
- interface: `internal/baseline/assets/setups/typescript.json`
- interface: `internal/baseline/testdata/catalog.diagnostics.golden.json`
- interface: `internal/baseline/testdata/catalog.digest`
- interface: `internal/baseline/testdata/catalog.normalized.json`
- interface: `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`
- interface: `docs/agents/setup-context.json`
- interface: `internal/baseline/upstream_setup_revision_test.go`
- interface: `internal/baseline/ai_memory_setup_test.go`
- interface: `internal/baseline/python_setup_test.go`
- interface: `internal/verifyselect/contracts_repository_test.go`
- creates: `internal/baseline/vendored_catalog_skills_docscontract_test.go`
- deletes: `.agents/skills/systematic-debugging/CREATION-LOG.md`
- deletes: `.agents/skills/systematic-debugging/test-academic.md`
- deletes: `.agents/skills/systematic-debugging/test-pressure-1.md`
- deletes: `.agents/skills/systematic-debugging/test-pressure-2.md`
- deletes: `.agents/skills/systematic-debugging/test-pressure-3.md`

## Consumers

- the Catalog Revision in `skills-lock.json`, the refreshed Setup Snapshots and the restored vendored folders — consumers: `./internal/baseline`, `./skills`, `./internal/verifyselect`, `./internal/docscontract`

## Verification

- `go test -tags docscontract -list '^TestEveryVendoredCatalogSkillIsTheTreeItsCatalogCommitPublishes$' ./internal/baseline | grep -qx TestEveryVendoredCatalogSkillIsTheTreeItsCatalogCommitPublishes || { echo 'missing test TestEveryVendoredCatalogSkillIsTheTreeItsCatalogCommitPublishes' >&2; exit 1; }; grep -q '"ref": "c8b2e8511a7ff89dbe9e87a56cc3c29a6acbcb84"' skills-lock.json || { echo 'skills-lock.json does not name the catalog commit' >&2; exit 1; }; ! grep -rq 86df8608726979064a8fe46003dfaf403ffe0087 internal --include=*_test.go || { echo 'a test still names the previous catalog commit' >&2; exit 1; }; go test -timeout 20m -count=1 ./internal/baseline ./skills && go test -timeout 20m -count=1 -tags docscontract ./internal/baseline -run '^TestEveryVendoredCatalogSkillIsTheTreeItsCatalogCommitPublishes$' && go test -timeout 20m -count=1 ./internal/verifyselect -run '^TestRepositoryContractTestsDeclareTheirRelevance$' && go test -timeout 10m -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSetupManifestIsCurrent$'` — expected: exit 0. Before this Task the contract does not exist, so the command fails at its first check; afterwards the lock names the commit, no test names the old one, the revision and tree tests read the lock, the new contract holds all 33 folders to their snapshot trees, and the relevance and Setup Manifest contracts pass.

## References

- `_prd.md` → Goals; Core Features 2 and 3; Success Metrics 2, 3, 4 and 6
- `_techspec.md` → System Architecture; Interfaces; API Contracts 4 and 5; Invariants 1, 3, 4 and 6; Testing Approach; Build Order 2
- ADR-0360; ADR-0224; ADR-0264; ADR-0191
