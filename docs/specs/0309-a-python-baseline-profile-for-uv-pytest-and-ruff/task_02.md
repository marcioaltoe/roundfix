---
task: task_02
spec: 0309-a-python-baseline-profile-for-uv-pytest-and-ruff
status: pending
type: backend
complexity: high
---

# Task 02: The python profile governs uv, ruff and pytest and dispatches the Python skills

## Overview

This Task ships the `python` Baseline Profile on task_01's `python` setup:
the `python` module with its guide, root block and nine dispatched skills (four
of them only when the repository has chosen them) and its gate rules for
uv, the offline lock check, ruff, ty and pytest, the templates, the
coverage entry, the core module's hierarchy and semantic-owner entries and the
profile itself. The sanctioned commands regenerate every derived file and this
repository's Setup Manifest, the existing tests that enumerate the catalog
learn the new entries, and two new tests hold the dispatch and the rendered
guide (`_prd.md` → Core Features 2-3; ADR-0322).

## Requirements

1. MUST answer the Backlog Entry
   [`references/2026-10-22-a-baseline-profile-for-python-repositories.md`](references/2026-10-22-a-baseline-profile-for-python-repositories.md)
   of 2026-10-22, its expected items 1, 2 and 3.
2. MUST create `internal/baseline/assets/modules/python.json` (`_techspec.md`
   → Data Model 3, the dispatch of API Contract 3 and the rules and clauses of
   API Contract 4 verbatim, module `version` 1 before the record command
   runs), `internal/baseline/assets/templates/root/python.md` and
   `internal/baseline/assets/templates/guides/python.md` (Data Model 4), add
   the two template entries to `internal/baseline/assets/templates/index.json`,
   the coverage entry to `internal/baseline/assets/coverage.json` (Data Model
   5), the two core entries to `internal/baseline/assets/modules/core.json`
   (Data Model 6; no clause, rule or other entry of the core module changes by
   hand), and `internal/baseline/assets/profiles/python.json` (Data Model 7:
   id `python`, title "Python repository (CLI, backend or worker)", no
   `cli-surface` module and no `rule.cli.*` rule).
3. MUST then run, in this order, and hand-edit none of their outputs:
   - `go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1`,
     which records the new `python` module and writes the next version into
     `core.json`, recording both in `internal/baseline/module-versions.json`;
   - `make baseline-digests`, which rewrites the derived files this Task
     declares;
   - `go build -buildvcs=false -o "$tmp/roundfix" ./cmd/roundfix` under the
     Run's environment, then `HOME="$(mktemp -d)" "$tmp/roundfix" baseline update --repo . --no-skills --no-history --yes --format text`
     twice, the Managed Refresh of this repository, which rewrites
     `docs/agents/setup-context.json` and then reports `File changes: 0`.
   A derived file the commands leave unchanged is not edited; a file they
   change that this Task does not declare is a stop-and-report condition.
4. MUST make the declared breaks of `_techspec.md` → Testing Approach (Build
   Order 2) and no other change to existing tests:
   `internal/baseline/catalog_test.go` (`TestEmbeddedCatalog`),
   `internal/baseline/derived_ownership_test.go` (`TestOutputsForCommand`),
   `internal/baseline/clause_characterization_test.go`
   (`characterizedBaselineForce`), `internal/cli/baseline_human_test.go`
   (`TestBaselineHumanProfileAdaptation`'s `answers` script selects `5`) and
   `internal/cli/baseline_release_gate_test.go`
   (`TestGuidanceCompositionJourney`'s list and its two fixture switches; the
   `python` formatter case runs nothing and says so in a comment).
5. MUST create `internal/baseline/python_profile_test.go` with
   `TestThePythonModuleDispatchesItsSkills` and
   `TestThePythonProfileRendersThePythonGuide` as `_techspec.md` → Testing
   Approach describes them, comparing guide text on whitespace-normalized
   content. Each top-level test calls `t.Parallel()` first.
6. MUST keep `skills/baseline_skill_contract_test.go`
   `TestNoPythonBaselineRuntime` passing unchanged: it guards Roundfix's own
   retired Python runtime, not an adopter's stack (`_techspec.md` → Current
   behavior 10).
