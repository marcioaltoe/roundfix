---
spec: 0309-a-python-baseline-profile-for-uv-pytest-and-ruff
prd: _prd.md
created: 2026-10-23
---

# A Python Baseline Profile for uv, pytest and ruff — Technical Spec

## Executive Summary

Three changes. First, the Baseline learns a second kind of derived Setup
Snapshot: a source `{"type": "common", "setups": [...]}` holds, in the first
component's order, the skills and activation bundles every component carries
with equal entries; the catalog validates it beside the composed kind and the
asset sync writes it after its components. Second, the catalog gains the
`python` common snapshot of `go` and `rust`, a `python` module with its guide,
root block, templates and coverage entry, the core module's hierarchy and
semantic-owner entries for them, and the `python` profile; the sanctioned
commands regenerate the derived files and this repository's Setup Manifest,
and the existing tests that enumerate profiles, setups, clauses or the
profile menu learn the new entries. Third, the user guide lists the profile
and the glossary defines **Common Setup Snapshot**.

The trade-off accepted: the profile carries no Python-specific skill until the
upstream repository publishes one; its guide carries the rules instead, and
the curation request is already filed outside this repository.

## Project Constraints

- Identifier strategy: not applicable — no project-owned identifier or record
  schema is minted; catalog identifiers follow the existing dotted naming.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path or credential;
  tests read the embedded catalog, temporary repositories and a synthetic
  asset-sync source. Source: `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — Build Order 1 follows ADR-0322 and
  ADR-0204 for the common snapshot and ADR-0191 for its single writer, with
  ADR-0206 and ADR-0246 standing through its components; Build Order 2
  follows ADR-0322 for the module, setup and profile, with ADR-0219 and
  ADR-0241 standing for repository-owned drafts of it; Build Order 3
  documents ADR-0322's term. All are named with their titles in `_prd.md` →
  Project Constraints. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AP approval of 2026-10-22 ("Preset de baseline Python"; authorization "same
  as cycle AO", carrying the cycle AJ grant "Mexer em Governed Paths
  declarados"). Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0309-a-python-baseline-profile-for-uv-pytest-and-ruff/_authorization.md`;
  bounded files: `docs/agents/setup-context.json`,
  `internal/baseline/assets/coverage.json`,
  `internal/baseline/assets/modules/core.json`,
  `internal/baseline/assets/modules/python.json`,
  `internal/baseline/assets/profiles/python.json`,
  `internal/baseline/assets/setups/python.json`,
  `internal/baseline/assets/templates/guides/python.md`,
  `internal/baseline/assets/templates/index.json`,
  `internal/baseline/assets/templates/root/python.md`,
  `internal/baseline/derived_ownership_test.go`,
  `internal/cli/baseline_human_test.go`,
  `internal/cli/baseline_release_gate_test.go`.

## Current behavior

Read and characterized on 2026-10-23 at `dde9d5f5`. Line numbers move when
other Specs merge; identifiers do not.

1. `internal/baseline/setup_composition.go`: `composeSetupSnapshot` builds the
   ordered union of components and `setupCompositionComponents` refuses fewer
   than two distinct identifiers, an unknown component and a component whose
   source type is `composed`.
2. `internal/baseline/catalog_validate.go` `validateSetups`: only a source of
   type `composed` is derived and checked, reporting
   `catalog.setup.composition.invalid`, `.conflict` and `.drift` (detail
   "skills or activation bundles differ from the component union").
3. `internal/baseline/assets_sync.go` `syncAssets`: setups of type `composed`
   are skipped in the upstream refresh loop and recomposed after it; any other
   type is refreshed from an upstream `setups/<id>.txt`.
4. Setup Snapshots: `go` (41 skills), `rust` (34), `typescript`, and the
   composed `go-cli-typescript-bun`. `go` and `rust` share thirty skills with
   deeply equal entries and both bundles `bundle.delivery` and `bundle.qa`.
5. Profiles: `go-cli-tui`, `go-cli-typescript-monorepo`, `rust-cli`,
   `standard-typescript-monorepo`. `rust-cli` selects `core`,
   `context-workflow`, `rust`, `cli-surface`, `autonomous-work`,
   `spec-workflow`, `external-triage`, `secondbrain`, `repository-extension`,
   declares formatter kind `none` and has no formatter fixture.
