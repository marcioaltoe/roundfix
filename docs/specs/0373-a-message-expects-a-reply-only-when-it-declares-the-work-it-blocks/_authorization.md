---
status: approved
granted: 2026-11-10
action: add the optional blocked_work key to the Follow-up Command's messages file with the declared basis that ranks such a message first without the judge, its tests, its command guide and the Roundfix follow-up reference; reword the three message clauses of the Baseline core module's cross-project continuity rule with their derived files and this repository's managed refresh; change the Roundfix skill's follow-up mode in the modes reference; and record the glossary entries, ADR-0383 and its amendment notices in ADR-0319 and ADR-0330
consuming: 0373-a-message-expects-a-reply-only-when-it-declares-the-work-it-blocks
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/follow-up.md
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
  - release
release_ceiling: minor
release_paths:
  - CHANGELOG.md
  - dist/npm/roundfix/package.json
  - internal/app/version.go
---

# Approved authority for Spec 0373

On 2026-11-10 the maintainer dictated, in the Roundfix Supervisor session
(`/roundfix --follow-up`, ~15:15 BRT), the request this Spec adopts
(`references/2026-11-10-a-message-expects-a-reply-only-when-the-sender-is-blocked.md`).
The request, verbatim (voice dictation; the first words are transcription
noise and are kept as received):

> Ana. Adicione e. Uma ideia que eu estou percebendo ao baseline: uma
> percepção que eu estou tendo ao baseline. Sempre que os projetos enviam algo
> para outro projeto usando AI memory, ele fica aguardando uma resposta. Essa
> resposta só é necessária se ele precisar da resposta para alguma coisa. Nova
> ação: se ele só está passando para ele algo que deve ser feito pelo outro
> projeto, isso não demanda uma resposta; isso está fazendo com que a gente
> tenha um fluxo de vai e volta e vigiar se foi feito ou não. Muito grande,
> sem necessidade. Então, a proposta seria uma regra para que só seja
> necessária uma resposta a uma mensagem enviada para outro projeto, se o
> projeto que enviou depender daquela resposta para executar ou continuar
> executando algo. E ele deve reforçar na mensagem que ele precisa da resposta
> para continuar o trabalho, e que se bloquear uma implementação, uma spec, ou
> o trabalho, deve ser tratado como urgência.

The maintainer then answered the operator's cycle BB questions with
AskUserQuestion (via `/roundfix --follow-up`, 2026-11-10 ~15:35 BRT). The
answers this Spec rests on are, verbatim:

- Question "Spec 0373 (sua ideia): a regra de resposta entre projetos fica só
  no texto da Baseline, ou o `roundfix follow-up` também passa a ordenar por
  ela?" → "Texto + ferramenta (Recommended)" — option text: "As três
  cláusulas de mensagem da Baseline e a skill mudam, e o `roundfix follow-up`
  põe no topo a mensagem que declara trabalho bloqueado, sem depender do juiz
  Jev. Uma Task a mais; o release provavelmente sai como minor (v0.138.0)."
- Question "Spec 0373: quando chega uma mensagem que declara trabalho
  bloqueado, o que \"tratar como urgência\" obriga o projeto que recebe a
  fazer?" → "Próximo da fila (Recommended)" — option text: "Passa à frente de
  tudo que está pendente e é atendida logo depois do que já está rodando.
  Nenhuma Run ou entrega em curso é interrompida."
- Question "Com que autorizações eu toco sem supervisão as três Specs deste
  ciclo: 0372 (review do Conexus), 0371 (relatório de QA, que você já decidiu
  às 06:20 pelo Jarvis) e 0373 (regra de resposta)?" → "Mesmas do BA, teto
  minor (Recommended)" — option text: "Escrever, entregar pela fila, merge e
  release de cada Spec assim que pronta (patch nas duas correções; a regra
  pode sair como minor), teto de 2 Tasks corretivas, eventos ao Jarvis, sem
  notas às caixas, e Makefile/CI/go.mod/lint só onde a Spec exigir."

