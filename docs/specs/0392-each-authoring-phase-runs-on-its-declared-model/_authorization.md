---
status: approved
granted: 2026-11-14
action: add the Authoring Model Declaration to Project Config, apply it when the Roundfix skills are installed in a project (a header in the four authoring skills and generated subagents under .claude/agents), report an unapplied declaration and user-level skill copies with the command that removes them, add roundfix skills remove, and change the four authoring skills and the Roundfix skill's setup, baseline and modes references to delegate each phase and its sweeps
consuming: 0392-each-authoring-phase-runs-on-its-declared-model
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/baseline.md
  - .agents/skills/roundfix/references/modes.md
  - .agents/skills/roundfix/references/setup.md
  - .agents/skills/write-idea/SKILL.md
  - .agents/skills/write-prd/SKILL.md
  - .agents/skills/write-tasks/SKILL.md
  - .agents/skills/write-techspec/SKILL.md
  - skills/roundfix/SKILL.md
  - skills/write-idea/SKILL.md
  - skills/write-prd/SKILL.md
  - skills/write-tasks/SKILL.md
  - skills/write-techspec/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
  - release
release_ceiling: minor
release_paths:
  - CHANGELOG.md
  - dist/npm/roundfix/package.json
  - internal/app/version.go
---

# Approved authority for Spec 0392

The Spec adopts two ai-memory messages, reported by indexus and reported by
pantheon on 2026-11-13, and the maintainer approved its Idea on 2026-11-14.
The maintainer's answers (AskUserQuestion via `/roundfix --follow-up`,
2026-11-14 ~06:45–08:05 BRT) are recorded whole in
[`references/2026-11-14-maintainer-answers.md`](references/2026-11-14-maintainer-answers.md).
The answers this record rests on are, verbatim:

- Scope. Question "O Indexus (mensagem 01a1d68d, pedido seu de 13/11) pede que
  `roundfix skills install --target project` e o Baseline escrevam `model` e
  `effort` no cabeçalho das skills e gerem um subagente de pesquisa em
  `.claude/agents/`, a partir de uma seção `skills` do `.roundfixrc.yml`. É
  feature (PRD + TechSpec + Tasks, release minor v0.144.0). O que faço com o
  pedido?" → "Spec completa" — option text: "Uma Spec com a seção `skills` no
  Project Config, o cabeçalho model/effort na instalação, o subagente de
  pesquisa só leitura e as skills de autoria passando a usá-lo. Cabe em 4
  Tasks + QA."
- The user-level copy notice. Question "Você respondeu "Spec pequena" para o
  ajuste do aviso de cópia de usuário (Pantheon). Pelo motivo que você
  lembrou, ela seria justamente uma Spec de 1 Task pagando um ciclo inteiro de
  Run, QA, review e release. O ajuste é do mesmo contexto da Spec do Indexus
  (instalação de skills e cópias de usuário), e ninguém está bloqueado por
  ele. Junto ou mantenho separado?" → "Juntar à Spec do Indexus (Recommended)"
  — option text: "O ajuste do aviso entra numa Task da Spec do Indexus que já
  mexe em `skills check`. Um ciclo de entrega e uma release (minor v0.144.0);
  a Spec adota as duas mensagens e continua em 4 Tasks + QA." The change
  itself was the earlier answer "Spec pequena (Recommended)" — option text:
  "Correção própria: o aviso passa a oferecer apagar a cópia de usuário, com o
  comando, ao lado do refresh. 1 Task + QA, release patch. A memória antiga é
  aposentada de qualquer forma."
- The shape. Question "Etapa de Ideia da 0392, pergunta 4 (a forma). Medi que
  o Claude Code 2.1.325 ignora o `model` do cabeçalho da skill quando é o
  agente que a invoca (só aplica quando você digita `/skill`), e o seu
  critério de sucesso inclui esse caminho. Qual forma a ideia toma?" → "Mais
  ambiciosa" — option text: "A híbrida, mais um subagente autor gerado com o
  modelo forte e a mesma declaração compilada para arquivos de agente do
  OpenCode e do Codex. Duas Specs." This Spec is the first of the two; the
  OpenCode and Codex agent files are Spec 0394.
- The draft. "Rascunho da Ideia 0392 (cada fase de autoria roda no modelo
  declarado): duas Specs, a 0392 para o Claude Code (declaração no Project
  Config, subagente autor e subagente de pesquisa gerados, cabeçalho na
  instalação, aviso de cópia de usuário) e a 0394 para OpenCode e Codex. […]
  Aprova?" → "Aprovado (Recommended)".
- Authorization. Question "As Specs que saírem das duas respostas acima
  precisam de autorização de push, merge e release. A entrega local só entra
  na fila com o Quaestor ocioso (regra do mini, 1 vaga de verificação). Qual
  autorização vale para elas?" → "Mesmas do BE (Recommended)" — option text:
  "Push, PR e merge; release após cada Spec (patch para correção, minor para
  feature); até 3 correções por Spec; Makefile, CI, go.mod e lint só quando
  uma correção exigir, descrito no relatório."

