---
spec: 0309-a-python-baseline-profile-for-uv-pytest-and-ruff
prd: _prd.md
created: 2026-10-23
---

# A Python Baseline Profile for uv, pytest and ruff — Technical Spec

## Executive Summary

Three changes, and no production Go code. First, the asset sync imports the
published upstream `python` preset at
`65d5f365cace27c4b468db1723213723b10c870b`, fetched from the local
`~/dev/skills` checkout and seeded the way Spec 0208 seeded its renamed
snapshots. Spec 0317, delivered first, has already moved `go`, `rust`,
`typescript` and the composed `go-cli-typescript-bun` to that commit, so the
sync writes only `python.json`, and a repository-owned profile never meets
two contracts for one skill (ADR-0241). The maintained asset-sync parity row
follows through its sanctioned transform; no skill of this repository is
restored. Second, the
catalog gains the `python` module with its guide, root block, templates,
coverage entry and nine dispatched skills, the core module's hierarchy and
semantic-owner entries for it, and the `python` profile; the sanctioned
commands regenerate the derived files. Third, the user guide lists the
profile.

The trade-off accepted: this Spec cannot be queued before Spec 0317 merges,
because its import relies on the other snapshots already naming the commit.

## Project Constraints

- Identifier strategy: not applicable — no project-owned identifier or record
  schema is minted; catalog identifiers follow the existing dotted naming.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path or credential;
  the sync reads a repository fetched from a local checkout, and tests read
  the embedded catalog, temporary repositories and a synthetic asset-sync
  source. Source: `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — Build Order 1 follows ADR-0322,
  ADR-0191, ADR-0206, ADR-0204, ADR-0241, ADR-0246 and ADR-0072; Build Order 2
  follows ADR-0322 for the module and profile, with ADR-0219 and ADR-0241
  standing for repository-owned drafts of it; Build Order 3 documents
  ADR-0322. All are named with their titles in `_prd.md` → Project
  Constraints. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AP approval of 2026-10-22 ("Preset de baseline Python"; authorization "same
  as cycle AO", carrying the cycle AJ grant "Mexer em Governed Paths
  declarados") and the decision of 2026-10-24 "Emendar a 0309 agora
  (Recomendado)". Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0309-a-python-baseline-profile-for-uv-pytest-and-ruff/_authorization.md`;
  bounded files:
  `docs/agents/setup-context.json`,
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

Read and characterized on 2026-10-24 at `64c56d7f`. Line numbers move when
other Specs merge; identifiers do not.

