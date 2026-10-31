---
task: task_02
spec: 0333-runs-load-the-dispatched-skills-and-the-baseline-follows-the-catalog
status: completed
type: backend
complexity: medium
---

# Task 02: The dispatch asks for person-started and polish skills only on request

## Overview

Five skills a Run should never start on its own become dispatched only when
the person asks: `architectural-analysis` and `grill-with-docs`, which only a
person can start, and `humanizer`, `writing-clearly-and-concisely` and
`web-design-guidelines`, which polish or review on request. `wcag-audit-patterns`
is kept for a formal audit. The frontend UI bundle splits into a product UI
bundle and a public pages bundle, so `interface-design` and `frontend-design`
no longer load together. The snapshots, derived files and this repository's
dispatch guide follow through the sanctioned commands. It applies the
maintainer's dispatch decisions of 2026-10-26 recorded in `_prd.md` →
Decisions.

## Requirements

1. MUST reword, as text that keeps each file's layout, the `when` of the
   eight triggers of `_techspec.md` → Data Model 3 in
   `internal/baseline/assets/modules/typescript.json`, `modules/frontend.json`
   and `modules/context-workflow.json`, verbatim; no other trigger of those
   modules changes, and `trigger.frontend.ui-ux-pro-max` keeps its text.
2. MUST change `internal/baseline/assets/skill-activations.json` to Data
   Model 2's frontend bundles: `bundle.frontend-ui-quality` holds
   `interaction-design`, `interface-design` and `fixing-accessibility`, and
   its trigger reads "Designing or reviewing product UI: application screens,
   forms, and dashboards."; add `bundle.frontend-public-pages` with
   `frontend-design` and `fixing-accessibility` and its activation
   `trigger.frontend.public-pages` (owner `frontend`, "Designing or reviewing
   marketing or other public pages."), each placed in `id` order. The file's
   `version` stays 3.
3. MUST run the module record command with its outputs, then `_techspec.md`
   → API Contract 2 in order: the `typescript.json` bundle assignment, the
   asset sync and its `--check`, the parity transform, `make
   baseline-digests` and the two managed refreshes (the second reports `File
   changes: 0`). This repository's `docs/agents/skill-dispatch.md` is
   rendered, never edited by hand.
4. MUST add `internal/baseline/skill_dispatch_review_test.go` with the
   `// Suite:`, `// Invariant:` and boundary header of its neighbours and
   `TestPersonStartedSkillsAreDispatchedOnlyWhenThePersonAsks` and
   `TestTheFrontendBundlesSeparateProductUIFromPublicPages` as `_techspec.md`
   → Testing Approach states, reading the embedded catalog. Each test must
   fail when one trigger loses "When the person asks" or one bundle regains
   `frontend-design` beside `interface-design`.

## Subtasks

- [ ] Reword the triggers and split the frontend bundle.
- [ ] Record the module versions, re-sync the snapshots and regenerate the derived files and this repository's dispatch guide.
- [ ] Add the two dispatch tests.

## Acceptance Criteria

- [ ] The five person-started or polish skills are dispatched only when the person asks and sit in no bundle.
- [ ] Product UI and public pages have separate bundles; `wcag-audit-patterns` and `web-design-guidelines` sit in no bundle.
- [ ] This repository's dispatch guide carries the new `grill-with-docs` trigger and a second refresh changes no file.

## Context

- instruction: `docs/adr/0246-a-retired-skill-leaves-the-baseline-and-the-adopter-removes-its-copy.md`
- interface: `internal/baseline/assets/modules/typescript.json`
- interface: `internal/baseline/assets/modules/frontend.json`
- interface: `internal/baseline/assets/modules/context-workflow.json`
- interface: `internal/baseline/assets/skill-activations.json`
- interface: `internal/baseline/module-versions.json`
- interface: `internal/baseline/assets/setups/typescript.json`
- interface: `internal/baseline/assets/setups/go-cli-typescript-bun.json`
- interface: `internal/baseline/assets/profiles/standard-typescript-monorepo.json`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/skill-dispatch.md`
- interface: `internal/baseline/testdata/catalog.diagnostics.golden.json`
- interface: `internal/baseline/testdata/catalog.digest`
- interface: `internal/baseline/testdata/catalog.normalized.json`
- interface: `internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json`
- interface: `internal/baseline/testdata/parity-corpus/v1/manifest.json`
- interface: `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`
- interface: `docs/agents/setup-context.json`
- interface: `docs/agents/skill-dispatch.md`
- creates: `internal/baseline/skill_dispatch_review_test.go`

## Consumers

- The activation bundles and triggers — consumers: `./internal/baseline`

## Verification

- `out="$(go test -count=1 -timeout 20m -v ./internal/baseline 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestPersonStartedSkillsAreDispatchedOnlyWhenThePersonAsks TestTheFrontendBundlesSeparateProductUIFromPublicPages TestTheTypeScriptSnapshotsCarryTheCatalogActivationBundles TestNoCatalogEntryNamesARetiredSkill; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < docs/agents/skill-dispatch.md | grep -qF -- 'When the person asks to stress-test a plan or decision against repository documents.' || { echo 'dispatch guide lacks the grill-with-docs trigger' >&2; exit 1; }` — expected: exit 0. Before this Task the two dispatch tests do not exist, so the first is missing from the passes. Afterwards the whole Baseline package passes with the reworded triggers and split bundles, the TypeScript snapshots still carry the catalog bundles, and this repository's rendered dispatch guide carries the new trigger.

## References

- `_prd.md` → Goals; Core Feature 3; Stories 1 and 4; Success Metric 2; Decisions
- `_techspec.md` → Data Models 2 and 3; API Contract 2; Testing Approach; Build Order 2

## Result

Skills loaded: implement-task, coding-guidelines, testing-boss, golang-testing, evidence-gate, setup-context-driven, roundfix, knowledge-workspace, tech-writer.

The eight dispatch triggers now use Data Model 3's exact text. The five
person-started or polish skills require a request and belong to no activation
bundle. Product UI uses `interaction-design`, `interface-design` and
`fixing-accessibility`; public pages use `frontend-design` and
`fixing-accessibility`. Neither audit nor web-guideline review joins a bundle.
The activation document remains version 3, with bundles and activations in ID
order. `ui-ux-pro-max` and every other module trigger retain their text.

Regeneration, in the required order:

- `go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1`
  exited 0. The recorder raised TypeScript 9 → 10, frontend 7 → 8 and context
  workflow 27 → 28, and wrote the module version record.
- Fetched catalog commit `86df8608726979064a8fe46003dfaf403ffe0087` from
  the authorized local checkout into a temporary repository and built the
  task binary with `go build -buildvcs=false`.
- Applied the prescribed `jq` bundle assignment to `typescript.json`, then
  ran `baseline assets sync --source-dir <fetched>/setups --format text`
  and its `--check` with an isolated HOME. Both exited 0; sync reported
  updates to the TypeScript and composed TypeScript snapshots, and check
  reported `setup-context-driven audit: ok`.
- Applied API Contract 1's parity transform with the composed-snapshot
  filter, then ran `make baseline-digests`: exit 0, `ok: true`. The command
  regenerated the declared formatter, profile, catalog, parity and plan
  artifacts.
- Rebuilt the task binary and ran
  `baseline update --repo . --no-skills --no-history --yes --format text`
  twice, each with a fresh HOME. The first rendered two files; the second
  reported `File changes: 0` and `Idempotence: verified`. The generated
  repository guide carries the new `grill-with-docs` trigger. Both refreshes
  reported the two existing nested fixture carriers as unchanged warnings.

Focused evidence:

- Before the catalog edits, both new tests failed on the five old trigger
  prefixes, the combined frontend bundle and its missing public-pages
  counterpart.
- `go test ./internal/baseline -run '^(TestPersonStartedSkillsAreDispatchedOnlyWhenThePersonAsks|TestTheFrontendBundlesSeparateProductUIFromPublicPages|TestTheTypeScriptSnapshotsCarryTheCatalogActivationBundles|TestNoCatalogEntryNamesARetiredSkill|TestEveryBaselineModuleVersionIsRecorded|TestCatalogCompatibility|TestCatalogDiagnosticCharacterization|TestBaselinePlanCharacterization|TestBaselineCompatibilityCorpus|TestFormatterComposition)$' -count=1 -v`
  exited 0; all ten selected tests passed. These cover acceptance criteria
  1 and 2, snapshot consistency and the regenerated artifacts.
- Six temporary Go overlays changed no repository source: five independently
  removed `When the person asks`, and one added `frontend-design` beside
  `interface-design`. Each focused test invocation exited 1 at its intended
  assertion, confirming both regression tests reject those changes.
- Parsed inspection compared the three modules with HEAD after allowing
  only the eight specified trigger changes and recorded version increments;
  it also confirmed activation version/order and the rendered guide text.
  `gofmt -d internal/baseline/skill_dispatch_review_test.go` produced no diff.
- `git diff --check` exited 0. The scope audit found all 22 changed or
  untracked paths declared by this Task, with its authored content and
  Daemon-written `in_progress` status preserved.

Declared Verification and repository Verification were not run; the Daemon
owns both and Task status. No other Task or Task Graph file was edited, and
no commit, push or pull request was made.
