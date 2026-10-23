---
task: task_02
spec: 0309-a-python-cli-baseline-profile-for-uv-pytest-and-ruff
status: pending
type: backend
complexity: high
---

# Task 02: The python-cli profile governs uv, ruff and pytest on the common setup of Go and Rust

## Overview

This Task ships the `python-cli` Baseline Profile: the `python` common setup
of `go` and `rust`, the `python` module with its guide and root block, the
templates, the coverage entry, the core module's hierarchy and semantic-owner
entries and the profile itself. The sanctioned commands regenerate every
derived file and this repository's Setup Manifest, the existing tests that
enumerate the catalog learn the new entries, and two new tests hold the setup
and the rendered guide (`_prd.md` → Core Features 2-4; ADR-0322).

## Requirements

1. MUST answer the Backlog Entry
   [`references/2026-10-22-a-baseline-profile-for-python-repositories.md`](references/2026-10-22-a-baseline-profile-for-python-repositories.md)
   of 2026-10-22, its expected items 1, 2, 3 and 5.
2. MUST create `internal/baseline/assets/setups/python.json` with exactly the
   bytes the asset sync writes for a `common` snapshot of `go` then `rust`
   (`_techspec.md` → API Contract 2, Data Model 1): the output of
   `commonSetupSnapshot("python", …)` over the embedded `go` and `rust`
   snapshots, decoded into the sync's snapshot struct and written with
   two-space indentation and a final newline. A temporary helper used to
   produce the bytes MUST NOT remain in the tree.
3. MUST create `internal/baseline/assets/modules/python.json` (Data Model 2,
   the rules and clauses of API Contract 3 verbatim, module `version` 1 before
   the record command runs), `internal/baseline/assets/templates/root/python.md`
   and `internal/baseline/assets/templates/guides/python.md` (Data Model 3),
   add the two template entries to `internal/baseline/assets/templates/index.json`,
   the coverage entry to `internal/baseline/assets/coverage.json` (Data Model
   4), the two core entries to `internal/baseline/assets/modules/core.json`
   (Data Model 5; no clause, rule or other entry of the core module changes by
   hand), and `internal/baseline/assets/profiles/python-cli.json` (Data Model
   6).
4. MUST then run, in this order, and hand-edit none of their outputs:
   - `go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1`,
     which records the new `python` module and writes the next version into
     `core.json`, recording both in `internal/baseline/module-versions.json`;
   - `make baseline-digests`, which rewrites the derived files this Task
     declares (measured on a prototype at `dde9d5f5`: the three catalog files
     and four plan characterization goldens under `internal/baseline/testdata/`);
   - `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text`,
     the Managed Refresh of this repository, which rewrites
     `docs/agents/setup-context.json`, then the same command again, which must
     report `File changes: 0`. Run both with `HOME` set to a temporary
     directory, so the refresh reads and writes no real Roundfix Home.
   A derived file the commands leave unchanged is not edited; a file they
   change that this Task does not declare is a stop-and-report condition.
5. MUST make the declared breaks of `_techspec.md` → Testing Approach (Build
   Order 2) and no other change to existing tests:
   `internal/baseline/catalog_test.go` (`TestEmbeddedCatalog`),
   `internal/baseline/derived_ownership_test.go` (`TestOutputsForCommand`),
   `internal/baseline/clause_characterization_test.go`
   (`characterizedBaselineForce`), `internal/baseline/assets_sync_test.go`
   (the three counts of four), `internal/cli/baseline_human_test.go`
   (`TestBaselineHumanProfileAdaptation`'s menu answer) and
   `internal/cli/baseline_release_gate_test.go`
   (`TestGuidanceCompositionJourney`'s list and its two fixture switches; the
   `python-cli` formatter case runs nothing and says so in a comment).
6. MUST create `internal/baseline/python_profile_test.go` with
   `TestThePythonSetupIsTheCommonPartOfGoAndRust` and
   `TestThePythonCliProfileRendersThePythonGuide` as `_techspec.md` → Testing
   Approach describes them, comparing guide text on whitespace-normalized
   content. Each top-level test calls `t.Parallel()` first.
7. MUST keep `skills/baseline_skill_contract_test.go`
   `TestNoPythonBaselineRuntime` passing unchanged: it guards Roundfix's own
   retired Python runtime, not an adopter's stack (`_techspec.md` → Current
   behavior 8).
