---
status: approved
granted: 2026-11-10
action: add a read-only reader of five user-level skill directories that names each Stale User Skill Copy, print it in roundfix skills check, roundfix doctor and the Update Check without changing an exit code, keep every internal/cli test away from the host's skill directories, and document it in the skills, doctor and update command guides, the Roundfix skill's setup and modes references, the glossary and ADR-0386
consuming: 0376-a-stale-user-level-skill-copy-is-reported-with-its-refresh-command
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/modes.md
  - .agents/skills/roundfix/references/setup.md
  - internal/cli/cli_test.go
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

# Approved authority for Spec 0376

On 2026-11-10 the operator recorded the open Backlog Entry this Spec adopts
(`references/2026-11-10-a-user-level-skill-copy-older-than-the-project-copy-is-loaded-unnoticed.md`).
The maintainer then answered the operator's cycle BB batch 2 questions with
AskUserQuestion (via `/roundfix --follow-up`, 2026-11-10 ~20:55–21:00 BRT).
The answers this Spec rests on are, verbatim:

- Question "Quais Backlog Entries viram Spec agora, depois da 0374? As que
  você não marcar ficam abertas no backlog." → all four selected, among them
  "Skills de usuário velhas (Recommended)" — option text: "O `doctor` e o
  `skills check` passam a apontar uma cópia de usuário de skill do Roundfix
  mais velha que a do binário, com o comando que a atualiza. Hoje cada
  release deixa as do vault velhas sem aviso."
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
`release`. The Spec is a `type: fix`, so the Release Ceiling is `patch` and
the release may change only the three version-only files in `release_paths`:
the changelog, the npm package manifest and the application version. No Task
performs it; the operator cuts it after the merge.

The governed set was measured on 2026-11-10 on the branch `docs/plan-0376`,
based on `b233761e` (re-measured after the rebase; the same set as on
`2b561a38`), with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` showed only this Spec's
  planning files, its adopted reference and ADR-0386 afterwards, and the
  probe file never existed in the tree.
- It covered every file the Tasks declare under `creates:`, `interface:` or
  `instruction:`, every file the scratch prototype of tasks 01 and 02
  changed or created, every planning file of this Spec, its adopted
  reference, the reference's former path under `docs/backlog/` and ADR-0386.
- The five paths above are governed: the Roundfix skill's `SKILL.md` pair
  and its canonical `modes` and `setup` references, which task_02 changes,
  and the CLI test helper file `internal/cli/cli_test.go`, whose
  `commandEnvironmentForTest` task_01 changes. The probe also named
  `.agents/skills/domain-modeling/SKILL.md`, which task_02 only reads; it is
  not granted.

## What is not governed

- Every other Go source and test file task_01 creates or changes under
  `skills/` and `internal/cli/`.
- The mirrors under `skills/roundfix/references/`,
  `skills/roundfix_modes_test.go` and
  `skills/testdata/owned-skill-versions.json`.
- `docs/user-guide/**`, `CONTEXT.md`, `docs/references/behavior-surfaces.json`
  and ADR-0386.
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
  than `roundfix`, any Roundfix reference other than `modes` and `setup`, or
  the `### QA settlement` section of any skill.
- No command, Task or QA row writes to a user-level skill directory or runs
  `roundfix skills install` with a user-level target.
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
