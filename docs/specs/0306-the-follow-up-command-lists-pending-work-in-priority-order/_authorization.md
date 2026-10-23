---
status: approved
granted: 2026-10-22
action: add the read-only roundfix follow-up command with its status view, two Secondbrain User Config keys and the follow-up-priority judgment, and describe them in a new command guide, a new Roundfix skill reference, the configuration guide and the glossary
consuming: 0306-the-follow-up-command-lists-pending-work-in-priority-order
paths:
  - .agents/skills/roundfix/references/follow-up.md
  - .agents/skills/roundfix/SKILL.md
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0306

On 2026-10-22 the maintainer approved cycle AP through AskUserQuestion; the
questions of about 15:25 BRT were answered by about 18:30 BRT. The answers
this Spec rests on were:

- scope (multi-select, all four): "Correções de adotantes (Recomendado)",
  "Ruído C2 e suite docscontract", "Modos follow-up A e B", "Preset de
  baseline Python"; this Spec is follow-up Spec A;
- order: "Correções, follow-up, Python (Recomendado)" (option text: "Segue a
  ordem que você definiu para o follow-up: primeiro o que trava outros
  projetos (Light Tier no Fluxus), depois o bom funcionamento, depois as
  features.");
- authorization: the same as cycle AO, that is the Governed Paths each Spec
  declares, and still not this repository's Makefile, `.github/workflows`,
  `go.mod`, `go.sum`, or lint and formatter configuration.

The follow-up design answers of 2026-10-22, at about 15:20 BRT, are:

- where: "Skill roundfix, fontes opcionais (Recomendado)" (option text: "Os
  modos entram na skill do produto. Cada fonte é detectada e opcional:
  Secondbrain pelo guia do baseline, ai-memory se o MCP existir, Jev se
  houver chave, bloco do GT3 só se as instruções do usuário o definirem. Sem
  a fonte, o modo segue e diz o que faltou.");
- code: "Comando CLI + skill (Recomendado)" (option text: "Novo comando só
  leitura (ex.: roundfix follow-up --format json) junta backlog, findings,
  specs, Runs, fila, PRs, git e o inbox do Secondbrain, pergunta ao Jev pelo
  internal/judge (model id no log) e ordena em código. A skill junta as
  mensagens do ai-memory (só o agente vê o MCP), conversa e pede as
  autorizações. Testável e igual nos três harnesses.");
- questions per call: "Até 3 no geral" (the maintainer's own answer; Spec
  0307 applies it);
- answers: "As respostas e autorizacoes devem ser parte da spec que será
  implementada, como o _autorizations.md e se não for virar uma spec deve
  virar um backlog ou finding e armazenado com as informaçoes e decisoes (sem
  autorizacoes) para virar uma spec no futuro." (the maintainer's own answer;
  Spec 0307 applies it);
- waiting for `default/jarvis`: "Sim, suas decisões bastam (Recomendado)".

The cycle AO rules that the authorization answer carries were approved on
2026-10-22, at about 00:55 BRT ("the same rules as cycle AN"), which carry
the cycle AJ grant of 2026-10-19, at about 16:20 BRT, "Mexer em Governed
Paths declarados" (option text: "'Concedo' para os caminhos que cada Spec
declarar (módulos Baseline, registros, skills, testes governados; fora:
Makefile, CI, go.mod, lint, que pedem autorização específica).") and "Merge
squash + release + push". The 2026-09-30 skills grant ("considere autorizado
a ajustar todas as skills se necessário") also stands. The judge's standing
authorization of 2026-10-05 covers sending this repository's planning
artifacts to TypeSafe during authoring; the product's own requests to the
judge are the maintainer's "pergunta ao Jev pelo internal/judge" above.

This record grants `implement`, `commit`, `push`, `pull_request` and `merge`.
It does not grant `release`: the operator releases.

The governed set was measured on the authoring branch `docs/plan-0306`, based
on `42f523c2`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` was identical before
  and after, and the probe file never existed in the tree.
- It covered every file the Tasks declare, every planning file of this Spec,
  its adopted reference and ADR-0319.
- Three of those files are governed, and all three are changed and listed
  above.

## Why each governed path is unavoidable

- `.agents/skills/roundfix/references/follow-up.md` is the one Roundfix
  reference of the new `follow-up` command family, which the Skill Coverage
  Map and the reference index require for a new command (task_01).
- `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md` carry the
  reference index row and the Roundfix skill's version fields, raised with
  its content (task_01).

## What is not governed

- The production and test sources task_02, task_03 and task_04 change in
  `internal/judge`, `internal/config` and `internal/cli`.
- The skill reference mirror `skills/roundfix/references/follow-up.md` and
  the layout test `skills/roundfix_layout_test.go`.
- `skills/testdata/owned-skill-versions.json`,
  `docs/references/skill-coverage.json` and
  `docs/references/behavior-surfaces.json`.
- `CONTEXT.md` and the user guides `docs/user-guide/commands/follow-up.md`,
  `docs/user-guide/commands.md`, `docs/user-guide/usage.md` and
  `docs/user-guide/configuration.md`.
- Every planning file of this Spec, its adopted reference and ADR-0319.

## Sanctioned regeneration

The repository-owned commands resolve the skill mirror, the skill version and
the Behavior Surface Record after the Tasks' edits. This declaration records
the regeneration that follows the approved edits and adds no source path.

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

No Baseline module changes, so no module version is recorded and `make
baseline-digests` has nothing to rewrite. No test is removed or renamed and
no platform set changes, so the coverage record is not re-recorded. No record
or mirror is hand-edited.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`, the lint,
  formatter or test-runner configuration, any Baseline module,
  `.roundfixrc.yml`, `docs/agents/`, `internal/store/`, `internal/delivery/`,
  `internal/releaseplan/` or the `### QA settlement` section of any skill.
- No change to `roundfix release plan`, `roundfix deliver status`,
  `roundfix runs list`, `roundfix spec judge` or the existing judgment kinds'
  questions and thresholds.
- No test, Verification command or QA row:
  - reads or writes the real `~/.roundfix` or the real Secondbrain;
  - reaches the network or runs the real `gh`;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository;
  - records an environment value.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