6. `internal/baseline/assets/modules/core.json`: the instruction hierarchy's
   `stack` level lists `root.typescript` … `root.go`, `root.rust`; the
   semantic owners list `guide.go` then `guide.rust`.
7. Tests that enumerate the catalog: `TestEmbeddedCatalog`
   (`catalog_test.go`, the four profile IDs), `TestOutputsForCommand`
   (`derived_ownership_test.go`, the `make baseline-digests` outputs),
   `characterizedBaselineForce` (`clause_characterization_test.go`, every
   clause's force), three asset-sync tests in `assets_sync_test.go` (four
   findings, four setups) and three `composed` skips in its helpers,
   `TestBaselineHumanProfileAdaptation` (`internal/cli/baseline_human_test.go`,
   menu answer `4` selects `standard-typescript-monorepo`) and
   `TestGuidanceCompositionJourney`
   (`internal/cli/baseline_release_gate_test.go`, the four profiles and two
   per-profile fixture switches that fail on an unknown profile).
8. `skills/baseline_skill_contract_test.go` `TestNoPythonBaselineRuntime`
   refuses `python3` and `Python fallback` in the user guide's
   `context-driven-development.md`, `commands.md`, the README and the
   `setup-context-driven` and `roundfix` skills: it guards Roundfix's own
   retired Python runtime, not an adopter's stack, and stays unchanged.
9. `docs/user-guide/context-driven-development.md` `### Profiles` names the
   four profiles.

## System Architecture

- `internal/baseline/setup_composition.go`, `catalog_validate.go`,
  `assets_sync.go`: the common kind.
- `internal/baseline/assets/setups/python.json` (new),
  `assets/modules/python.json` (new), `assets/templates/root/python.md` and
  `assets/templates/guides/python.md` (new), `assets/templates/index.json`,
  `assets/coverage.json`, `assets/modules/core.json`,
  `assets/profiles/python.json` (new).
- Derived files, the module version record and this repository's Setup
  Manifest: regenerated only by the sanctioned commands.
- Tests: `internal/baseline/common_setup_test.go` and
  `internal/baseline/python_profile_test.go` (new), and the existing tests of
  Current behavior 7.
- `docs/user-guide/context-driven-development.md`, `CONTEXT.md`,
  `docs/references/skill-coverage.json`, `docs/references/behavior-surfaces.json`.

No package, import edge or exported symbol changes.

## Implementation Design

### Interfaces

```go
// internal/baseline/setup_composition.go (unexported)
func derivedSetupSource(setup document) (kind string, derived bool) // "composed" or "common"
func deriveSetupSnapshot(id, kind string, components []document) (document, error)
func commonSetupSnapshot(id string, components []document) (document, error)
```

### Data Models

1. A common Setup Snapshot is a Setup Snapshot document whose `source` is
   `{"setups": ["<first>", "<second>", ...], "type": "common"}`, serialized by
   the asset sync's snapshot struct (field order `schemaVersion`, `id`,
   `version`, `source`, `digest`, `skills`, `activationBundles`), with
   `version` `0.0.1` and `digest` the canonical SHA-256 of `{"skills": ...,
   "activationBundles": ...}` (of the skills alone when no bundle is shared),
   as the composed kind computes it.
2. The `python` module: `schemaVersion` `setup-context-driven/module-v3`,
   `kind` `language-runtime`, `title` "Python repository guidance",
   `dependsOn` `["core"]`, empty `requiredSkills`, `skillDispatch`,
   `requiredDecisions` and `repositoryExtensions`; root block `root.python`
   (version 1, template `template.root.python`, reference
   `reference.root.python` with token `reference.python`, ownership `setup`,
   managed id `guide.python`); supporting guide `guide.python` (version 1,
   path `docs/agents/python.md`, template `template.guide.python`, rules in
   API Contract 3's order); its own `version` is the one the record command
   writes.
3. Templates: `template.root.python` (version 1, kind `root-block`, path
   `root/python.md`, tokens `["reference.python"]`) whose content is "### Python",
   a blank line, and "- Python implementation rules are mandatory:
   {{reference.python}}."; `template.guide.python` (version 1, kind `guide`,
   path `guides/python.md`, tokens `["artifact.rules"]`) whose content is the
   heading "# Python", a blank line, the paragraph "These rules govern the
   repository's Python project: its package, its entry points and its tests,
   managed with uv. Code in another language follows its own guide.", a blank
   line and `{{artifact.rules}}`. Each is inserted in `templates/index.json`
   right before the matching Rust entry.
