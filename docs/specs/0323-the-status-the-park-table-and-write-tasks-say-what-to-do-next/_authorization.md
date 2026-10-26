---
status: approved
granted: 2026-10-26
action: make roundfix follow-up --status name the manual release classification count and command with no blank field, make the deliver Park Class table list blockers only and state the secret-found next action for that blocker alone under a documentation contract, and add the shared-test-file rule to the write-tasks skill and the Recorded Path glossary entry
consuming: 0323-the-status-the-park-table-and-write-tasks-say-what-to-do-next
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/deliver.md
  - .agents/skills/roundfix/references/follow-up.md
  - .agents/skills/write-tasks/SKILL.md
  - skills/roundfix/SKILL.md
  - skills/write-tasks/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0323

On 2026-10-26, at about 15:00 BRT, the maintainer answered the operator's
cycle AQ proposal in chat. The answers this Spec rests on are, verbatim:

- "Sim, pode fazer a higiene."
- "Vamos seguir a sua sugestão de sequencia para a spec e os demais itens.
  Talvez seja importante usar o jev para validar."

The sequence those answers approve is: hygiene (the inbox and the Baseline
snapshot check against `marcioaltoe/skills` `d12d7452` and later) → the Spec
Check probe fix → the cycle AP triage (the shared build cache in `gc` first)
→ skill and tool studies. This Spec belongs to the cycle AP triage, batch 2
of cycle AQ, and delivers after Specs 0319, 0320, 0321 and 0322.

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

The governed set was measured on the authoring branch `docs/plan-0323`, based
on `aa8b3db7`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` was identical before
  and after, and the probe file never existed in the tree.
- It covered every file the Tasks declare under `interface:` or
  `instruction:`, every planning file of this Spec and its three adopted
  Backlog Entry copies.
- `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/deliver.md`,
  `.agents/skills/roundfix/references/follow-up.md`,
  `.agents/skills/write-tasks/SKILL.md`, `skills/roundfix/SKILL.md` and
  `skills/write-tasks/SKILL.md` are governed and listed in `paths`.

## What is not governed

- `internal/cli/follow_up_status.go`, `internal/cli/follow_up_status_test.go`,
  `internal/docscontract/command_documentation_test.go`.
- `skills/roundfix/references/deliver.md`,
  `skills/roundfix/references/follow-up.md`,
  `skills/testdata/owned-skill-versions.json`.
- `docs/user-guide/commands/deliver.md`,
  `docs/user-guide/commands/follow-up.md`, `CONTEXT.md`,
  `docs/references/behavior-surfaces.json`.
- Every planning file of this Spec and its adopted Backlog Entry copies.

## Sanctioned regeneration

- `make skills-sync`, which copies the canonical `roundfix` and
  `write-tasks` skills to `skills/`.
- `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
  outputs: `skills/testdata/owned-skill-versions.json`,
  `.agents/skills/roundfix/SKILL.md`, `.agents/skills/write-tasks/SKILL.md`.
- `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`,
  outputs: `docs/references/behavior-surfaces.json`.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`, the lint,
  formatter or test-runner configuration, any Baseline module or derived
  Baseline file, `docs/agents/`, `.roundfixrc.yml`, any ADR, any skill other
  than `roundfix` and `write-tasks`, or the `### QA settlement` section of any
  skill.
- No change to the Release Plan, the Delivery Queue's Park Classes, blockers
  or next actions, or the Wave collision rule and its evidence.
- No test, Verification command or QA row:
  - reads or writes the real `~/.roundfix`;
  - starts an Agent session;
  - reaches the network;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository;
  - records an environment value.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
