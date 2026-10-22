---
status: approved
granted: 2026-10-22
action: make a light Task whose light runtime offers no mode for the requested access policy skip the Light Tier with the reason code access_mode_unavailable, make any other refusal of a light session before Agent work yield to the category's Preferred Selection, add a Doctor light-tier line that warns about the combination, and say so in the configuration and Doctor guides, the Roundfix runtime and setup references and the glossary
consuming: 0301-the-light-tier-yields-to-the-category-profile
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/runtime.md
  - .agents/skills/roundfix/references/setup.md
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0301

On 2026-10-22 the maintainer approved cycle AP through AskUserQuestion
(questions of 15:25 BRT, answered by about 18:30 BRT). The answers this Spec
rests on were:

- scope (multi-select, all four): "Correções de adotantes (Recomendado)",
  "Ruído C2 e suite docscontract", "Modos follow-up A e B", "Preset de
  baseline Python"; this Spec is the first adopter correction and delivers
  first in batch 1, after Spec 0300;
- order "Correções, follow-up, Python (Recomendado)" (option text: "Segue a
  ordem que você definiu para o follow-up: primeiro o que trava outros
  projetos (Light Tier no Fluxus), depois o bom funcionamento, depois as
  features.");
- Light Tier "Pular no planejamento + escalar (Recomendado)" (option text:
  "Quando o acesso pedido é full-access e o runtime light não o atende, a
  Task nem escolhe o Light Tier (evento light_tier_skipped, motivo
  access_mode_unavailable). Qualquer outra recusa da sessão light antes do
  trabalho escala para o perfil normal da categoria. O doctor avisa a
  combinação.");
- authorization: the same as cycle AO, that is the Governed Paths each Spec
  declares, and still not this repository's Makefile, `.github/workflows`,
  `go.mod`, `go.sum`, or lint and formatter configuration.

The cycle AO rules that answer carries were approved on 2026-10-22, at about
00:55 BRT ("the same rules as cycle AN"), which carry the cycle AN approval of
2026-10-21, at about 01:00 BRT ("the same rules as cycle AM"), the cycle AM
approval of 2026-10-20, at about 18:15–18:35 BRT ("the same rules as cycle
AL"), the cycle AL approval of 2026-10-20, at about 13:10 BRT ("Mesmas regras
do AK (Recomendado)"), and the cycle AJ grant of 2026-10-19, at about 16:20
BRT, "Mexer em Governed Paths declarados" (option text: "'Concedo' para os
caminhos que cada Spec declarar (módulos Baseline, registros, skills, testes
governados; fora: Makefile, CI, go.mod, lint, que pedem autorização
específica).") and "Merge squash + release + push". The 2026-09-30 skills
grant ("considere autorizado a ajustar todas as skills se necessário") also
stands.

This record grants `implement`, `commit`, `push`, `pull_request` and `merge`.
It does not grant `release`: the operator releases. No Task, Verification
command or QA row reads or writes the real `~/.roundfix`, starts a real ACP
runtime, sends a prompt to OpenRouter or any provider, reaches the network or
reads another repository.

The governed set was measured on the authoring branch `docs/plan-0301`, based
on `42f523c2`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` was identical before
  and after, and the probe file never existed in the tree.
- It covered every file the Tasks declare, every planning file of this Spec,
  its adopted reference, ADR-0314, `docs/references/skill-coverage.json` and
  `CONTEXT.md`.
- Four of those files are governed, and all four are changed and listed
  above.

## Why each governed path is unavoidable

- `.agents/skills/roundfix/references/runtime.md` lists the Light Tier's skip
  reason codes and its escalation, and gains `access_mode_unavailable` and the
  pre-work yield (task_01).
- `.agents/skills/roundfix/references/setup.md` describes every Doctor line
  and finding code, and gains `light-tier:` and `DR-LIGHT-ACCESS-MODE`
  (task_01).
- `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md` carry the
  Roundfix skill's version fields, raised with its content (task_01).

## What is not governed

- The production and test sources task_02 and task_03 change or create in
  `internal/lighttier`, `internal/daemon` and `internal/cli`.
- The skill reference mirrors `skills/roundfix/references/runtime.md` and
  `skills/roundfix/references/setup.md`.
- `skills/testdata/owned-skill-versions.json` and
  `docs/references/behavior-surfaces.json`.
- `CONTEXT.md`, `docs/user-guide/configuration.md` and
  `docs/user-guide/commands/doctor.md`.
- Every planning file of this Spec, its adopted reference and ADR-0314.

## Sanctioned regeneration

The repository-owned commands resolve the skill mirrors, the skill versions
and the Behavior Surface Record after task_01's edits. This declaration
records the regeneration that follows the approved edits and adds no source
path.

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

No Baseline module changes, so no module version is recorded. An owned skill
edit leaves the derived Baseline artifacts byte-identical, so `make
baseline-digests` has nothing to rewrite. No test is removed or renamed and no
platform set changes, so the coverage record is not re-recorded. No command
surface is added, so `docs/references/skill-coverage.json` does not change.
No record or mirror is hand-edited.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`, the lint,
  formatter or test-runner configuration, any Baseline module,
  `.roundfixrc.yml`, `docs/agents/setup-context.json`, the Run Database
  schema or the `### QA settlement` section of any skill.
- No change to `internal/agent`, `roundfix profiles validate`, the full-access
  mode choice, the Light Spend Log, the ceiling, the escalation after a
  failed Verification or the selection fallback of non-light Tasks.
- No test, Verification command or QA row:
  - reads or writes the real `~/.roundfix`;
  - starts a real ACP runtime or adapter;
  - reaches the network;
  - sends a request to a provider, OpenRouter, the judge or a reviewer;
  - reads another repository;
  - records an environment value.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