4. `coverage.json`: `{"id": "coverage.python", "version": 1, "description":
   "Python work runs through uv, keeps ruff and pytest gates honest, and keeps
   tests hermetic.", "delegationAliases": ["python", "python guidance"]}` right
   before `coverage.rust`.
5. `core.json`: `root.python` between `root.go` and `root.rust` in the `stack`
   level, and `{"managedId": "guide.python", "title": "Python guidance",
   "classifications": ["python"]}` between the Go and Rust semantic owners. No
   clause of the core module changes.
6. `profiles/python.json`: `rust-cli.json` with `id` `python`, `title`
   "Python repository (CLI, backend or worker)", `setup` `python`, `modules`
   `["core", "context-workflow", "python", "autonomous-work",
   "spec-workflow", "external-triage", "secondbrain",
   "repository-extension"]` (the `rust-cli` list with `python` in place of
   `rust` and without `cli-surface`), and `requiredRules` with
   `rule.python.uv-environment`, `rule.python.ruff-and-pytest-gates` and
   `rule.python.hermetic-tests` in place of the three Rust rules and without
   `rule.cli.output-contract` and `rule.cli.non-interactive`; every other
   field equal. The profile selects no application-surface module: it serves
   a command-line tool, a backend service or a worker alike, and a
   repository adds the rules of its own surface in its repository-owned
   `docs/agents/specific-repository.md`, the `repository-extension` module's
   carrier (maintainer decision of 2026-10-24).

### API Contracts

1. API Contract: the common derivation. `commonSetupSnapshot` keeps each skill
   of the first component, in its order, whose name every other component
   carries with a deeply equal entry, likewise each activation bundle by `id`;
   a name another component carries with a different entry is dropped, not a
   conflict. Components without a shared skill are an error.
   `setupCompositionComponents` refuses a component of either derived kind
   (message "component %q is <kind>") and keeps its other refusals.
   `validateSetups` checks every derived snapshot through `derivedSetupSource`
   and `deriveSetupSnapshot` with the existing codes; the drift detail for a
   common snapshot is exactly "skills or activation bundles differ from the
   components' common part", and the composed detail keeps its bytes. The
   asset sync skips both derived kinds in its upstream loop and derives each
   after it, writing the result through the existing transaction.
2. API Contract: the `python` setup is the asset sync's derivation of the
   common part of `go` then `rust`, written with the sync's serialization;
   the Task creates it with those bytes (the Run cannot reach an upstream
   skills checkout), and the catalog's drift and digest checks hold it.