"Mesmas do BE" keeps the cycle BE rules of 2026-11-13 and, through them, the
cycle BD, BC, BB, BA, AZ, AY, AX, AW, AV, AU, AT and AS answers, which release
to the Runs the owned skills and changes of skills that reach the adopters,
the cycle AJ grant of 2026-10-19 "Mexer em Governed Paths declarados" (option
text: "'Concedo' para os caminhos que cada Spec declarar (módulos Baseline,
registros, skills, testes governados; fora: Makefile, CI, go.mod, lint, que
pedem autorização específica).") and the 2026-09-30 grant of every skill
("considere autorizado a ajustar todas as skills se necessário").

The Makefile, `.github/workflows`, `go.mod`, `go.sum` and the lint and
formatter configuration are granted in this cycle only for a correction that
requires one. This Spec's design needs none of them, and it changes no
Baseline module, template or profile.

This record grants `implement`, `commit`, `push`, `pull_request`, `merge` and
`release`. The Spec is a `type: feat`, so the Release Ceiling is `minor` and
the release may change only the three version-only files in `release_paths`:
the changelog, the npm package manifest and the application version. No Task
performs it; the operator cuts it after the merge.

The governed set was measured on 2026-11-14 on the branch `docs/plan-0392`,
based on `251ed81e`, with `speccheck.GovernedPath`; the rebase on `f17949de`
added only the plan of Spec 0391 and moved none of these paths.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` showed only this Spec's
  planning files and ADR-0401, ADR-0402 and ADR-0405 afterwards, and the probe
  file never existed in the tree.
- It covered every file the TechSpec's Build Order names, the Go files it
  expects the Tasks to create, every planning file of this Spec and the three
  ADRs.
- It was run again on 2026-11-14 at `f17949de`, when the Task Graph was
  written, over every path the five Task files declare under `interface:`,
  `creates:` and `instruction:`, the Task files, `_tasks.md` and the planning
  files: sixty-nine paths with the two controls below. The governed set did
  not move.
- The thirteen paths above are governed: the `SKILL.md` pairs of the Roundfix
  skill and of the four authoring skills, and the Roundfix skill's canonical
  `baseline`, `modes` and `setup` references. Task 03 changes the `setup`
  reference and the Roundfix `SKILL.md` pair; task 04 changes all thirteen.
  The probe also named `.agents/skills/domain-modeling/SKILL.md`, which a
  Task only reads, and the two controls `internal/cli/cli_test.go` and
  `docs/references/coverage-record.json`, which no Task declares or changes;
  none of the three is granted.

## What is not governed

- Every Go source and test file under `internal/config`, `skills` and
  `internal/cli` that the Tasks declare.
- The mirrors under `skills/roundfix/references/` and
  `skills/testdata/owned-skill-versions.json`.
- `docs/user-guide/**`, `CONTEXT.md`, `docs/references/behavior-surfaces.json`,
  `docs/references/skill-coverage.json`, ADR-0401, ADR-0402, ADR-0405 and the
  amendment notice in ADR-0386.
- Every planning file of this Spec and its references.

## Sanctioned regeneration

The repository-owned commands write the skill mirrors and the generated
records. This declaration records the regeneration that follows the approved
edits and adds no source path.

```yaml
command: make skills-sync
```

```yaml
command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
outputs:
  - skills/testdata/owned-skill-versions.json
  - .agents/skills/roundfix/SKILL.md
  - skills/roundfix/SKILL.md
  - .agents/skills/write-idea/SKILL.md
  - skills/write-idea/SKILL.md
  - .agents/skills/write-prd/SKILL.md
  - skills/write-prd/SKILL.md
  - .agents/skills/write-tasks/SKILL.md
  - skills/write-tasks/SKILL.md
  - .agents/skills/write-techspec/SKILL.md
  - skills/write-techspec/SKILL.md
```

```yaml
command: go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage
outputs:
  - docs/references/behavior-surfaces.json
```

No mirror, version field, digest or record is hand-edited.

## Limits

- This repository's `.roundfixrc.yml` is not changed: the Roundfix repository
  does not declare the section, because its `.agents/skills` are the source of
  the shipped bundle and the installed owner binary reads the default
  branch's Project Config.
- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`, the lint,
  formatter or test-runner configuration, `skills-lock.json`, any Vendored
  Skill under `.agents/skills`, any Baseline module, template, profile or
  Setup Snapshot, any guide under `docs/agents/`,
  `docs/references/coverage-record.json`, any Archive Record, any owned skill
  other than `roundfix`, `write-idea`, `write-prd`, `write-techspec` and
  `write-tasks`, any Roundfix reference other than `baseline`, `modes` and
  `setup`, or the `### QA settlement` section of any skill.
- No shipped `SKILL.md`, canonical or mirrored, gains a `model` or `effort`
  key; the header exists only in a repository that declares it.
- No command, Task or QA row writes to a user-level skill directory of the
  host, runs `roundfix skills install` with a user-level target or runs
  `roundfix skills remove --apply` outside a temporary home.
- No test, Verification command or QA row:
  - starts Claude Code, another Agent harness, an Agent session or an ACP
    adapter;
  - reaches the network, calls ai-memory or runs the real `gh`;
  - reads or writes the real `~/.roundfix` or the real home directory;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository;
  - records an environment value or writes a real-looking secret literal;
  - signals a process it did not start, or ends a process by name.
- No Task sends an ai-memory message or a Spec Event Announcement; the
  operator does.