1. `internal/baseline/assets_sync.go` `syncAssets` lists every
   `assets/setups/*.json`, rebuilds each non-composed snapshot from the
   upstream `setups/<id>.txt` of one clean checkout (one `source.ref` for all),
   keeping Roundfix-owned entries with their minimum versions and the
   activation bundles from the current file, recomposes the composed ones, and
   validates the catalog before writing through its transaction. A snapshot
   that has no file is never created, so a new one starts as a seed (Spec
   0208: "Each new snapshot is seeded as a byte copy of its predecessor with
   only `"id"` and `"source.path"` changed").
2. After Spec 0317 the snapshots `go`, `rust`, `typescript` and the composed
   `go-cli-typescript-bun` name `marcioaltoe/skills` at
   `65d5f365cace27c4b468db1723213723b10c870b`, hold the six `ai-memory-*`
   skills and no `handoff`, and `council` and `the-fool` stay out
   (ADR-0246). Upstream added `python.txt` at `32e4143e` (47 entries) and, at
   `65d5f36`, replaced its `handoff` with the six ai-memory skills (52
   entries).
3. `internal/baseline/skills_restore.go` `repositorySnapshotContracts` refuses
   a skill whose contracts (source, ref, tree digest) differ between embedded
   snapshots (`restore.snapshot-conflict`);
   `TestEmbeddedSetupSnapshotsAgreeOnEveryExternalSkill` holds it for the
   whole catalog.
4. `internal/baseline/repository_skills_snapshot_test.go`
   `TestThisRepositoryHoldsItsRequiredSkillsAtTheSetupSnapshot` and the two
   `SetupSnapshotDifferences` tests compare this repository's
   `.agents/skills/` with the `go` snapshot.
5. `internal/baseline/assets_sync_test.go`
   `TestAssetsSyncCompatibilityMatchesMaintainedPythonContract` compares every
   upstream snapshot with `testdata/parity-corpus/v1/fixtures/asset-sync.json`
   (rows `go`, `rust`, `typescript`); Spec 0208's jq transform rebuilds that
   row from the snapshots, and `make baseline-digests` fills its digests.
6. Profiles `go-cli-tui`, `go-cli-typescript-monorepo`, `rust-cli`,
   `standard-typescript-monorepo`. `rust-cli` selects `core`,
   `context-workflow`, `rust`, `cli-surface`, `autonomous-work`,
   `spec-workflow`, `external-triage`, `secondbrain`, `repository-extension`
   and declares formatter kind `none`.
7. `catalog_validate.go`: one module per dispatched skill
   (`catalog.skill.dispatch.id.duplicate`); a profile's required and
   dispatched skills must be in its setup.
   `TestEveryModuleRequiresEverySkillItDispatches` requires every dispatched
   skill to be required. The `backend` module dispatches `backend-docs`,
   `external-api-adapters`, `integration-contract-testing` and
   `security-best-practices`; `cli-surface` dispatches `agentic-cli-design`.
8. `internal/baseline/assets/modules/core.json`: the instruction hierarchy's
   `stack` level lists `root.go`, `root.rust`; the semantic owners list
   `guide.go` then `guide.rust`.
9. Tests that enumerate the catalog: `TestEmbeddedCatalog` (`catalog_test.go`,
   four profiles), `TestOutputsForCommand` (`derived_ownership_test.go`),
   `characterizedBaselineForce` (`clause_characterization_test.go`), four
   counts of four in `assets_sync_test.go`,
   `TestBaselineHumanProfileAdaptation` (`internal/cli/baseline_human_test.go`,
   the menu answer `4` in its `answers` script selects
   `standard-typescript-monorepo`) and `TestGuidanceCompositionJourney`
   (`internal/cli/baseline_release_gate_test.go`, four profiles and two
   per-profile fixture switches that fail on an unknown profile).
10. `skills/baseline_skill_contract_test.go` `TestNoPythonBaselineRuntime`
    refuses `python3` and `Python fallback` in Roundfix's own documents and
    skills: it guards Roundfix's retired Python runtime, not an adopter's
    stack, and stays unchanged. The upstream Streamlit skill mentions
    `python3`, but it is never copied into this repository.
11. `docs/user-guide/context-driven-development.md` `### Profiles` names the
    four profiles.

## System Architecture

- `internal/baseline/assets/setups/*.json`: the five snapshots, written by the
  asset sync (`python.json` seeded first).
- `internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json`: the
  sanctioned transform.
- `internal/baseline/assets/modules/python.json`, `templates/root/python.md`,
  `templates/guides/python.md`, `profiles/python.json` (new);
  `templates/index.json`, `coverage.json`, `modules/core.json`.
- Derived files, the module version record and this repository's Setup
  Manifest: regenerated only by the sanctioned commands.
- Tests: `internal/baseline/python_setup_test.go` and
  `internal/baseline/python_profile_test.go` (new), and the existing tests of
  Current behavior 5 and 9.
- `docs/user-guide/context-driven-development.md`,
  `docs/references/skill-coverage.json`, `docs/references/behavior-surfaces.json`.

No Go production file, package, import edge or exported symbol changes.

## Implementation Design

### Interfaces

None. The Spec changes assets, generated records, tests and documentation.

### Data Models

1. The seed: `internal/baseline/assets/setups/python.json` starts as the
   bytes of `rust.json` with exactly two replacements, `"id": "rust"` →
   `"id": "python"` and `"path": "setups/rust.txt"` → `"path":
   "setups/python.txt"`; the sync rewrites everything else and keeps the
   seed's Roundfix-owned minimum versions and its bundles `bundle.delivery`
   and `bundle.qa`.
2. The parity transform: Spec 0208's jq program, unchanged except that it
   skips composed snapshots (`.manifest.setups = [$s[] | select(.source.type
   != "composed") | row(.)]`), so the row holds `go`, `python`, `rust` and
   `typescript`, as it held the upstream snapshots before. The program text is
   API Contract 1's.
3. The `python` module: `schemaVersion` `setup-context-driven/module-v3`,
   `kind` `language-runtime`, `title` "Python repository guidance",
   `dependsOn` `["core"]`, `requiredSkills` the nine skills of API Contract 3
   in its order, `skillDispatch` one entry per skill with one trigger
   `trigger.python.<skill>`, empty `requiredDecisions` and
   `repositoryExtensions`; root block `root.python` (version 1, template
   `template.root.python`, reference `reference.root.python` with token
   `reference.python`, ownership `setup`, managed id `guide.python`);
   supporting guide `guide.python` (version 1, path `docs/agents/python.md`,
   template `template.guide.python`, rules in API Contract 4's order); module
   `version` 1 before the record command writes its own.
4. Templates: `template.root.python` (version 1, kind `root-block`, path
   `root/python.md`, tokens `["reference.python"]`) whose content is "###
   Python", a blank line, and "- Python implementation rules are mandatory:
   {{reference.python}}."; `template.guide.python` (version 1, kind `guide`,
   path `guides/python.md`, tokens `["artifact.rules"]`) whose content is the
   heading "# Python", a blank line, the paragraph "These rules govern the
   repository's Python project: its package, its entry points and its tests,
   managed with uv. Code in another language follows its own guide.", a blank
   line and `{{artifact.rules}}`. Each is inserted in `templates/index.json`
   right before the matching Rust entry.
5. `coverage.json`: `{"id": "coverage.python", "version": 1, "description":
   "Python work runs through uv, keeps ruff and pytest gates honest, and keeps
   tests hermetic.", "delegationAliases": ["python", "python guidance"]}` right
   before `coverage.rust`.
6. `core.json`: `root.python` between `root.go` and `root.rust` in the `stack`
   level, and `{"managedId": "guide.python", "title": "Python guidance",
   "classifications": ["python"]}` between the Go and Rust semantic owners. No
   clause of the core module changes.
7. `profiles/python.json`: `rust-cli.json` with `id` `python`, `title`
   "Python repository (CLI, backend or worker)", `setup` `python`, `modules`
   `["core", "context-workflow", "python", "autonomous-work",
   "spec-workflow", "external-triage", "secondbrain",
   "repository-extension"]`, and `requiredRules` with
   `rule.python.uv-environment`, `rule.python.verification-gates` and
   `rule.python.hermetic-tests` in place of the three Rust rules and without
   `rule.cli.output-contract` and `rule.cli.non-interactive`; every other
   field equal.

### API Contracts

1. API Contract: the import (task_01), run from the repository root in this
   order, each command verbatim, stopping on the first failure:

   ```bash
   sha=65d5f365cace27c4b468db1723213723b10c870b
   src="$HOME/dev/skills"
   git -C "$src" cat-file -e "$sha^{commit}" || { echo "stop: $src lacks $sha" >&2; exit 1; }
   grep -qF -- "\"ref\": \"$sha\"" internal/baseline/assets/setups/go.json || { echo "stop: Spec 0317 has not moved the setups to $sha" >&2; exit 1; }
   tmp="$(mktemp -d)"
   git init --quiet "$tmp/skills"
   git -C "$tmp/skills" fetch --quiet "$src" "$sha"
   git -C "$tmp/skills" checkout --quiet --detach "$sha"
   git -C "$tmp/skills" remote add origin https://github.com/marcioaltoe/skills.git
   go build -buildvcs=false -o "$tmp/roundfix" ./cmd/roundfix
   HOME="$(mktemp -d)" "$tmp/roundfix" baseline assets sync --source-dir "$tmp/skills/setups" --format text
   HOME="$(mktemp -d)" "$tmp/roundfix" baseline assets sync --source-dir "$tmp/skills/setups" --check --format text
   ```

   The first sync reports `skills.setup-snapshot.updated` for the five
   snapshots; the check exits `0`. Then the parity transform, with this
   program in a temporary file and the fixture rewritten through a temporary
   copy:

   ```jq
   .input.sourceRevision as $rev
   | def row($s): {activationBundles: $s.activationBundles, digest: "", id: $s.id,
       schemaVersion: $s.schemaVersion, version: $s.version,
       source: {path: $s.source.path, ref: $rev, repository: "example/skills", type: "github"},
       skills: [$s.skills[] | if .source.type == "github"
         then {name, path, source: {path: .source.path, ref: $rev, repository: "example/skills", type: "github"}, treeDigest: ""}
         else {minimumVersion, name, path, source: {name: .source.name, type: "repo"}} end]}
       | if .activationBundles == null then del(.activationBundles) else . end;
     .manifest.setups = [$s[] | select(.source.type != "composed") | row(.)]
   | .managedEntryLedger = [.manifest.setups[] | {digest: "", id}]
   | .normalizedOutput.result.findings = [.manifest.setups[] | {
       action: "Review the snapshot diff and run asset validation.",
       code: "skills.setup-snapshot.updated", managedId: ("setup." + .id),
       message: "Setup snapshot was synchronized from the canonical source.",
       path: ("assets/setups/" + .id + ".json"), severity: "info"}]
   | .normalizedOutput.result.summary.info = (.manifest.setups | length)
   ```

   run as `jq --slurpfile s <(jq -s 'sort_by(.id)[]'
   internal/baseline/assets/setups/*.json) -f <program> <fixture>`. Then
   `make baseline-digests`; then, with a fresh `HOME` each time and the
   binary rebuilt from the tree, the managed refresh `"$tmp/roundfix"
   baseline update --repo . --no-skills --no-history --yes --format text`
   twice (the second reports `File changes: 0`). No skill is restored: the
   `python` setup adds no skill this repository's `go-cli-tui` profile
   requires, and the other snapshots keep their bytes, so
   `TestThisRepositoryHoldsItsRequiredSkillsAtTheSetupSnapshot` holds without
   a restore (re-amendment probe, `_prd.md` → Validation).
2. API Contract: the `python` snapshot's source is `{"type": "github",
   "repository": "marcioaltoe/skills", "ref":
   "65d5f365cace27c4b468db1723213723b10c870b", "path": "setups/python.txt"}`
   and its 52 skill names, in order, are `find-skills`, `knowledge-workspace`,
   `setup-context-driven`, `brainstorming`, `business-analyst`,
   `context7-cli`, `exa-web-search`, `grill-with-docs`, `grilling`,
   `write-idea`, `write-prd`, `agentic-cli-design`, `domain-modeling`,
   `write-techspec`, `write-tasks`, `coding-guidelines`, `implement-spec`,
   `implement-task`, `no-workarounds`, `testing-boss`, `typesafe-ai`,
   `evidence-gate`, `qa-gate`, `systematic-debugging`, `roundfix`,
   `archive-spec`, `conventional-commits`, `crafting-effective-readmes`,
   `github-pr-workflow`, `ai-memory-durable-pages`, `ai-memory-handoff`,
   `ai-memory-learning-maintenance`, `ai-memory-messaging`,
   `ai-memory-retrieval`, `ai-memory-routing-install`, `tech-writer`, `uv`,
   `ruff`, `ty`, `pytest`, `typer`, `fastapi`, `python-services`,
   `python-file-pipelines`, `backend-docs`, `external-api-adapters`,
   `integration-contract-testing`, `security-best-practices`, `pdf`,
   `pdftk-server`, `docker`, `developing-with-streamlit`. Every upstream snapshot and each of its
   `github` skills names the same `ref`.
3. API Contract: the module's dispatch, each skill required and dispatched
   by exactly one trigger, in this order:
   - `uv`: "Changing the Python project's dependencies, lockfile, Python
     version or environment, or running a project tool through uv."
   - `ruff`: "Formatting or linting Python code, fixing a ruff finding, or
     changing the ruff configuration."
   - `ty`: "Type checking Python code with ty, fixing a ty diagnostic, or
     changing the ty configuration."
   - `pytest`: "Writing, selecting, or reviewing Python tests, fixtures,
     parametrization, or pytest configuration."
   - `python-services`: "Structuring a Python command-line tool, HTTP
     backend, or worker around shared application services: its entry
     points, configuration, retries, cancellation, or durable work."
   - `typer`: "Only when the repository's command-line interface already
     uses Typer: adding or changing its commands, options, or their tests."
   - `fastapi`: "Only when the repository's HTTP service already uses
     FastAPI: adding or changing its routes, dependencies, Pydantic models,
     or their tests."
   - `python-file-pipelines`: "Only when the repository ingests or
     reconciles files or generates reports in a Python pipeline, such as
     CSV, XML, ZIP, SQLite, XLSX, or DOCX processing."
   - `developing-with-streamlit`: "Only when the repository ships a
     Streamlit internal tool or prototype the maintainer asked for: creating
     or changing that Streamlit app."
4. API Contract: the Python guide's rules and clauses, in this order, each
   rule at version 1:
   - `rule.python.uv-environment`, coverage `coverage.python`, `coverage.dependencies`:
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
     - `clause.python.check-the-lock-offline`, mandatory: "The repository
       Verification checks the lockfile with `uv lock --check --offline`: it
       compares `uv.lock` with `pyproject.toml` from local files alone, fails
       when they no longer match, and needs no network or cache, so it runs
       inside a Run as in CI."
   - `rule.python.verification-gates`, coverage `coverage.python`, `coverage.verification-integrity`:
     - `clause.python.iterate-with-focused-commands`, mandatory: "Use focused
       commands such as `uv run pytest <node id>`, `uv run ruff check <path>`
       and `uv run ty check <path>` while iterating, then run the selected
       repository Verification."
     - `clause.python.ruff-findings-fail-the-gate`, prohibited: "Do not hide a
       ruff finding: never pass `--exit-zero`, add a blanket `# noqa`, or
       widen the ruff configuration to make a finding disappear. A `# noqa`
       names the rule code it silences and the reason."
     - `clause.python.type-check-with-ty`, mandatory: "The repository
       Verification type-checks the package and its tests with `uv run ty
       check --error-on-warning`, ty being a development dependency. Do not
       hide a ty diagnostic: never pass `--exit-zero` or
       `--exit-zero-on-warning`, add a blanket `# ty: ignore` or `# type:
       ignore`, or widen the ty configuration to make a diagnostic disappear.
       A suppression names the rule it silences and the reason."
     - `clause.python.select-a-new-test-by-its-node-id`, mandatory: "A Task
       Verification selects a new test by its pytest node id, such as
       `tests/test_report.py::test_totals`, and never passes that test's file
       or directory in the same command. pytest collects a node id once when
       another argument already contains it and exits 0 even before the test
       exists, so the command could never fail before the work; a missing node
       id alone exits 4 and an empty selection exits 5. `roundfix spec check`
       reports a command that names a test both ways."
   - `rule.python.hermetic-tests`, coverage `coverage.python`, `coverage.verification`:
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
       reads, or run `uv run ty check --python-platform <platform>` for each
       platform it names, because a test run on the host executes only the
       host's branch."
     - `clause.python.keep-entry-points-thin`, mandatory: "Keep each entry
       point thin, whether a console script, a web application factory or a
       worker's start-up: it reads its input and calls behavior that lives in
       the package. Test behavior through the package and the entry point's
       observable contract: for a console script, standard output, standard
       error and the exit code it raises with `SystemExit`; for a web
       application or a worker, the responses or messages it produces."
