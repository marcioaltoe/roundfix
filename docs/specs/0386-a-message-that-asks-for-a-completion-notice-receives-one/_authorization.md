---
status: approved
granted: 2026-11-12
action: reword the sender and receiving message clauses of the Baseline core module's cross-project continuity rule so a message that asks to be told when its work is done receives one completion notice after the delivery, with their derived files and this repository's managed refresh; change the Roundfix skill's follow-up mode and release reference to keep the request in the message's record and send the notice; and record the glossary changes and ADR-0396's amendment notices in ADR-0390, ADR-0383 and ADR-0330
consuming: 0386-a-message-that-asks-for-a-completion-notice-receives-one
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/modes.md
  - .agents/skills/roundfix/references/release.md
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

# Approved authority for Spec 0386

The Spec adopts the maintainer's instruction of 2026-11-12, given in chat at
about 18:40 BRT in reply to the operator's cycle BD report; requested by the
maintainer. Both texts are copied verbatim under `references/`.

- The operator's question: "O Pantheon e o Drive pediram aviso quando a
  correção sair, e o Drive diz que o pedido do retorno veio de você. A regra
  da 0380 proíbe aviso de conclusão a mensagem que não declarou trabalho
  bloqueado, e nenhuma das duas declarou. Mando os dois avisos quando as Specs
  saírem, ou mantenho a regra?"
- The maintainer's instruction: "Adicione exceção a regra para quando houver
  uma solicitação de aviso."

The operator takes the instruction as authorization on the cycle BD terms,
whose answer of 2026-11-12 (AskUserQuestion via `/roundfix --follow-up`,
~17:05–17:20 BRT) was "A → C → B, mesmas do BC (Recommended)" — option text:
"Push, merge e release logo após cada merge, como no ciclo BC. A primeiro
porque trava dois adotantes; C antes de B porque o disco do mini está abaixo
dos 80 GB." Those terms keep the cycle BC, BB, BA, AZ, AY, AX, AW, AV, AU, AT
and AS answers, which release to the Runs the Baseline source
(`internal/baseline/assets/**`) with its fixtures and goldens, the guides
under `docs/agents/`, the owned skills, and changes of Baseline and skills
that reach the adopters, the cycle AJ grant of 2026-10-19 "Mexer em Governed
Paths declarados" (option text: "'Concedo' para os caminhos que cada Spec
declarar (módulos Baseline, registros, skills, testes governados; fora:
Makefile, CI, go.mod, lint, que pedem autorização específica).") and the
2026-09-30 grant of every skill ("considere autorizado a ajustar todas as
skills se necessário").

The Makefile, `.github/workflows`, `go.mod`, `go.sum` and the lint and
formatter configuration are not expected for this Spec, and it needs none of
them.

This record grants `implement`, `commit`, `push`, `pull_request`, `merge` and
`release`. The Spec is a `type: feat`, so the Release Ceiling is `minor` and
the release may change only the three version-only files in `release_paths`:
the changelog, the npm package manifest and the application version. No Task
performs it; the operator cuts it right after the merge.

The governed set was measured on 2026-11-12 on the branch `docs/plan-0386`,
based on `73421db2`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` showed only this Spec's
  planning files and ADR-0396 afterwards, and the probe file never existed in
  the tree.
- It covered every file the Tasks declare under `interface:` or
  `instruction:`, every file the prototype of tasks 01 and 02 changed or
  created (the derived Baseline files included), every planning file of this
  Spec, its adopted reference and ADR-0396.
- The nine paths above are governed: the core module, the Standard TypeScript
  profile and its formatter golden of the agent instructions guide, this
  repository's agent instructions guide and its setup manifest (task_01), and
  the Roundfix skill's `SKILL.md` pair and its canonical `modes` and `release`
  references (task_02). The probe also reported the `domain-modeling` skill,
  which task_02 only reads; it is not granted.

## What is not governed

- Every Go test file tasks 01 and 02 change under `internal/baseline` and
  `skills`.
- The mirrors under `skills/roundfix/references/`,
  `skills/testdata/owned-skill-versions.json`,
  `internal/baseline/module-versions.json` and the catalog snapshots and plan
  goldens under `internal/baseline/testdata/`.
- `CONTEXT.md`, ADR-0396 and the amendment notices in ADR-0390, ADR-0383 and
  ADR-0330.
- Every planning file of this Spec and its adopted reference.

## Sanctioned regeneration

The repository-owned commands write the module version record, the derived
Baseline files, this repository's managed guide and setup manifest, the skill
mirrors and the owned-skill version record. This declaration records the
regeneration that follows the approved edits and adds no source path.

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

No mirror, version field, digest or record is hand-edited.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`, the lint,
  formatter or test-runner configuration, `.roundfixrc.yml`,
  `skills-lock.json`, any Vendored Skill under `.agents/skills` (the
  `ai-memory-messaging` skill included), any Baseline module other than
  `core`, any Source Baseline, Baseline template or Setup Snapshot, any guide
  under `docs/agents/` other than through the managed refresh, any page under
  `docs/user-guide/`, `docs/references/coverage-record.json`, any Archive
  Record, any file under
  `internal/cli/`, any skill other than `roundfix`, any Roundfix reference
  other than `modes` and `release`, or the `### QA settlement` section of any
  skill.
- No clause is added, removed or renamed, and no code reads a message body or
  a record for a request.
- No test, Verification command or QA row:
  - reaches the network, calls ai-memory or runs the real `gh`;
  - reads or writes the real `~/.roundfix` or the real home directory;
  - starts an Agent session or an ACP adapter;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository;
  - records an environment value or writes a real-looking secret literal;
  - signals a process it did not start, or ends a process by name.
- No Task sends an ai-memory message, a completion notice or a Spec Event
  Announcement; the operator does.
