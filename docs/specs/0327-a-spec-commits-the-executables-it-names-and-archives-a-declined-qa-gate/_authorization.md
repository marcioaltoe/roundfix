---
status: approved
granted: 2026-10-26
action: let a Task commit keep the untracked executables a Spec's authorization lists under executables and refuse every other one, also inside a new directory, in Task and Batch commits; let roundfix archive accept a Task Graph that declares qa declined once every Task is completed; make roundfix settle name each worktree it removes; and document the three rules in the write-tasks, roundfix and archive-spec skills, the user guides and CONTEXT.md
consuming: 0327-a-spec-commits-the-executables-it-names-and-archives-a-declined-qa-gate
paths:
  - .agents/skills/archive-spec/SKILL.md
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/archive.md
  - .agents/skills/roundfix/references/implement.md
  - .agents/skills/roundfix/references/settle.md
  - .agents/skills/roundfix/references/spec-delivery.md
  - .agents/skills/write-tasks/SKILL.md
  - .agents/skills/write-tasks/references/task-template.md
  - internal/spec/archive.go
  - skills/archive-spec/SKILL.md
  - skills/roundfix/SKILL.md
  - skills/write-tasks/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0327

On 2026-10-26, at about 22:35 BRT, the maintainer answered the operator's
`/roundfix follow-up` questions on Quaestor's report (ai-memory message
`01a17abd-20fe-7213-ab5c-f1a1bdb043fe` from
`default/gesttione-solutions-quaestor`, 2026-10-26 21:22 BRT) with
AskUserQuestion. The answers this Spec rests on are, verbatim:

- Grouping: "Duas Specs (Recomendado)" — 0327 Run lifecycle (executable
  scripts, settle exit 128, archive with `qa: declined`); 0328 Baseline
  profile-change frictions.
- Executables: "_authorization.md nomeia os executáveis (Recomendado)" —
  option text: "O commit da Task aceita 0755 só nos caminhos que a
  autorização da Spec lista; todo o resto continua recusado. O write-tasks e a
  recusa passam a citar essa regra."
- Authorization: "Mesmas regras do ciclo AQ (Recomendado)" — option text:
  "implement, commit, push, PR e merge nos Governed Paths que cada Spec
  declarar; fora Makefile, .github/workflows, go.mod/go.sum e lint. A release
  continua a meu critério, pelos critérios gravados."

The cycle AQ rules those answers keep are the maintainer's chat answers of
2026-10-26, at about 15:00 BRT, "Sim, pode fazer a higiene." and "Vamos seguir
a sua sugestão de sequencia para a spec e os demais itens. Talvez seja
importante usar o jev para validar.", over the cycle AJ grant of 2026-10-19,
at about 16:20 BRT, "Mexer em Governed Paths declarados" (option text:
"'Concedo' para os caminhos que cada Spec declarar (módulos Baseline,
registros, skills, testes governados; fora: Makefile, CI, go.mod, lint, que
pedem autorização específica).") and the 2026-09-30 grant of every skill
("considere autorizado a ajustar todas as skills se necessário"). Still not
authorized: this repository's Makefile, `.github/workflows`, `go.mod`,
`go.sum`, and lint and formatter configuration.

This record grants `implement`, `commit`, `push`, `pull_request` and `merge`.
It does not grant `release`: the operator chooses the release point. This
Spec delivers after Spec 0326.

The governed set was measured on 2026-10-26 on the branch `docs/plan-0327`,
based on `6a2630e5`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` showed no change from
  it, and the probe file never existed in the tree.
- It covered every file the Tasks declare under `creates:` or `interface:`,
  every planning file of this Spec, its adopted Backlog Entry copy, ADR-0340
  and the two open Backlog Entries the planning change writes.
- The twelve paths above are governed: the canonical `write-tasks`,
  `roundfix` and `archive-spec` skill files the Tasks change, the three
  mirrored `SKILL.md` files whose version fields rise, and
  `internal/spec/archive.go`, a historically bounded path.

## What is not governed

- `internal/authorization/authorization.go`,
  `internal/authorization/executables_test.go`,
  `internal/daemon/task_engine.go`, `internal/daemon/engine.go`,
  `internal/daemon/task_engine_test.go`,
  `internal/daemon/task_engine_executables_test.go`.
- `internal/cli/archive.go`, `internal/cli/settle.go`,
  `internal/cli/settle_test.go`, `internal/cli/archive_declined_qa_test.go`,
  `internal/cli/archive_pull_request_test.go`.
- The mirrored references under `skills/roundfix/references/` and
  `skills/write-tasks/references/task-template.md`,
  `skills/testdata/owned-skill-versions.json`.
- `docs/user-guide/commands/archive.md`, `settle.md` and `implement.md`,
  `CONTEXT.md`, `docs/references/behavior-surfaces.json`.
- ADR-0340, the two open Backlog Entries, every planning file of this Spec and
  its adopted Backlog Entry copy.

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
  - .agents/skills/write-tasks/SKILL.md
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/archive-spec/SKILL.md
  - skills/write-tasks/SKILL.md
  - skills/roundfix/SKILL.md
  - skills/archive-spec/SKILL.md
```

```yaml
command: go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage
outputs:
  - docs/references/behavior-surfaces.json
```

No mirror, version field or record is hand-edited. No top-level test is
removed or renamed, so the coverage record is not re-recorded.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`, the lint,
  formatter or test-runner configuration, `.roundfixrc.yml`, any Baseline
  module, profile, template or derived Baseline file, any guide under
  `docs/agents/`, any ADR other than ADR-0340 (written by the planning
  change), the Archive Record schema, or the `### QA settlement` section of
  any skill.
- No test, Verification command or QA row:
  - reads or writes the real `~/.roundfix`;
  - starts an Agent session;
  - reaches the network;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository, Quaestor's included;
  - records an environment value;
  - signals a process it did not start, or ends a process by name.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
