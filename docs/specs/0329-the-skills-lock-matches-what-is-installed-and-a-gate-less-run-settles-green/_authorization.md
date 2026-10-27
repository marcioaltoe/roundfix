---
status: approved
granted: 2026-10-27
action: change the owned-skill bundle install so it rewrites the skills-lock.json hash of each installed owned skill, make every Task settle on the repository Verification with or without a QA gate, add the Integrated Repository Verification for a Run without a QA gate, and update the guides, the roundfix and write-tasks skills, the glossary and the Behavior Surface record
consuming: 0329-the-skills-lock-matches-what-is-installed-and-a-gate-less-run-settles-green
paths:
  - .agents/skills/roundfix/SKILL.md
  - skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/baseline.md
  - .agents/skills/roundfix/references/implement.md
  - .agents/skills/roundfix/references/setup.md
  - .agents/skills/write-tasks/SKILL.md
  - skills/write-tasks/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0329

On 2026-10-27, at about 11:40 BRT, the maintainer answered the cycle AR
questions (AskUserQuestion) the operator asked in the `/roundfix
--follow-up` mode, after an inventory with `--messages` and the Jev model
`typesafe/jev-1.13-20260917`. The sources include Quaestor's v0.125.0 proof,
ai-memory message `01a17d0a-cd60-7191-a724-2988affea2d0` from
`default/gesttione-solutions-quaestor` (2026-10-27 08:07 BRT). Verbatim:

- Scope: "A, B e C (Recommended)" — option text: "Lock + Clean vermelho;
  família do settle + aceite da QA; profile proof sob carga. Todas são
  health e duas se repetem em adotantes." This Spec is A.
- QA acceptance: "Casar por linha (Recommended)" — option text: "Uma
  declaração cobre cada criterion que nomeia; o accept confere linha a
  linha, não a contagem. write-prd/qa-gate passam a dizer isso." This answer
  governs Spec 0330.
- Authorization: "Mesmas regras do AQ (Recommended)" — option text: "Push de
  branches, PR dos planos, entrega e merge pela Delivery Queue, eu escolho o
  ponto de release e aviso os adotantes por ai-memory."

The cycle AQ rules are the ones of 2026-10-26, about 15:00 BRT, when the
maintainer wrote "Sim, pode fazer a higiene." and "Vamos seguir a sua
sugestão de sequencia para a spec e os demais itens. Talvez seja importante
usar o jev para validar.": the Governed Paths each Spec declares, and still
not this repository's Makefile, `.github/workflows`, `go.mod`, `go.sum`, or
lint and formatter configuration. They carry the cycle AJ grant of
2026-10-19, at about 16:20 BRT, "Mexer em Governed Paths declarados" (option
text: "'Concedo' para os caminhos que cada Spec declarar (módulos Baseline,
registros, skills, testes governados; fora: Makefile, CI, go.mod, lint, que
pedem autorização específica).") and the 2026-09-30 grants of every skill
("considere autorizado a ajustar todas as skills se necessário") and of
Baseline sources and guides ("Autorizar os dois").

This record grants `implement`, `commit`, `push`, `pull_request` and `merge`.
It does not grant `release`: the release point stays the maintainer's, by
the recorded criteria.

The governed set was measured on 2026-10-27 on the branch `docs/plan-0329`,
based on `5e9475bf`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` was identical before
  and after, and the probe file never existed in the tree.
- It covered every file the Tasks declare under `creates:`, `interface:` or
  `instruction:`, the new ADR, the two amended ADRs and every planning file
  of this Spec.
- The seven paths above are governed and changed by task_01, and are listed
  in `paths`. `.agents/skills/domain-modeling/SKILL.md` is governed but only
  read, so it is not listed. Every other declared file is not governed,
  among them the production and test files under `internal/baseline/`,
  `internal/cli/` and `internal/daemon/`, the four user guide pages, the
  mirrored references `skills/roundfix/references/baseline.md`,
  `skills/roundfix/references/implement.md` and
  `skills/roundfix/references/setup.md`,
  `skills/testdata/owned-skill-versions.json`, `CONTEXT.md`,
  `docs/references/behavior-surfaces.json` and the ADRs.

## Why each governed path is needed

- `.agents/skills/roundfix/SKILL.md`, its mirror and the canonical
  references `baseline.md`, `implement.md` and `setup.md` carry the Roundfix
  owned skill, whose Settlement Checks section and `implement`, `baseline
  update` and `skills install` references must describe the changed
  behavior (`docs/agents/cli.md`) and whose
  version rises with it.
- `.agents/skills/write-tasks/SKILL.md` and its mirror carry the authoring
  rule that today scopes a green repository Verification to gated graphs.

## Sanctioned regeneration

The owned-skill version record is rewritten only by its record command,
after `make skills-sync` copies the canonical skills to their mirrors:

```yaml
command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions -count=1
outputs:
  - skills/testdata/owned-skill-versions.json
  - .agents/skills/roundfix/SKILL.md
  - skills/roundfix/SKILL.md
  - .agents/skills/write-tasks/SKILL.md
  - skills/write-tasks/SKILL.md
```

The Behavior Surface record is rewritten only by its record command, run
after the last guide edit:

```yaml
command: go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1
outputs:
  - docs/references/behavior-surfaces.json
```

The repository coverage record is rewritten only by its record command,
after task_03 replaces the old gate-less settlement test (amended
2026-10-27 after the second delivery Run's settlement failed on
`TestCoverageEquivalence`):

```yaml
command: go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record -count=1
outputs:
  - docs/references/coverage-record.json
```

No version field, fingerprint or digest is hand-edited.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`, the lint,
  formatter or test-runner configuration, `.roundfixrc.yml`, any Baseline
  module, profile, template, Setup Snapshot, coverage or retention asset,
  `docs/agents/`, this repository's `skills-lock.json`, any installed
  external skill, or the `### QA settlement` section of any skill.
- No file of another Spec's folder.
- No test, Verification command or QA row:
  - reads or writes the real `~/.roundfix`;
  - reaches the network;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository or `~/dev/skills`;
  - records an environment value.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
