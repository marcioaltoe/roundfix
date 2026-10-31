---
status: approved
granted: 2026-10-31
action: make the Spec Check's untasked-skills rule and the Behavior Surface fingerprint read one file set through one shared predicate, so a Go test file never asks for a skills Task while a source file still does, and say so in the spec command guide and the Roundfix skill's spec reference
consuming: 0336-the-untasked-skills-check-reads-the-files-the-fingerprint-reads
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
---

# Approved authority for Spec 0336

On 2026-10-31 the maintainer answered the operator's cycle AS batch 2
questions with AskUserQuestion. The answers this Spec rests on are, verbatim:

- Scope: "Entra no AS como 0336" — option text: "Planejo a Spec 0336 agora,
  ela entra na fila depois da 0335 e sai na mesma v0.127.0."
- Tooling: "Makefile, .github/workflows, go.mod / go.sum, Config de
  lint/formatter" — question: "Hoje o grant do ciclo AS proíbe que os Runs
  mexam em Makefile (fora da linha OWNED_SKILLS da 0333), .github/workflows,
  go.mod/go.sum e config de lint/formatter. Para nenhuma Spec travar por
  autorização, quais destes eu libero aos Runs e aos planos de 0332 a 0336,
  sempre que o desenho da Spec pedir?"
- Technical failure: "Consertar e seguir (Recommended)" — option text:
  "Diagnostico, emendo o plano ou crio uma Spec bloqueadora e entrego, sem
  teto fixo de corretivas. Só pulo a Spec se o conserto exigir uma decisão de
  produto ou arquitetura que contradiga o PRD aprovado; aí ela sai do ciclo e
  a release vai com as demais, ou espera se for a 0333."

Those answers keep the cycle AS answers of 2026-10-30, verbatim:

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
a ajustar todas as skills se necessário").

The Tooling answer of 2026-10-31 makes the Makefile, `.github/workflows`,
`go.mod`, `go.sum` and the lint and formatter configuration grantable to this
Spec when its design needs them. This design needs none of them, so none is
listed in `paths`.

This record grants `implement`, `commit`, `push`, `pull_request` and `merge`.
It does not grant `release`: the operator chooses the release point, which
the cycle AS answers place after the cycle's Specs, in v0.127.0. This Spec
delivers in its own queue after Spec 0335 merges.

The governed set was measured on 2026-10-31 on the branch `docs/plan-0336`,
based on `0c27cdc9`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` showed no change from
  it, and the probe file never existed in the tree.
- It covered every file the Tasks declare under `creates:` or `interface:`,
  every planning file of this Spec and its adopted Backlog Entry, at both its
  former and its adopted path.
- The three paths above are governed: the canonical Roundfix `spec`
  reference and the Roundfix `SKILL.md` pair whose version fields rise with
  it.

## What is not governed

- `internal/skillcoverage/skillcoverage.go`,
  `internal/skillcoverage/skillcoverage_test.go`,
  `internal/speccheck/skills_declaration_test.go`,
  `internal/docscontract/skill_coverage_test.go`.
- The mirrored reference `skills/roundfix/references/spec.md` and
  `skills/testdata/owned-skill-versions.json`.
- `docs/user-guide/commands/spec.md`, `docs/references/behavior-surfaces.json`.
- Every planning file of this Spec and its adopted Backlog Entry.

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

No mirror, version field or record is hand-edited. No top-level test is
removed or renamed, so the coverage record is not re-recorded.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`, the lint,
  formatter or test-runner configuration, `.roundfixrc.yml`, any Baseline
  module, profile, template or derived Baseline file, any guide under
  `docs/agents/`, `docs/references/skill-coverage.json`, `CONTEXT.md`, any
  ADR, any Archive Record, or the `### QA settlement` section of any skill.
- No test, Verification command or QA row:
  - reads or writes the real `~/.roundfix`;
  - starts an Agent session;
  - reaches the network;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository;
  - records an environment value;
  - signals a process it did not start, or ends a process by name.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
