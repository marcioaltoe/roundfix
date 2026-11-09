---
task: task_01
spec: 0370-the-inngest-capability-takes-the-seven-official-inngest-skills
status: pending
type: backend
complexity: medium
---

# Task 01: The backend module takes the seven Inngest skills, inngest is retired and the snapshots follow the catalog

## Overview

The catalog replaced `inngest` with seven official `inngest-*` skills at
commit `e45466d2`, but the `backend` module still requires and dispatches
`inngest` and the Catalog Skill Bump stops at the mismatch. This Task answers
the report adopted in
[`references/2026-11-09-marcioaltoe-skills-inngest-replaced-by-seven-official-skills.md`](references/2026-11-09-marcioaltoe-skills-inngest-replaced-by-seven-official-skills.md):
the module requires and dispatches the seven, `inngest` becomes a Retired
Skill, and the snapshots, the vendored catalog skills and every derived file
follow through the repository's own commands (`_prd.md` → Core Features
1-3).

## Requirements

1. MUST answer the marcioaltoe-skills report of 2026-11-09 (ai-memory
   message 01a1c01f-b6ad, adopted in `references/`), its Request section.
2. MUST first check, read-only, that the operator's clone
   `/Users/marcio/.roundfix-operator/scratch/spec-0370/skills` holds commit
   `e45466d2b9e15480eb286a4c19b513494a867bce`
   (`git -C /Users/marcio/.roundfix-operator/scratch/spec-0370/skills cat-file -e e45466d2b9e15480eb286a4c19b513494a867bce^{commit}`);
   when it does not, MUST stop and report before writing anything. MUST NOT
   check out, commit, fetch into or otherwise write that clone, MUST NOT
   read or write `/Users/marcio/dev/skills`, and MUST NOT reach a network
   remote.
3. MUST edit `internal/baseline/assets/modules/backend.json`: in
   `requiredSkills`, replace `"inngest"` with `"inngest-cli"`,
   `"inngest-durable-functions"`, `"inngest-events"`,
   `"inngest-flow-control"`, `"inngest-middleware"`, `"inngest-setup"` and
   `"inngest-steps"`, in that order and in the same place; in
   `skillDispatch`, replace the `inngest` entry with seven entries in the same
   order, each with exactly one trigger whose `id` is
   `trigger.backend.<skill>` and whose `when` is the text `_techspec.md` →
   API Contract 1 gives, byte for byte. No other line of the module changes
   by hand; its version line changes only through the record command.
4. MUST add `"inngest": {},` to `retiredSkills` in
   `internal/baseline/retired_skills.go`, in lexical position after
   `"implement-spec"`, and add `"inngest"` in the same position to the
   expected list of `TestNoCatalogEntryNamesARetiredSkill` in
   `internal/baseline/retired_skills_test.go`; MUST change no other
   production Go file.
5. MUST run, in order, the commands of `_techspec.md` → API Contract 3:
   the module record command; the Catalog Skill Bump with
   `GIT_ALLOW_PROTOCOL=file`, `SKILLS_CATALOG` set to the clone of
   requirement 2 and `SKILLS_REF=e45466d2b9e15480eb286a4c19b513494a867bce`
   (exit `0`); the asset-sync parity transform with the quoted `jq` program,
   through temporary files under the Run's own temporary directory;
   `make baseline-digests`; and the managed refresh
   `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text`
   twice, the second reporting `File changes: 0`. Every command that runs
   Roundfix directly gets a fresh temporary `HOME`, with `GOMODCACHE` and
   `GOCACHE` kept at their `go env` values, removed afterwards.
6. MUST hand-edit no snapshot, vendored skill file, lock entry, derived file,
   digest, parity fixture or Setup Manifest value. A file these commands
   change that this Task does not declare is a stop-and-report condition.
7. MUST create `internal/baseline/inngest_skills_test.go`, with a
   `// Suite:`, `// Invariant:`, `// Boundary IN:` and `// Boundary OUT:`
   header like its neighbours, holding
   `TestTheBackendModuleTakesTheSevenOfficialInngestSkills`, which calls
   `t.Parallel()` first and, from `LoadEmbeddedCatalog()`:
   - decodes the `backend` module and fails when `requiredSkills` or any
     `skillDispatch` entry names `inngest`, when a skill of the seven is not
     required, or when it is not dispatched exactly once with exactly the one
     trigger and `when` text of API Contract 1;
   - for `standard-typescript-monorepo` and `go-cli-typescript-monorepo`,
     calls `loadRestoreProfile` and fails when `inngest` is in its `Skills`
     or `RequiredSkills`, or when a skill of the seven is not in
     `RequiredSkills` or is not in `Skills` with `Source.Path`
     `skills/05-implementation-loop/<skill>` and a non-empty `TreeDigest`;
   - fails when `RetiredSkills()` does not hold `inngest`.

   The test MUST NOT name a Catalog Revision or a tree digest literal
   (ADR-0360).
