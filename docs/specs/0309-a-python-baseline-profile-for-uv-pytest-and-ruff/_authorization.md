---
status: approved
granted: 2026-10-22
action: import the published upstream python preset through the asset sync at the commit every other upstream Setup Snapshot already names, ship the python module with its uv, ruff and pytest guide and its dispatched Python skills and the python Baseline Profile, which fixes no application surface, regenerate the derived Baseline files and this repository's Setup Manifest through the sanctioned commands, adjust the tests that enumerate the catalog, and document the profile
consuming: 0309-a-python-baseline-profile-for-uv-pytest-and-ruff
paths:
  - docs/agents/setup-context.json
  - internal/baseline/assets/coverage.json
  - internal/baseline/assets/modules/core.json
  - internal/baseline/assets/modules/python.json
  - internal/baseline/assets/profiles/python.json
  - internal/baseline/assets/setups/python.json
  - internal/baseline/assets/templates/guides/python.md
  - internal/baseline/assets/templates/index.json
  - internal/baseline/assets/templates/root/python.md
  - internal/baseline/derived_ownership_test.go
  - internal/cli/baseline_human_test.go
  - internal/cli/baseline_release_gate_test.go
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---


# Approved authority for Spec 0309

On 2026-10-22 the maintainer approved cycle AP through AskUserQuestion; the
questions of about 15:25 BRT were answered by about 18:30 BRT. The answers
this Spec rests on were:

- scope (multi-select, all four): "Correções de adotantes (Recomendado)",
  "Ruído C2 e suite docscontract", "Modos follow-up A e B", "Preset de
  baseline Python"; this Spec is the first half of the Python preset;
