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
were picked by hand. On 2026-10-24 the upstream skills repository published a
curated Python preset, `setups/python.txt` at commit `32e4143e`, and at
`65d5f36` replaced its `handoff` skill with six ai-memory skills, leaving 52.
After this Spec Roundfix imports that preset at `65d5f36`, the commit Spec
0317 moves every other upstream snapshot to, as the `python` Setup Snapshot
through its asset sync, and ships the `python` Baseline Profile,
which serves a Python command-line tool, backend service or worker alike: it
brings the language and the workflow modules and fixes no application
surface. Its guide states the uv, ruff, ty and pytest rules with force,
including an offline lock check and a ty type check in the gate, and its
module dispatches the uv, ruff, ty, pytest and Python-services skills always,
and the Typer, FastAPI, file-pipeline and Streamlit skills only when the
repository has chosen them.

## Sources

This Spec owns one source, adopted at birth as a promoted Backlog Entry under
`references/` (see `references/_index.md`):

1. [`references/2026-10-22-a-baseline-profile-for-python-repositories.md`](references/2026-10-22-a-baseline-profile-for-python-repositories.md),
   type `feat`, from Quaestor's ai-memory request of 2026-10-21 at 19:02 BRT
   and the maintainer's decisions of 2026-10-22 and 2026-10-24. This Spec
   delivers its expected items 1, 2, 3 and 5; Spec 0310 delivers item 4, the
   generated Makefile, Project Config and CI workflow, and links this copy.

No other open Backlog Entry or unresolved Finding shares this context: at the
amendment base `64c56d7f` (2026-10-24) `docs/backlog/` and `docs/findings/`
held none about Python or the setup snapshots.

## Prerequisites