8. MUST NOT change the Source Baseline corpus, manifest, accounting or index,
   `internal/baseline/assets/retention/**`, the formatter fixtures, any other
   module, setup or profile by hand, any skill, any user guide page,
   `.roundfixrc.yml`, the Makefile, or any Go production file.

## Subtasks

- [ ] Write the setup, module, templates, index, coverage, core and profile entries.
- [ ] Record the module versions, run `make baseline-digests` and the Managed Refresh twice.
- [ ] Make the declared test changes and add the two tests.

## Acceptance Criteria

- [ ] The catalog lists five profiles and five setups, and the `python` setup equals the common part of `go` and `rust`.
- [ ] A greenfield `python-cli` plan is ready and renders the Python guide and root pointer, and the release journey passes for `python-cli`.
- [ ] A second Managed Refresh of this repository reports no file change.

## Consumers

- the `python-cli` profile and `python` module — consumers: `./internal/baseline`, `./internal/cli`

## Context

- instruction: `docs/adr/0322-a-python-profile-takes-the-skills-the-go-and-rust-setups-share.md`
- instruction: `internal/baseline/assets/profiles/rust-cli.json`
- instruction: `internal/baseline/assets/modules/rust.json`
- instruction: `skills/baseline_skill_contract_test.go`
- creates: `internal/baseline/assets/setups/python.json`
- creates: `internal/baseline/assets/modules/python.json`
- creates: `internal/baseline/assets/templates/root/python.md`
- creates: `internal/baseline/assets/templates/guides/python.md`
- creates: `internal/baseline/assets/profiles/python-cli.json`
- interface: `internal/baseline/assets/templates/index.json`
- interface: `internal/baseline/assets/coverage.json`
- interface: `internal/baseline/assets/modules/core.json`
- interface: `internal/baseline/module-versions.json`
- interface: `internal/baseline/testdata/catalog.diagnostics.golden.json`
- interface: `internal/baseline/testdata/catalog.digest`
- interface: `internal/baseline/testdata/catalog.normalized.json`
- interface: `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`
- interface: `docs/agents/setup-context.json`
- interface: `internal/baseline/catalog_test.go`
- interface: `internal/baseline/derived_ownership_test.go`
- interface: `internal/baseline/clause_characterization_test.go`
- interface: `internal/baseline/assets_sync_test.go`
- interface: `internal/cli/baseline_human_test.go`
- interface: `internal/cli/baseline_release_gate_test.go`
- creates: `internal/baseline/python_profile_test.go`

## Verification

- `for name in TestThePythonSetupIsTheCommonPartOfGoAndRust TestThePythonCliProfileRendersThePythonGuide; do go test -list "^${name}\$" ./internal/baseline | grep -qx "$name" || { printf 'missing test %s\n' "$name" >&2; exit 1; }; done; go test -count=1 ./internal/baseline` — expected: exit 0. Before this Task neither test exists, so the command fails at its first check. Afterwards both exist and the whole package passes: the catalog, the module record, the derived-file characterizations, the clause force characterization, the asset sync and the derived-output enumeration.
- `test -f internal/baseline/assets/profiles/python-cli.json || { echo 'missing python-cli profile' >&2; exit 1; }; grep -qF -- '"python-cli",' internal/cli/baseline_release_gate_test.go || { echo 'release journey lacks python-cli' >&2; exit 1; }; go test -count=1 -run '^(TestGuidanceCompositionJourney|TestBaselineHumanProfileAdaptation)$' ./internal/cli && go test -count=1 -run '^TestNoPythonBaselineRuntime$' ./skills && GOOS=windows go vet ./internal/baseline ./internal/cli` — expected: exit 0. Before this Task the profile does not exist, so the command fails at its first check. Afterwards the release journey plans, applies and re-plans every built-in profile including `python-cli`, the human menu selects the right profile, the retired-runtime guard still passes, and both packages' tests compile for Windows.

## References

- `_prd.md` → Goals; Core Features 2-4; Success Metrics 1, 2 and 4; User Stories 1-3
- `_techspec.md` → Current behavior 4-8; Data Models 1-6; API Contracts 2-4; Invariants 3-5; Testing Approach; Build Order 2
- ADR-0322
