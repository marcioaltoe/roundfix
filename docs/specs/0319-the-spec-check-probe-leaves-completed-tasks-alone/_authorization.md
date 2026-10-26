---
status: approved
granted: 2026-10-26
action: make roundfix spec check --run-verification probe only the Tasks a Run would still execute, through one exported Daemon function the Daemon's scheduler and QA routing also read, report each completed Task's command as not probed with its reason in text, JSON and the verdict line, and document the rule in the spec guide, the Roundfix and qa-gate skills, ADR-0148 and CONTEXT.md
consuming: 0319-the-spec-check-probe-leaves-completed-tasks-alone
paths:
  - .agents/skills/qa-gate/SKILL.md
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/spec.md
  - skills/qa-gate/SKILL.md
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0319

On 2026-10-26, at about 15:00 BRT, the maintainer answered the operator's
cycle AQ proposal in chat. The answers this Spec rests on are, verbatim:

- "Sim, pode fazer a higiene."
- "Vamos seguir a sua sugestão de sequencia para a spec e os demais itens.
  Talvez seja importante usar o jev para validar."

The sequence those answers approve is: hygiene (the inbox and the Baseline
snapshot check against `marcioaltoe/skills` `d12d7452` and later) → this Spec
Check probe fix → the cycle AP triage (the shared build cache in `gc` first)
→ skill and tool studies. The delivery order of the cycle's Specs is 0319 →
0320, independent of each other.

The cycle AQ rules carry the cycle AP rules unchanged: the Governed Paths each
Spec declares, under the cycle AJ grant of 2026-10-19, at about 16:20 BRT,
"Mexer em Governed Paths declarados" (option text: "'Concedo' para os caminhos
que cada Spec declarar (módulos Baseline, registros, skills, testes
governados; fora: Makefile, CI, go.mod, lint, que pedem autorização
específica).") and "Merge squash + release + push", and the skills
authorization "considere autorizado a ajustar todas as skills se necessário".
Still not authorized: this repository's Makefile, `.github/workflows`,
`go.mod`, `go.sum`, and lint and formatter configuration.

This record grants `implement`, `commit`, `push`, `pull_request` and `merge`.
It does not grant `release`: the operator chooses the release point.

The governed set was measured on the authoring branch `docs/plan-0319`, based
on `86cb6b73`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` was identical before
  and after, and the probe file never existed in the tree.
- It covered every file the Tasks declare under `interface:` or `creates:`,
  every planning file of this Spec and its adopted Backlog Entry copy.
- `.agents/skills/qa-gate/SKILL.md`, `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/spec.md`, `skills/qa-gate/SKILL.md` and
  `skills/roundfix/SKILL.md` are governed and listed in `paths`.

## What is not governed

- `internal/daemon/verification_probe.go`, `internal/daemon/task_engine.go`,
  `internal/daemon/task_awaits_work_test.go`, `internal/cli/spec_check.go`,
  `internal/cli/spec_check_completed_tasks_test.go`,
  `internal/speccheck/report.go`, `internal/speccheck/report_test.go`.
- `skills/roundfix/references/spec.md`,
  `skills/testdata/owned-skill-versions.json`.
- `docs/user-guide/commands/spec.md`, `CONTEXT.md`,
  `docs/adr/0148-one-verification-prober-refuses-vacuous-tasks-before-agent-and-at-authoring.md`,
  `docs/references/behavior-surfaces.json`.
- Every planning file of this Spec and its adopted Backlog Entry copy.

## Sanctioned regeneration

- `make skills-sync`, which copies the canonical `roundfix` and `qa-gate`
  skills to `skills/`.
- `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
  outputs: `skills/testdata/owned-skill-versions.json`,
  `.agents/skills/roundfix/SKILL.md`, `.agents/skills/qa-gate/SKILL.md`.
- `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`,
  outputs: `docs/references/behavior-surfaces.json`.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`, the lint,
  formatter or test-runner configuration, any Baseline module or derived
  Baseline file, `docs/agents/`, `.roundfixrc.yml`, any skill other than
  `roundfix` and `qa-gate`, any ADR other than ADR-0148, or the
  `### QA settlement` section of any skill.
- No change to the prober's classification, the disposable checkout, or which
  Tasks the Daemon runs, probes and gates.
- No test, Verification command or QA row:
  - reads or writes the real `~/.roundfix`;
  - starts an Agent session;
  - reaches the network;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository;
  - records an environment value.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