"Mesmas do BA" keeps the cycle BA rules of 2026-11-09, and through them the
cycle AZ, AY, AX, AW, AV, AU, AT and AS answers, which release to the Runs the
Baseline source (`internal/baseline/assets/**`) with its fixtures and goldens,
the guides under `docs/agents/`, the owned skills, and changes of Baseline and
skills that reach the adopters, the cycle AJ grant of 2026-10-19 "Mexer em
Governed Paths declarados" (option text: "'Concedo' para os caminhos que cada
Spec declarar (módulos Baseline, registros, skills, testes governados; fora:
Makefile, CI, go.mod, lint, que pedem autorização específica).") and the
2026-09-30 grant of every skill ("considere autorizado a ajustar todas as
skills se necessário").

The Makefile, `.github/workflows`, `go.mod`, `go.sum` and the lint and
formatter configuration are grantable to this Spec only where its design
needs them. It needs none of them.

This record grants `implement`, `commit`, `push`, `pull_request`, `merge` and
`release`. The Spec is a `type: feat`, so the Release Ceiling is `minor` and
the release may change only the three version-only files in `release_paths`:
the changelog, the npm package manifest and the application version. No Task
performs it; the operator cuts it right after the merge.

The governed set was measured on 2026-11-10 on the branch `docs/plan-0373`,
based on `9fda0b65`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` showed only this Spec's
  planning files afterwards, and the probe file never existed in the tree.
- It covered every file the Tasks declare under `creates:`, `interface:` or
  `instruction:`, every file the scratch prototype of tasks 01 to 03 changed
  or created (the derived Baseline files included), every planning file of
  this Spec, its adopted reference and ADR-0383.
- The nine paths above are governed: the Roundfix skill's `SKILL.md` pair and
  its canonical `follow-up` and `modes` references (tasks 01 and 03), and the
  core module, the Standard TypeScript profile and its formatter golden of
  the agent instructions guide, this repository's agent instructions guide and
  its setup manifest (task_02). The probe also reported the
  `domain-modeling` skill, which task_03 only reads; it is not granted.

## What is not governed

- Every Go source and test file tasks 01 to 03 create or change under
  `internal/cli`, `internal/baseline` and `skills`.
- The mirrors under `skills/roundfix/references/`,
  `skills/testdata/owned-skill-versions.json`,
  `internal/baseline/module-versions.json` and the catalog snapshots and plan
  goldens under `internal/baseline/testdata/`.
- `docs/user-guide/**`, `CONTEXT.md`, `docs/references/behavior-surfaces.json`,
  ADR-0383 and the amendment notices in ADR-0319 and ADR-0330.
- Every planning file of this Spec and its adopted reference.

## Sanctioned regeneration

The repository-owned commands write the module version record, the derived
Baseline files, this repository's managed guide and setup manifest, the skill
mirrors and the generated records. This declaration records the regeneration
that follows the approved edits and adds no source path.

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
command: go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text
outputs:
  - docs/agents/agent-instructions.md
  - docs/agents/setup-context.json
```

```yaml
command: make skills-sync
```

```yaml
command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
outputs:
  - skills/testdata/owned-skill-versions.json
  - .agents/skills/roundfix/SKILL.md
  - skills/roundfix/SKILL.md
```

```yaml
command: go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage
outputs:
  - docs/references/behavior-surfaces.json
```

No mirror, version field, digest or record is hand-edited.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`, the lint,
  formatter or test-runner configuration, `.roundfixrc.yml`,
  `skills-lock.json`, any Vendored Skill under `.agents/skills` (the
  `ai-memory-messaging` skill included), any Baseline module other than
  `core`, any Source Baseline, Baseline template or Setup Snapshot, any guide
  under `docs/agents/` other than through the managed refresh,
  `docs/references/coverage-record.json`, any Archive Record, any skill other
  than `roundfix`, any Roundfix reference other than `follow-up` and `modes`,
  `internal/judge`, or the `### QA settlement` section of any skill.
- No clause is added, removed or renamed, and no code reads a message body
  for a marker phrase.
- No test, Verification command or QA row:
  - reaches the network, calls ai-memory or runs the real `gh`;
  - reads or writes the real `~/.roundfix` or the real home directory;
  - starts an Agent session or an ACP adapter;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository;
  - records an environment value or writes a real-looking secret literal;
  - signals a process it did not start, or ends a process by name.
- No Task sends an ai-memory message or a Spec Event Announcement; the
  operator does.