5. API Contract: rendering. A greenfield plan for the `python` profile with
   the usual decisions in a repository whose Makefile declares `verify` and
   `verify-incremental` is `ready`, renders `docs/agents/python.md` with the
   twelve clauses, `docs/agents/skill-dispatch.md` naming the nine skills and the
   root pointer "Python implementation rules are mandatory", renders no
   `docs/agents/go.md`, `docs/agents/rust.md`, `docs/agents/tui.md` or
   `docs/agents/cli.md`, records `python` in its Setup Manifest, applies, and
   re-plans with no file change.

### Surface Transcripts

None. No command, flag, exit code or output line changes; the profile menu
gains one entry.

### Invariants

```text
1. A Setup Snapshot holds no list of its own choosing; only the asset sync writes it, from one upstream commit for every upstream snapshot; this Spec's sync changes only python.json.
2. The embedded snapshots agree on every external skill's contract.
3. Existing profiles and clauses keep their bytes; only setups, derived files, the module record, this repository's skill copies and lock, and its Setup Manifest change, through the sanctioned commands.
4. The python profile selects no application-surface module and renders no Go, Rust, TUI or CLI guide.
5. No Roundfix code or test runs a Python tool.
```

## Coverage Map

Deviations from the PRD: None.

- Goals 1 and 2, Core Feature 1, Success Metrics 1, 2 and 5, User Story 4, Invariants 1-3 → API Contracts 1 and 2
- Goals 3-5, Core Features 2 and 3, Success Metrics 3 and 4, User Stories 1-3, Invariants 4-5 → API Contracts 3-5
- Core Feature 4, Success Metric 6 → Build Order 3

