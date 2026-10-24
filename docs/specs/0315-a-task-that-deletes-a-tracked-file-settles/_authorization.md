---
status: approved
granted: 2026-10-24
action: make the three readers that list tracked files and then read them from the working tree skip a tracked path the working tree no longer holds, so a gated Task whose declared deletion is only on disk settles through the Spec Consistency settlement check, while every other read failure, the deletes check, link findings and the authorization check keep judging the deletion
consuming: 0315-a-task-that-deletes-a-tracked-file-settles
paths: []
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0315

On 2026-10-24 the maintainer ran `/roundfix --follow-up` and, at about 11:15
BRT, answered the cycle AP batch 5 questions through AskUserQuestion. The
answer this Spec rests on is:

- deletes fix: "Spec 0315 após a 0314 (Recomendado)" (option text: "Correção
  pequena, 1–2 Tasks + QA, entregue depois da auditoria dos testes e antes
  da 0307. Avisa o Oraculum quando sair.").

The same answers order the batch: Specs 0311 to 0314, then 0315, then 0316,
then 0307, 0309, 0310 and 0317. The operator tells oraculum when this Spec's
release ships.

The cycle AP rules that answer carries were approved on 2026-10-22 (questions
of 15:25 BRT, answered by about 18:30 BRT): "the same as cycle AO", the
Governed Paths each Spec declares, and still not this repository's Makefile,
`.github/workflows`, `go.mod`, `go.sum`, or lint and formatter configuration.
Those rules carry the cycle AJ grant of 2026-10-19, at about 16:20 BRT,
"Mexer em Governed Paths declarados" (option text: "'Concedo' para os caminhos
que cada Spec declarar (módulos Baseline, registros, skills, testes
governados; fora: Makefile, CI, go.mod, lint, que pedem autorização
específica).") and "Merge squash + release + push".

This record grants `implement`, `commit`, `push`, `pull_request` and `merge`.
It does not grant `release`: the operator releases.

The governed set was measured on the authoring branch `docs/plan-0315-0316`,
based on `64c56d7f`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` was identical before
  and after, and the probe file never existed in the tree.
- It covered every file the Tasks declare under `interface:` or `creates:`,
  every planning file of this Spec and its adopted Inbox Entry copy.
- None of them is governed, so `paths:` is empty.

## What is not governed

- `internal/spec/citation_sources.go`, `internal/spec/history_entries.go`,
  `internal/spec/stale_citations.go`,
  `internal/spec/citation_sources_working_tree_test.go` and
  `internal/daemon/task_deletes_settlement_test.go`.
- Every planning file of this Spec and its adopted Inbox Entry copy.

## Sanctioned regeneration

None. No guide, skill, Baseline module, record or mirror changes, so no
regeneration command runs.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`, the lint,
  formatter or test-runner configuration, any Baseline module,
  `.roundfixrc.yml`, `docs/agents/setup-context.json`, any guide, skill or
  glossary entry, or the `### QA settlement` section of any skill.
- No change to the deletes check, the scope record of undeclared paths, the
  authorization settlement check or the Governed Path rule.
- No test, Verification command or QA row:
  - reads or writes the real `~/.roundfix`;
  - reaches the network;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository;
  - records an environment value.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
