---
task: task_01
spec: 0326-the-baseline-setups-stop-sending-agents-to-the-secondbrain-inbox
status: pending
type: backend
complexity: medium
---

# Task 01: Every upstream setup and this repository carry the knowledge-workspace of b56fc5ff

## Overview

The five bundled Setup Snapshots pin `knowledge-workspace` at upstream commit
`8d71fade`, a tree that still sends agents to the retired Secondbrain inbox;
the published `b56fc5ffcad1ab366fbc03839a6ce8c662401178` fixes it. This Task
imports that commit through the asset sync, the only writer of a Setup
Snapshot (ADR-0191), restores this repository's installed
`knowledge-workspace` and its lock entry through the confirmed restore
command (ADR-0264), regenerates the derived Baseline files and the Setup
Manifest through their sanctioned commands, and moves the shared revision
test constant with a new test that pins the corrected tree
(`_prd.md` → Core Features 1-4).

## Requirements

1. MUST answer the Backlog Entry
   [`references/2026-10-26-the-knowledge-workspace-skill-still-writes-to-the-secondbrain-inbox.md`](references/2026-10-26-the-knowledge-workspace-skill-still-writes-to-the-secondbrain-inbox.md)
   of 2026-10-26, its Expected section.
2. MUST first check, read-only, that the local clone
   `/Users/marcio/.roundfix-operator/scratch/spec-0326/skills` holds commit
   `b56fc5ffcad1ab366fbc03839a6ce8c662401178`
   (`git -C /Users/marcio/.roundfix-operator/scratch/spec-0326/skills cat-file -e b56fc5ffcad1ab366fbc03839a6ce8c662401178^{commit}`);
   when it does not, MUST stop and report before writing anything. MUST NOT
   check out, commit, fetch into or otherwise write that clone, MUST NOT
   read or write `/Users/marcio/dev/skills`, and MUST NOT fetch from or
   otherwise reach a network remote.
3. MUST build the binary under the Run's environment before any command
   that sets a fresh `HOME`, so the build reads the existing module cache.
4. MUST run, in order, the commands of `_techspec.md` → API Contracts 1, 2
   and 3: the fetch into an empty temporary repository with the provenance
   remote added without contacting it, the sync (exit `0`), its `--check`
   (exit `0`), the restore preview (exit `3`, exactly the two planned
   changes of API Contract 2) and its confirmed apply (exit `0`),
   `make baseline-digests`, and the managed refresh twice, the second
   reporting `File changes: 0`. Every command that runs Roundfix gets a
   fresh temporary `HOME`.
5. MUST hand-edit no snapshot, installed skill, lock entry, derived file,
   digest or Setup Manifest value. A file these commands change that this
   Task does not declare is a stop-and-report condition.
6. MUST change `internal/baseline/upstream_setup_revision_test.go` as
   `_techspec.md` → Interfaces and API Contracts 4 and 5 state: move
   `upstreamSetupRevision`, add `knowledgeWorkspaceTreeDigest` and
   `TestEverySetupCarriesTheKnowledgeWorkspaceThatRetiredTheSecondbrainInbox`
   (the top-level test and each subtest call `t.Parallel()` first), and
   extend the `Invariant:` header line; every other line keeps its bytes.
7. MUST NOT change any other test, Go production file, Baseline module,
   profile, template, Source Baseline, retention file, coverage entry,
   installed skill other than `knowledge-workspace`, lock entry other than
   `knowledge-workspace`'s, `skills/recommended.txt`, user guide page, ADR,
   `CONTEXT.md`, `.roundfixrc.yml`, the Makefile or `go.mod`/`go.sum`.

## Subtasks

- [ ] Check the commit, fetch it into an empty repository and run the sync and its check.
- [ ] Restore this repository's knowledge-workspace from the fetched checkout.
- [ ] Regenerate the derived files and refresh the Setup Manifest twice.
- [ ] Move the revision constant and add the knowledge-workspace tree test.

## Acceptance Criteria

- [ ] Every upstream snapshot and every `github` skill it pins names `b56fc5ff`, and all five carry the `knowledge-workspace` tree `490af37d…`.
- [ ] This repository's installed `knowledge-workspace` reads `version: 0.6.2`, no longer names the Secondbrain inbox, and matches its lock entry and the `go` snapshot.
- [ ] The derived files and the Setup Manifest are current, and a second managed refresh reports no file change.
- [ ] No test names `8d71fade`; the new tree test and the four revision tests read the shared constant and pass.

