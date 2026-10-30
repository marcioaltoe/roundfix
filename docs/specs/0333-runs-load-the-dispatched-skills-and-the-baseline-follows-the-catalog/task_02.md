---
task: task_02
spec: 0333-runs-load-the-dispatched-skills-and-the-baseline-follows-the-catalog
status: pending
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
