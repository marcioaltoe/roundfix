---
status: approved
granted: 2026-11-12
action: read skills-lock.json versions 1 and 2 in every lock reader, write and compare each computedHash in the hash algorithm of the lock's own version without changing that version, refuse any other version by name with roundfix upgrade as the next action in Doctor and the Baseline update, and document it in the doctor, skills and baseline guides, the Roundfix skill's setup, baseline and runtime references, the glossary and ADR-0393
consuming: 0383-roundfix-reads-a-version-2-skills-lock-and-keeps-its-version
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/baseline.md
  - .agents/skills/roundfix/references/runtime.md
  - .agents/skills/roundfix/references/setup.md
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
  - release
release_ceiling: patch
release_paths:
  - CHANGELOG.md
  - dist/npm/roundfix/package.json
  - internal/app/version.go
---

# Approved authority for Spec 0383

On 2026-11-12 vortex and pantheon reported that Roundfix 0.141.1 refuses the
`skills-lock.json` version 2 that the skills CLI 1.8 writes
(`references/2026-11-12-vortex-skills-lock-version-2-refused.md`, reported by
vortex; `references/2026-11-12-pantheon-skills-lock-version-2-refused.md`,
reported by pantheon). The maintainer answered the operator's cycle BD
questions with AskUserQuestion (via `/roundfix --follow-up`, 2026-11-12
~17:05–17:20 BRT). The answers this Spec rests on are, verbatim:

- Question "Três Specs novas saem das mensagens de hoje: A = aceitar o
  skills-lock.json versão 2 (Vortex e Pantheon travados no update), B = o
  dono da fila não morre mais ao esperar a Machine Verification Capacity +
  âncora 'path:line.' aceita, C = cache Go do roundfix fora do disco interno
  do mini (Drive). Que autorizações e ordem?" → "A → C → B, mesmas do BC
  (Recommended)" — option text: "Push, merge e release logo após cada merge,
  como no ciclo BC. A primeiro porque trava dois adotantes; C antes de B
  porque o disco do mini está abaixo dos 80 GB."
- Context only, nothing to build: question "A skill roundfix rebaixada no
  Vortex vem do catálogo marcioaltoe/skills, que sincroniza com o Roundfix
  uma vez por dia (hoje 04:52 BRT, antes da v0.141.1 sair às 05:20). Disparo
  o 'Sync upstream skills' do catálogo (workflow_dispatch, auto-merge no
  marcioaltoe/skills)?" → free answer "Já pedi para o projeto skills disparar
  a sincronização".

This Spec is A, delivered first. "Mesmas do BC" keeps the cycle BC rules of
2026-11-11 ("Após a 0379 (Recommended)" — "Mesmas autorizações do ciclo BB
(push, merge, release, teto minor)…", and "Mesmas do BB, após E
(Recommended)"), and through them the cycle BB, BA, AZ, AY, AX, AW, AV, AU,
AT and AS answers, which release to the Runs the owned skills and changes of
skills that reach the adopters, the cycle AJ grant of 2026-10-19 "Mexer em
Governed Paths declarados" (option text: "'Concedo' para os caminhos que cada
Spec declarar (módulos Baseline, registros, skills, testes governados; fora:
Makefile, CI, go.mod, lint, que pedem autorização específica).") and the
2026-09-30 grant of every skill ("considere autorizado a ajustar todas as
skills se necessário").

The Makefile, `.github/workflows`, `go.mod`, `go.sum` and the lint and
formatter configuration are grantable in this cycle only where a design
needs them. This Spec needs none of them.

This record grants `implement`, `commit`, `push`, `pull_request`, `merge` and
`release`. The Spec is a `type: fix`, so the Release Ceiling is `patch` and
the release may change only the three version-only files in `release_paths`:
the changelog, the npm package manifest and the application version. No Task
performs it; the operator cuts it after the merge.

The governed set was measured on 2026-11-12 on the branch `docs/plan-0383`,
based on `d79b30b4`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` showed only this Spec's
  folder and ADR-0393 afterwards, and the probe file never existed in the
  tree.
- It covered every file the Tasks declare under `creates:`, `interface:` or
  `instruction:`, every file the scratch prototype of tasks 01 to 03 changed
  or created, every planning file of this Spec, its adopted references and
  ADR-0393.
- Six paths are governed. The five above are the Roundfix skill's `SKILL.md`
  pair and its canonical `setup`, `baseline` and `runtime` references, which
  task_03 changes. The sixth, `.agents/skills/domain-modeling/SKILL.md`, is
  only read by task_03 (`instruction:`) and is not granted.

## What is not governed

- Every Go source and test file tasks 01 and 02 create or change under
  `internal/skillhash/`, `skills/`, `internal/baseline/` and `internal/cli/`.
- The mirrors under `skills/roundfix/references/` and
  `skills/testdata/owned-skill-versions.json`.
- `docs/user-guide/**`, `CONTEXT.md`, `docs/references/behavior-surfaces.json`
  and ADR-0393.
- Every planning file of this Spec and its adopted references.

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
  module, template, profile, asset or Setup Snapshot (including
  `internal/baseline/assets/lock-hash-compatibility-v1.json`), any guide under
  `docs/agents/`, `docs/references/coverage-record.json`, any Archive Record,
  any skill other than `roundfix`, any Roundfix reference other than `setup`,
  `baseline` and `runtime`, any ADR other than ADR-0393, or the
  `### QA settlement` section of any skill.
- No test, Verification command or QA row:
  - reaches the network, calls ai-memory or runs the real `gh`;
  - runs the skills CLI (`bunx skills`, `npx skills`) or `npm`;
  - reads or writes the real `~/.roundfix`;
  - starts an Agent session or an ACP adapter;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository;
  - records an environment value or writes a real-looking secret literal;
  - signals a process it did not start, or ends a process by name.
- No Task sends an ai-memory message or a Spec Event Announcement; the
  operator does.
