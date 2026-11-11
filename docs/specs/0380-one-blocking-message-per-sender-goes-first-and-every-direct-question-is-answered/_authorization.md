---
status: approved
granted: 2026-11-11
action: give the Follow-up Command's declared basis to one message entry per sender with the precedence_held_by key on every later declared entry of that sender, its tests, its command guide and the Roundfix follow-up reference; reword the sender and receiving message clauses of the Baseline core module's cross-project continuity rule so every direct question is answered and the declared precedence is bounded per sender, with their derived files and this repository's managed refresh; change the Roundfix skill's follow-up mode in the modes reference; and record the glossary changes and ADR-0390's amendment notices in ADR-0383, ADR-0319 and ADR-0330
consuming: 0380-one-blocking-message-per-sender-goes-first-and-every-direct-question-is-answered
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

# Approved authority for Spec 0380

The Spec adopts two ai-memory messages of 2026-11-11, both findings on the
cross-project message rule of Spec 0373: the Conexus message
`01a1c9ec-5406-7170-bc20-0a69c687b8ae` (a question in a message that declares
no blocked work has no rule that answers it; reported by conexus) and the
Quaestor message `01a1cbb3-63ed-7be3-9cca-6e806ef76e54` (a body any sender
writes can put every message ahead of all pending work; reported by
quaestor). Both are copied verbatim under `references/`.

The maintainer answered the operator's cycle BC questions with
AskUserQuestion (via `/roundfix --follow-up`, 2026-11-11 ~14:45–14:50 BRT).
The answers this Spec rests on are, verbatim:

- Question "Specs novas D (mensagens: pergunta não bloqueante recebe
  resposta; precedência de 'bloqueado' limitada) e E (review pré-PR revisa
  só o delta e com contexto depois de merge da main): como autorizar?" →
  "Após a 0379 (Recommended)" — option text: "Mesmas autorizações do ciclo
  BB (push, merge, release, teto minor); entregues depois da fila
  0377→0379, D antes de E."
- Question "Spec D (mensagens entre projetos): como limitar a precedência de
  uma mensagem que declara trabalho bloqueado, já que o corpo não é
  confiável (achado do Quaestor)?" → "Uma por remetente + juízo
  (Recommended)" — option text: "Só a mensagem bloqueante mais antiga de
  cada remetente passa à frente; as demais seguem a ordem normal. O texto
  diz que a precedência não dispensa o juízo do maintainer. Pergunta direta
  não bloqueante recebe resposta quando existir."

"Mesmas autorizações do ciclo BB" keeps the cycle BB answers of 2026-11-10
("Mesmas do BA, teto minor (Recommended)": "Escrever, entregar pela fila,
merge e release de cada Spec assim que pronta (patch nas duas correções; a
regra pode sair como minor), teto de 2 Tasks corretivas, eventos ao Jarvis,
sem notas às caixas, e Makefile/CI/go.mod/lint só onde a Spec exigir.") and,
through them, the cycle BA, AZ, AY, AX, AW, AV, AU, AT and AS answers, which
release to the Runs the Baseline source (`internal/baseline/assets/**`) with
its fixtures and goldens, the guides under `docs/agents/`, the owned skills,
and changes of Baseline and skills that reach the adopters, the cycle AJ
grant of 2026-10-19 "Mexer em Governed Paths declarados" (option text:
"'Concedo' para os caminhos que cada Spec declarar (módulos Baseline,
registros, skills, testes governados; fora: Makefile, CI, go.mod, lint, que
pedem autorização específica).") and the 2026-09-30 grant of every skill
("considere autorizado a ajustar todas as skills se necessário").

The Makefile, `.github/workflows`, `go.mod`, `go.sum` and the lint and
formatter configuration are grantable to this Spec only where its design
needs them. It needs none of them.

This record grants `implement`, `commit`, `push`, `pull_request`, `merge` and
`release`. The Spec is a `type: feat`, so the Release Ceiling is `minor` and
the release may change only the three version-only files in `release_paths`:
the changelog, the npm package manifest and the application version. No Task
performs it; the operator cuts it right after the merge.

The governed set was measured on 2026-11-11 on the branch `docs/plan-0380`,
based on `013af3b4`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` showed only this Spec's
  planning files and ADR-0390 afterwards, and the probe file never existed in
  the tree.
- It covered every file the Tasks declare under `creates:`, `interface:` or
  `instruction:`, every file the scratch prototype of tasks 01 to 03 changed
  or created (the derived Baseline files included), every planning file of
  this Spec, its adopted references and ADR-0390.
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
  ADR-0390 and the amendment notices in ADR-0383, ADR-0319 and ADR-0330.
- Every planning file of this Spec and its adopted references.

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
  `internal/judge`, `internal/cli/follow_up_inventory.go`, or the `### QA
  settlement` section of any skill.
- No clause is added, removed or renamed, and no code reads a message body
  for a marker phrase or a sender.
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
