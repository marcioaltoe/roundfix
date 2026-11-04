---
task: task_02
spec: 0350-the-vendored-catalog-skills-follow-the-catalog-through-a-pinned-bump
status: completed
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

## Result

Skills loaded: implement-task, coding-guidelines, golang-testing, testing-boss, tech-writer, evidence-gate, roundfix, systematic-debugging.

The generated catalog bump and this Task's test changes are present. The
Verification Feedback repair added the missing documentation contract and
updated the six setup tests and the relevance declaration. The second-bump
refusal described below remains a blocker in task_01's Makefile target;
that file is outside this Task's edit slice. Task status remains Daemon-owned.

### Implementation

- `upstreamSetupRevision(t)` reads the existing structured lock reader and
  requires every `marcioaltoe/skills` entry to name one full lowercase
  40-character commit. It rejects branch, short, mixed and absent refs.
  All call sites in the three declared setup test files use the helper.
- The three existing tree tests retain their names and revision checks and
  reject the trees replaced by Specs 0320, 0333 and 0326. Their constants are
  the replaced trees from the TechSpec; none pins a current snapshot tree.
  The three files' `Invariant:` headers describe the new contracts.
- The new `docscontract` test is `//verify:always`, calls `t.Parallel()`
  first, checks every lock entry in sorted order against the owned names
  and embedded tree digests, and compares each installed tree using
  `inspectRestoreTarget` and `portableRestoreDigest`. Differences name the
  skill and `make skills-third-party-bump`. It also checks every upstream
  snapshot and every `github` skill revision against the lock.
- The relevance test accepts both CLI and Baseline documentation contracts
  as always selected and requires three live documentation contract files.

### Bump evidence and unresolved refusal

- Before any repository write, the required read-only command
  `git -C /Users/marcio/.roundfix-operator/scratch/spec-0350/skills cat-file -e c8b2e8511a7ff89dbe9e87a56cc3c29a6acbcb84^{commit}`
  exited 0. The operator clone was only read; no network transport was used.
- The first
  `GIT_ALLOW_PROTOCOL=file make skills-third-party-bump SKILLS_CATALOG=/Users/marcio/.roundfix-operator/scratch/spec-0350/skills SKILLS_REF=c8b2e8511a7ff89dbe9e87a56cc3c29a6acbcb84`
  exited 0. Its last line was exactly:
  `skills-third-party-bump: 33 catalog skills and every upstream Setup Snapshot pinned to c8b2e8511a7ff89dbe9e87a56cc3c29a6acbcb84`.
- Inspection of the actual first restore preview found exactly 33
  `update-lock-entry` operations and the five declared `systematic-debugging`
  removals, with no other operations. Its Plan Digest was
  `26028d9f266ea7aa595fdcb8f385cf1b3d989ea95720a3e8fc536b0f1b333855`.
  The target's folder comparisons and lock-ref checks reached success.
- The second invocation of the same command exited 2 from `make`:
  restore reported no changes and returned 0, but the target unconditionally
  requires preview exit 3. Its refusal was
  `restore preview exited 0; expected 3 with a Plan Digest`.
  `git status --porcelain` captured immediately before and after the second
  invocation was byte-identical. The repeat did not reach the target's
  post-restore checks, regeneration or managed refresh. Requirement 3 caused
  the initial stop; Verification Feedback subsequently requested the missing
  Task-local implementation. No Makefile edit or third bump was attempted.
- Local transcripts: `/private/tmp/spec-0350-task02-bump-1.log` and
  `/private/tmp/spec-0350-task02-bump-2.log`. These temporary logs are
  supplementary; the relevant command results are recorded above.

### Focused checks

The Daemon's attempt-1 diagnostic artifact was read at
`/Users/marcio/.roundfix/artifacts/339f8dac2b687a04/runs/run_20261104T202718Z_36f747a05c5e57d3/verification/batch-002-attempt-1.log`.
It identifies the absent vendored contract. Its body is not reproduced here.

The following focused implementation check exited 0 for both packages:

```sh
GOPROXY=off GOSUMDB=off go test -count=1 -timeout 5m -tags docscontract ./internal/baseline ./internal/verifyselect -run '^(TestTheTypeScriptSetupsCarryTheRepublishedOnioncryTree|TestTheTypeScriptSetupsCarryTheCompletedUiUxProMaxTree|TestEverySetupCarriesTheKnowledgeWorkspaceThatRetiredTheSecondbrainInbox|TestEveryUpstreamSetupFollowsTheAiMemoryRevision|TestThePythonSetupIsThePublishedUpstreamPreset|TestEveryUpstreamSetupFollowsThePythonPresetRevision|TestEveryVendoredCatalogSkillIsTheTreeItsCatalogCommitPublishes|TestRepositoryContractTestsDeclareTheirRelevance)$'
```

For the negative probes, the updated baseline tests were compiled once with
`GOPROXY=off GOSUMDB=off go test -c -tags docscontract -o /private/tmp/spec-0350-task02-baseline.test ./internal/baseline`.
That binary ran only selected tests against a temporary Git fixture holding
copies of the lock, skills, Setup Manifest and repository profiles. No probe
edited a vendored folder or lock in the assigned worktree.

- Healthy fixture: the new contract passed for all 33 skills.
- Re-hashed vendored edit: appending text to the fixture's
  `ai-memory-retrieval/SKILL.md` and updating its `computedHash` caused the
  new contract to exit 1, naming `ai-memory-retrieval` and the bump command.
  The hash computation was first checked against the untouched lock value.
  The existing profile-required snapshot test still passed on this fixture,
  demonstrating the gap the new contract closes. Restoring the original
  fixture bytes made the new contract pass again.
- Lock helper: branch `main`, short ref, a different full commit in one
  entry, and no catalog entry each made the selected Python setup test exit
  1 with the corresponding helper diagnostic. The actual failing transcripts
  were inspected; no suite-guard boundary violation occurred.
- Supplementary probe logs:
  `/private/tmp/spec-0350-task02-probes-oxg9_8qf/`.
- Source inspection found no previous catalog commit in internal test files
  and no overlap between the three replaced-tree constants and current
  snapshot tree digests. `gofmt -l` reported no changed test file.

### Acceptance evidence and remaining limitation

1. JSON inspection confirms all 33 catalog lock entries, all five upstream
   snapshots and every `github` skill they carry name the requested full
   commit. The five removed files are absent. The first successful target
   compared every installed folder with its fetched catalog tree. No owned
   skill changed.
2. The second bump preserved the exact porcelain status, but refused rather
   than finishing successfully. The task_01 target needs to handle a current
   restore preview; this Task does not repair its Makefile.
3. The six setup tests, new vendored contract and relevance declaration pass
   the focused check above, without hand-pinned current revisions or trees.
4. The isolated re-hashed vendored-edit probe fails the new contract for the
   intended tree mismatch, and its healthy/restored controls pass.

No snapshot, lock entry, vendored file, derived artifact or Setup Manifest
value was hand-edited in the assigned worktree. The first target ran its
sanctioned digest-recording checks and refreshes. No declared Verification
command or repository Verification was rerun; the Daemon owns that sequence.
The shell sandbox could not start because its writable root contains a
symlink; the permission reviewer allowed commands from the physical
worktree. No commit, push, pull request, Task Graph edit or other Task edit
was performed.

## Carry-forward provenance

- Source Run: `run_20261104T202718Z_36f747a05c5e57d3`
- Source commit: `0b4d886a8a291bbacd6b39cfc9731adeef1eadad`