None in `requires`. The operator delivers it after Specs 0311 to 0317 and
0307 (cycle AP order 0311–0314 → 0315 → 0316 → 0317 → 0307 → 0309 → 0310,
by the maintainer's batch 6 answer quoted in `_authorization.md`). task_01's
Verification depends on Spec 0317's result: it expects the `go`, `rust`,
`typescript` and composed snapshots already at
`65d5f365cace27c4b468db1723213723b10c870b`, and stops and reports when they
are not, so this Spec must not be queued before Spec 0317 merges. Specs 0307
and 0317 also edit the core module and change the Baseline digests; this
Spec's core edit touches other entries, names no module version or digest
value, and regenerates every derived file through the sanctioned commands on
the tree it starts on.

The import needs the upstream commit
`65d5f365cace27c4b468db1723213723b10c870b` in the local checkout
`~/dev/skills`, from which task_01 fetches it without reaching a network
remote; task_01 stops, writing nothing, when it is absent.

## Project Constraints

- Identifier strategy: not applicable — no project-owned identifier or record
  schema is minted; the new catalog identifiers (the `python` setup, module
  and profile, `rule.python.*`, `clause.python.*`, `trigger.python.*`,
  `coverage.python`) follow the existing dotted naming. Profiles, modules and
  setups are separate catalog namespaces (`internal/baseline/catalog.go`
  `Catalog` keeps the maps `profiles`, `modules` and `setups`, and
  `catalog_load.go` `readCollection` checks duplicates within one
  collection), as the `go` and `rust` setups and modules already share their
  identifiers. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, credential or
  recipient is added; the asset sync and the skill restore read a local clone
  of a local checkout, and tests read the embedded catalog, temporary
  repositories and a synthetic asset-sync source; the sync reads a local
  repository fetched from a local checkout. Source:
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0322, "A Python profile imports
  the published upstream Python preset", which this Spec records; ADR-0191,
  "A setup snapshot follows its upstream by name", which keeps the asset sync
  the only writer of a snapshot ("the asset sync is its only writer");
  ADR-0206, "The Baseline takes upstream setup names and drops a removed
  skill", whose seed-then-sync path imports the new snapshot; ADR-0204, "A
  composed profile takes a setup composed from upstream setups by name",
  whose composed snapshot, composed from `go` and `typescript`, this Spec's
  sync leaves unchanged; ADR-0241, "A repository
  profile takes its skill contracts from the embedded Setup Snapshots", whose
  agreement rule requires the `python` snapshot at the commit the other
  upstream snapshots already name; ADR-0219, "A profile draft binds to the closest
  built-in profile it adapts", which binds a repository-owned draft of the
  `python` profile like any built-in profile; ADR-0246, "A retired skill
  leaves the Baseline and the adopter removes its copy", whose retired skills
  (`council`, `the-fool`) the sync keeps out; and ADR-0072, "Baseline Go
  cutover preserves Python contracts", whose maintained asset-sync parity row
  the sanctioned transform follows. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AP approval of 2026-10-22 through AskUserQuestion, scope "Preset de baseline
  Python", authorization "same as cycle AO", which carries the cycle AJ grant
  "Mexer em Governed Paths declarados" ("'Concedo' para os caminhos que cada
  Spec declarar (módulos Baseline, registros, skills, testes governados; fora:
  Makefile, CI, go.mod, lint, que pedem autorização específica)."), the
  preset's reach, the maintainer's own answer of the same day, and the
  decision of 2026-10-24, "Emendar a 0309 agora (Recomendado)".
  Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
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
  No change to this repository's Makefile, workflows, root Go module, lint,
  formatter or Project Config.

## Goals

- The `python` Setup Snapshot is the published upstream preset, imported by
  the asset sync from `65d5f36`, with no list of Roundfix's choosing.
- The `python` snapshot names the upstream commit every other upstream
  snapshot already names, so a repository-owned profile never meets two
  contracts for one skill.
- A Python repository can select a built-in `python` Baseline Profile and
  receive a ready plan whose managed guidance includes
  `docs/agents/python.md`; the profile fixes no application surface.
- The Python guide states, as Baseline clauses with force, the uv, ruff and
  pytest rules the request lists, the offline lock check and the ty type
  check every Python gate runs, the node-id Verification rule and the
  operating-system branch rule.
- The `python` module dispatches uv, ruff, ty, pytest and Python services
  always, and Typer, FastAPI, file pipelines and Streamlit only when the
  repository has chosen them.

## User Stories

1. As a maintainer of a Python repository, I want to pick a built-in Python
   profile, so that I stop maintaining a hand-made profile without Python
   rules or skills.
2. As an Agent in a Python repository, I want the uv, ruff and pytest rules in
   a managed guide and the matching skills dispatched, so that I do not guess
   from pip habits or hide a lint finding.
3. As an Agent in a repository that never chose FastAPI or Streamlit, I want
   those skills not dispatched to me, so that I do not introduce a framework
   nobody chose.
4. As the Baseline maintainer, I want the Python skill list to come from the
   curated upstream preset through the asset sync, so that it never drifts
   silently.

## Core Features

1. The import. A seed `python.json` (a copy of `rust.json` with only `id` and
   `source.path` changed) and one asset-sync refresh from a repository
   fetched at `65d5f36` write the `python` snapshot and leave `go`, `rust`,
   `typescript` and the composed `go-cli-typescript-bun`, already at that
   commit, unchanged; the maintained asset-sync parity row follows through its
   sanctioned transform. No skill of this repository is restored: the
   `python` setup adds no skill its `go-cli-tui` profile requires.
2. The `python` module: a guide `docs/agents/python.md` with three rules and
   twelve clauses (uv environment and the offline lock check; ruff, ty and
   pytest gates, including the node-id rule; hermetic tests, outside data,
   platform branches and thin entry points), a root block pointing at it, a
   coverage entry, nine required and dispatched skills of which four carry a
   trigger that begins "Only when the repository", and the core module's
   hierarchy and semantic-owner entries for it.
3. The `python` profile: the `rust-cli` module set with `python` in place of
   `rust` and without the `cli-surface` module and its two `rule.cli.*`
   required rules, on the `python` setup, with no formatter; its title is
   "Python repository (CLI, backend or worker)".
4. Documentation: the user guide's profile list names `python` and says what
   it governs, where its skills come from and which it dispatches only when
   chosen.

## User Experience

`roundfix baseline` offers `python` in its profile menu between
`go-cli-typescript-monorepo` and `rust-cli`; `roundfix baseline plan --profile
python` with the usual decisions produces a ready plan in a repository that
declares `make verify` and `make verify-incremental`. The adopter's
`docs/agents/python.md`, `docs/agents/skill-dispatch.md` and root
instructions name the Python rules and skills. Adopters of the other built-in
profiles see no setup change from this Spec; Spec 0317 already moved their
snapshots.

## Non-Goals / Out of Scope

- Generating the Makefile, Project Config or CI workflow: Spec 0310.
- Adding ty to an adopter's `pyproject.toml`: the guide requires ty as a
  development dependency of the gate, and the repository adds it with
  `uv add --group dev ty`; Spec 0310's Starter Makefile runs it.
- Dispatching the preset's backend, document and container skills
  (`backend-docs`, `external-api-adapters`, `integration-contract-testing`,
  `security-best-practices`, `docker`, `pdf`, `pdftk-server`) or
  `agentic-cli-design`: they are installed with the setup; the catalog lets
  one module dispatch a skill and the `backend` module already dispatches four
  of them.
- Installing `find-skills` into this repository: no module requires it.
- Refreshing `go`, `rust` and `typescript`, the ai-memory skills and the
  retirement of `handoff`: Spec 0317 delivers them before this Spec.
- Running any Python tool from Roundfix or its tests.

## Success Metrics

1. Success Metric: the embedded `python` setup's source is
   `marcioaltoe/skills` `setups/python.txt` at
   `65d5f365cace27c4b468db1723213723b10c870b`, and its 52 skill names equal
   the published preset in order, the six `ai-memory-*` skills included and
   `handoff` absent; every upstream setup and each of its external skills
   names that commit.
2. Success Metric: `roundfix baseline assets sync --source-dir <clone>/setups
   --check` over a clone of upstream at that commit exits `0` after the Task,
   and the embedded snapshots agree on every external skill
   (`TestEmbeddedSetupSnapshotsAgreeOnEveryExternalSkill`); the sync leaves
   the four other snapshots byte-identical to the Delivery Base.
3. Success Metric: the embedded catalog lists the profiles `go-cli-tui`,
   `go-cli-typescript-monorepo`, `python`, `rust-cli` and
   `standard-typescript-monorepo`, and `roundfix baseline plan --profile
   python` in a temporary repository declaring the two Make targets returns a
   ready plan that renders `docs/agents/python.md` with the twelve clauses and a
   skill-dispatch guide naming the nine Python skills, renders no Go, Rust,
   TUI or CLI guide, and whose re-plan after apply has no file change.
4. Success Metric: the `python` module requires exactly `uv`, `ruff`, `ty`,
   `pytest`, `python-services`, `typer`, `fastapi`,
   `python-file-pipelines` and `developing-with-streamlit`, and the trigger of
   exactly the last four begins "Only when the repository"; the guide names
   `uv lock --check --offline` and `uv run ty check --error-on-warning` as
   part of the gate; each check fails when the matching asset is sabotaged.
5. Success Metric: this repository still holds every skill it requires at
   the `go` snapshot (`TestThisRepositoryHoldsItsRequiredSkillsAtTheSetupSnapshot`),
   `.agents/skills/` and `skills-lock.json` are byte-identical to the Delivery
   Base, and a second managed refresh reports no file change.
6. Success Metric: the user guide names `python`, says it fixes no
   application surface and which skills it dispatches only when chosen, and
   neither the user guide nor the Roundfix skill contains `python3` or
   `Python fallback`.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write.

- The upstream preset itself: `setups/python.txt` of `marcioaltoe/skills`,
  added at commit `32e4143e1b018005a65069afa408cec2c6e73a83` (pull request
  #130, "feat: add a unified Python setup and document tooling", 2026-10-24
  09:43 BRT) with 47 entries and changed at
  `65d5f365cace27c4b468db1723213723b10c870b` (pull request #131,
  "refactor: replace workspace and continuity skills"), read at authoring in
  a scratch repository: 52 entries, `handoff` replaced by the six
  `ai-memory-*` skills, the first 36 equal, in order, to the entries
  `setups/go.txt` and `setups/rust.txt` share at that commit; at `32e4143e` the registry `skills-registry.json` names
  `microsoft/debugpy` `.claude/skills/pytest` as the pytest skill's upstream
  and Astral's plugin repository as the uv, ruff and ty skills' upstream.
  The QA compares the audited `python` snapshot's names and source with this
  description, and records the upstream file itself as `blocked (environment:
  network prohibited by Spec authorization)` when no local clone holds the
  commit.
- Quaestor's report through the coordinator (ai-memory, 2026-10-24 about
  11:18 BRT): its adopted skill set was the 47 skills of `32e4143e` at its
  commit `d17c98a`, `handoff` included, before upstream's `65d5f36`, and its CLI rules live in its own
  `docs/agents/specific-repository.md`, the extension this profile leaves
  for surface rules.
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

- criterion: a real Python repository adopts the `python` profile with `roundfix baseline update`, restores the 52 skills and runs its own `make verify` under the new guide
  reason: it needs another repository, its installed skills and Python tooling, which a Run sandbox and this Spec's authorization exclude
  satisfied-by: Quaestor's first managed refresh after the release, which the operator announces to `default/marcioaltoe-quaestor`, recorded as this Spec's unproven step

## Validation

The decisions were checked on 2026-10-23 and, for the amendment, on
2026-10-24 against these sources.

- Characterization at `dde9d5f5` (2026-10-23): four profiles and four Setup
  Snapshots, one composed; `rust-cli` declares formatter kind `none`; the
  Baseline writes only guidance and the Setup Manifest.
- Amendment prototype of 2026-10-24, in scratch worktrees of `64c56d7f` and a
  scratch clone of upstream at `32e4143e`:
  - Importing `python` alone at the new commit, with the other snapshots at
    `b3c45a4`, failed `TestEmbeddedSetupSnapshotsAgreeOnEveryExternalSkill`,
    `TestSetupSnapshotDifferencesCompareARepositoryProfile` and
    `TestRestoreSkillsAcceptsARepositoryProfile` with "requires skill
    \"agentic-cli-design\", whose contracts disagree in Setup Snapshots
    \"go\" and \"python\"": six skills shared with `go`, `rust` and
    `typescript` (`knowledge-workspace`, `context7-cli`, `grilling`,
    `typesafe-ai`, `systematic-debugging`, `handoff`) changed bytes upstream.
    The `go`, `rust` and `typescript` refresh was therefore needed; since
    the batch 6 order, Spec 0317 performs it at `65d5f36` (ADR-0322).
  - The seed and one full sync wrote the five snapshots byte-identically in
    two independent worktrees; `--check` then exited `0`. The `python`
    snapshot holds 47 skills, 13 of them Roundfix-owned entries and both
    activation bundles carried from the seed.
  - `make baseline-digests` rewrote the three catalog files, the four plan
    characterization goldens and the parity `v1/manifest.json`; the managed
    refresh rewrote only `docs/agents/setup-context.json`, and a second
    refresh reported `File changes: 0`.
  - `baseline update --no-skills` left seven required skills differing from
    the refreshed `go` snapshot; `baseline skills restore --profile go-cli-tui`
    with those seven `--skill` flags and the confirmed plan digest changed
    ten files under `.agents/skills/` (eight rewritten,
    `.agents/skills/context7-cli/references/skills.md` deleted and
    `.agents/skills/golang-cli/assets/examples/go.mod` added) and seven
    `skills-lock.json` entries.
  - With the module, profile and declared test changes, `go test` passed for
    `./internal/baseline`, `./internal/cli`, `./internal/spec`,
    `./internal/speccheck`, `./internal/daemon`, `./internal/store`,
    `./internal/verifyselect`, `./internal/worktree`, `./skills` and
    `-tags docscontract ./internal/docscontract`.
  - The sabotages `_techspec.md` → Testing Approach lists each failed its
    named test.
- Re-amendment probe of 2026-10-24 (cycle AP batch 6), in a scratch worktree
  of `b5e45c63` after Spec 0317's import had moved `go`, `rust`,
  `typescript` and the composed snapshot to `65d5f36`: the seed and one sync
  from a repository fetched at that commit wrote only `python.json` (52
  skills), left the four other snapshots byte-identical, and
  `TestEmbeddedSetupSnapshotsAgreeOnEveryExternalSkill` and
  `TestThisRepositoryHoldsItsRequiredSkillsAtTheSetupSnapshot` passed with no
  skill restore. The earlier prototype's restore of seven skills therefore
  moves to Spec 0317, which restores them at `65d5f36`.
- Catalog facts: `catalog_validate.go` lets one module dispatch a skill
  (`catalog.skill.dispatch.id.duplicate`), and
  `TestEveryModuleRequiresEverySkillItDispatches` requires every dispatched
  skill to be required, so the four optional skills are required (installed
  with the setup) and dispatched conditionally; a decision effect can
  activate a module or select a template but has no skill-dispatch field
  (`catalog_validate.go`, the effect fields `when`, `activateModules`,
  `includeArtifacts`, `excludeArtifacts`, `selectTemplates`,
  `renderBindings`, `requireDecisions`).
- Secondbrain, read only: `wiki/index.md` (entity `quaestor`) and the
  2026-10-23 `qmd query`, which returned Quaestor's
  `projects/quaestor/mirror/docs/adr/0004-roundfix-e-baseline.md` and
  Roundfix's decision that repository profiles cannot add modules; the
  curation request
  `inbox/skills/2026-10-22-setup-python-e-skills-python-uv-pytest-ruff.md`,
  which upstream answered with pull request #130, was read, not written.
- Exa: uv's CLI reference (https://docs.astral.sh/uv/reference/cli/:
  "`--frozen` ... Instead of checking if the lockfile is up-to-date, uses the
  versions in the lockfile as the source of truth", "`--locked` ... Requires
  that the lockfile is up-to-date", "`--offline` ... uv will only use locally
  cached data") and "Locking and syncing", whose "Checking the lockfile"
  section says `uv lock --check` "is equivalent to the `--locked` flag for
  other commands" and asserts the lockfile "would remain unchanged after a
  resolution". The checking section adds: "uv will check if it matches the
  project metadata" (Context7, `astral-sh/uv`
  `docs/concepts/projects/sync.md`), and uv's settings resolve `--offline` to
  "only uses the cache and local sources" (`crates/uv-cli/src/settings.rs`).