## Consumers

- the refreshed Setup Snapshots, the restored skill and the upstream revision constant — consumers: `./internal/baseline`, `./skills`, `./internal/docscontract`

## Context

- instruction: `docs/adr/0191-a-setup-snapshot-follows-its-upstream-by-name.md`
- instruction: `docs/adr/0264-roundfix-never-replaces-a-vendored-skill-it-cannot-prove-older.md`
- instruction: `internal/baseline/assets_sync.go`
- instruction: `internal/baseline/skills_restore.go`
- instruction: `internal/baseline/repository_skills_snapshot_test.go`
- instruction: `internal/baseline/skills_snapshot_profile_test.go`
- instruction: `skills/repository_test.go`
- instruction: `internal/docscontract/publicdocs_test.go`
- interface: `internal/baseline/upstream_setup_revision_test.go`
- interface: `internal/baseline/assets/setups/go.json`
- interface: `internal/baseline/assets/setups/go-cli-typescript-bun.json`
- interface: `internal/baseline/assets/setups/python.json`
- interface: `internal/baseline/assets/setups/rust.json`
- interface: `internal/baseline/assets/setups/typescript.json`
- interface: `.agents/skills/knowledge-workspace/SKILL.md`
- interface: `skills-lock.json`
- interface: `internal/baseline/testdata/catalog.diagnostics.golden.json`
- interface: `internal/baseline/testdata/catalog.digest`
- interface: `internal/baseline/testdata/catalog.normalized.json`
- interface: `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`
- interface: `docs/agents/setup-context.json`

## Verification

- `go test -list '^TestEverySetupCarriesTheKnowledgeWorkspaceThatRetiredTheSecondbrainInbox$' ./internal/baseline | grep -qx TestEverySetupCarriesTheKnowledgeWorkspaceThatRetiredTheSecondbrainInbox || { echo 'missing test TestEverySetupCarriesTheKnowledgeWorkspaceThatRetiredTheSecondbrainInbox' >&2; exit 1; }; go test -timeout 20m -count=1 ./internal/baseline ./skills` — expected: exit 0. Before this Task the test does not exist, so the command fails at its first check. Afterwards all five snapshots carry the corrected `knowledge-workspace` tree at the new revision, the four revision tests read the moved constant, this repository's installed skill matches its lock and the `go` snapshot, and the catalog, derived-file and snapshot-agreement tests pass.
- `for f in go rust typescript python go-cli-typescript-bun; do grep -qF -- '"ref": "b56fc5ffcad1ab366fbc03839a6ce8c662401178"' "internal/baseline/assets/setups/$f.json" || { printf '%s is not at b56fc5ff\n' "$f" >&2; exit 1; }; if grep -qF -- 8d71fade68c54693bedbd917371f0514311af157 "internal/baseline/assets/setups/$f.json"; then printf '%s keeps 8d71fade\n' "$f" >&2; exit 1; fi; done; if grep -rlF --include='*_test.go' -- 8d71fade68c54693bedbd917371f0514311af157 internal/baseline; then echo 'a test still pins 8d71fade' >&2; exit 1; fi; tr -s '[:space:]' ' ' < .agents/skills/knowledge-workspace/SKILL.md | grep -qF -- 'version: 0.6.2' || { echo 'installed knowledge-workspace is not 0.6.2' >&2; exit 1; }; if tr -s '[:space:]' ' ' < .agents/skills/knowledge-workspace/SKILL.md | grep -qF -- 'dev/secondbrain/inbox/'; then echo 'installed knowledge-workspace still writes to the Secondbrain inbox' >&2; exit 1; fi; go test -timeout 20m -count=1 -tags docscontract -run '^TestTheSetupManifestIsCurrent$' ./internal/docscontract` — expected: exit 0. Before this Task every snapshot names `8d71fade`, so the command fails at its first check. Afterwards every snapshot names only the new commit, no test pins the old one, the installed skill is the corrected 0.6.2, and the Setup Manifest is current.

## References

- `_prd.md` → Goals; Core Features 1-4; Success Metrics 1-4
- `_techspec.md` → Current behavior 1-7; Interfaces; API Contracts 1-5; Invariants 1-5; Testing Approach; Build Order 1
- ADR-0191; ADR-0264
