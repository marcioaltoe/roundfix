---
status: approved
granted: 2026-10-24
action: add the internal/sharedcache trim engine, the worktree.shared_caches configuration key and a GC Command section that trims each declared Shared Build Cache above its size by removing only idle Go build cache entries, declare this repository's Makefile cache in .roundfixrc.yml, and document the key and section in the GC and configuration guides, the Roundfix skill's storage reference and the glossary
consuming: 0314-gc-trims-the-shared-build-cache
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/storage.md
  - .roundfixrc.yml
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0314

On 2026-10-23, at about 23:55 BRT, the maintainer asked: "Vejo que
frequentemente temos problemas com os testes implementados para roundfix.
Precisamos trabalhar em testes mais eficientes e rápidos. Isso deve ser uma
spec de performance para os próximos passos de evolução de roundfix,
saneamento dos testes e foco de performance para garantir a qualidade."

On 2026-10-24, at about 00:05 BRT, the maintainer answered cycle AP batch 4
through AskUserQuestion:

- order: "Logo após a 0306 (Recomendado)" (option text: "Entra antes da
  0307, da 0309 e da 0310: cada Spec que vem depois roda a mesma suíte, e os
  testes instáveis continuariam custando Runs. Eu escrevo o plano enquanto a
  0306 termina.");
- scope: "Auditoria completa" (option text: "Inventário de toda a suíte:
  testes lentos, redundantes, instáveis e com insumos ocultos, com o plano de
  2 a 3 Specs ordenadas."); this Spec bounds the shared build cache the
  audit measured;
- authorization: "Sim, só Makefile e workflows (Recomendado)" (option text:
  "Autoriza o Makefile, .github/workflows e flags de teste, limitado ao que a
  Spec declarar (ex.: rodar o docscontract uma vez por job, limitar o
  GOCACHE, paralelismo). go.mod e go.sum continuam proibidos."). This Spec
  declares no Makefile, workflow or test-flag change.

On 2026-10-24, at about 01:50 BRT, the maintainer answered the authoring
questions of this batch through AskUserQuestion; the answer this Spec rests
on is:

- build cache: "roundfix gc recolhe" (option text: "Vira uma mudança do
  produto: o comando gc passa a recolher esse cache."); the Makefile trim
  first planned in Spec 0312 was dropped, and the reclamation is this change
  to the GC Command. The same answers order the batch 0311, 0312, 0313,
  0314, all before 0307.

The cycle AP rules that answer carries were approved on 2026-10-22 (questions
of 15:25 BRT, answered by about 18:30 BRT): the Governed Paths each Spec
declares, with the cycle AJ grant of 2026-10-19, at about 16:20 BRT, "Mexer
em Governed Paths declarados" (option text: "'Concedo' para os caminhos que
cada Spec declarar (módulos Baseline, registros, skills, testes governados;
fora: Makefile, CI, go.mod, lint, que pedem autorização específica).") and
"Merge squash + release + push".

This record grants `implement`, `commit`, `push`, `pull_request` and `merge`.
It does not grant `release`: the operator releases.

The governed set was measured on the authoring branch `docs/plan-0311`, based
on `56dd1a76`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` was identical before
  and after, and the probe file never existed in the tree.
- It covered every file the Tasks declare, every planning file of this Spec
  and ADR-0327.
- `.roundfixrc.yml`, `.agents/skills/roundfix/references/storage.md` and
  both `SKILL.md` files of the Roundfix skill are governed.

## What is not governed

- `internal/sharedcache/*`, `internal/config/config.go`, the new
  `internal/config/worktree_shared_caches_test.go`, `internal/cli/gc.go`,
  the new `internal/cli/gc_shared_cache.go` and its test.
- `docs/user-guide/commands/gc.md`, `docs/user-guide/configuration.md`,
  `skills/roundfix/references/storage.md`,
  `skills/testdata/owned-skill-versions.json`, `CONTEXT.md` and
  `docs/references/behavior-surfaces.json`.
- `docs/adr/0327-gc-trims-a-declared-shared-build-cache.md`, every file the
  Tasks read as an instruction and every planning file of this Spec.

## Sanctioned regeneration

- `make skills-sync`, which copies the canonical Roundfix skill to
  `skills/roundfix/`.
- `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
  outputs: `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`,
  `skills/testdata/owned-skill-versions.json`.
- `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`,
  outputs: `docs/references/behavior-surfaces.json`.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`, the lint,
  formatter or test-runner configuration, any Baseline module,
  `docs/agents/setup-context.json`, any skill other than the Roundfix
  skill's storage reference and version fields, or the `### QA settlement`
  section of any skill.
- No test, Verification command or QA row:
  - reads or writes `~/.roundfix`;
  - measures or trims the host's real Go build caches;
  - reaches the network;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository;
  - records an environment value.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
