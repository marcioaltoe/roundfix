---
spec: 0309-a-python-baseline-profile-for-uv-pytest-and-ruff
status: active
created: 2026-10-23
surfaces: [backend, docs]
---

# A Python Baseline Profile for uv, pytest and ruff

A maintainer who adopts the Context-Driven Baseline in a Python repository
today finds no built-in profile for it: Quaestor had to write a
repository-owned profile that removes the `rust` module from a built-in one,
so its agents read no Baseline rule about uv, ruff or pytest, and its skills
were picked by hand. After this Spec Roundfix ships the `python` Baseline
Profile, which serves a Python command-line tool, backend service or worker
alike: it brings the language and the workflow modules and fixes no
application surface. Its guide states the uv, ruff and pytest rules with force, including
how a Task Verification selects a new test, how tests stay hermetic and how an
operating-system branch is exercised, and the profile takes the
language-neutral skills that the upstream Go and Rust setups share through a
new **Common Setup Snapshot**, so an adopter restores no Go or Rust skill.

## Sources

This Spec owns one source, adopted at birth as a promoted Backlog Entry under
`references/` (see `references/_index.md`):

1. [`references/2026-10-22-a-baseline-profile-for-python-repositories.md`](references/2026-10-22-a-baseline-profile-for-python-repositories.md),
   type `feat`, from Quaestor's ai-memory request of 2026-10-21 at 19:02 BRT
   and the maintainer's decisions of 2026-10-22. This Spec delivers its
   expected items 1, 2, 3 and 5; Spec 0310 delivers item 4, the generated
   Makefile, Project Config and CI workflow, and links this copy.

No other open Backlog Entry or unresolved Finding shares this context: at the
authoring base `dde9d5f5` (2026-10-23) `docs/backlog/` held no entry and
`docs/findings/` held none.

## Prerequisites

None in `requires`: no Verification command of this Spec runs code of another
active Spec. This Spec is authored against `dde9d5f5` (v0.113.0 with the plans
of Specs 0303 to 0307 merged), and the operator delivers it after Specs 0303 to
0307. Spec 0307 also edits the core module and changes the Baseline digests;
this Spec's core edit touches other entries (the instruction hierarchy and the
semantic owners), names no module version or digest value, and regenerates
every derived file through the sanctioned commands on the tree it starts on.
Spec 0302, already delivered, gives `roundfix spec check` the finding for a
pytest node id passed together with an argument that contains it; the guide's
node-id clause names that behavior.

## Project Constraints

