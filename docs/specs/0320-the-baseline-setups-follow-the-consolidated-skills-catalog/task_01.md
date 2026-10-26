---
task: task_01
spec: 0320-the-baseline-setups-follow-the-consolidated-skills-catalog
status: pending
type: backend
complexity: medium
---

# Task 01: Every upstream setup follows the published commit 8d71fade

## Overview

The five bundled Setup Snapshots name upstream commit `65d5f365`, one
catalog behind the published `8d71fade68c54693bedbd917371f0514311af157`.
This Task imports that commit through the asset sync, the only writer of a
Setup Snapshot (ADR-0191), regenerates the derived Baseline files and this
repository's Setup Manifest through their sanctioned commands, and moves the
tests that pinned the old revision to one shared constant, with a new test
that pins the republished `onioncry` tree (`_prd.md` → Core Features 1-3).

## Requirements

1. MUST answer the Backlog Entry
   [`references/2026-10-26-the-baseline-setups-trail-the-consolidated-skills-catalog.md`](references/2026-10-26-the-baseline-setups-trail-the-consolidated-skills-catalog.md)
   of 2026-10-26, its Expected section.
2. MUST first check, read-only, that the local checkout `/Users/marcio/dev/skills`
   holds commit `8d71fade68c54693bedbd917371f0514311af157`
   (`git -C /Users/marcio/dev/skills cat-file -e 8d71fade68c54693bedbd917371f0514311af157^{commit}`);
   when it does not, MUST stop and report before writing anything. MUST NOT
   check out, commit, fetch into or otherwise write that checkout, and MUST
   NOT fetch from or otherwise reach a network remote.
3. MUST run, in order, the commands of `_techspec.md` → API Contracts 1 and
   2: the fetch into an empty temporary repository with the provenance
   remote added without contacting it, the binary build, the sync (exit
   `0`), its `--check` (exit `0`), `make baseline-digests`, and the managed
   refresh twice, the second reporting `File changes: 0`. Every command that
   runs Roundfix gets a fresh temporary `HOME`.
4. MUST hand-edit no snapshot, derived file, digest or Setup Manifest value.
   A file these commands change that this Task does not declare is a
   stop-and-report condition.
5. MUST create `internal/baseline/upstream_setup_revision_test.go` with the
   constants and the test of `_techspec.md` → Interfaces and API Contract 4,
   with the `Suite:`/`Invariant:`/`Boundary IN:`/`Boundary OUT:` header; the
   top-level test and each subtest call `t.Parallel()` first.
6. MUST make the declared break of `_techspec.md` → API Contract 3 in
   `internal/baseline/ai_memory_setup_test.go` and
   `internal/baseline/python_setup_test.go`, and no other change to existing
   tests.
7. MUST NOT change any Go production file, Baseline module, profile,
   template, coverage entry, installed skill, `skills-lock.json`,
   `skills/recommended.txt`, user guide page, ADR, `.roundfixrc.yml`, the
   Makefile or `go.mod`/`go.sum`.

## Subtasks

- [ ] Check the commit, fetch it into an empty repository and run the sync and its check.
- [ ] Regenerate the derived files and refresh the Setup Manifest twice.
- [ ] Add the revision test file and move the two tests to its constant.

## Acceptance Criteria

- [ ] Every upstream snapshot and every `github` skill it pins names `8d71fade`, and both TypeScript snapshots carry the `onioncry` tree `be502b02…`.
- [ ] The derived files and the Setup Manifest are current, and a second managed refresh reports no file change.
- [ ] No test names `65d5f365`; the three revision tests and the new tree test read one constant and pass.

## Consumers

- the refreshed Setup Snapshots and the upstream revision constant — consumers: `./internal/baseline`, `./skills`, `./internal/docscontract`

## Context

- instruction: `docs/adr/0191-a-setup-snapshot-follows-its-upstream-by-name.md`
- instruction: `internal/baseline/assets_sync.go`
- instruction: `internal/baseline/repository_skills_snapshot_test.go`
- instruction: `internal/baseline/skills_snapshot_profile_test.go`
- instruction: `internal/docscontract/publicdocs_test.go`
- creates: `internal/baseline/upstream_setup_revision_test.go`
- interface: `internal/baseline/ai_memory_setup_test.go`
- interface: `internal/baseline/python_setup_test.go`
- interface: `internal/baseline/assets/setups/go.json`
- interface: `internal/baseline/assets/setups/go-cli-typescript-bun.json`
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

## Verification

- `go test -list '^TestTheTypeScriptSetupsCarryTheRepublishedOnioncryTree$' ./internal/baseline | grep -qx TestTheTypeScriptSetupsCarryTheRepublishedOnioncryTree || { echo 'missing test TestTheTypeScriptSetupsCarryTheRepublishedOnioncryTree' >&2; exit 1; }; go test -timeout 20m -count=1 ./internal/baseline ./skills` — expected: exit 0. Before this Task the test does not exist, so the command fails at its first check. Afterwards it exists, both TypeScript snapshots carry the republished `onioncry` tree at the new revision, the three revision tests read the shared constant and pass, and the catalog, derived-file, snapshot-agreement and installed-skill tests pass.
- `for f in go rust typescript python go-cli-typescript-bun; do grep -qF -- '"ref": "8d71fade68c54693bedbd917371f0514311af157"' "internal/baseline/assets/setups/$f.json" || { printf '%s is not at 8d71fade\n' "$f" >&2; exit 1; }; if grep -qF -- 65d5f365cace27c4b468db1723213723b10c870b "internal/baseline/assets/setups/$f.json"; then printf '%s keeps 65d5f365\n' "$f" >&2; exit 1; fi; done; if grep -rlF --include='*_test.go' -- 65d5f365cace27c4b468db1723213723b10c870b internal/baseline; then echo 'a test still pins 65d5f365' >&2; exit 1; fi; go test -timeout 20m -count=1 -tags docscontract -run '^TestTheSetupManifestIsCurrent$' ./internal/docscontract` — expected: exit 0. Before this Task every snapshot names `65d5f365`, so the command fails at its first check. Afterwards every snapshot names only the new commit, no test pins the old one, and the Setup Manifest is current.

## References

- `_prd.md` → Goals; Core Features 1-3; Success Metrics 1-4
- `_techspec.md` → Current behavior 1-7; Interfaces; API Contracts 1-4; Invariants 1-5; Testing Approach; Build Order 1
- ADR-0191