- uv offline lock-check probe (2026-10-24, uv 0.13.5, scratch project with
  `iniconfig` and a `dev` group with `pytest`, `UV_CACHE_DIR` set to a fresh
  empty directory each time): `uv lock --check --offline` on an up-to-date
  lock printed "Resolved 7 packages in 2ms" and exited `0` with an empty
  cache, also with a `hatchling` backend and with `dynamic = ["version"]`; with
  a dependency added to `pyproject.toml` and not locked it exited `1`, with
  "The lockfile at `uv.lock` needs to be updated, but `--check` was
  provided." when the cache held the package and with "Packages were
  unavailable because the network was disabled" when it did not. The check of
  a current lock needs neither the network nor a cache, so a Run's gate can
  run it; a stale lock fails the gate either way, and CI's `uv sync --locked`
  names the cause with the network.
- ty's documentation (Context7, `astral-sh/ty` `docs/reference/cli.md` and
  `docs/reference/exit-codes.md`): "`--error-on-warning`: ty will exit with
  `1` if it finds any violations with severity `warning` or higher", it
  "cannot be supplied at the same time as `--exit-zero` or
  `--exit-zero-on-warning`", and `--python-platform` is the "Target platform
  to assume when resolving types", which "is used to specialize the type of
  `sys.platform`".
