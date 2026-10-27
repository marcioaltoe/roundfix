---
status: approved
granted: 2026-10-27
action: make QA acceptance match each declared-blocked row to a declaration naming its criterion, with one declaration naming several criteria; make roundfix settle refuse inside a Run or Task Worktree, follow the Spec plan its Run's branch carries and refuse the untracked executables a Task commit refuses; and document the rules in the qa-gate, write-prd and roundfix skills, the user guides, the QA prompt and CONTEXT.md
consuming: 0330-settle-follows-the-branch-plan-and-declared-rows-match-their-criteria
paths:
  - .agents/skills/qa-gate/SKILL.md
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/settle.md
  - .agents/skills/roundfix/references/spec-delivery.md
  - .agents/skills/write-prd/SKILL.md
  - .agents/skills/write-prd/references/prd-template.md
  - skills/qa-gate/SKILL.md
  - skills/roundfix/SKILL.md
  - skills/write-prd/SKILL.md
  - skills/write-prd/references/prd-template.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0330

On 2026-10-27, at about 11:40 BRT, the maintainer answered the operator's
`/roundfix follow-up` questions (inventory with `--messages`, Jev
`typesafe/jev-1.13-20260917`), which included Quaestor's proof of v0.125.0
(ai-memory message `01a17d0a-cd60-7191-a724-2988affea2d0` from
`default/gesttione-solutions-quaestor`, 2026-10-27 08:07 BRT), with
AskUserQuestion. The answers this Spec rests on are, verbatim:

- Scope: "A, B e C (Recommended)" — option text: "Lock + Clean vermelho;
  família do settle + aceite da QA; profile proof sob carga. Todas são health
  e duas se repetem em adotantes."
- QA acceptance: "Casar por linha (Recommended)" — option text: "Uma
  declaração cobre cada criterion que nomeia; o accept confere linha a linha,
  não a contagem. write-prd/qa-gate passam a dizer isso."
- Authorization: "Mesmas regras do AQ (Recommended)" — option text: "Push de
  branches, PR dos planos, entrega e merge pela Delivery Queue, eu escolho o
  ponto de release e aviso os adotantes por ai-memory."

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
Spec delivers after Spec 0329 and before Spec 0331.

The governed set was measured on 2026-10-27 on the branch `docs/plan-0330`,
based on `5e9475bf`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` showed no change from
  it, and the probe file never existed in the tree.
- It covered every file the Tasks declare under `creates:` or `interface:`,
  every planning file of this Spec, its adopted Backlog Entries, the removed
  `docs/backlog/` paths and ADR-0343.
- The ten paths above are governed: the canonical `qa-gate`, `write-prd`
  and `roundfix` skill files the Tasks change, the mirrored `SKILL.md` files
  and PRD template whose bytes rise with them. The Tasks read but never
  change `internal/spec/archive_test.go`, a historically bounded path.

## What is not governed

- `internal/spec/spec.go`, `internal/spec/qa.go`, `internal/spec/spec_test.go`,
  `internal/spec/qa_declared_rows_test.go`,
  `internal/speccheck/mechanical.go`, `internal/agent/spec_prompt.go`,
  `internal/agent/spec_prompt_test.go`.
- `internal/cli/settle.go`, `internal/cli/qa_report_declared_rows_test.go`,
  `internal/cli/settle_checkout_test.go`, `internal/cli/settle_plan_test.go`,
  `internal/cli/settle_executables_test.go`, `internal/daemon/task_engine.go`,
  `internal/daemon/task_engine_executables_test.go`.
- The mirrored references `skills/roundfix/references/settle.md` and
  `skills/roundfix/references/spec-delivery.md`,
  `skills/testdata/owned-skill-versions.json`.
- `docs/user-guide/commands/qa-report.md`, `docs/user-guide/commands/settle.md`,
  `CONTEXT.md`, `docs/references/behavior-surfaces.json`.
- ADR-0343, every planning file of this Spec and its adopted Backlog Entries.

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
  - .agents/skills/qa-gate/SKILL.md
  - .agents/skills/write-prd/SKILL.md
  - .agents/skills/roundfix/SKILL.md
  - skills/qa-gate/SKILL.md
  - skills/write-prd/SKILL.md
  - skills/roundfix/SKILL.md
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
  `docs/agents/`, any ADR other than ADR-0343 (written by the planning
  change), the Archive Record schema, any Archive Record, or the
  `### QA settlement` section of any skill.
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
