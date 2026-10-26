---
status: approved
granted: 2026-10-26
action: make the unreaped-group detach fixture test wait for its member's zombie state before its group probe, prove the zombie reading rejects an exiting member, and widen the Hard Link rule test to every reference and every module test file, changing test files only
consuming: 0324-the-group-probe-waits-for-a-zombie-and-the-hard-link-rule-sees-all-forms
paths: []
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0324

On 2026-10-26, at about 15:00 BRT, the maintainer answered the operator's
cycle AQ proposal in chat. The answers this Spec rests on are, verbatim:

- "Sim, pode fazer a higiene."
- "Vamos seguir a sua sugestão de sequencia para a spec e os demais itens.
  Talvez seja importante usar o jev para validar."

The sequence those answers approve is: hygiene (the inbox and the Baseline
snapshot check against `marcioaltoe/skills` `d12d7452` and later) → the Spec
Check probe fix → the cycle AP triage (the shared build cache in `gc` first)
→ skill and tool studies. This Spec belongs to the cycle AP triage (batch 2
of cycle AQ: triage items 9 and 12) and delivers after Specs 0319 to 0322,
independent of them.

The cycle AQ rules carry the cycle AP rules unchanged: the Governed Paths each
Spec declares, under the cycle AJ grant of 2026-10-19, at about 16:20 BRT,
"Mexer em Governed Paths declarados" (option text: "'Concedo' para os caminhos
que cada Spec declarar (módulos Baseline, registros, skills, testes
governados; fora: Makefile, CI, go.mod, lint, que pedem autorização
específica).") and "Merge squash + release + push". Still not authorized:
this repository's Makefile, `.github/workflows`, `go.mod`, `go.sum`, and lint
and formatter configuration.

This record grants `implement`, `commit`, `push`, `pull_request` and `merge`.
It does not grant `release`: the operator chooses the release point.

The governed set was measured on the authoring branch `docs/plan-0324`,
based on `aa8b3db7`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` was identical before
  and after, and the probe file never existed in the tree.
- It covered every file the Tasks declare under `interface:` or
  `instruction:`, every planning file of this Spec, its adopted Backlog
  Entry copies and `docs/references/2026-10-24-test-suite-audit.md`.
- None is governed, so `paths` is empty.

## What is not governed

- `internal/cli/implement_detach_teardown_test.go`,
  `internal/cli/detach_fixture_group_darwin_test.go`,
  `internal/cli/detach_fixture_group_linux_test.go`,
  `internal/cli/detach_fixture_group_other_test.go`,
  `internal/testfixture/hard_link_rule_test.go`.
- `docs/references/2026-10-24-test-suite-audit.md`.
- Every planning file of this Spec and its adopted Backlog Entry copies.

## Sanctioned regeneration

None.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`, the lint,
  formatter or test-runner configuration, any production file, any Baseline
  module or derived Baseline file, `docs/agents/`, `CONTEXT.md`,
  `.roundfixrc.yml`, any skill, any ADR, or the `### QA settlement` section
  of any skill.
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
