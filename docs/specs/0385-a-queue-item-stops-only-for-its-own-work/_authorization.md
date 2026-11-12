---
status: approved
granted: 2026-11-12
action: give the Delivery Queue workflow the Queue Owner's log, keep the owner record of an owner that stops with an item left, read a review anchor that ends a sentence, copy the reviewed checkout's .ai-memory.toml into the review validator's sealed working directory, and document them in the deliver and review guides, the Roundfix skill's deliver and review references, the glossary and ADR-0395
consuming: 0385-a-queue-item-stops-only-for-its-own-work
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/deliver.md
  - .agents/skills/roundfix/references/review.md
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

# Approved authority for Spec 0385

On 2026-11-12 quaestor reported, in ai-memory message
`01a1d01d-ac6c-71a0-a548-0398a33c1ba6`, that the Delivery Owner panics at the
repository gate while it waits for the Machine Verification Capacity
(`references/2026-11-12-quaestor-the-delivery-owner-panics-at-the-gate.md`),
and pantheon reported, in ai-memory message
`01a1d01a-fa8c-7211-bd88-5023f39563e1`, that the review validator becomes a
stray ai-memory project
(`references/2026-11-12-pantheon-the-review-validator-becomes-a-stray-ai-memory-project.md`).
The operator recorded the owner crash and the `path:line.` anchor as Backlog
Entries, which this Spec adopts. The maintainer then answered the operator's
cycle BD questions with AskUserQuestion (via `/roundfix --follow-up`,
2026-11-12 ~17:05–17:20 BRT). The answers this Spec rests on are, verbatim:

- Question "Três Specs novas saem das mensagens de hoje: A = aceitar o
  skills-lock.json versão 2 (Vortex e Pantheon travados no update), B = o dono
  da fila não morre mais ao esperar a Machine Verification Capacity + âncora
  'path:line.' aceita, C = cache Go do roundfix fora do disco interno do mini
  (Drive). Que autorizações e ordem?" → "A → C → B, mesmas do BC
  (Recommended)" — option text: "Push, merge e release logo após cada merge,
  como no ciclo BC. A primeiro porque trava dois adotantes; C antes de B porque
  o disco do mini está abaixo dos 80 GB." This Spec is B, delivered after 0383
  and 0384.
- Question "Pantheon achou que o validador de revisão roda o Codex em
  /var/folders/…/roundfix-review-validator-<n>, criando projetos soltos no
  ai-memory. Já há rota em ~/.ai-memory.toml mandando /var/folders para o
  Roundfix. O que faço com o achado?" → "Juntar à Spec B" — option text: "Vira
  uma quarta Task da Spec B (fila e revisão)."

"Mesmas do BC" keeps the cycle BC rules of 2026-11-11 ("Mesmas do BB, após E
(Recommended)" — option text: "Push, merge, release, teto minor; entregue
depois de D e E, fechando o ciclo."), and through them the cycle BB, BA, AZ,
AY, AX, AW, AV, AU, AT and AS answers, which release to the Runs the owned
skills and changes of skills that reach the adopters, the cycle AJ grant of
2026-10-19 "Mexer em Governed Paths declarados" (option text: "'Concedo' para
os caminhos que cada Spec declarar (módulos Baseline, registros, skills,
testes governados; fora: Makefile, CI, go.mod, lint, que pedem autorização
específica).") and the 2026-09-30 grant of every skill ("considere autorizado
a ajustar todas as skills se necessário").

This Spec needs no Makefile, `.github/workflows`, `go.mod`, `go.sum`, lint or
formatter configuration change, and no Baseline module, template or profile
change.

This record grants `implement`, `commit`, `push`, `pull_request`, `merge` and
`release`. The Spec is a `type: fix`, so the Release Ceiling is `patch` and
the release may change only the three version-only files in `release_paths`:
the changelog, the npm package manifest and the application version. No Task
performs it; the operator cuts it after the merge.

The governed set was measured on 2026-11-12 on the branch `docs/plan-0385`,
based on `d79b30b4`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` showed only this Spec's
  planning files, its adopted references and ADR-0395 afterwards, and the
  probe file never existed in the tree.
- It covered every file the Tasks declare under `creates:`, `interface:` or
  `instruction:`, every file the scratch prototype of tasks 01 to 04 changed or
  created, every planning file of this Spec, its adopted references, ADR-0395
  and `docs/agents/setup-context.json`.
- The four paths above are governed: the Roundfix skill's `SKILL.md` pair and
  its canonical `deliver` and `review` references, which task_04 changes. The
  probe also named `.agents/skills/domain-modeling/SKILL.md`, which task_01
  only reads, and `docs/agents/setup-context.json`, which no Task changes;
  neither is granted.

## What is not governed

- Every Go source and test file tasks 01 to 03 create or change under
  `internal/cli/` and `internal/agent/`.
- The mirrors under `skills/roundfix/references/` and
  `skills/testdata/owned-skill-versions.json`.
- `docs/user-guide/**`, `CONTEXT.md`, `docs/references/behavior-surfaces.json`
  and ADR-0395.
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
  module, template, profile, Setup Snapshot or plan golden, any guide under
  `docs/agents/`, `docs/agents/setup-context.json`,
  `docs/references/coverage-record.json`, any Archive Record, any skill other
  than `roundfix`, any Roundfix reference other than `deliver` and `review`,
  `internal/baselineacp`, `internal/delivery`, or the `### QA settlement`
  section of any skill.
- No test, Verification command or QA row:
  - reaches the network, calls ai-memory or runs the real `gh`;
  - reads or writes the real `~/.roundfix` or the real home directory;
  - starts an Agent session or an ACP adapter;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository;
  - records an environment value or writes a real-looking secret literal;
  - signals a process it did not start, or ends a process by name;
  - waits on elapsed time instead of the event it asserts.
- No Task sends an ai-memory message or a Spec Event Announcement; the
  operator does.
