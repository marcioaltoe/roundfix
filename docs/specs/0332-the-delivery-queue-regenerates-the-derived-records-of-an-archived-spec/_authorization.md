---
status: approved
granted: 2026-10-30
action: make the coverage record writer name its regeneration command to the suite guard; make the regeneration contract model an archived item sanctioned only by its Archive Record, with a stale coverage record and an unsanctioned case; add a queue-level archived-item conflict fixture; record the repository rule and the glossary revision; re-record the coverage record with its command
consuming: 0332-the-delivery-queue-regenerates-the-derived-records-of-an-archived-spec
paths:
  - internal/spec/coverage_test.go
  - docs/agents/specific-repository.md
  - docs/references/coverage-record.json
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0332

On 2026-10-30 the maintainer answered the cycle AS questions
(AskUserQuestion). The operator asked them after a follow-up inventory
(`roundfix follow-up --messages`, judge `typesafe/jev-1.13-20260917`) and a
direct Jev priority check (`jev-1.14.0`). Verbatim:

- Plan: "Aprovar como está (Recommended)" — option text: "Quatro Specs na
  ordem validada pelo Jev; a 0332 entra primeiro para que a fila resolva
  sozinha conflitos de registros gerados durante a 0333."
- Authorization: "Mesmas regras do AR (Recommended)" — option text: "Push de
  branches, PR dos planos, entrega e merge pela Delivery Queue, eu escolho o
  ponto de release e aviso os adotantes por ai-memory."
- Release: "Uma no fim do ciclo (Recommended)" — option text: "v0.127.0
  depois das quatro Specs; os adotantes recebem a revisão de skills e as
  correções juntas."

The cycle AR rules are those of 2026-10-27, about 11:40 BRT. They carry
forward the cycle AQ rules of 2026-10-26 ("Sim, pode fazer a higiene.").
Those rules allow the Governed Paths each Spec declares. They still do not
allow this repository's Makefile, `.github/workflows`, `go.mod`, `go.sum`, or
lint and formatter configuration. They also carry two older grants:

- the cycle AJ grant of 2026-10-19, "Mexer em Governed Paths declarados"
  (option text: "'Concedo' para os caminhos que cada Spec declarar (módulos
  Baseline, registros, skills, testes governados; fora: Makefile, CI, go.mod,
  lint, que pedem autorização específica).");
- the grant for agent guides, "Autorizar os dois" (Baseline source and
  `docs/agents/*`).

This record grants `implement`, `commit`, `push`, `pull_request` and `merge`.
It does not grant `release`, which stays at the maintainer's discretion.

The governed set was measured on 2026-10-30 on the branch `docs/plan-0332`,
based on `35de7f74`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` was identical before
  and after, and the probe file never existed in the tree.
- It covered every file the Tasks declare under `interface:`, `creates:` or
  `instruction:`, the new ADR, the adopted Backlog Entry and every planning
  file of this Spec.
- The three paths above are governed and listed in `paths`.
- `.roundfixrc.yml` and `.agents/skills/domain-modeling/SKILL.md` are
  governed too, but task_01 only reads them as instructions and no Task
  changes them.
- Every other declared file is not governed:
  - `internal/config/regeneration_declared_test.go`;
  - `internal/cli/deliver_derived_archived_item_test.go`;
  - the existing conflict tests task_02 reads;
  - `docs/references/skill-coverage.json`;
  - `CONTEXT.md`;
  - the ADR;
  - the planning files.

## Why each governed path is needed

- `internal/spec/coverage_test.go` holds the coverage record writer, which
  must name its command to the suite guard (ADR-0345).
- `docs/agents/specific-repository.md` holds the repository rules. The rule
  that a derived-record writer names its command belongs there.
- `docs/references/coverage-record.json` is re-recorded after task_02 adds
  two top-level tests.

## Sanctioned regeneration

The repository coverage record is rewritten only by its record command, as
the last step of task_02:

```yaml
command: go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record -count=1
outputs:
  - docs/references/coverage-record.json
```

No record is hand-edited.

## Limits

This Spec changes none of the following:

- the Makefile, the CI workflows, `go.mod`, `go.sum`, or the lint, formatter
  or test-runner configuration;
- `.roundfixrc.yml`, `internal/suiteguard/`, `internal/suiteguardcontract/`
  or `internal/cli/deliver_workflow.go`;
- any Baseline module, profile or snapshot;
- any skill or `skills-lock.json`;
- the `### QA settlement` section of any skill;
- any file in another Spec's folder.

No test, Verification command or QA row:

- reads or writes the real `~/.roundfix`;
- reaches the network or pushes to a remote other than a local bare
  repository;
- sends a request to a provider, the judge or a reviewer;
- reads another repository;
- records an environment value.

There is no release, tag or deployment. Verification stays Daemon-owned, and
Task status stays Daemon-written.
