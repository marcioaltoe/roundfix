---
status: approved
granted: 2026-10-22
action: let a Baseline module declare Starter Files, give the python module its Makefile, Project Config and CI workflow starters, add the confirmation-gated roundfix baseline starters command that writes them only where absent, name it in the profile alignment, regenerate the derived Baseline files and this repository's Setup Manifest through the sanctioned commands, and document the command in the guides, the Roundfix skill and the glossary
consuming: 0310-the-python-cli-profile-writes-its-starter-files
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/baseline.md
  - docs/agents/setup-context.json
  - internal/baseline/assets/modules/python.json
  - internal/baseline/assets/templates/index.json
  - internal/baseline/assets/templates/starters/python/Makefile
  - internal/baseline/assets/templates/starters/python/roundfixrc.yml
  - internal/baseline/assets/templates/starters/python/verify.yml
  - internal/cli/baseline_documentation_contract_test.go
  - internal/docscontract/publicdocs_test.go
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0310

On 2026-10-22 the maintainer approved cycle AP through AskUserQuestion; the
questions of about 15:25 BRT were answered by about 18:30 BRT. The answers
this Spec rests on were:

- scope (multi-select, all four): "Correções de adotantes (Recomendado)",
  "Ruído C2 e suite docscontract", "Modos follow-up A e B", "Preset de
  baseline Python"; this Spec is the second half of the Python preset, split
  from Spec 0309 because the whole preset does not fit four implementation
  Tasks;
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
Secondbrain, runs uv, pytest, ruff, a starter Makefile target or a workflow,
or reads another repository.

The governed set was measured on the authoring branch `docs/plan-0309`, based
on `dde9d5f5` (2026-10-23), with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` was identical before
  and after, and the probe file never existed in the tree.
- It covered every file the Tasks declare, every planning file of this Spec
  and ADR-0323. Eleven of those files are governed, all listed above. The
  derived files were measured by a prototype in a scratch worktree, on top of
  Spec 0309's prototype, that applied the declarations and ran the sanctioned
  commands below in order.

## Why each governed path is unavoidable

- The three files under `internal/baseline/assets/templates/starters/python/`
  are the Starter File templates; they are Baseline assets rendered into an
  adopter on its maintainer's confirmation, not this repository's own
  Makefile, Project Config or workflow (task_01).
- `internal/baseline/assets/templates/index.json` registers them, and
  `internal/baseline/assets/modules/python.json` declares them and takes the
  module version the record command writes (task_01).
- `docs/agents/setup-context.json` is rewritten by this repository's Managed
  Refresh after the catalog changes (task_01).
- `internal/cli/baseline_documentation_contract_test.go` and
  `internal/docscontract/publicdocs_test.go` hold the Baseline command
  documentation contracts the new subcommand joins (task_03).
- `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md` and
  `.agents/skills/roundfix/references/baseline.md` name the command, carry
  the version fields and describe the command (task_03).

## What is not governed

- `internal/baseline/catalog_load.go`, `catalog_validate.go`,
  `starter_files.go`, `starters.go`, `profile_alignment.go`,
  `internal/cli/baseline_starters.go`, `baseline_profile.go`, `cli.go`,
  `internal/baseline/module-versions.json`, the catalog files and plan
  characterization goldens under `internal/baseline/testdata/`, and the new
  tests.
- `skills/roundfix/references/baseline.md`,
  `skills/testdata/owned-skill-versions.json`,
  `docs/user-guide/commands/baseline.md`, `CONTEXT.md`,
  `docs/references/skill-coverage.json`,
  `docs/references/behavior-surfaces.json`, ADR-0323 and every planning file
  of this Spec.

## Sanctioned regeneration

The repository-owned commands resolve the Baseline's derived files, the skill
mirrors and the generated records. This declaration records the regeneration
that follows the approved edits and adds no source path.

The Module Version Record step writes the `python` module's version line and
the record; suiteguard accepts those writes only for this declared command:

```yaml
command: go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1
outputs:
  - internal/baseline/module-versions.json
  - internal/baseline/assets/modules/python.json
```

```yaml
command: make baseline-digests
```

```yaml
command: make skills-sync
```

```yaml
command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
outputs:
  - skills/testdata/owned-skill-versions.json
```

```yaml
command: go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1
outputs:
  - docs/references/behavior-surfaces.json
```

`go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text`
renders this repository's Setup Manifest, and a second refresh reports no file
change. No digest pin, golden, generated guide, mirror or record is
hand-edited. No test is removed or renamed, so the coverage record is not
re-recorded.

## Limits

- No change to this repository's Makefile, CI workflows, `go.mod`, `go.sum`,
  lint, formatter or test-runner configuration, or `.roundfixrc.yml`; to any
  module other than `python`, any profile or setup, the Source Baseline
  corpus, manifest, accounting or index, `internal/baseline/assets/retention/**`
  or the formatter fixtures; to any owned skill other than `roundfix`; or to
  the `### QA settlement` section of any skill.
- No test, Verification command or QA row:
  - reads or writes the real `~/.roundfix`;
  - reaches the network;
  - sends a request to a provider, the judge or a reviewer;
  - reads or writes the Secondbrain;
  - reads another repository;
  - runs uv, pytest, ruff, a starter Makefile target or a workflow;
  - records an environment value.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
