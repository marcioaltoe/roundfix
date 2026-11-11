---
status: approved
granted: 2026-11-10
action: add the internal/corebudget package that computes the Verification Core Budget; export it as GOMAXPROCS and ROUNDFIX_VERIFICATION_CORES to every Verification and Agent Session of an Implement Run, to the Delivery Queue's repository gate and to the Linux Lane container; print it in a verification-cores line of doctor; and document it in the configuration, doctor and deliver guides, the Roundfix skill's implement, deliver and setup references, the glossary and ADR-0385
consuming: 0375-a-verification-core-budget-keeps-concurrent-verifications-light
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/deliver.md
  - .agents/skills/roundfix/references/implement.md
  - .agents/skills/roundfix/references/setup.md
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

# Approved authority for Spec 0375

On 2026-11-10 the vortex project sent this repository the maintainer's
request, ai-memory message 01a1c7d9-d4ac-7553-9f7e-e1cf756c6ac3 ("para
roundfix resolver isso no setup dele"), which Spec 0374 adopted and this Spec
links. The maintainer then answered the operator's cycle BB batch 2 questions
with AskUserQuestion (via `/roundfix --follow-up`, 2026-11-10 ~20:55–21:00
BRT). The answers this Spec rests on are, verbatim:

- Question "Pedido do Vortex (Spec 0374): verificações de projetos diferentes
  rodando juntas nesta máquina derrubam testes sensíveis a tempo. Hoje o
  Roundfix só limita Runs inteiras (`runs.max_active`), e você não quer isso.
  Que mecanismo o Roundfix ganha?" → "Os dois" — option text: "Capacidade da
  máquina para serializar a verificação e, além dela, o orçamento exportado
  para as verificações que ainda rodarem juntas. Spec maior, provavelmente
  dividida em duas." The two options it combines: "Capacidade da máquina" —
  "Uma chave na config de usuário limita quantas verificações rodam ao mesmo
  tempo na máquina, somando todos os repositórios. Os agentes seguem em
  paralelo; só a etapa de verificação espera a vez, e `deliver status` e
  `runs list` mostram quem espera. Vale para qualquer stack (Go, Vitest,
  Bun)."; "Orçamento de núcleos" — "O Roundfix divide os núcleos entre as Runs
  ativas e exporta o número (GOMAXPROCS e uma variável própria) para os
  comandos de verificação. O Go obedece sozinho; Vitest e Bun só se cada
  projeto ler a variável. As verificações continuam simultâneas, só mais
  leves."
- Question "Até a 0374 existir, a verificação da 0373 do Roundfix vai rodar
  esta noite em paralelo com a fila do Vortex. O que faço com a fila do
  Roundfix?" → custom answer "GOMAXPROCS=4 até o final das implementações e
  depois usar 6", confirmed by "Isso (Recommended)" to "GOMAXPROCS=4 em toda
  Run do Roundfix até a última Spec deste ciclo (0373 a 0379) ser entregue, e
  depois 6 como valor fixo da fila do Roundfix nesta máquina, até o orçamento
  de núcleos da 0375 assumir."
- Question "Com que autorizações eu toco sem supervisão as seis Specs novas:
  0374 (capacidade da máquina), 0375 (orçamento de núcleos), 0376 (skills de
  usuário velhas), 0377 (portão do CI), 0378 (digest do bump de skills) e
  0379 (marcador \"não se aplica\")?" → "Mesmas do BB, teto minor
  (Recommended)" — option text: "Escrever, entregar pela fila, merge e
  release de cada Spec assim que pronta (patch nas correções, minor nas
  features; 0374 e 0375 saem juntas), teto de 2 Tasks corretivas, CI na 0377
  e Makefile na 0378 só nos arquivos que elas exigirem, eventos ao Jarvis e
  resposta ao Vortex quando a 0374 e a 0375 saírem."

"Mesmas do BB" keeps the cycle BB rules of 2026-11-10 ("Mesmas do BA, teto
minor (Recommended)"), and through them the cycle BA, AZ, AY, AX, AW, AV, AU,
AT and AS answers, which release to the Runs the owned skills and changes of
skills that reach the adopters, the cycle AJ grant of 2026-10-19 "Mexer em
Governed Paths declarados" (option text: "'Concedo' para os caminhos que cada
Spec declarar (módulos Baseline, registros, skills, testes governados; fora:
Makefile, CI, go.mod, lint, que pedem autorização específica).") and the
2026-09-30 grant of every skill ("considere autorizado a ajustar todas as
skills se necessário").

The Makefile, `.github/workflows`, `go.mod`, `go.sum` and the lint and
formatter configuration are grantable in this cycle only to Specs 0377 and
0378. This Spec needs none of them, and it changes no Baseline module.

This record grants `implement`, `commit`, `push`, `pull_request`, `merge` and
`release`. The Spec is a `type: feat`, so the Release Ceiling is `minor` and
the release may change only the three version-only files in `release_paths`:
the changelog, the npm package manifest and the application version. No Task
performs it; the operator cuts it after the merge, together with Spec 0374
("0374 e 0375 saem juntas"), and its changelog section names
`GOMAXPROCS` and `ROUNDFIX_VERIFICATION_CORES` and the half-of-the-cores
default.

The governed set was measured on 2026-11-10 on the branch `docs/plan-0375`,
based on `b233761e`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` showed only this Spec's
  planning files and ADR-0385 afterwards, and the probe file never existed in
  the tree.
- It covered every file the Tasks declare under `creates:`, `interface:` or
  `instruction:`, every file the scratch prototype of tasks 01 to 03 changed
  or created, every planning file of this Spec and ADR-0385.
- The five paths above are governed: the Roundfix skill's `SKILL.md` pair and
  its canonical `deliver`, `implement` and `setup` references. task_02
  changes the `implement` reference and the `SKILL.md` pair's version fields;
  task_03 the `deliver` and `setup` references and the version fields again.

## What is not governed

- Every Go source and test file tasks 01 to 03 create or change under
  `internal/corebudget`, `internal/daemon`, `internal/agent` and
  `internal/cli`.
- The mirrors under `skills/roundfix/references/` and
  `skills/testdata/owned-skill-versions.json`.
- `docs/user-guide/**`, `CONTEXT.md`, `docs/references/behavior-surfaces.json`
  and ADR-0385.
- Every planning file of this Spec.

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
  `skills-lock.json`, any Vendored Skill under `.agents/skills`, any Baseline
  module, template, profile or Setup Snapshot, any guide under `docs/agents/`,
  `docs/references/coverage-record.json`, any Archive Record, any skill other
  than `roundfix`, any Roundfix reference other than `deliver`, `implement`
  and `setup`, or the `### QA settlement` section of any skill.
- No User Config or Project Config key is added, and nothing is stored.
- No test, Verification command or QA row:
  - reaches the network, calls ai-memory or runs the real `gh` or Docker;
  - reads or writes the real `~/.roundfix` or the real home directory;
  - starts an Agent Session or an ACP adapter;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository;
  - records an environment value other than a budget number a test prints,
    or writes a real-looking secret literal;
  - signals a process it did not start, or ends a process by name.
- No Task sends an ai-memory message or a Spec Event Announcement; the
  operator does.