## Absence Evidence

1. The Python profile, module and setup — command: `ls internal/baseline/assets/profiles internal/baseline/assets/modules internal/baseline/assets/setups | grep -i python`; result: no match at `64c56d7f`.
2. The python setup file — command: `test -e internal/baseline/assets/setups/python.json`; result: exit 1 at `b5e45c63`.

## ADR Alternatives

- ADR-0322: none of its rejected alternatives is taken: no common or composed derivation, no single-snapshot import at a second commit, no per-tool decision.
- ADR-0191: its rule stands; the sync writes every snapshot, and the only hand edit is the seed's two values, which the sync overwrites.
- ADR-0206: no rejected alternative is taken; the snapshot takes the upstream list's name `python`, seeded as that ADR's Spec seeded renamed snapshots.
- ADR-0204: no rejected alternative is taken; the composed snapshot keeps the bytes Spec 0317's sync gave it.
- ADR-0241: no rejected alternative is taken; its agreement rule is why the python snapshot takes the commit the others already name.
- ADR-0219: no rejected alternative is taken; a repository-owned draft adapted from the `python` profile binds to it like to any built-in profile.
- ADR-0246: no rejected alternative is taken; `council` and `the-fool` stay out of every snapshot.
- ADR-0072: no rejected alternative is taken; the frozen parity row changes only through the sanctioned transform, and the other frozen fixtures keep their bytes.

