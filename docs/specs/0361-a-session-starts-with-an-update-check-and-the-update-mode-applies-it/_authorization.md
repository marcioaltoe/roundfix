---
status: approved
granted: 2026-11-06
action: add the update.mode Project Config key, the read-only roundfix update check command with its daily cached release listing, its guide and the Roundfix setup reference, the Roundfix skill's update Skill Mode in the modes reference and SKILL.md, the mandatory session-start clause in the Baseline core module with its derived files and this repository's managed refresh, the glossary entries and ADR-0371
consuming: 0361-a-session-starts-with-an-update-check-and-the-update-mode-applies-it
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/modes.md
  - .agents/skills/roundfix/references/setup.md
  - docs/agents/agent-instructions.md
  - docs/agents/setup-context.json
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/agent-instructions.md
  - internal/baseline/assets/modules/core.json
  - internal/baseline/assets/profiles/standard-typescript-monorepo.json
  - internal/cli/cli_test.go
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

# Approved authority for Spec 0361

On 2026-11-06 the maintainer answered the operator's cycle AZ batch 4
questions with AskUserQuestion, after jarvis relayed his request to replace
the release notes in the projects' inboxes with an update check at session
start
(`references/2026-11-06-an-update-check-at-session-start-replaces-the-release-notes-in-the-inboxes.md`).
The answers this Spec rests on are, verbatim:

- Question "O Jarvis repassou seu pedido de 06/11: parar de mandar as notas de
  release às caixas (havia 30 paradas) e trocar por uma verificação de
  atualização na Baseline, com /roundfix --update e escolha por repositório
  entre automático e manual. Como seguir?" → "Spec agora, padrão avisar
  (Recommended)" — option text: "Paro as notas já. A Baseline ganha uma
  checagem no início da sessão: compara a versão do Roundfix e da Baseline com
  a última publicada e mostra o que muda. /roundfix --update aplica. Uma chave
  por repositório liga o modo automático; o padrão é só avisar. Quando sair,
  cancelo as notas pendentes nas caixas."
- Rules: "Mesmas do AZ (Recommended)" — option text: "Sem supervisão, merge
  pela fila, conserto e retry quando parar. Release assim que cada Spec ficar
  pronta; resposta ao OnionCry, ao sql-library e ao Jarvis."
- The maintainer's standing instruction of 2026-11-05: "Quero que faça os
  releases hoje assim que estiverem prontos. E se existe no baseline regra
  para um release por dia por algum motivo deve ser removida."

"Mesmas do AZ" keeps the cycle AZ rules of 2026-11-05, and through them the
cycle AY, AX, AW, AV, AU, AT and AS answers, which release to the Runs the
Baseline source (`internal/baseline/assets/**`), the guides under
`docs/agents/`, the skills, and changes of Baseline and skills that reach the
adopters ("mudanças de Baseline e skills que chegam aos adotantes
liberadas"), the cycle AJ grant of 2026-10-19 "Mexer em Governed Paths
declarados" (option text: "'Concedo' para os caminhos que cada Spec declarar
(módulos Baseline, registros, skills, testes governados; fora: Makefile, CI,
go.mod, lint, que pedem autorização específica).") and the 2026-09-30 grant of
every skill ("considere autorizado a ajustar todas as skills se necessário").

The Makefile, `.github/workflows`, `go.mod`, `go.sum` and the lint and
formatter configuration are grantable to this Spec when its design needs them.
It needs none of them.

This record grants `implement`, `commit`, `push`, `pull_request`, `merge` and
`release`. The Spec is a `type: feat`, so the Release Ceiling is `minor` and
the release may change only the three version-only files in `release_paths`:
the changelog, the npm package manifest and the application version. No Task
performs it; the operator cuts it right after the merge.

The governed set was measured on 2026-11-06 on the branch `docs/plan-0361`,
based on `9a285d77`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` showed no change from
  it, and the probe file never existed in the tree.
- It covered every file the Tasks declare under `creates:`, `interface:` or
  `instruction:`, every derived file the task_04 prototype rewrote, every
  planning file of this Spec, its adopted reference and ADR-0371.
- The ten paths above are governed: the Roundfix skill's `SKILL.md` pair and
  canonical `setup` and `modes` references (tasks 02 and 03), the CLI help
  test `internal/cli/cli_test.go` (task_02), and the core module, the
  Standard TypeScript profile and its formatter golden of the agent
  instructions guide, this repository's agent instructions guide and its
  setup manifest (task_04).

## What is not governed

- Every Go source and test file tasks 01 to 04 create or change under
  `internal/config`, `internal/app`, `internal/cli` (other than
  `cli_test.go`), `internal/baseline` and `skills`.
- The mirrors under `skills/roundfix/references/`,
  `skills/testdata/owned-skill-versions.json`,
  `internal/baseline/module-versions.json` and the catalog snapshots and plan
  goldens under `internal/baseline/testdata/`.
- `docs/user-guide/**`, `CONTEXT.md`, `docs/references/behavior-surfaces.json`,
  `docs/references/skill-coverage.json` and ADR-0371.
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
  `skills-lock.json`, any Vendored Skill under `.agents/skills`, any Baseline
  module other than `core`, any Baseline template or Setup Snapshot, any
  guide under `docs/agents/` other than through the managed refresh,
  `docs/references/coverage-record.json`, any Archive Record, any skill other
  than `roundfix`, any Roundfix reference other than `setup` and `modes`, or
  the `### QA settlement` section of any skill.
- No change to `roundfix upgrade`, its Version Freshness line, `roundfix
  baseline update`, the release workflow or the npm packages.
- No Task writes a harness hook into `.claude/`, `.codex/` or `.opencode/`.
- No test, Verification command or QA row:
  - reaches the network or runs the real `gh`;
  - reads or writes the real `~/.roundfix` or the real home directory;
  - starts an Agent session or an ACP adapter;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository;
  - records an environment value or writes a real-looking secret literal;
  - signals a process it did not start, or ends a process by name.
- No Task releases, tags or deploys; the release is the operator's, under the
  Release Ceiling above. Verification stays Daemon-owned, and Task status
  stays Daemon-written.
