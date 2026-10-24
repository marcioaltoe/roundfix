---
status: approved
granted: 2026-10-24
action: remove disable-model-invocation from the Roundfix-owned setup-context-driven skill and its shipped mirror, raise and record its version, and state in the skill and the Baseline guide that starting it writes nothing because roundfix baseline still writes only after the person approves the exact Plan Digest, leaving the Baseline's person-only clause and every other locked skill unchanged
consuming: 0316-an-agent-may-start-the-setup-skill
paths:
  - .agents/skills/setup-context-driven/SKILL.md
  - skills/setup-context-driven/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0316

On 2026-10-24 the maintainer ran `/roundfix --follow-up` and, at about 11:15
BRT, answered the cycle AP batch 5 questions through AskUserQuestion. The
answer this Spec rests on is:

- setup skill: "Tirar a trava da skill (Recomendado)" (option text: "Remove
  `disable-model-invocation: true` da skill setup-context-driven e o bloqueio
  no skill-dispatch; as confirmações do próprio comando (Plan Digest,
  --confirm-plan) continuam. Spec pequena, junto da 0315 ou logo depois.").

The request it answers came from the maintainer through Quaestor (ai-memory
message `01a16dbe`, 2026-10-24), "preservando as confirmações do próprio
comando para planos e operações quando aplicáveis"; the adopted Backlog Entry
records it. The same answers order the batch: Specs 0311 to 0314, then 0315,
then 0316, then 0307, 0309, 0310 and 0317.

The "bloqueio no skill-dispatch" is the Baseline's generic clause that tells
an Agent to ask the person to run a skill whose metadata turns off model
invocation. It names no skill, so removing the field from this skill lifts the
block for it; the clause stays for the skills that keep the field. This Spec
therefore changes no Baseline module.

The cycle AP rules that answer carries were approved on 2026-10-22 (questions
of 15:25 BRT, answered by about 18:30 BRT): "the same as cycle AO", the
Governed Paths each Spec declares, and still not this repository's Makefile,
`.github/workflows`, `go.mod`, `go.sum`, or lint and formatter configuration.
Those rules carry the cycle AJ grant of 2026-10-19, at about 16:20 BRT,
"Mexer em Governed Paths declarados" (option text: "'Concedo' para os caminhos
que cada Spec declarar (módulos Baseline, registros, skills, testes
governados; fora: Makefile, CI, go.mod, lint, que pedem autorização
específica).") and "Merge squash + release + push", and the skills
authorization "considere autorizado a ajustar todas as skills se necessário".

This record grants `implement`, `commit`, `push`, `pull_request` and `merge`.
It does not grant `release`: the operator releases.

The governed set was measured on the authoring branch `docs/plan-0315-0316`,
based on `64c56d7f`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` was identical before
  and after, and the probe file never existed in the tree.
- It covered every file the Task declares under `interface:` or `creates:`,
  every planning file of this Spec and its adopted Backlog Entry copy.
- `.agents/skills/setup-context-driven/SKILL.md` and
  `skills/setup-context-driven/SKILL.md` are governed and listed in `paths`.

## What is not governed

- `skills/testdata/owned-skill-versions.json`,
  `skills/setup_skill_invocation_test.go`,
  `docs/user-guide/context-driven-development.md` and
  `docs/references/behavior-surfaces.json`.
- Every planning file of this Spec and its adopted Backlog Entry copy.

## Sanctioned regeneration

- `make skills-sync`, which copies the canonical skill to
  `skills/setup-context-driven/`.
- `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
  outputs: `skills/testdata/owned-skill-versions.json`.
- `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`,
  outputs: `docs/references/behavior-surfaces.json`.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`, the lint,
  formatter or test-runner configuration, any Baseline module or derived
  Baseline file, `docs/agents/`, `.roundfixrc.yml`, `CONTEXT.md`, any Go
  production file, any skill other than `setup-context-driven`, or the
  `### QA settlement` section of any skill.
- No change to `roundfix baseline`'s confirmations: the Plan Digest,
  `--confirm-plan`, `--yes` and the interactive approval stay.
- No test, Verification command or QA row:
  - reads or writes the real `~/.roundfix`;
  - starts an Agent session;
  - reaches the network;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository;
  - records an environment value.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
