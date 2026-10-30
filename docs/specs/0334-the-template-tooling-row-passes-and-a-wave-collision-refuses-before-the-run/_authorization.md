---
status: approved
granted: 2026-10-30
action: make the Spec Check read a Tooling authority row that declares no protected tooling mutation and cites an authorization record with an explicitly empty paths list as no mutation, say in the write-prd and write-techspec templates how that row ends, keep a Completed Task out of every Wave Collision, refuse a Wave Collision in the Implement Command's Preflight Validation before a Run exists, and document the rules in the Roundfix skill, the user guides and CONTEXT.md
consuming: 0334-the-template-tooling-row-passes-and-a-wave-collision-refuses-before-the-run
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/implement.md
  - .agents/skills/roundfix/references/spec.md
  - .agents/skills/write-prd/SKILL.md
  - .agents/skills/write-prd/references/prd-template.md
  - .agents/skills/write-techspec/SKILL.md
  - .agents/skills/write-techspec/references/techspec-template.md
  - internal/speccheck/constraints.go
  - skills/roundfix/SKILL.md
  - skills/write-prd/SKILL.md
  - skills/write-prd/references/prd-template.md
  - skills/write-techspec/SKILL.md
  - skills/write-techspec/references/techspec-template.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0334

On 2026-10-30 the maintainer answered the operator's cycle AS questions
(after a `roundfix follow-up --messages` inventory, judge
`typesafe/jev-1.13-20260917`, and a direct Jev priority check, `jev-1.14.0`)
with AskUserQuestion. The answers this Spec rests on are, verbatim:

- Plan: "Aprovar como está (Recommended)" — option text: "Quatro Specs na
  ordem validada pelo Jev; a 0332 entra primeiro para que a fila resolva
  sozinha conflitos de registros gerados durante a 0333."
- Authorization: "Mesmas regras do AR (Recommended)" — option text: "Push de
  branches, PR dos planos, entrega e merge pela Delivery Queue, eu escolho o
  ponto de release e aviso os adotantes por ai-memory."
- Release: "Uma no fim do ciclo (Recommended)" — option text: "v0.127.0
  depois das quatro Specs; os adotantes recebem a revisão de skills e as
  correções juntas."

The cycle AR rules those answers keep are the maintainer's answers of
2026-10-27, at about 11:40 BRT, among them "Mesmas regras do AQ
(Recommended)", which keep the cycle AQ answers of 2026-10-26, at about
15:00 BRT, "Sim, pode fazer a higiene." and "Vamos seguir a sua sugestão de
sequencia para a spec e os demais itens. Talvez seja importante usar o jev
para validar.", over the cycle AJ grant of 2026-10-19, at about 16:20 BRT,
"Mexer em Governed Paths declarados" (option text: "'Concedo' para os
caminhos que cada Spec declarar (módulos Baseline, registros, skills, testes
governados; fora: Makefile, CI, go.mod, lint, que pedem autorização
específica).") and the 2026-09-30 grant of every skill ("considere autorizado
a ajustar todas as skills se necessário"). Still not authorized: this
repository's Makefile, `.github/workflows`, `go.mod`, `go.sum`, and lint and
formatter configuration.

This record grants `implement`, `commit`, `push`, `pull_request` and `merge`.
It does not grant `release`: the operator chooses the release point, which
the cycle AS answer places after the four Specs, in v0.127.0. This Spec
delivers third in the cycle, after Specs 0332 and 0333 and before Spec 0335.

The governed set was measured on 2026-10-30 on the branch `docs/plan-0334`,
based on `35de7f74`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` showed no change from
  it, and the probe file never existed in the tree.
- It covered every file the Tasks declare under `creates:` or `interface:`,
  every planning file of this Spec, its adopted Backlog Entries and ADR-0347.
- The thirteen paths above are governed:
  - `internal/speccheck/constraints.go`, a path an earlier authorization
    bounded;
  - the canonical `write-prd` and `write-techspec` templates and the four
    `SKILL.md` files whose versions rise with them, with the two mirrored
    templates;
  - the canonical Roundfix `implement` and `spec` references and the
    Roundfix `SKILL.md` pair.

## What is not governed

- `internal/speccheck/report_test.go`,
  `internal/docscontract/command_documentation_test.go`.
- `internal/spec/collision.go`, `internal/spec/collision_test.go`,
  `internal/daemon/task_engine.go`,
  `internal/daemon/qa_gate_wave_collision_test.go`,
  `internal/cli/implement.go`, `internal/cli/implement_wave_collision_test.go`.
- The mirrored references `skills/roundfix/references/implement.md` and
  `skills/roundfix/references/spec.md`,
  `skills/testdata/owned-skill-versions.json`.
- `docs/user-guide/commands/implement.md`, `docs/user-guide/commands/spec.md`,
  `CONTEXT.md`, `docs/references/behavior-surfaces.json`.
- ADR-0347, every planning file of this Spec and its adopted Backlog Entries.

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
  - .agents/skills/write-prd/SKILL.md
  - .agents/skills/write-techspec/SKILL.md
  - .agents/skills/roundfix/SKILL.md
  - skills/write-prd/SKILL.md
  - skills/write-techspec/SKILL.md
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
  `docs/agents/`, any ADR other than ADR-0347 (written by the planning
  change), any Archive Record, or the `### QA settlement` section of any
  skill.
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
