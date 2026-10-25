---
status: approved
granted: 2026-10-25
action: make the roundfix baseline update preview list each external skill the adopted modules require that is absent from .agents/skills or skills-lock.json, or no longer matches its lock entry, with its restore command, and report plan_ready instead of current; make roundfix skills check fail on the same skills inside an adopted repository; pin roundfix doctor's existing failure with a test; and document the result in the guides, the Roundfix skill references and CONTEXT.md
consuming: 0318-the-baseline-reports-a-required-skill-that-is-not-installed
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/baseline.md
  - .agents/skills/roundfix/references/setup.md
  - internal/cli/cli_test.go
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0318

On 2026-10-25 the operator relayed to the maintainer the Fluxus message
`01a16ea1` (ai-memory, 2026-10-24 12:57 BRT), and at about 09:05 BRT the
maintainer answered the cycle AP batch 7 questions through AskUserQuestion.
The answers this Spec rests on are:

- Order: "0317 logo após a 0313 (Recomendado)" (option text: "Ordem nova:
  0313 → 0317 → 0314 → 0315 → 0316 → 0307 → 0309 → 0310. A 0317 sai na v0.121
  e desbloqueia o Fluxus em cerca de 3–5 h; a 0315 (Oraculum, que tem
  contorno) passa para a v0.123.").
- Detection: "Spec 0318 logo após a 0317 (Recomendado)" (option text:
  "Correção pequena: `baseline update`, `doctor` e `skills check` acusam
  quando o dispatch exige uma skill que não está instalada. 1–2 Tasks + QA.
  Não atrasa a 0317.").

The final delivery order is 0313 → 0317 → 0318 → 0314 → 0315 → 0316 → 0307 →
0309 → 0310.

The cycle AP rules those answers carry were approved on 2026-10-22 (questions
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

The governed set was measured on the authoring branch `docs/plan-0318`, based
on `f3c39140`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` was identical before
  and after, and the probe file never existed in the tree.
- It covered every file the Tasks declare under `interface:` or `creates:`,
  every planning file of this Spec and its adopted Backlog Entry copy.
- `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/baseline.md`,
  `.agents/skills/roundfix/references/setup.md`,
  `internal/cli/cli_test.go` and `skills/roundfix/SKILL.md` are governed and
  listed in `paths`.

## What is not governed

- `internal/cli/baseline_update.go`, `internal/cli/cli.go`,
  `internal/cli/required_skills.go`,
  `internal/cli/skills_check_repository.go`,
  `internal/cli/baseline_update_missing_skills_test.go`,
  `internal/cli/baseline_update_retired_skills_test.go`,
  `internal/cli/skills_check_repository_test.go`,
  `internal/cli/doctor_test.go`.
- `skills/roundfix/references/baseline.md`,
  `skills/roundfix/references/setup.md`,
  `skills/testdata/owned-skill-versions.json`.
- `docs/user-guide/commands/baseline.md`,
  `docs/user-guide/commands/skills.md`, `CONTEXT.md`,
  `docs/references/behavior-surfaces.json`.
- Every planning file of this Spec and its adopted Backlog Entry copy.

## Sanctioned regeneration

- `make skills-sync`, which copies the canonical Roundfix skill to
  `skills/roundfix/`.
- `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
  outputs: `skills/testdata/owned-skill-versions.json`,
  `.agents/skills/roundfix/SKILL.md`.
- `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`,
  outputs: `docs/references/behavior-surfaces.json`.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`, the lint,
  formatter or test-runner configuration, any Baseline module or derived
  Baseline file, `docs/agents/`, `.roundfixrc.yml`, `internal/cli/doctor.go`,
  any skill other than `roundfix`, or the `### QA settlement` section of any
  skill.
- No change to `--yes`, `--confirm-plan`, `--no-skills`, the applied skills
  stage, `skills.differing` or Doctor's output.
- No test, Verification command or QA row:
  - reads or writes the real `~/.roundfix`;
  - reads this repository's or the host's installed skills as test input;
  - starts an Agent session;
  - reaches the network;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository;
  - records an environment value.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
