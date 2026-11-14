---
status: approved
granted: 2026-11-14
action: make roundfix spec check --run-verification prepare its disposable checkout with the worktree.copy files and the Worktree Bootstrap, as the Daemon prepares a Task Worktree, once before the first probed command, stop without probing when the bootstrap fails, and document it in the spec guide, the Roundfix skill's spec reference, the glossary and the ADR-0148 and ADR-0034 amendment notices
consuming: 0396-the-verification-probe-prepares-its-tree-as-the-daemon-does
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/spec.md
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

# Approved authority for Spec 0396

On 2026-11-14 indexus reported, in ai-memory message
`01a1d9be-ad8c-7f33-a039-0d8d4e60fef7`, that `roundfix spec check
--run-verification` called a command "honest" that the Daemon refused minutes
later as vacuous. The message is copied verbatim in this Spec's `references/`
(`references/2026-11-14-indexus-the-verification-probe-called-honest-a-command-the-daemon-refused-as-vacuous.md`).
The maintainer then answered the operator with AskUserQuestion
(2026-11-14 ~11:35 BRT). The answers this Spec rests on are, verbatim:

- Question "Mensagem do Indexus: o `roundfix spec check --run-verification`
  chamou de \"honest\" um comando que o Daemon depois recusou como vacuous
  (provável causa: a checagem roda sem as dependências instaladas). O
  Indexus contornou e não está bloqueado. O que fazemos?" → "Spec fix 0396
  (Recommended)" — option text: "Reproduzo a causa, escrevo a TechSpec no
  Opus (sem Fable) e entrego depois da 0393 como patch v0.145.1, com o
  Quaestor ocioso."
- Question "Se a Spec 0396 do Indexus for aberta, como a checagem `spec
  check --run-verification` deve se corrigir? (O Jev ficou indeciso
  exatamente nesta escolha.)" → "Igualar ao Daemon (Recommended)" — option
  text: "A checagem prepara as dependências como o Daemon prepara (ex.: bun
  install) antes de rodar o comando. Mais lenta, porém o veredito bate com o
  da entrega."
- Question "Para a Spec 0396 do Indexus (uma correção com o defeito medido),
  posso pular a etapa de Ideia e ir direto à TechSpec?" → "Pular a Ideia
  (Recommended)" — option text: "Correção com evidência medida; a forma já
  foi decidida na pergunta anterior. TechSpec direto, no Opus."
- Question "Autorizações para a Spec 0396 (Indexus) e para a 0397 (limite do
  Claude): valem as mesmas do ciclo BE?" → "Mesmas do BE (Recommended)" —
  option text: "Push, PR, merge, release depois de cada Spec, até 3
  correções; Makefile/CI/go.mod/lint só se a correção exigir. Com o Quaestor
  ocioso, 1 vaga."

"Mesmas do BE" carries the cycle BE answers of 2026-11-13 ~08:20 BRT: a
patch release after each Spec ("Patch após cada Spec (Recommended)") and up
to three corrections per Spec with Makefile, CI workflow, `go.mod` and lint
changes granted when a correction requires one ("Até 3 por Spec + Makefile/CI
(Recommended)"). The cycle BE rules keep the earlier cycles' grants of owned
skills and of changes to skills that reach the adopters, the cycle AJ grant of
2026-10-19 "Mexer em Governed Paths declarados" (option text: "'Concedo' para
os caminhos que cada Spec declarar (módulos Baseline, registros, skills,
testes governados; fora: Makefile, CI, go.mod, lint, que pedem autorização
específica).") and the 2026-09-30 grant of every skill ("considere autorizado
a ajustar todas as skills se necessário").

This Spec's design needs no Makefile, `.github/workflows`, `go.mod`, `go.sum`,
lint or formatter configuration change, and no Baseline module, template or
profile change. The answers would grant a Makefile, workflow, `go.mod` or lint
change if a correction required one; none is listed here, and any such
change must first be added to `paths` on the default branch with its reason.

This record grants `implement`, `commit`, `push`, `pull_request`, `merge` and
`release`. The Spec is a `type: fix`, so the Release Ceiling is `patch` and
the release may change only the three version-only files in `release_paths`:
the changelog, the npm package manifest and the application version. No Task
performs it; the operator cuts it after the merge.

The governed set was measured on 2026-11-14 on the branch `docs/plan-0396`,
based on `4f8a5b7d`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` showed only this Spec's
  planning files and ADR-0406 afterwards, and the probe file never existed in
  the tree.
- It covered every file the Tasks declare under `creates:`, `interface:` or
  `instruction:`, every file the scratch prototype of task_01 changed or
  created, every planning file of this Spec, its adopted reference, ADR-0406
  and `docs/agents/setup-context.json`.
- The three paths above are governed: the Roundfix skill's `SKILL.md` pair
  and its canonical `spec` reference, which task_01 changes. The probe also
  named `docs/agents/setup-context.json`, which no Task changes; it is not
  granted.

## What is not governed

- `internal/cli/spec_check.go`, `internal/cli/spec_check_bootstrap_test.go`
  and `internal/worktree/worktree.go`, which task_01 changes or creates.
- The mirror `skills/roundfix/references/spec.md` and
  `skills/testdata/owned-skill-versions.json`.
- `docs/user-guide/commands/spec.md`, `docs/references/behavior-surfaces.json`,
  `CONTEXT.md`, ADR-0148, ADR-0034 and ADR-0406.
- Every planning file of this Spec and its adopted reference.

## Sanctioned regeneration

The repository-owned commands write the skill mirror and the generated
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
  `docs/references/coverage-record.json`, `docs/references/skill-coverage.json`,
  any Archive Record, any skill other than `roundfix`, any Roundfix reference
  other than `spec`, or the `### QA settlement` section of any skill.
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