3. API Contract: the Python guide's rules and clauses, in this order, each
   rule at version 1:
   - `rule.python.uv-environment`, coverage `coverage.python`,
     `coverage.dependencies`:
     - `clause.python.change-dependencies-through-uv`, mandatory: "Change
       dependencies only through uv: `uv add`, `uv add --group dev`, `uv
       remove` and `uv lock --upgrade-package`. Commit `pyproject.toml` and
       `uv.lock` together, and never install into the project environment with
       `pip`."
     - `clause.python.run-tools-through-uv`, mandatory: "Run Python, pytest,
       ruff and every other project tool through `uv run`, so each uses the
       locked project environment and the Python version the project pins,
       never a globally installed interpreter or tool."
     - `clause.python.sync-the-locked-environment`, mandatory: "Create a
       worktree's or CI job's environment with `uv sync --frozen` and the
       development group the gate needs, so it installs exactly what `uv.lock`
       records; a step that must also refuse a lockfile that no longer matches
       `pyproject.toml` uses `--locked` instead."
   - `rule.python.ruff-and-pytest-gates`, coverage `coverage.python`,
     `coverage.verification-integrity`:
     - `clause.python.iterate-with-focused-commands`, mandatory: "Use focused
       commands such as `uv run pytest <node id>` and `uv run ruff check
       <path>` while iterating, then run the selected repository
       Verification."
     - `clause.python.ruff-findings-fail-the-gate`, prohibited: "Do not hide a
       ruff finding: never pass `--exit-zero`, add a blanket `# noqa`, or
       widen the ruff configuration to make a finding disappear. A `# noqa`
       names the rule code it silences and the reason."
     - `clause.python.select-a-new-test-by-its-node-id`, mandatory: "A Task
       Verification selects a new test by its pytest node id, such as
       `tests/test_report.py::test_totals`, and never passes that test's file
       or directory in the same command. pytest collects a node id once when
       another argument already contains it and exits 0 even before the test
       exists, so the command could never fail before the work; a missing node
       id alone exits 4 and an empty selection exits 5. `roundfix spec check`
       reports a command that names a test both ways."
   - `rule.python.hermetic-tests`, coverage `coverage.python`,
     `coverage.verification`:
     - `clause.python.keep-tests-hermetic`, mandatory: "Keep tests hermetic:
       set or clear with `monkeypatch` every environment variable the code
       under test reads, write only under `tmp_path`, and never read the
       host's credentials, home directory, network or tool state."
     - `clause.python.outside-data-skips-by-name`, mandatory: "A test that
       needs data from outside the repository, such as a mounted drive or a
       credentialed service, calls `pytest.skip` with a reason that names the
       environment variable or path it lacks, so the selected Verification
       passes on a machine without that data. A Task Verification never
       depends on such a test."
     - `clause.python.exercise-every-platform-branch`, mandatory: "When a
       change touches code that branches on the operating system, such as on
       `sys.platform`, `os.name` or `platform.system()`, exercise every branch
       it names: test each one with `monkeypatch` over the value the code
       reads, or run the repository's type checker for each platform with its
       platform option, because a test run on the host executes only the
       host's branch."
     - `clause.python.keep-entry-points-thin`, mandatory: "Keep each entry
       point thin, whether a console script, a web application factory or a
       worker's start-up: it reads its input and calls behavior that lives in
       the package. Test behavior through the package and the entry point's
       observable contract: for a console script, standard output, standard
       error and the exit code it raises with `SystemExit`; for a web
       application or a worker, the responses or messages it produces."
4. API Contract: rendering. A greenfield plan for the `python` profile with the usual
   decisions in a repository whose Makefile declares `verify` and
   `verify-incremental` is `ready`, renders `docs/agents/python.md` with the
   ten clauses and the root pointer "Python implementation rules are
   mandatory", renders no `docs/agents/go.md`, `docs/agents/rust.md`,
   `docs/agents/tui.md` or `docs/agents/cli.md`, records `python` in its Setup Manifest, applies,
   and re-plans with no file change. Profile alignment, decisions and
   capabilities behave as for `rust-cli`, without the CLI surface's
   `agentic-cli-design` dispatch.

### Surface Transcripts

None. No command, flag, exit code or output line changes; the profile menu
gains one entry.

### Invariants

```text
1. A derived Setup Snapshot holds no list of its own choosing; the catalog refuses one that differs from its derivation.
2. No derived snapshot is a component of another.
3. Existing setups, profiles and clauses keep their bytes; only derived files, the module record and this repository's Setup Manifest change through the sanctioned commands.
4. The python profile restores no skill outside the python setup, the python setup holds no Go, TUI or Rust skill, and the profile selects no application-surface module.
5. No Roundfix code or test runs a Python tool.
```

## Coverage Map

Deviations from the PRD: None.

- Goal 3, Core Features 1 and 2, Success Metrics 2 and 3, User Story 4, Invariants 1, 2 and 4 → API Contracts 1 and 2
- Goals 1 and 2, Core Features 3 and 4, Success Metrics 1 and 4, User Stories 1-3, Invariant 5 → API Contracts 3 and 4
- Goal 4, Success Metric 1, Invariant 3 → Data Models 5 and 6, Testing Approach (declared breaks)
- Core Feature 5, Success Metric 5 → Build Order 3

## Absence Evidence

1. The common kind — command: `grep -rn '"common"' internal/baseline/*.go`; result: no match at `dde9d5f5`.
2. The Python profile, module and setup — command: `ls internal/baseline/assets/profiles internal/baseline/assets/modules internal/baseline/assets/setups | grep -i python`; result: no match at `dde9d5f5`.
3. The glossary term — command: `grep -n "Common Setup Snapshot" CONTEXT.md`; result: no match at `dde9d5f5`.

