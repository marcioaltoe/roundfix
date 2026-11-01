---
status: approved
granted: 2026-11-01
action: make roundfix spec check --run-verification leave the QA gate's derived Verification unprobed through the predicate the Daemon uses to route the gate, make the SC-QA-VERIFICATION-AUTHORED finding print the exact expected Verification block in text and JSON, and teach the write-tasks skill and the Roundfix skill's spec reference to copy it
consuming: 0337-the-spec-check-leaves-the-qa-gate-unprobed-and-prints-its-expected-verification
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/spec.md
  - .agents/skills/write-tasks/SKILL.md
  - .agents/skills/write-tasks/references/task-template.md
  - internal/speccheck/coherence.go
  - skills/roundfix/SKILL.md
  - skills/write-tasks/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0337

On 2026-11-01 the maintainer answered the operator's cycle AT questions with
AskUserQuestion, in the `/roundfix follow-up` mode. The answers this Spec
rests on are, verbatim:

- Scope: "0337 primeiro, release já" — option text: "Só a 0337, que
  desbloqueia Fluxus e Quaestor, com release logo depois; as outras três num
  ciclo seguinte."
- Derived QA Verification: "spec check mostra (Recommended)" — option text:
  "O achado SC-QA-VERIFICATION-AUTHORED passa a imprimir o bloco exato
  esperado, e o template do write-tasks ensina a copiar dele. Nenhum comando
  novo. Junto: o spec check --run-verification deixa de medir a Verification
  derivada da QA (correção do Fluxus)."
- Authorization: "Mesmas do AS (Recommended)" — option text: "Push e PR do
  plano, entrega e merge pela fila, Makefile/workflows/go.mod/lint quando o
  desenho pedir, 'Consertar e seguir' em falha técnica, release v0.128.0 logo
  após o merge e aviso aos adotantes (Fluxus e Quaestor, com resposta às
  mensagens deles). As outras três mensagens viram Backlog Entries com as suas
  decisões, para o próximo ciclo."

Those answers keep the cycle AS rules: the answers of 2026-10-31 "Consertar e
seguir (Recommended)" (option text: "Diagnostico, emendo o plano ou crio uma
Spec bloqueadora e entrego, sem teto fixo de corretivas. Só pulo a Spec se o
conserto exigir uma decisão de produto ou arquitetura que contradiga o PRD
aprovado; aí ela sai do ciclo e a release vai com as demais, ou espera se for
a 0333.") and the Tooling answer "Makefile, .github/workflows, go.mod /
go.sum, Config de lint/formatter", and the answer of 2026-10-30 "Mesmas regras
do AR (Recommended)" (option text: "Push de branches, PR dos planos, entrega e
merge pela Delivery Queue, eu escolho o ponto de release e aviso os adotantes
por ai-memory."). Through them they keep the cycle AQ answers of 2026-10-26
("Sim, pode fazer a higiene." and "Vamos seguir a sua sugestão de sequencia
para a spec e os demais itens. Talvez seja importante usar o jev para
validar."), the cycle AJ grant of 2026-10-19 "Mexer em Governed Paths
declarados" (option text: "'Concedo' para os caminhos que cada Spec declarar
(módulos Baseline, registros, skills, testes governados; fora: Makefile, CI,
go.mod, lint, que pedem autorização específica).") and the 2026-09-30 grant of
every skill ("considere autorizado a ajustar todas as skills se necessário").

The Makefile, `.github/workflows`, `go.mod`, `go.sum` and the lint and
formatter configuration are grantable to this Spec when its design needs
them. This design needs none of them, so none is listed in `paths`.

This record grants `implement`, `commit`, `push`, `pull_request` and `merge`.
It does not grant `release`: the release is the operator's, after the merge,
as the Authorization answer states.

The governed set was measured on 2026-11-01 on the branch `docs/plan-0337`,
based on `340c9004`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` showed no change from
  it, and the probe file never existed in the tree.
- It covered every file the Tasks declare under `creates:` or `interface:`,
  every planning file of this Spec and its adopted references.
- The seven paths above are governed: the canonical Roundfix `spec`
  reference, the canonical `write-tasks` task template, the `SKILL.md` pairs
  of both skills whose version fields rise, and the Spec Check's coherence
  detectors.

## What is not governed

- `internal/cli/spec_check.go`, `internal/cli/spec_check_completed_tasks_test.go`,
  `internal/daemon/verification_probe.go`, `internal/daemon/task_engine.go`,
  `internal/speccheck/report.go`, `internal/speccheck/report_test.go`,
  `internal/speccheck/coherence_test.go`.
- The mirrors `skills/roundfix/references/spec.md` and
  `skills/write-tasks/references/task-template.md`, and
  `skills/testdata/owned-skill-versions.json`.
- `docs/user-guide/commands/spec.md`, `docs/references/behavior-surfaces.json`,
  `CONTEXT.md`,
  `docs/adr/0148-one-verification-prober-refuses-vacuous-tasks-before-agent-and-at-authoring.md`.
- Every planning file of this Spec and its adopted references.

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
  - .agents/skills/write-tasks/SKILL.md
  - skills/roundfix/SKILL.md
  - skills/write-tasks/SKILL.md
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
  `docs/agents/`, `docs/references/skill-coverage.json`, any Archive Record,
  or the `### QA settlement` section of any skill.
- No change to `spec.DerivedQAVerification`, `spec.QAReportEligibility`, the
  QA gate's audited head or the Daemon's QA step `runQAGate`;
  `taskPlanWithoutQAGate` only reads the shared predicate.
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
