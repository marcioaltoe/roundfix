---
task: task_01
spec: 0350-the-vendored-catalog-skills-follow-the-catalog-through-a-pinned-bump
status: pending
type: infra
complexity: medium
---

# Task 01: The Catalog Skill Bump pins the vendored catalog skills to one catalog commit

## Overview

Add the Makefile target `skills-third-party-bump` that `_techspec.md` → API
Contracts 1-3 define and retire `skills-update`. The target composes the
existing asset sync, confirmed restore, `make baseline-digests` and managed
refresh; its own shell resolves and fetches one catalog commit, refuses what
it must not pin and checks the result. A fixture test runs the real Makefile
target against a local fake catalog with recording stubs, so the Task is
verifiable without the network and without moving any real skill.

## Requirements

1. MUST answer the Backlog Entry
   [`references/2026-11-04-the-vendored-catalog-skills-trail-the-catalog-and-name-branch-refs.md`](references/2026-11-04-the-vendored-catalog-skills-trail-the-catalog-and-name-branch-refs.md),
   its Expected section, for the repeatable update and the owned-skill rule.
2. MUST add `skills-third-party-bump` with a `##` help line, the variables
   `SKILLS_CATALOG` (default `https://github.com/marcioaltoe/skills.git`),
   `SKILLS_CATALOG_ORIGIN` and `SKILLS_REF` (default empty), and the
   `.PHONY` entry, behaving exactly as `_techspec.md` → API Contracts 1, 2
   and 3 state: the refusals before any command and their message fragments,
   the step order with two builds, the post-restore checks, the regeneration
   and the final stdout line. Only `git ls-remote` and `git fetch` may use the
   network.
3. MUST compute the catalog set from the `skills-lock.json` entries whose
   `source` is `marcioaltoe/skills`, cross-check it with the `.agents/skills`
   folders outside `OWNED_SKILLS`, and refuse an owned name in it; MUST read
   the profile from `docs/agents/setup-context.json`; MUST create its
   temporary directory under `$TMPDIR` (falling back to `/tmp`) and remove it
   on every exit; MUST run every `roundfix` call with `HOME` set to an empty
   directory inside it; MUST call `baseline-digests` through `$(MAKE) -f` this
   Makefile, so the target works from any working directory.
4. MUST remove the `skills-update` target and its `.PHONY` mention, if any.
5. MUST create `internal/verifyselect/skills_bump_makefile_test.go` with
   `TestTheSkillsBumpPinsTheCatalogMainCommit` and
   `TestTheSkillsBumpRefusesWhatItMustNotPin`, built as `_techspec.md` →
   Testing Approach states: `make -f <repository Makefile>
   skills-third-party-bump` run in a temporary repository fixture, a local
   fake catalog Git repository with branch `main`, `GO=` a recording shell
   stub that copies a recording `roundfix` stub on `build -o`,
   `DERIVED_DIGEST_PATHS=` an empty directory, `BASELINE_DIGEST_STEPS=probe:Probe`,
   and an environment holding only `PATH`, `TMPDIR`, `HOME` (all under the
   test's temporary directory), `GIT_ALLOW_PROTOCOL=file` and
   `GIT_CONFIG_NOSYSTEM=1`. The first test MUST leave `SKILLS_REF` unset and
   assert the final line names the fake `main` commit, exactly two builds,
   the sync with `--source-dir .../setups` and its `--check`, the restore of
   exactly the lock's catalog names in sorted order (never an owned name),
   the confirmed restore with the stub's Plan Digest, the `baseline-digests`
   probe and the managed refresh. The second MUST have one parallel subtest
   per refusal: an owned skill in the lock, `SKILLS_REF=main`, a commit the
   catalog lacks, a vendored folder without a lock entry (each asserting no
   stub call was recorded), a vendored folder that differs from the catalog
   and an entry left at another ref (each asserting no regeneration call).
   Both tests MUST call `t.Parallel()` first and reuse the package helpers
   `writeFile`, `runGit` and `findRepositoryRoot`.
6. MUST run, as the last step, `go run -buildvcs=false ./cmd/roundfix baseline
   update --repo . --no-skills --no-history --yes --format text` with a fresh
   temporary `HOME` whose module and build caches are the Run's own
   (`GOMODCACHE` and `GOCACHE` set from `go env` before `HOME` changes, and
   the temporary `HOME` removed afterwards), so `docs/agents/setup-context.json`
   records the new Makefile digest; a second run MUST report
   `File changes: 0`.
7. MUST NOT run the target against the real catalog, change any vendored
   skill, `skills-lock.json` or Setup Snapshot, or change any production
   `.go` file.

## Subtasks

- [ ] Add the target, its variables and its help line; retire `skills-update`.
- [ ] Write the fixture test with the fake catalog and recording stubs.
- [ ] Refresh the Setup Manifest after the last Makefile edit.

## Acceptance Criteria

- [ ] With no `SKILLS_REF`, the target pins the fake catalog's `main` commit and runs sync, check, two builds, the confirmed restore of the lock's catalog names, `baseline-digests` and the refresh.
- [ ] Each refusal of API Contract 1 runs no command; each check of API Contract 3 stops before regeneration.
- [ ] `make skills-update` no longer exists, and the Setup Manifest is current.

## Context

- instruction: `docs/adr/0360-the-vendored-catalog-skills-follow-the-catalog-through-one-pinned-commit.md`
- instruction: `internal/verifyselect/verify_docs_makefile_test.go`
- instruction: `internal/verifyselect/verifyselect_test.go`
- instruction: `skills/baseline_skill_contract_test.go`
- interface: `Makefile`
- interface: `docs/agents/setup-context.json`
- creates: `internal/verifyselect/skills_bump_makefile_test.go`

## Consumers

- the `skills-third-party-bump` target and the Makefile declaration digest — consumers: `./internal/verifyselect`, `./internal/docscontract`

## Verification

- `for name in TestTheSkillsBumpPinsTheCatalogMainCommit TestTheSkillsBumpRefusesWhatItMustNotPin; do go test -list "^$name\$" ./internal/verifyselect | grep -qx "$name" || { echo "missing test $name" >&2; exit 1; }; done; ! grep -q '^skills-update:' Makefile || { echo 'skills-update target still present' >&2; exit 1; }; go test -timeout 20m -count=1 ./internal/verifyselect && go test -timeout 10m -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSetupManifestIsCurrent$'` — expected: exit 0. Before this Task neither test exists, so the command fails at its first check; afterwards the two fixture tests and every other Makefile contract test of the package pass, and the Setup Manifest records the new Makefile.

## References

- `_prd.md` → Goals; Core Feature 1; User Experience; Success Metric 1; Success Metric 5
- `_techspec.md` → API Contracts 1-3; Invariants 2, 4 and 5; Testing Approach; Build Order 1
- ADR-0360