## ADR Alternatives

- ADR-0322: none of its rejected alternatives is taken; the setup is neither a union of `go` and `rust`, nor hand-written, nor the `go` list.
- ADR-0204: its rejected alternatives (a profile naming several setups; deriving at load time) are not taken; the common snapshot is a file the sync writes, read by the existing readers.
- ADR-0191: its rule stands; the sync remains the only writer, and the Task writes the `python` snapshot with the sync's own derivation bytes, which the catalog checks.
- ADR-0206: no rejected alternative is taken; the `python` snapshot names its components by their upstream setup names, and a skill dropped upstream leaves it with its components.
- ADR-0219: no rejected alternative is taken; a repository-owned draft adapted from the `python` profile binds to it like to any built-in profile.
- ADR-0241: no rejected alternative is taken; a repository-owned profile that selects the `python` module takes its skill contracts from the embedded snapshots, of which `python` is one.
- ADR-0246: no rejected alternative is taken; a retired skill leaves the common part when it leaves its components, and no module requires it.

## Integration Points

The embedded Baseline catalog and its derived files; the asset sync over a
synthetic source; the repository's Managed Refresh. No test, Verification
command or QA row reads the upstream skills checkout, runs a Python tool,
reaches the network, reads the real `~/.roundfix`, reads the Secondbrain or
reads another repository.

## Testing Approach

Every new top-level test and subtest calls `t.Parallel()` first and passes no
real-time duration; the new files use no platform API and compile for Windows.

- `internal/baseline/common_setup_test.go` (new, Build Order 1):
  - `TestACommonSetupKeepsOnlyTheSkillsEveryComponentCarries`: explicit
    components `first` (skills `a`, `c`, `differs`, `b`; bundles `shared`,
    `only-first`) and `second` (skills `b`, `differs` with another path, `c`,
    `d`; bundles `last`, `shared` equal) give skills `c`, `b` and bundle
    `shared`, source type `common` with setups `first`, `second`, and the
    canonical digest; components without a shared skill return an error.
  - `TestACommonSetupThatDiffersFromItsComponentsIsReported`: on
    `cloneEmbeddedAssets`, a fixture `setups/common-example.json` built by the
    derivation from `go` and `rust` loads; with its first skill dropped it
    reports `catalog.setup.composition.drift` whose detail contains
    "components' common part"; with components `["go"]` or `["go",
    "go-cli-typescript-bun"]` it reports `catalog.setup.composition.invalid`.
  - `TestTheAssetSyncWritesACommonSetupAfterItsComponents`: in an asset-sync
    target and synthetic source, a stale `common-example` snapshot (first
    skill dropped) is rewritten by a refresh to the common part of the
    refreshed `go` and `rust`, with source type `common`.
- `internal/baseline/python_profile_test.go` (new, Build Order 2):
  - `TestThePythonSetupIsTheCommonPartOfGoAndRust`: the embedded `python`
    setup has source type `common` with `go`, `rust`, equals the derivation,
    holds none of the toolchain skills of Success Metric 2 and holds the three
    language-neutral ones.
  - `TestThePythonProfileRendersThePythonGuide`: through
    `newPlanRepository` and `planTestDecisions` with greenfield preservation,
    API Contract 4's plan, guide phrases (on whitespace-normalized text),
    absent guides (including `docs/agents/cli.md`), root pointer and
    manifest profile.
- Declared breaks (Build Order 1): the three `composed` skips in
  `assets_sync_test.go` (the template check, the maintained parity row and
  `buildAssetsSyncSource`) skip both derived kinds; their assertions keep
  their meaning.