8. MUST create `internal/cli/baseline_update_inngest_migration_test.go`,
   with the same four-line header, holding two top-level tests that call
   `t.Parallel()` first:
   - `TestTheBackendModuleRequiresTheSevenInngestSkills`: in `t.TempDir()`,
     write `docs/agents/setup-context.json` as
     `{"modules":["typescript","backend"]}` and call
     `resolveExternalSkillRequirement`; fail unless it reports a manifest
     and its list holds each of the seven and not `inngest`.
   - `TestBaselineUpdateApprovedApplyDeletesTheReplacedInngestSkill`, with
     subtests `unedited` and `edited`: on `newBaselineUpdateRepository` with
     `installRequiredExternalSkills`, write `.agents/skills/inngest/SKILL.md`,
     record its `skills.SkillFolderHash` as the `computedHash` of an
     `inngest` lock entry, link `.claude/skills/inngest` to
     `../../.agents/skills/inngest`, and commit. `unedited`: the text preview
     exits `exitUnverified` and contains the deletion line of `_techspec.md`
     → API Contract 2; the JSON preview is `plan_ready` with
     `Skills.Retired` equal to one `delete` entry with both paths, the lock
     entry and the hash; the preview changes no byte; the apply with
     `--confirm-plan <digest>` is `verified` with that entry in
     `Skills.Deleted`, both paths are gone, and the lock equals the earlier
     one without `inngest`. `edited`: after the hash is recorded, rewrite the
     `SKILL.md`, commit; the JSON preview is `current` with one `keep` entry,
     reason `edited`; the text output of `--yes` contains the kept line of
     API Contract 2; no byte of the repository changes.
9. MUST update every other test that pins the Retired Skill list, the
   `backend` module's skills or the snapshot contents on the default branch
   when the Task runs, if any; at `cfc08683` only
   `TestNoCatalogEntryNamesARetiredSkill` does.
10. MUST NOT change any guide, skill, ADR, `CONTEXT.md`, the Makefile, a
    workflow, `go.mod`, `go.sum`, any Baseline module other than `backend`,
    any template, Source Baseline or retention file, or any frozen parity
    fixture other than `v1/fixtures/asset-sync.json` and `v1/manifest.json`;
    task_02 documents the change.

## Subtasks

- [ ] Check the clone, edit the module and the Retired Skill set, and record the module version.
- [ ] Run the bump, the parity transform, the digests and the two managed refreshes.
- [ ] Add the two test files and extend the Retired Skill list test.

## Acceptance Criteria

- [ ] The `backend` module requires and dispatches the seven skills with the prescribed triggers, and `inngest` is a Retired Skill.
- [ ] All five snapshots, the lock and the asset-sync row name `e45466d2`, the two TypeScript snapshots pin the seven and no `inngest`, and a second managed refresh changes nothing.
- [ ] An unedited `inngest` copy is deleted by the approved update, and an edited one is kept.

## Context

