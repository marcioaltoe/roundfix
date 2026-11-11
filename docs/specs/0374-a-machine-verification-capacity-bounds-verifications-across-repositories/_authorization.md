---
status: approved
granted: 2026-11-10
action: add the runs.max_verifying User Config key and the Machine Verification Capacity entry store in the Roundfix Home; make every Daemon Verification and the Delivery Queue's repository gate take a machine unit, with the Run Budget paused during the wait; print the capacity in runs list, deliver status and doctor; and document it in the configuration and command guides, the Roundfix skill's implement, runs, deliver and setup references, the glossary, ADR-0384 and the amendment notice in ADR-0056
consuming: 0374-a-machine-verification-capacity-bounds-verifications-across-repositories
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/deliver.md
  - .agents/skills/roundfix/references/implement.md
  - .agents/skills/roundfix/references/runs.md
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

# Approved authority for Spec 0374

On 2026-11-10 the vortex project sent this repository the maintainer's
request, ai-memory message 01a1c7d9-d4ac-7553-9f7e-e1cf756c6ac3 ("para
roundfix resolver isso no setup dele"), which this Spec adopts
(`references/2026-11-10-vortex-verifications-of-different-repositories-collide-on-one-machine.md`).
The maintainer then answered the operator's cycle BB batch 2 questions with
AskUserQuestion (via `/roundfix --follow-up`, 2026-11-10 ~20:55–21:00 BRT).
The answers this Spec rests on are, verbatim:

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
- Question "A Capacidade da máquina (Spec 0374) é uma chave da config de
  usuário. Que valor ela tem por padrão para quem atualizar o Roundfix sem
  configurá-la?" → "1 por padrão (Recommended)" — option text: "Toda máquina
  passa a verificar um repositório por vez, e os agentes seguem em paralelo.
  O `setup` grava o valor na config de usuário. Muda o comportamento de todos
  os adotantes logo que atualizarem."
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
0378. This Spec needs none of them.

This record grants `implement`, `commit`, `push`, `pull_request`, `merge` and
`release`. The Spec is a `type: feat`, so the Release Ceiling is `minor` and
the release may change only the three version-only files in `release_paths`:
the changelog, the npm package manifest and the application version. No Task
performs it; the operator cuts it after the merge, together with the core
budget Spec ("0374 e 0375 saem juntas"), and its changelog section says that
every adopter's Verifications become one at a time per machine after the
update.

The governed set was measured on 2026-11-10 on the branch `docs/plan-0374`,
based on `60ac8142`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` showed only this Spec's
  planning files and ADR-0384 afterwards, and the probe file never existed in
  the tree.
- It covered every file the Tasks declare under `creates:`, `interface:` or
  `instruction:`, every file the scratch prototype of tasks 01 to 03 changed
  or created, every planning file of this Spec, its adopted reference and
  ADR-0384.
- The six paths above are governed: the Roundfix skill's `SKILL.md` pair and
  its canonical `deliver`, `implement`, `runs` and `setup` references:
  task_02 changes the `implement` reference and the `SKILL.md` pair's
  version fields, task_03 the `deliver`, `runs` and `setup` references and
  the version fields again.

## What is not governed

- Every Go source and test file tasks 01 to 03 create or change under
  `internal/store`, `internal/config`, `internal/daemon` and `internal/cli`.
- The mirrors under `skills/roundfix/references/` and
  `skills/testdata/owned-skill-versions.json`.
- `docs/user-guide/**`, `CONTEXT.md`, `docs/references/behavior-surfaces.json`,
  ADR-0384 and the amendment notice in ADR-0056.
- Every planning file of this Spec and its adopted reference.

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
  than `roundfix`, any Roundfix reference other than `deliver`, `implement`,
  `runs` and `setup`, or the `### QA settlement` section of any skill.
- No Run Database table, column or schema version is added, and no entry is
  removed by its age.
- No environment variable is exported to a Verification command; that is the
  core budget Spec's.
- No test, Verification command or QA row:
  - reaches the network, calls ai-memory or runs the real `gh` or Docker;
  - reads or writes the real `~/.roundfix` or the real home directory;
  - starts an Agent session or an ACP adapter;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository;
  - records an environment value or writes a real-looking secret literal;
  - signals a process it did not start, or ends a process by name.
- No Task sends an ai-memory message or a Spec Event Announcement; the
  operator does.