- Declared breaks (Build Order 2), each moving to the five-profile, five-setup
  catalog: `TestEmbeddedCatalog` lists `python`; `TestOutputsForCommand`
  adds `internal/baseline/assets/setups/python.json` and
  `internal/baseline/assets/profiles/python.json` to its enumeration;
  `characterizedBaselineForce` lists the ten Python clauses with API Contract
  3's force; `TestAssetsSyncCheckIsReadOnlyAndReportsDrift`,
  `TestBaselineAssetsSyncRefreshProducesCanonicalTreeAndIsIdempotent` and
  `TestAssetsSyncCompatibilityMatchesMaintainedPythonContract` count five
  setups and findings; `TestBaselineHumanProfileAdaptation` answers `5` for
  `standard-typescript-monorepo`; `TestGuidanceCompositionJourney` lists
  `python` and its two fixture switches gain a `python` case (a
  `pyproject.toml` and an empty `src/example/__init__.py`; no formatter run).
  No other existing test changes; the prototype passed every other test of
  `internal/baseline`, `internal/cli`, `skills`, `internal/spec` and
  `internal/speccheck` with these changes.
- Prototype record (2026-10-23, scratch worktree at `dde9d5f5`, with the
  earlier profile shape `python-cli` and its `cli-surface` module; the
  2026-10-24 amendment was checked by reading, `_prd.md` → Validation): the
  five new tests failed to exist before and passed after. These sabotages each failed
  the named test: the common derivation keeping a same-named entry without
  comparing it → `TestACommonSetupKeepsOnlyTheSkillsEveryComponentCarries`;
  the common drift detail replaced by the composed one →
  `TestACommonSetupThatDiffersFromItsComponentsIsReported/skill_dropped`; the
  asset sync's upstream loop skipping only `composed` →
  `TestTheAssetSyncWritesACommonSetupAfterItsComponents`; the node-id clause's
  "selects a new test by its pytest node id" reworded →
  `TestThePythonProfileRendersThePythonGuide`; the first skill removed from
  `setups/python.json` → catalog validation, failing
  `TestThePythonSetupIsTheCommonPartOfGoAndRust`.

## Build Order

1. The common kind: `derivedSetupSource`, `deriveSetupSnapshot`,
   `commonSetupSnapshot`, the refusal of derived components, the validator,
   the asset sync, the helper skips and `common_setup_test.go` (no
   dependency).
2. The Python profile: the `python` setup, module, templates and index
   entries, coverage entry, core entries and profile; the module record
   command, `make baseline-digests` and the Managed Refresh twice; the
   declared breaks and `python_profile_test.go` (depends on 1, whose validator
   accepts the common kind and whose helper skips the sync tests need).
3. Documentation: the user guide's profile list and paragraph, **Common Setup
   Snapshot** in the glossary, the Coverage Review and the Behavior Surface
   re-record last (depends on 2, which ships what it describes).
4. QA gate (depends on: 1, 2, 3).

No Spec is listed in `requires`. Spec 0307 delivers first and changes the
core module and the derived digests; Build Order 2 regenerates on the tree it
starts on and names no version or digest value.

## Risks & Considerations

- A later upstream change can shrink the common part: a skill only one of
  `go` and `rust` keeps leaves the `python` setup at the next asset sync. The
  sync's diff shows it; a skill a module requires would fail the catalog
  check, and the `python` module requires none.
- The Managed Refresh of an adopter with the `python` profile writes
  `docs/agents/python.md`; an existing repository-owned Python guide at that
  path follows the Baseline's ordinary carrier and preservation rules.
- Rebasing over Spec 0307 conflicts in `core.json`'s version line, the module
  record and the derived goldens; the Task resolves them by rerunning the
  sanctioned commands on the rebased tree, never by hand.

## Glossary

- adds: **Common Setup Snapshot**

## Decisions

- The setup is a common snapshot of `go` and `rust`; see ADR-0322.
- The module dispatches no skill until upstream publishes Python skills.
- The guide's node-id clause names `roundfix spec check`'s existing finding by
  behavior, not by code.
- The starter files are Spec 0310's.
- The profile is named `python` and fixes no application surface: no
  `cli-surface` module and no `rule.cli.*` rules, the hermetic-tests rule
  covers `coverage.verification`, and the entry-point clause speaks of every
  kind of entry point (maintainer decisions of 2026-10-24, quoted in
  `_prd.md` → Decisions and `_authorization.md`).
- `agentic-cli-design` stays in the `python` setup, because the setup is the
  derived common part of `go` and `rust`, but no module of the profile
  requires or dispatches it; a command-line adopter names it, with its own
  surface rules, in its repository-owned `docs/agents/specific-repository.md`
  (the `repository-extension` module).