## Integration Points

The embedded Baseline catalog and its derived files; the asset sync over a
repository fetched from `~/dev/skills`; the repository's Managed Refresh. No
test, Verification command or QA row reaches the network, reads the real
`~/.roundfix`, reads the Secondbrain, reads another repository or runs a
Python tool. Only task_01's import reads `~/dev/skills`, fetching one commit
into an empty local repository.

## Testing Approach

Every new top-level test calls `t.Parallel()` first and passes no real-time
duration; the new files use no platform API and compile for Windows.

- `internal/baseline/python_setup_test.go` (new, Build Order 1):
  - `TestThePythonSetupIsThePublishedUpstreamPreset`: the embedded `python`
    setup has API Contract 2's source and its 52 names in order.
  - `TestEveryUpstreamSetupFollowsThePythonPresetRevision`: `go`, `rust`,
    `typescript` and `python` and each of their `github` skills name
    `65d5f365cace27c4b468db1723213723b10c870b`, and `python` holds
    `ai-memory-handoff` and no `handoff`.
- `internal/baseline/python_profile_test.go` (new, Build Order 2):
  - `TestThePythonModuleDispatchesItsSkills`: the module's `requiredSkills`
    equal API Contract 3's nine in order; each dispatch has one trigger, and
    a trigger begins "Only when the repository" exactly for `typer`,
    `fastapi`, `python-file-pipelines` and `developing-with-streamlit`.
  - `TestThePythonProfileRendersThePythonGuide`: through `newPlanRepository`
    and `planTestDecisions` with greenfield preservation, API Contract 5's
    ready plan, eight guide clauses with their force, the lock-check and ty clauses among them (on whitespace-normalized
    text), the nine skills in the skill-dispatch guide, the root pointer, and
    no Go, Rust, TUI or CLI guide among the postimages.