7. MUST NOT change the setup snapshots, the Source Baseline corpus, manifest,
   accounting or index, `internal/baseline/assets/retention/**`, the formatter
   fixtures, any module other than `core` and `python`, any existing profile,
   any skill, any user guide page, `.roundfixrc.yml`, the Makefile or any Go
   production file.

## Subtasks

- [ ] Write the module, templates, index, coverage, core and profile entries.
- [ ] Record the module versions, run `make baseline-digests` and the Managed Refresh twice.
- [ ] Make the declared test changes and add the two tests.

## Acceptance Criteria

- [ ] The catalog lists five profiles, and the `python` module requires and dispatches the nine Python skills, four of them only when chosen, and its guide names the offline lock check and ty in the gate.
- [ ] A greenfield `python` plan is ready and renders the Python guide, the nine skills in the dispatch guide and the root pointer, and no Go, Rust, TUI or CLI guide; the release journey passes for `python`.
- [ ] A second Managed Refresh of this repository reports no file change.

## Consumers

- the `python` profile and `python` module — consumers: `./internal/baseline`, `./internal/cli`

## Context

- instruction: `docs/adr/0322-a-python-profile-imports-the-published-upstream-python-preset.md`
- instruction: `internal/baseline/assets/profiles/rust-cli.json`
- instruction: `internal/baseline/assets/modules/rust.json`
- instruction: `internal/baseline/stack_force_rust_test.go`
- instruction: `skills/baseline_skill_contract_test.go`
- creates: `internal/baseline/assets/modules/python.json`
- creates: `internal/baseline/assets/templates/root/python.md`
- creates: `internal/baseline/assets/templates/guides/python.md`
- creates: `internal/baseline/assets/profiles/python.json`
- creates: `internal/baseline/python_profile_test.go`
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
- interface: `internal/cli/baseline_human_test.go`
- interface: `internal/cli/baseline_release_gate_test.go`

## Verification

- `for name in TestThePythonModuleDispatchesItsSkills TestThePythonProfileRendersThePythonGuide; do go test -list "^${name}\$" ./internal/baseline | grep -qx "$name" || { printf 'missing test %s\n' "$name" >&2; exit 1; }; done; go test -count=1 ./internal/baseline` — expected: exit 0. Before this Task neither test exists, so the command fails at its first check. Afterwards both exist and the whole package passes: the catalog, the module record, the dispatch rule that every dispatched skill is required, the derived-file characterizations and the clause force characterization.
- `test -f internal/baseline/assets/profiles/python.json || { echo 'missing python profile' >&2; exit 1; }; if grep -qE -- 'cli-surface|rule[.]cli[.]' internal/baseline/assets/profiles/python.json; then echo 'python profile selects the CLI surface' >&2; exit 1; fi; grep -qF -- '"python",' internal/cli/baseline_release_gate_test.go || { echo 'release journey lacks python' >&2; exit 1; }; go test -count=1 -run '^(TestGuidanceCompositionJourney|TestBaselineHumanProfileAdaptation)$' ./internal/cli && go test -count=1 -run '^TestNoPythonBaselineRuntime$' ./skills` — expected: exit 0. Before this Task the profile does not exist, so the command fails at its first check. Afterwards the profile exists without the CLI surface, the release journey plans, applies and re-plans every built-in profile including `python`, the human menu selects the right profile, and the retired-runtime guard still passes.

## References

- `_prd.md` → Goals 3-5; Core Features 2-3; Success Metrics 3 and 4; User Stories 1-3
- `_techspec.md` → Current behavior 6-10; Data Models 3-7; API Contracts 3-5; Invariants 3-5; Testing Approach; Build Order 2
- ADR-0322
