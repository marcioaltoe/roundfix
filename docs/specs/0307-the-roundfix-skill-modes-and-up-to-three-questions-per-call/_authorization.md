---
status: approved
granted: 2026-10-22
action: reword the Baseline clause on cheap user decisions to allow up to three independent questions per call, regenerate the derived Baseline files and this repository's guides through the sanctioned commands, and give the Roundfix skill its follow-up, status and unattended Skill Modes with their reference, metadata, tests and glossary term
consuming: 0307-the-roundfix-skill-modes-and-up-to-three-questions-per-call
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/modes.md
  - docs/agents/agent-instructions.md
  - docs/agents/setup-context.json
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/agent-instructions.md
  - internal/baseline/assets/modules/core.json
  - internal/baseline/assets/profiles/standard-typescript-monorepo.json
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0307

On 2026-10-22 the maintainer approved cycle AP through AskUserQuestion; the
questions of about 15:25 BRT were answered by about 18:30 BRT. The answers
this Spec rests on were:

- scope (multi-select, all four): "Correções de adotantes (Recomendado)",
  "Ruído C2 e suite docscontract", "Modos follow-up A e B", "Preset de
  baseline Python"; this Spec is follow-up Spec B;
- order: "Correções, follow-up, Python (Recomendado)" (option text: "Segue a
  ordem que você definiu para o follow-up: primeiro o que trava outros
  projetos (Light Tier no Fluxus), depois o bom funcionamento, depois as
  features.");
- authorization: the same as cycle AO, that is the Governed Paths each Spec
  declares, and still not this repository's Makefile, `.github/workflows`,
  `go.mod`, `go.sum`, or lint and formatter configuration.

The follow-up design was settled with the maintainer the same day, through
AskUserQuestion at about 15:20 BRT, after his request relayed at 15:09 and
15:12 BRT:

- where: "Skill roundfix, fontes opcionais (Recomendado)" (option text: "Os
  modos entram na skill do produto. Cada fonte é detectada e opcional:
  Secondbrain pelo guia do baseline, ai-memory se o MCP existir, Jev se houver
  chave, bloco do GT3 só se as instruções do usuário o definirem. Sem a fonte,
  o modo segue e diz o que faltou.");
- code: "Comando CLI + skill (Recomendado)" (option text: "Novo comando só
  leitura (ex.: roundfix follow-up --format json) junta backlog, findings,
  specs, Runs, fila, PRs, git e o inbox do Secondbrain, pergunta ao Jev pelo
  internal/judge (model id no log) e ordena em código. A skill junta as
  mensagens do ai-memory (só o agente vê o MCP), conversa e pede as
  autorizações. Testável e igual nos três harnesses."); the command is Spec
  0306's;
- questions per call: "Até 3 no geral", his own answer, which replaces the
  one-question rule of `clause.core.ask-user-answerable-decisions` for every
  adopter;
- answers: "As respostas e autorizacoes devem ser parte da spec que será
  implementada, como o _autorizations.md e se não for virar uma spec deve
  virar um backlog ou finding e armazenado com as informaçoes e decisoes (sem
  autorizacoes) para virar uma spec no futuro.", his own answer;
- "Sim, suas decisões bastam (Recomendado)": the operator does not wait for
  the assistant's reply.

On 2026-10-24, at about 11:35 BRT, before this Spec entered the Delivery
Queue, the maintainer amended it through AskUserQuestion. The question was "A
skill deve aceitar o modo sem os traços (`/roundfix follow-up`, `/roundfix
status`, `/roundfix unattended`) além de `--follow-up`?", and the answer was
"Aceitar as duas formas (Recomendado)" (option text: "Emendo a 0307 antes de
ela entrar na fila: a primeira palavra depois do nome da skill, com ou sem
`--`, escolhe o modo nos três harnesses. Sem custo de versão."). The
amendment changes the skill text, the modes reference, the contracts, the
tests and the QA rows that describe how a mode is chosen; it adds no Governed
Path, operation or grant, so the governed set and the operations below stand.

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

This record grants `implement`, `commit`, `push`, `pull_request` and `merge`.
It does not grant `release`: the operator releases. No Task, Verification
command or QA row reads or writes the real `~/.roundfix`, reaches the network,
sends a request to a provider, the judge or a reviewer, reads or writes the
Secondbrain, calls ai-memory, or reads another repository.

The governed set was measured on the authoring branch `docs/plan-0307`, based
on v0.109.0, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` was identical before
  and after, and the probe file never existed in the tree.
- It covered every file the Tasks declare, every planning file of this Spec,
  ADR-0320, `CONTEXT.md`,
  `docs/references/behavior-surfaces.json` and
  `docs/references/skill-coverage.json`.
- Eight of those files are governed, all listed above. The derived files among
  them were measured by a prototype in a scratch worktree that applied the
  reworded clause and ran the sanctioned commands below in order.

## Why each governed path is unavoidable

- `internal/baseline/assets/modules/core.json` holds the reworded clause, its
  rule version and the module version the record command writes (task_01).
- The Standard TypeScript golden `docs/agents/agent-instructions.md` and the
  Standard TypeScript profile are rewritten by `make baseline-digests` from
  the core edit (task_01).
- `docs/agents/agent-instructions.md` and `docs/agents/setup-context.json`
  are rewritten by this repository's Managed Refresh after the core edit
  (task_01).
- `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md` carry the
  modes section, the metadata, the index row and the version fields
  (task_02, task_03).
- `.agents/skills/roundfix/references/modes.md` is the new modes reference
  (task_02, task_03).

## What is not governed

- `internal/baseline/module-versions.json`, the catalog files and plan
  characterization goldens under `internal/baseline/testdata/`, and the new
  test file `internal/baseline/question_batch_clause_test.go`.
- `skills/roundfix/references/modes.md`, `skills/roundfix_layout_test.go`,
  `skills/roundfix_modes_test.go` and
  `skills/testdata/owned-skill-versions.json`.
- `CONTEXT.md`, ADR-0320 and every planning file of this Spec.

## Sanctioned regeneration

The repository-owned commands resolve the Baseline's derived files, the skill
mirrors and the generated records. This declaration records the regeneration
that follows the approved edits and adds no source path.

The Module Version Record step writes the core module's version line and the
record; suiteguard accepts those writes only for this declared command:

```yaml
command: go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1
outputs:
  - internal/baseline/module-versions.json
  - internal/baseline/assets/modules/core.json
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

`go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text`
renders this repository's guides and Setup Manifest, and a second refresh
reports no file change. No digest pin, golden, generated guide, mirror or
record is hand-edited. No user guide page changes, so the Behavior Surface
Record is not re-recorded; no command surface is added, so
`docs/references/skill-coverage.json` does not change; no test is removed or
renamed, so the coverage record is not re-recorded.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`, the lint,
  formatter or test-runner configuration, `.roundfixrc.yml`, any module other
  than `core`, the Source Baseline corpus, manifest, accounting or index,
  `internal/baseline/assets/retention/**`, the Setup Snapshots, or the
  `### QA settlement` section of any skill.
- No change to an owned skill other than `roundfix`, to the follow-up
  command, its guide or its reference, to any user guide page, or to Go
  production code.
- No test, Verification command or QA row:
  - reads or writes the real `~/.roundfix`;
  - reaches the network;
  - sends a request to a provider, the judge or a reviewer;
  - reads or writes the Secondbrain or calls ai-memory;
  - reads another repository;
  - records an environment value.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