- Declared breaks (Build Order 1): `assets_sync_test.go` counts five setups
  and findings in `TestAssetsSyncCheckIsReadOnlyAndReportsDrift`,
  `TestBaselineAssetsSyncRefreshProducesCanonicalTreeAndIsIdempotent` (two
  counts) and `TestAssetsSyncCompatibilityMatchesMaintainedPythonContract`;
  `TestOutputsForCommand` adds `internal/baseline/assets/setups/python.json`.
- Declared breaks (Build Order 2): `TestEmbeddedCatalog` lists `python`;
  `TestOutputsForCommand` adds `internal/baseline/assets/profiles/python.json`;
  `characterizedBaselineForce` lists the ten Python clauses with API Contract
  4's force; `TestBaselineHumanProfileAdaptation` answers `5` in its
  `answers` script; `TestGuidanceCompositionJourney` lists `python` and its
  two fixture switches gain a `python` case (a `pyproject.toml` and an empty
  `src/example/__init__.py`; no formatter run).
- Prototype record (2026-10-24, scratch worktrees of `64c56d7f`): after Build
  Order 1 alone, `./internal/baseline` and `./skills` passed; after Build
  Orders 1 and 2, every package of `_prd.md` → Validation passed. These
  sabotages each failed the named test:
  - the first skill removed from `setups/python.json` →
    `TestThePythonSetupIsThePublishedUpstreamPreset`;
  - "Only when the repository's command-line interface already uses Typer: "
    removed from the `typer` trigger, or "Only when the repository has chosen
    ty: " prefixed to the `ty` trigger →
    `TestThePythonModuleDispatchesItsSkills`;
  - `ty` removed from `requiredSkills` →
    `TestEveryModuleRequiresEverySkillItDispatches` ("python dispatches
    unrequired skill ty");
  - `--offline` removed from the lock-check clause, or `--error-on-warning`
    from the ty clause → `TestThePythonProfileRendersThePythonGuide`;
  - "selects a new test by its pytest node id" reworded →
    `TestThePythonProfileRendersThePythonGuide`;
  - `cli-surface` added after `python` in the profile's modules →
    `TestThePythonProfileRendersThePythonGuide` ("python plan renders
    docs/agents/cli.md").

## Build Order

1. The import: the seed, the sync and its check, the parity transform,
   `make baseline-digests`, the managed refresh twice, the Build Order 1
   declared breaks and `python_setup_test.go`. It depends on Spec 0317
   having merged, which moves the other upstream snapshots to the same
   commit; the import's first commands stop when it has not.
2. The Python profile: module, templates and index entries, coverage entry,
   core entries and profile; the module record command, `make
   baseline-digests` and the managed refresh twice; the Build Order 2 declared
   breaks and `python_profile_test.go` (depends on 1, whose `python` setup the
   profile names and whose derived files it regenerates).
3. Documentation: the user guide's profile list and paragraph, the Coverage
   Review and the Behavior Surface re-record last (depends on 2).
4. QA gate (depends on: 1, 2, 3).

No Spec is listed in `requires`; the operator queues this Spec after Spec
0317 merges (cycle AP batch 6 order).

## Risks & Considerations

- Fleet effect: none beyond the new profile; Spec 0317 already moved the
  other snapshots and its release note covers the changed skills.
- The local checkout may lack the commit; task_01 stops before any write.
- Building `bin/roundfix` or running `go run` under a fresh `HOME` moves the
  Go module cache into that home and downloads modules; API Contract 1
  builds the binary once under the Run's environment and gives only the
  binary a fresh `HOME`.
- Rebasing over Specs 0317 and 0307 conflicts in `core.json`'s version line, the module
  record and the derived goldens; the Task resolves them by rerunning the
  sanctioned commands, never by hand.

## Glossary

None.

## Decisions

- The `python` setup is the imported upstream preset, at the commit Spec 0317
  moved every other upstream snapshot to; see ADR-0322.
- Optional skills are required and dispatched with a trigger beginning "Only
  when the repository"; see ADR-0322.
- Every Python gate runs ty and the lock check (maintainer, 2026-10-24 about
  12:30 BRT, "ty e lock check sempre"). The guide's
  `clause.python.type-check-with-ty` names `uv run ty check
  --error-on-warning` over the package and its tests and forbids hiding a
  diagnostic, as the ruff clause does; the platform clause names ty's
  `--python-platform`; `ty` is dispatched always.
- The lock check is `uv lock --check --offline`
  (`clause.python.check-the-lock-offline`). uv decides whether a lock is
  current by comparing it with the project metadata, and the probe in
  `_prd.md` → Validation shows a current lock passes with no network and an
  empty cache, while a stale one fails with or without a cache. So the gate
  runs it inside a Run without the network, and CI keeps `uv sync --locked`,
  which names the cause with the network. Spec 0310's Starter Makefile carries
  both targets.
- The guide's node-id clause names `roundfix spec check`'s existing finding by
  behavior, not by code.
- The starter files are Spec 0310's.
- The profile is named `python` and fixes no application surface (maintainer
  decisions of 2026-10-24).