- order: "Correções, follow-up, Python (Recomendado)" (option text: "Segue a
  ordem que você definiu para o follow-up: primeiro o que trava outros
  projetos (Light Tier no Fluxus), depois o bom funcionamento, depois as
  features."), so the preset is delivered last;
- the preset's reach, the maintainer's own answer: "Perfil completo e pode, se
  necessário, solicitar ou requerer novas skills recomendadas ao repo skills
  para que tenhamos a curadoria e ter o melhor resultado no uso do python, uv
  e linters e outros toolings especificos do setup.";
- authorization: the same as cycle AO, that is the Governed Paths each Spec
  declares, and still not this repository's Makefile, `.github/workflows`,
  `go.mod`, `go.sum`, or lint and formatter configuration.

Three further cycle AP decisions were taken by the maintainer through
AskUserQuestion on 2026-10-23, at about 16:15 BRT, on the questions the
authoring report raised:

- Starter Files: "Comando starters (Recomendado)" (option text: "Keep the
  separate confirmation-gated `roundfix baseline starters` command, which
  writes the Makefile, `.roundfixrc.yml` and `.github/workflows/verify.yml`
  only where absent and only after the exact preview is confirmed (Spec
  0310).");
- lockfile flags: "--frozen e --locked (Recomendado)" (option text: "The
  generated Makefile runs `uv run --frozen`, so a Run never rewrites the
  lockfile or needs the network, and CI runs `uv sync --locked`, which
  refuses a lockfile that no longer matches `pyproject.toml`.");
- Python skills: "Entregar agora (Recomendado)" (option text: "Ship now with
  the rules in the Python guide, and adopt the Python skills in a later Spec
  once the curation request in the skills repository is answered.").

Two further cycle AP decisions were taken by the maintainer on 2026-10-24,
after the plans of Specs 0309 and 0310 merged and before either started
delivery:

- at about 09:10 BRT, in the maintainer's own words: "mude o nome do profile
  de python-cli para python (vai ser para cli, backend e workers)";
- through AskUserQuestion at about 09:15 BRT: "Nenhuma superfície fixa
  (Recomendado)" (option text: "O perfil traz só a linguagem (uv, pytest,
  ruff) e os módulos de fluxo (core, specs, autônomo, secondbrain, extensão).
  Serve igual para CLI, backend e workers; regras de uma superfície
  específica entram pela extensão do repositório. Muda pouco o plano e mantém
  0309/0310 na v0.123/v0.124.").

The profile identifier is therefore `python`, the Spec folders no longer say
`python-cli`, and the profile selects no `cli-surface` module and no
`rule.cli.*` rule. These decisions add no Governed Path and grant nothing
beyond the operations below.

One further cycle AP decision was taken by the maintainer on 2026-10-24,
through AskUserQuestion at about 11:15 BRT (cycle AP batch 5), after upstream
published its Python preset:

- the Python preset: "Emendar a 0309 agora (Recomendado)" — "Antes de a 0309
  entrar na fila: importa o preset publicado (commit 32e4143e) e o módulo
  python passa a despachar uv, ruff e pytest para as skills novas; ty, typer,
  fastapi, pipelines e streamlit ficam só quando escolhidos. Mantém a
  v0.123."

Importing the preset needs every upstream Setup Snapshot at one commit,
because a repository-owned profile refuses two contracts for one skill
(ADR-0241, measured in `_prd.md` → Validation). A first version of this
amendment moved the four existing snapshots and restored ten
`.agents/skills/` files itself; the batch 6 order below gives that move to
Spec 0317, so this record no longer bounds those paths, and this Spec
imports only the `python` snapshot, at the commit the others already name.

Three further answers were given by the maintainer through AskUserQuestion on
2026-10-24 at about 12:30 BRT, on the questions of the amendment report:

- Go/Rust: "Manter na 0309 (Recomendado)" — "A 0309 importa o python e
  atualiza Go e Rust juntos, pelo comando de sincronização já usado (Spec
  0208). A nota de release explica o `baseline skills restore --skill` para
  os adotantes.";
- ty / lock: "ty e lock check sempre" — "Todo repositório Python recebe ty no
  verify e `uv lock --check`, igual ao Quaestor." This supersedes "ty ... só
  quando escolhidos" for ty only: the `python` module dispatches ty always,
  the guide names ty and `uv lock --check --offline` as part of the gate,
  Spec 0310's Starter Makefile carries `lock-check` and `typecheck` in `make
  verify`, and CI keeps `uv sync --locked`; Typer, FastAPI, file pipelines
  and Streamlit stay conditional;
- optional skills: the report's recommended option (a), the conditional
  triggers (relayed by the coordinator by its letter, without the option
  label).

These answers add no Governed Path and grant nothing beyond the operations
below.

On 2026-10-24, at about 13:05 BRT, the maintainer wrote, in his own words:
"Lembrando que solicitei em algum momento para matar o inbox de secondbrain
e usar somente as mensagens do ai-memory". Through AskUserQuestion at about
13:10 BRT (cycle AP batch 6) he chose the order "Antes da 0307
(Recomendado)" (option text: "A 0317 entra logo após a 0316. A 0307 é
emendada para usar só o ai-memory (sem triagem de inbox), e a 0317 também
tira a fonte secondbrain-inbox do comando follow-up. O setup Python vai uma
versão para frente (v0.127)."). Spec 0317 moves `go`, `rust`, `typescript`
and the composed snapshot to upstream commit
`65d5f365cace27c4b468db1723213723b10c870b`, where upstream replaced
`handoff` with six `ai-memory-*` skills in every preset, `python` included,
and restores this repository's changed skills. This Spec therefore imports
`python` at that commit, no longer refreshes the other snapshots and
restores no skill: the `python` setup adds no skill this repository's
`go-cli-tui` profile requires. That answer supersedes the 12:30 BRT "Manter
na 0309" answer for the refresh only; it removes Governed Paths and adds
none.

The cycle AO rules carried by the authorization answer were approved on
2026-10-22, at about 00:55 BRT ("the same rules as cycle AN"), which carry the
cycle AN approval of 2026-10-21, at about 01:00 BRT, the cycle AM, AL and AK
approvals, and the cycle AJ grant of 2026-10-19, at about 16:20 BRT, "Mexer em
Governed Paths declarados" (option text: "'Concedo' para os caminhos que cada
Spec declarar (módulos Baseline, registros, skills, testes governados; fora:
Makefile, CI, go.mod, lint, que pedem autorização específica).") and "Merge
squash + release + push". The 2026-09-30 grants also stand: of the skills
"considere autorizado a ajustar todas as skills se necessário", and of
Baseline sources, guides and `.roundfixrc.yml` "Autorizar os dois".

The skills curation request the reach allowed was filed by the operator as
the Secondbrain Inbox Entry
`inbox/skills/2026-10-22-setup-python-e-skills-python-uv-pytest-ruff.md`;
upstream answered it with pull request #130. No Task, Verification command or
QA row writes it or reads it.

This record grants `implement`, `commit`, `push`, `pull_request` and `merge`.
It does not grant `release`: the operator releases. No Task, Verification
command or QA row reads or writes the real `~/.roundfix`, reaches the network,
sends a request to a provider, the judge or a reviewer, reads or writes the
Secondbrain, runs a Python tool, or reads another repository. Only task_01
reads the local checkout `~/dev/skills`, fetching one commit from it into an
empty local repository, and stops when it lacks the commit.

The governed set was measured on 2026-10-24 on the branch
`docs/plan-0309-python-preset`, based on `64c56d7f`, with
`speccheck.GovernedPath`, through a `go test -overlay` probe in
`internal/speccheck` that wrote nothing to the repository (`git status
--porcelain` identical before and after; the probe file never existed in the
tree). It covered every file the Tasks declare, every planning file of this
Spec and ADR-0322. Twenty-nine of those files were governed then; the batch 6
amendment removed the fourteen the Tasks no longer change (four upstream
snapshots and ten `.agents/skills/` files), leaving the twelve listed above,
which the Tasks change, and three the Tasks only read
(`internal/baseline/assets/modules/rust.json`,
`internal/baseline/assets/profiles/rust-cli.json`,
`skills/baseline_skill_contract_test.go`), which are not granted. The changed
set was measured by prototypes in scratch worktrees that ran the sanctioned
commands below in order.

## Why each governed path is unavoidable

- `internal/baseline/assets/setups/python.json` is the imported preset, seeded
  and then written by the asset sync, which leaves `go.json`, `rust.json`,
  `typescript.json` and `go-cli-typescript-bun.json` as Spec 0317 wrote them
  (task_01).
- `internal/baseline/assets/modules/python.json`,
  `internal/baseline/assets/templates/root/python.md`,
  `internal/baseline/assets/templates/guides/python.md` and
  `internal/baseline/assets/profiles/python.json` are the new module,
  templates and profile (task_02).
- `internal/baseline/assets/templates/index.json` and
  `internal/baseline/assets/coverage.json` register the templates and the
  coverage entry; `internal/baseline/assets/modules/core.json` places the root
  block in the hierarchy and the guide among the semantic owners, and takes
  the module version the record command writes (task_02).
- `docs/agents/setup-context.json` is rewritten by this repository's Managed
  Refresh after the catalog changes (task_01, task_02).
- `internal/baseline/derived_ownership_test.go` enumerates the sanctioned
  outputs and learns the new setup (task_01) and profile (task_02);
  `internal/cli/baseline_human_test.go` and
  `internal/cli/baseline_release_gate_test.go` enumerate the profiles
  (task_02).

## What is not governed

- `internal/baseline/module-versions.json`, the catalog
  files, plan characterization goldens and the parity files
  `v1/fixtures/asset-sync.json` and `v1/manifest.json` under
  `internal/baseline/testdata/`, and the tests `python_setup_test.go`,
  `python_profile_test.go`, `catalog_test.go`,
  `clause_characterization_test.go` and `assets_sync_test.go`.
- `docs/user-guide/context-driven-development.md`,
  `docs/references/skill-coverage.json`,
  `docs/references/behavior-surfaces.json`, ADR-0322 and every planning file
  of this Spec.

## Sanctioned regeneration

The repository-owned commands resolve the Baseline's derived files and the
generated records. This declaration records the regeneration that follows the
approved edits and adds no source path.

The asset sync writes the `python` setup snapshot from a repository fetched
from the local checkout at `65d5f365cace27c4b468db1723213723b10c870b`, after
the seed of `_techspec.md` → Data Model 1; its `--check` must then exit `0`,
and the other snapshots keep their bytes:

```yaml
command: roundfix baseline assets sync --source-dir <fetched>/setups --format text
outputs:
  - internal/baseline/assets/setups/python.json
```

The parity row follows through Spec 0208's jq transform, quoted with its
composed-snapshot filter in `_techspec.md` → API Contract 1, and the
regeneration fills its digests.

The Module Version Record step records the new module and writes the core
module's version line; suiteguard accepts those writes only for this declared
command:

```yaml
command: go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1
outputs:
  - internal/baseline/module-versions.json
  - internal/baseline/assets/modules/core.json
  - internal/baseline/assets/modules/python.json
```

```yaml
command: make baseline-digests
```

```yaml
command: go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1
outputs:
  - docs/references/behavior-surfaces.json
```

`roundfix baseline update --repo . --no-skills --no-history --yes --format text`,
run from a binary built under the Run's environment and with a fresh `HOME`,
renders this repository's Setup Manifest, and a second refresh reports no
file change. No digest pin, golden, generated guide, snapshot beyond the
seed's two values, skill file or lock entry is hand-edited, and no skill is
restored. No test is removed
or renamed, so the coverage record is not re-recorded.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`, the lint,
  formatter or test-runner configuration, `.roundfixrc.yml`, any module other
  than `core` and the new `python`, the existing profiles, the Source Baseline
  corpus, manifest, accounting or index, `internal/baseline/assets/retention/**`,
  the formatter fixtures, any parity fixture other than the two named above,
  any owned skill, any file under `.agents/skills/`, `skills-lock.json`, any
Setup Snapshot other than `python.json`, any Go
  production file, or the `### QA settlement` section of any skill.
- No test, Verification command or QA row:
  - reads or writes the real `~/.roundfix`;
  - reaches the network;
  - sends a request to a provider, the judge or a reviewer;
  - reads or writes the Secondbrain;
  - reads `~/dev/skills` (only task_01's import does) or another repository;
  - runs a Python tool;
  - records an environment value.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