- instruction: `docs/specs/0370-the-inngest-capability-takes-the-seven-official-inngest-skills/references/2026-11-09-marcioaltoe-skills-inngest-replaced-by-seven-official-skills.md`
- instruction: `docs/adr/0380-a-skill-the-catalog-replaces-is-followed-to-its-replacements-and-retired.md`
- instruction: `docs/adr/0350-an-approved-baseline-update-deletes-the-retired-skills-it-proves-unedited.md`
- instruction: `docs/adr/0360-the-vendored-catalog-skills-follow-the-catalog-through-one-pinned-commit.md`
- instruction: `internal/baseline/skills_restore.go`
- instruction: `internal/baseline/retired_skills_installed.go`
- instruction: `internal/cli/baseline_update_retired_skills_test.go`
- interface: `internal/baseline/assets/modules/backend.json`
- interface: `internal/baseline/module-versions.json`
- interface: `internal/baseline/retired_skills.go`
- interface: `internal/baseline/retired_skills_test.go`
- creates: `internal/baseline/inngest_skills_test.go`
- creates: `internal/cli/baseline_update_inngest_migration_test.go`
- interface: `internal/baseline/assets/setups/go-cli-typescript-bun.json`
- interface: `internal/baseline/assets/setups/go.json`
- interface: `internal/baseline/assets/setups/python.json`
- interface: `internal/baseline/assets/setups/rust.json`
- interface: `internal/baseline/assets/setups/typescript.json`
- interface: `skills-lock.json`
- interface: `.agents/skills/context7-cli/references/setup.md`
- interface: `internal/baseline/assets/profiles/standard-typescript-monorepo.json`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/skill-dispatch.md`
- interface: `internal/baseline/testdata/catalog.diagnostics.golden.json`
- interface: `internal/baseline/testdata/catalog.digest`
- interface: `internal/baseline/testdata/catalog.normalized.json`
- interface: `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`
- interface: `internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json`
- interface: `internal/baseline/testdata/parity-corpus/v1/manifest.json`
- interface: `docs/agents/setup-context.json`

## Consumers

- The `backend` module's required skills and dispatch triggers — consumers: `./internal/baseline`, `./internal/cli`
- The Retired Skill set with `inngest` — consumers: `./internal/baseline`, `./internal/cli`
- The Setup Snapshots, `skills-lock.json` and the vendored catalog skills at `e45466d2` — consumers: `./internal/baseline`, `./skills`, `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^TestTheBackendModuleTakesTheSevenOfficialInngestSkills$' ./internal/baseline 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestTheBackendModuleTakesTheSevenOfficialInngestSkills ' || { echo 'missing pass: TestTheBackendModuleTakesTheSevenOfficialInngestSkills' >&2; exit 1; }; go test -count=1 -timeout 20m ./internal/baseline ./skills` — expected: exit 0. Before this Task the test does not exist, so its pass line is missing; afterwards it passes with the whole Baseline and skills suites, the snapshot, parity and Retired Skill tests included.
- `out="$(go test -count=1 -timeout 20m -v -run '^(TestTheBackendModuleRequiresTheSevenInngestSkills|TestBaselineUpdateApprovedApplyDeletesTheReplacedInngestSkill)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTheBackendModuleRequiresTheSevenInngestSkills TestBaselineUpdateApprovedApplyDeletesTheReplacedInngestSkill TestBaselineUpdateApprovedApplyDeletesTheReplacedInngestSkill/unedited TestBaselineUpdateApprovedApplyDeletesTheReplacedInngestSkill/edited; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 30m -run 'Baseline' ./internal/cli` — expected: exit 0. Before this Task neither test exists; afterwards both pass with every Baseline test of the CLI package.
- `for f in go rust typescript python go-cli-typescript-bun; do p="internal/baseline/assets/setups/$f.json"; grep -qF -- '"ref": "e45466d2b9e15480eb286a4c19b513494a867bce"' "$p" || { printf '%s is not at e45466d2\n' "$f" >&2; exit 1; }; if grep -qF -- c8b2e8511a7ff89dbe9e87a56cc3c29a6acbcb84 "$p"; then printf '%s keeps c8b2e85\n' "$f" >&2; exit 1; fi; done; for f in typescript go-cli-typescript-bun; do p="internal/baseline/assets/setups/$f.json"; if grep -qF -- '"skills/05-implementation-loop/inngest"' "$p"; then printf '%s pins inngest\n' "$f" >&2; exit 1; fi; for s in inngest-cli inngest-durable-functions inngest-events inngest-flow-control inngest-middleware inngest-setup inngest-steps; do grep -qF -- "\"skills/05-implementation-loop/$s\"" "$p" || { printf '%s lacks %s\n' "$f" "$s" >&2; exit 1; }; done; done; if grep -qF -- c8b2e8511a7ff89dbe9e87a56cc3c29a6acbcb84 skills-lock.json; then echo 'skills-lock.json keeps c8b2e85' >&2; exit 1; fi; if grep -qF -- '"name": "inngest",' internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json; then echo 'the asset-sync parity row names inngest' >&2; exit 1; fi; go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract -run '^TestTheSetupManifestIsCurrent$'` — expected: exit 0. Before this Task every snapshot names `c8b2e85`, so the first check fails; afterwards the snapshots, the lock and the parity row follow `e45466d2` and the Setup Manifest is current.

## References

- `_prd.md` → Goals; User Stories 1-3; Core Features 1-3; Success Metrics 1-3
- `_techspec.md` → Interfaces; API Contracts 1-3; Testing Approach; Build Order 1
- ADR-0380; ADR-0350; ADR-0360; ADR-0191
