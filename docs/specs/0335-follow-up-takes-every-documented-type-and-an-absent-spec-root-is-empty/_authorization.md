---
status: approved
granted: 2026-10-30
action: make the Spec Root resolver report an absent in-repository root as a value, switch the enumerating readers to it and keep the named-Spec refusals; give the follow-up message parser, its error and its ordering fallback one type table; name the JSON fields and types in both follow-up guides and pin them with a documentation contract test; update the configuration guide, the Roundfix skill references, the glossary and the Behavior Surface record
consuming: 0335-follow-up-takes-every-documented-type-and-an-absent-spec-root-is-empty
paths:
  - .agents/skills/roundfix/references/archive.md
  - .agents/skills/roundfix/references/baseline.md
  - .agents/skills/roundfix/references/deliver.md
  - .agents/skills/roundfix/references/follow-up.md
  - .agents/skills/roundfix/references/modes.md
  - .agents/skills/roundfix/references/spec.md
  - .agents/skills/roundfix/SKILL.md
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0335

On 2026-10-30 the maintainer answered the cycle AS questions
(AskUserQuestion), asked after a follow-up inventory read with `roundfix
follow-up --messages` and judge `typesafe/jev-1.13-20260917`, and a direct Jev
priority check (`jev-1.14.0`). Verbatim:

- Plan: "Aprovar como está (Recommended)" — option text: "Quatro Specs na
  ordem validada pelo Jev; a 0332 entra primeiro para que a fila resolva
  sozinha conflitos de registros gerados durante a 0333." This Spec is the
  fourth, delivered after 0332, 0333 and 0334.
- Authorization: "Mesmas regras do AR (Recommended)" — option text: "Push de
  branches, PR dos planos, entrega e merge pela Delivery Queue, eu escolho o
  ponto de release e aviso os adotantes por ai-memory."
- Release: "Uma no fim do ciclo (Recommended)" — option text: "v0.127.0
  depois das quatro Specs; os adotantes recebem a revisão de skills e as
  correções juntas."

The cycle AR rules are the ones of 2026-10-27, about 11:40 BRT ("Mesmas
regras do AQ (Recommended)"), which carry the cycle AQ rules of 2026-10-26:
the Governed Paths each Spec declares, and still not this repository's
Makefile, `.github/workflows`, `go.mod`, `go.sum`, or lint and formatter
configuration. They carry the cycle AJ grant of 2026-10-19, at about 16:20
BRT, "Mexer em Governed Paths declarados" (option text: "'Concedo' para os
caminhos que cada Spec declarar (módulos Baseline, registros, skills, testes
governados; fora: Makefile, CI, go.mod, lint, que pedem autorização
específica).") and the 2026-09-30 grant of every skill ("considere
autorizado a ajustar todas as skills se necessário").

This record grants `implement`, `commit`, `push`, `pull_request` and `merge`.
It does not grant `release`: the release stays at the maintainer's
discretion, by the recorded criteria.

The governed set was measured on 2026-10-30 on the branch `docs/plan-0335`,
based on `35de7f74`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` was identical before
  and after, and the probe file never existed in the tree.
- It covered every file the Tasks declare under `interface:`, `creates:` or
  `instruction:`, the new ADR and every planning file of this Spec.
- The eight paths above are governed and listed in `paths`. The
  `domain-modeling` skill is governed too, but task_01 only reads it as an
  instruction and no Task changes it. Every other declared file is not
  governed, among them the production files and tests under
  `internal/config/` and `internal/cli/`, both user guide pages, the
  `skills/roundfix/references/` mirrors,
  `skills/testdata/owned-skill-versions.json`, `CONTEXT.md` and
  `docs/references/behavior-surfaces.json`.

## Why each governed path is needed

- `.agents/skills/roundfix/references/spec.md`, `deliver.md`, `archive.md`
  and `baseline.md` describe `spec check`, `deliver plan` and `status`,
  `history sanitize` and the history plan of `baseline update`, which stop
  refusing an absent Spec Root (`docs/agents/cli.md`).
- `.agents/skills/roundfix/references/follow-up.md` describes the accepted
  message types and the JSON fields this Spec pins.
- `.agents/skills/roundfix/references/modes.md` tells the agent which types
  it may write into the messages file.
- The two Roundfix `SKILL.md` files carry the owned skill's version, which
  rises with its content.

## Sanctioned regeneration

The owned-skill version record is rewritten only by its record command,
after `make skills-sync` copies the canonical skill to its mirror:

```yaml
command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions -count=1
outputs:
  - skills/testdata/owned-skill-versions.json
  - .agents/skills/roundfix/SKILL.md
  - skills/roundfix/SKILL.md
```

The Behavior Surface record is rewritten only by its record command, run
after the last guide edit:

```yaml
command: go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1
outputs:
  - docs/references/behavior-surfaces.json
```

No version field, fingerprint or digest is hand-edited. No test is removed
or renamed, so the repository coverage record is not re-recorded.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`, the lint,
  formatter or test-runner configuration, `.roundfixrc.yml`, any Baseline
  module, profile, template or snapshot asset, `docs/agents/`,
  `skills-lock.json`, any installed external skill, or the `### QA
  settlement` section of any skill.
- No change to `internal/cli/archive.go` or `internal/cli/cancel.go`, and no
  placeholder file written by any command.
- No file of another Spec's folder.
- No test, Verification command or QA row:
  - reads or writes the real `~/.roundfix`;
  - reaches the network;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository or `~/dev/skills`;
  - records an environment value.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