- Quaestor's tooling, read through `gh api` without cloning:
  `marcioaltoe/quaestor@e6b390d` (`Makefile`, `.roundfixrc.yml`,
  `.github/workflows/verify.yml`, `pyproject.toml`) and the `Makefile` at
  `d17c98a`, whose `verify` runs `lock-check` (`uv lock --check`),
  `fmt-check`, `lint` and `test` with `uv run --locked`; no client data was
  read.

## Glossary

None.

## Skills

- unchanged: guide docs/user-guide/context-driven-development.md — the only change names the new built-in `python` Baseline Profile in the profile list and its skills; the covering skills describe how a Baseline is adopted and refreshed, not which built-in profiles or setups exist, so no skill text changes.

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
  plano e mantém 0309/0310 na v0.123/v0.124.").
- 2026-10-24, about 11:15 BRT (AskUserQuestion): the preset — "Emendar a 0309
  agora (Recomendado)" — "Antes de a 0309 entrar na fila: importa o preset
  publicado (commit 32e4143e) e o módulo python passa a despachar uv, ruff e
  pytest para as skills novas; ty, typer, fastapi, pipelines e streamlit ficam
  só quando escolhidos. Mantém a v0.123."
- 2026-10-24, about 12:30 BRT (AskUserQuestion), on the questions of the
  amendment report:
  - Go and Rust: "Manter na 0309 (Recomendado)" — "A 0309 importa o python e
    atualiza Go e Rust juntos, pelo comando de sincronização já usado (Spec
    0208). A nota de release explica o `baseline skills restore --skill` para
    os adotantes.";
  - ty and the lock check: "ty e lock check sempre" — "Todo repositório
    Python recebe ty no verify e `uv lock --check`, igual ao Quaestor." This
    supersedes "ty ... só quando escolhidos" for ty only: the module dispatches
    ty always and the guide names ty and the lock check as part of the gate;
    Typer, FastAPI, file pipelines and Streamlit stay conditional;
  - optional skills: the conditional triggers (the report's recommendation).
- 2026-10-24, about 13:10 BRT (AskUserQuestion, cycle AP batch 6): order
  "Antes da 0307 (Recomendado)" — "A 0317 entra logo após a 0316. A 0307 é
  emendada para usar só o ai-memory (sem triagem de inbox), e a 0317 também
  tira a fonte secondbrain-inbox do comando follow-up. O setup Python vai uma
  versão para frente (v0.127)." Spec 0317 therefore moves `go`, `rust` and
  `typescript` to `65d5f36` before this Spec, which supersedes the 12:30 BRT
  "Manter na 0309" answer for the refresh: this Spec imports only `python`,
  at the same commit.

Autonomous route: no person present; decisions taken from the operator's
briefings and the characterization above, with no further interview:

- The `python` setup is the imported preset, at the commit every other
  upstream snapshot already names after Spec 0317; see ADR-0322. The Common Setup Snapshot of the first version
  is dropped.
- "Só quando escolhidos" is expressed by the dispatch trigger, the mechanism
  the catalog has; a decision per tool would need four modules.
- The lock check runs as `uv lock --check --offline`, proven by the probe
  above; see `_techspec.md` → Decisions.
- The profile is split from the starter files (Spec 0310) because the whole
  preset does not fit four implementation Tasks.

## Open Questions

None.
