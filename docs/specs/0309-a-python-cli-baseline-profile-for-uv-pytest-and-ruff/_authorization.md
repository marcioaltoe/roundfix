---
status: approved
granted: 2026-10-22
action: add a Common Setup Snapshot kind to the Baseline, ship the python common setup of the go and rust setups, the python module with its uv, ruff and pytest guide and the python-cli Baseline Profile, regenerate the derived Baseline files and this repository's Setup Manifest through the sanctioned commands, adjust the tests that enumerate the catalog, and document the profile with its glossary term
consuming: 0309-a-python-cli-baseline-profile-for-uv-pytest-and-ruff
paths:
  - docs/agents/setup-context.json
  - internal/baseline/assets/coverage.json
  - internal/baseline/assets/modules/core.json
  - internal/baseline/assets/modules/python.json
  - internal/baseline/assets/profiles/python-cli.json
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

The skills curation request the reach allows is already filed by the operator
as the Secondbrain Inbox Entry
`inbox/skills/2026-10-22-setup-python-e-skills-python-uv-pytest-ruff.md`;
no Task, Verification command or QA row writes it or reads it.

This record grants `implement`, `commit`, `push`, `pull_request` and `merge`.
It does not grant `release`: the operator releases. No Task, Verification
command or QA row reads or writes the real `~/.roundfix`, reaches the network,
sends a request to a provider, the judge or a reviewer, reads or writes the
Secondbrain, reads the upstream skills checkout, runs a Python tool, or reads
another repository.

The governed set was measured on the authoring branch `docs/plan-0309`, based
on `dde9d5f5` (2026-10-23), with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` was identical before
  and after, and the probe file never existed in the tree.
- It covered every file the Tasks declare, every planning file of this Spec
  and ADR-0322. Fifteen of those files are governed: the twelve listed above,
  which the Tasks change, and three the Tasks only read
  (`internal/baseline/assets/modules/rust.json`,
  `internal/baseline/assets/profiles/rust-cli.json`,
  `skills/baseline_skill_contract_test.go`), which are not granted. The
  derived files were measured by a prototype in a scratch worktree that
  applied the edits and ran the sanctioned commands below in order.

## Why each governed path is unavoidable

- `internal/baseline/assets/setups/python.json`,
  `internal/baseline/assets/modules/python.json`,
  `internal/baseline/assets/templates/root/python.md`,
  `internal/baseline/assets/templates/guides/python.md` and
  `internal/baseline/assets/profiles/python-cli.json` are the new setup,
  module, templates and profile (task_02).
- `internal/baseline/assets/templates/index.json` and
  `internal/baseline/assets/coverage.json` register the templates and the
  coverage entry; `internal/baseline/assets/modules/core.json` places the root
  block in the hierarchy and the guide among the semantic owners, and takes
  the module version the record command writes (task_02).
- `docs/agents/setup-context.json` is rewritten by this repository's Managed
  Refresh after the catalog changes (task_02).
- `internal/baseline/derived_ownership_test.go`,
  `internal/cli/baseline_human_test.go` and
  `internal/cli/baseline_release_gate_test.go` enumerate the catalog's setups
  and profiles and learn the new entries (task_02).

## What is not governed

- `internal/baseline/setup_composition.go`, `catalog_validate.go`,
  `assets_sync.go`, `internal/baseline/module-versions.json`, the catalog files
  and plan characterization goldens under `internal/baseline/testdata/`, and
  the tests `common_setup_test.go`, `python_profile_test.go`,
  `catalog_test.go`, `clause_characterization_test.go` and
  `assets_sync_test.go`.
- `docs/user-guide/context-driven-development.md`, `CONTEXT.md`,
  `docs/references/skill-coverage.json`,
  `docs/references/behavior-surfaces.json`, ADR-0322 and every planning file
  of this Spec.

## Sanctioned regeneration

The repository-owned commands resolve the Baseline's derived files and the
generated records. This declaration records the regeneration that follows the
approved edits and adds no source path.

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

`go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text`
renders this repository's Setup Manifest, and a second refresh reports no file
change. No digest pin, golden, generated guide or record is hand-edited. No
test is removed or renamed, so the coverage record is not re-recorded.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`, the lint,
  formatter or test-runner configuration, `.roundfixrc.yml`, any module other
  than `core` and the new `python`, the existing setups and profiles, the
  Source Baseline corpus, manifest, accounting or index,
  `internal/baseline/assets/retention/**`, the formatter fixtures, any skill,
  or the `### QA settlement` section of any skill.
- No test, Verification command or QA row:
  - reads or writes the real `~/.roundfix`;
  - reaches the network;
  - sends a request to a provider, the judge or a reviewer;
  - reads or writes the Secondbrain;
  - reads the upstream skills checkout or another repository;
  - runs a Python tool;
  - records an environment value.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