- Identifier strategy: not applicable — no project-owned identifier or record
  schema is minted; the new catalog identifiers (the `python` setup, module
  and profile, `rule.python.*`, `clause.python.*`, `coverage.python`) follow
  the existing dotted naming. Profiles, modules and setups are separate
  catalog namespaces (`internal/baseline/catalog.go` `Catalog` keeps the maps
  `profiles`, `modules` and `setups`, and `catalog_load.go` `readCollection`
  checks duplicates within one collection), as the `go` and `rust` setups and
  modules already share their identifiers. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, credential or
  recipient is added; tests read the embedded catalog, temporary repositories
  and a synthetic asset-sync source. Source:
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0322, "A Python profile takes the
  skills the Go and Rust setups share", which this Spec records; ADR-0204, "A
  composed profile takes a setup composed from upstream setups by name", whose
  derived-snapshot rules the common kind follows ("holds no list of its own
  choosing"); and ADR-0191, "A setup snapshot follows its upstream by name",
  which keeps the asset sync the only writer of a snapshot ("the asset sync is
  its only writer"); ADR-0206, "The Baseline takes upstream setup names and
  drops a removed skill", whose names the `python` snapshot's components use;
  ADR-0219, "A profile draft binds to the closest built-in profile it
  adapts", which binds a repository-owned draft of the `python` profile like any
  built-in profile; ADR-0241, "A repository profile takes its skill contracts
  from the embedded Setup Snapshots", which includes the `python` snapshot;
  and ADR-0246, "A retired skill leaves the Baseline and the
  adopter removes its copy", whose removal reaches the common part through its
  components. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AP approval of 2026-10-22 through AskUserQuestion, scope "Preset de baseline
  Python", authorization "same as cycle AO", which carries the cycle AJ grant
  "Mexer em Governed Paths declarados" ("'Concedo' para os caminhos que cada
  Spec declarar (módulos Baseline, registros, skills, testes governados; fora:
  Makefile, CI, go.mod, lint, que pedem autorização específica)."), and the
  preset's reach, the maintainer's own answer of the same day.
  Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
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
  `internal/cli/baseline_release_gate_test.go`. No change to this repository's Makefile,
  workflows, Go module, lint, formatter or Project Config.

## Goals

- A Python repository can select a built-in `python` Baseline Profile and
  receive a ready plan whose managed guidance includes `docs/agents/python.md`.
- The same profile serves a command-line tool, a backend service or a worker:
  it selects no application-surface module, and a repository adds the rules
  of its own surface through the repository extension.
- The Python guide states, as Baseline clauses with force, the uv, ruff and
  pytest rules the request lists, the node-id Verification rule and the
  operating-system branch rule.
- The profile's setup holds only skills the upstream Go and Rust lists share,
  derived and checked like a composed setup, so no Go, TUI or Rust skill is
  restored into a Python repository.
- Existing profiles, setups and adopters see no change except the catalog
  digest.

## User Stories

1. As a maintainer of a Python repository, I want to pick a built-in Python
   profile, so that I stop maintaining a hand-made profile without Python rules.
2. As an Agent in a Python repository, I want the uv, ruff and pytest rules in
   a managed guide, so that I do not guess from pip habits or hide a lint
   finding.
3. As a Spec author, I want the guide to say how a new pytest test is selected
   in a Verification, so that the command fails before the work.
4. As the Baseline maintainer, I want the Python setup derived from the
   upstream lists and checked like a composed setup, so that it never drifts
   silently.

## Core Features

1. The common setup kind. A Setup Snapshot whose source is `{"type":
   "common", "setups": [...]}` holds, in the first component's order, each
   skill and activation bundle every component carries with an equal entry;
   the catalog refuses one that differs, has fewer than two components, or
   names an unknown or derived component; the asset sync writes it after its
   components.
2. The `python` setup: the common part of `go` and `rust`.
3. The `python` module: a guide `docs/agents/python.md` with three rules and
   ten clauses (uv environment; ruff and pytest gates, including the node-id
   rule; hermetic tests, outside data, platform branches and thin entry
   points),
   a root block pointing at it, a coverage entry, and the core module's
   hierarchy and semantic-owner entries for it.
4. The `python` profile: the `rust-cli` module set with `python` in place
   of `rust` and without the `cli-surface` module and its two `rule.cli.*`
   required rules, on the `python` setup, with no formatter; its title is
   "Python repository (CLI, backend or worker)".
5. Documentation and glossary: the user guide's profile list names
   `python` and says what it governs and why it dispatches no Python skill
   yet; the glossary defines **Common Setup Snapshot**.

## User Experience

`roundfix baseline` offers `python` in its profile menu between
`go-cli-typescript-monorepo` and `rust-cli`; `roundfix baseline plan --profile
python` with the usual decisions produces a ready plan in a repository
that declares `make verify` and `make verify-incremental`. The adopter's
`docs/agents/python.md` and root instructions name the Python rules.

## Non-Goals / Out of Scope

- Generating the Makefile, Project Config or CI workflow: Spec 0310.
- Python-specific skills: the upstream skills repository has none; the
  curation request
  `~/dev/secondbrain/inbox/skills/2026-10-22-setup-python-e-skills-python-uv-pytest-ruff.md`
  is already filed by the operator, and a later Spec adopts what it produces.
- A Python formatter fixture: the profile declares no formatter, as `rust-cli`
  declares none.
- Running any Python tool from Roundfix or its tests.
- A type-checker rule: the request names none; the platform clause mentions a
  type checker's platform option only as one way to exercise a branch.

## Success Metrics

1. Success Metric: the embedded catalog lists the profiles `go-cli-tui`,
   `go-cli-typescript-monorepo`, `python`, `rust-cli` and
   `standard-typescript-monorepo`, and `roundfix baseline plan --profile
   python` in a temporary repository declaring the two Make targets
   returns a ready plan that renders `docs/agents/python.md`, renders no Go,
   Rust, TUI or CLI guide (`docs/agents/cli.md`), and whose re-plan after
   apply has no file change; the profile selects no `application-surface`
   module.
2. Success Metric: the embedded `python` setup equals the common part of `go`
   and `rust` computed by the derivation, holds no `golang-*`, `rust*`,
   `clap-rust`, `cut-release`, `bubbletea` or `tui-design` skill, and holds
   `roundfix`, `testing-boss` and `agentic-cli-design`.
3. Success Metric: the common kind keeps only equal shared entries in the
   first component's order, the catalog reports a drifted common setup with
   the detail "skills or activation bundles differ from the components'
   common part", refuses one component or a derived component, and the asset
   sync rewrites a stale common setup after its components; each check fails
   when the matching code is sabotaged.
4. Success Metric: the Python guide carries the ten clauses with their
   enforcement, and the clause force characterization lists them.
5. Success Metric: the user guide names `python` and says it fixes no
   application surface, the glossary defines
   **Common Setup Snapshot**, and neither the user guide nor the Roundfix
   skill contains `python3` or `Python fallback`.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write.

- The upstream Setup Snapshots `internal/baseline/assets/setups/go.json` and
  `internal/baseline/assets/setups/rust.json`, written by the asset sync from
  `marcioaltoe/skills` at `b3c45a4` long before this Spec and unchanged by
  it: the skills they share with equal entries are the `python` setup's
  content, and the QA counts them on the Delivery Base. Quaestor's own
  decision record (`docs/adr/0004-roundfix-e-baseline.md` in its repository,
  2026-10-19, read through the Secondbrain mirror at authoring) reports the
  same number for its hand-made profile: "`roundfix doctor` confirma
  `skills: ok (30 required: 13 Roundfix-owned, 17 external)`", and states the
  gap this Spec closes: "Para ter um guia de Python gerido pelo baseline, o
  Roundfix precisa ganhar um módulo `python` no catálogo embutido; perfis do
  repositório não podem declarar módulos novos."
- uv's documentation, "Locking and syncing"
  (https://docs.astral.sh/uv/concepts/projects/sync/), quoted: "To use the
  lockfile without checking if it is up-to-date, use the `--frozen` option"
  and "If the lockfile is not up-to-date, uv will raise an error instead of
  updating the lockfile." — the distinction the sync clause states.
- pytest's documentation, "Exit codes"
  (https://docs.pytest.org/en/stable/reference/exit-codes.html): exit code 4
  is a command-line usage error and exit code 5 means no tests were collected
  — the codes the node-id clause names.

## Unreachable Acceptance

- criterion: a real Python repository adopts the `python` profile with `roundfix baseline update`, restores the thirty skills and runs its own `make verify` under the new guide
  reason: it needs another repository, its installed skills and Python tooling, which a Run sandbox and this Spec's authorization exclude
  satisfied-by: Quaestor's first managed refresh after the release, which the operator announces to `default/marcioaltoe-quaestor`, recorded as this Spec's unproven step

## Validation

The decisions were checked on 2026-10-23 against these sources.

- Characterization at `dde9d5f5`. The catalog lists four profiles and four
  Setup Snapshots, one composed (`go-cli-typescript-bun`, ADR-0204); `go` and
  `rust` share exactly thirty skills with equal entries and both bundles
  (`bundle.delivery`, `bundle.qa`); `rust-cli` declares formatter kind `none`
  and no formatter fixture. The Baseline writes only guidance and the Setup
  Manifest. A prototype in a scratch worktree added the common kind, the
  `python` setup, module, templates, coverage entry, core entries and the
  profile, ran the module record command, `make baseline-digests` and the
  Managed Refresh twice, and measured the derived files and the existing tests
  the TechSpec lists; the `internal/baseline`, `internal/cli`, `skills`,
  `internal/spec` and `internal/speccheck` suites then passed, the release
  journey passed for the profile, and the sabotages the TechSpec lists each
  failed the named test. That prototype carried the earlier shape, the
  identifier `python-cli` with the `cli-surface` module.
- Amendment of 2026-10-24 (the maintainer's decisions in Decisions), checked
  by reading at `87094aa9`: no Go code reads a module's `kind`
  (`grep -rn "application-surface" internal --include='*.go'` has no match),
  so a profile without an application-surface module is valid;
  `validateProfiles` in `catalog_validate.go` only requires each module's
  required and dispatched skills to be in the setup, never the reverse, so
  `agentic-cli-design` stays in the derived `python` setup without a module
  that dispatches it; the coverage catalog is global, so the hermetic-tests
  rule may name `coverage.verification`; and the existing tests that mention
  `docs/agents/cli.md` (`internal/cli/baseline_release_gate_test.go`) use the
  `go-cli-tui` profile only. The menu position and every declared break of
  the TechSpec are unchanged, because `python` sorts where `python-cli`
  sorted.
- Secondbrain, read only. Read `wiki/index.md` (entity `quaestor`) and ran
  `qmd query "Python Baseline profile uv pytest ruff Quaestor" --all --files
  --min-score 0.3`, which returned Quaestor's
  `projects/quaestor/mirror/docs/adr/0004-roundfix-e-baseline.md` (quoted
  above: it confirms the gap and the thirty hand-installed skills) and
  Roundfix's mirrored decisions that Baseline execution belongs to the CLI and
  that custom profiles are repository-owned and cannot add modules, which is
  why the profile must be built in. The curation request
  `inbox/skills/2026-10-22-setup-python-e-skills-python-uv-pytest-ruff.md` was
  read, not written.
- Exa: uv's CLI reference (https://docs.astral.sh/uv/reference/cli/: "`--frozen`
  ... Instead of checking if the lockfile is up-to-date, uses the versions in
  the lockfile as the source of truth" and "`--locked` ... Requires that the
  lockfile is up-to-date"), uv's locking page above and pytest's exit codes
  page above. Challenging: Quaestor bootstraps with `--frozen`, which never
  notices a lockfile that no longer matches `pyproject.toml`, so the clause
  keeps `--frozen` for creating an environment and names `--locked` for a step
  that must refuse a stale lock, instead of prescribing one flag everywhere.
- Quaestor's tooling at `marcioaltoe/quaestor@e6b390d`, read through `gh api`
  without cloning: `Makefile`, `.roundfixrc.yml`,
  `.github/workflows/verify.yml` and `pyproject.toml`, which confirm the
  request's commands; no client data was read.

## Glossary

- adds: **Common Setup Snapshot**

## Skills

- unchanged: guide docs/user-guide/context-driven-development.md — the only change names the new built-in `python` Baseline Profile in the profile list; the covering skills describe how a Baseline is adopted and refreshed, not which built-in profiles exist, so no skill text changes.

## Decisions

- 2026-10-22: Scope of cycle AP (AskUserQuestion, questions of about 15:25 BRT
  answered by about 18:30 BRT) — "Preset de baseline Python", delivered last
  ("Correções, follow-up, Python (Recomendado)").
- 2026-10-22: Reach of the preset — "Perfil completo e pode, se necessário,
  solicitar ou requerer novas skills recomendadas ao repo skills para que
  tenhamos a curadoria e ter o melhor resultado no uso do python, uv e linters
  e outros toolings especificos do setup." (the maintainer's own answer).
- 2026-10-24, about 09:10 BRT: the profile's name — "mude o nome do profile
  de python-cli para python (vai ser para cli, backend e workers)" (the
  maintainer's own words).
- 2026-10-24, about 09:15 BRT (AskUserQuestion): the profile's surface —
  "Nenhuma superfície fixa (Recomendado)" (option text: "O perfil traz só a
  linguagem (uv, pytest, ruff) e os módulos de fluxo (core, specs, autônomo,
  secondbrain, extensão). Serve igual para CLI, backend e workers; regras de
  uma superfície específica entram pela extensão do repositório. Muda pouco o
  plano e mantém 0309/0310 na v0.123/v0.124."). The profile therefore drops
  the `cli-surface` module and the `rule.cli.*` rules, the hermetic-tests rule
  covers `coverage.verification` instead of `coverage.cli`, and the
  entry-point clause speaks of any entry point (console script, web
  application or worker), not only a console script.

Autonomous route: no person present; decisions taken from the operator's
cycle AP batch 3 briefing and the characterization above, with no further
interview:

- The profile takes a Common Setup Snapshot of `go` and `rust` rather than a
  union, a hand-written list or the `go` list; see ADR-0322.
- The `python` module requires and dispatches no skill until the upstream
  repository publishes Python skills; the guide carries the rules.
- The profile is split from the starter files (Spec 0310) because the whole
  preset does not fit four implementation Tasks.

## Open Questions

None.
