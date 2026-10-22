---
status: approved
granted: 2026-10-22
action: make the Spec Consistency Check report a pytest node id that another collection argument of the same invocation contains as SC-VERIFY-SUBSUMED-SELECTOR, make every QA Report the Daemon closes without an Agent turn record the commit it audited as build, and say so in the spec guide, the Roundfix spec reference, the write-tasks and qa-gate skills, the Context-driven development guide and the glossary
consuming: 0302-a-subsumed-test-selector-is-flagged-and-a-qa-refusal-names-its-commit
paths:
  - .agents/skills/qa-gate/SKILL.md
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/spec.md
  - .agents/skills/write-tasks/SKILL.md
  - internal/speccheck/coherence.go
  - internal/speccheck/constraints.go
  - skills/qa-gate/SKILL.md
  - skills/roundfix/SKILL.md
  - skills/write-tasks/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0302

On 2026-10-22 the maintainer approved cycle AP through AskUserQuestion; the
questions of 15:25 BRT were answered by about 18:30 BRT. The answers this
Spec rests on were:

- scope (multi-select, all four) "Correções de adotantes (Recomendado)",
  "Ruído C2 e suite docscontract", "Modos follow-up A e B", "Preset de
  baseline Python"; this Spec is one of the adopter fixes of batch 1 and
  delivers second, after Spec 0301;
- order "Correções, follow-up, Python (Recomendado)" (option text: "Segue a
  ordem que você definiu para o follow-up: primeiro o que trava outros
  projetos (Light Tier no Fluxus), depois o bom funcionamento, depois as
  features.");
- authorization: the same as cycle AO, that is the Governed Paths each Spec
  declares, and still not this repository's Makefile, `.github/workflows`,
  `go.mod`, `go.sum`, or lint and formatter configuration.

This Spec's design is the operator's cycle AP briefing for Spec 0302: a
structural check over the written command rather than a richer vacuous
diagnosis, language-agnostic in its naming, with the rule added to the
`write-tasks` skill; and a refusal report that records `build` as a full
report does, with reports already written keeping their bytes.

The cycle AO rules that answer carries were approved on 2026-10-22, at about
00:55 BRT ("the same rules as cycle AN"), which carry the cycle AN approval
of 2026-10-21, at about 01:00 BRT ("the same rules as cycle AM"), the cycle
AM approval of 2026-10-20, at about 18:15–18:35 BRT ("the same rules as cycle
AL"), the cycle AL approval of 2026-10-20, at about 13:10 BRT ("Mesmas regras
do AK (Recomendado)"), and the cycle AJ grant of 2026-10-19, at about 16:20
BRT, "Mexer em Governed Paths declarados" (option text: "'Concedo' para os
caminhos que cada Spec declarar (módulos Baseline, registros, skills, testes
governados; fora: Makefile, CI, go.mod, lint, que pedem autorização
específica).") and "Merge squash + release + push". The 2026-09-30 skills
grant ("considere autorizado a ajustar todas as skills se necessário") also
stands.

This record grants `implement`, `commit`, `push`, `pull_request` and `merge`.
It does not grant `release`: the operator releases. No Task, Verification
command or QA row runs pytest, reads or writes the real `~/.roundfix`,
reaches the network or reads another repository.

The governed set was measured on the authoring branch `docs/plan-0302`, based
on `42f523c2`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` was identical before
  and after, and the probe file never existed in the tree.
- It covered every file the Tasks declare, every planning file of this Spec,
  its two adopted references and the Backlog Entry path one of them left,
  ADR-0315, `docs/references/skill-coverage.json` and `CONTEXT.md`.
- Nine of those files are governed, and all nine are changed and listed
  above.

## Why each governed path is unavoidable

- `internal/speccheck/constraints.go` runs the Verification detectors of
  `Check`, and gains the call of the new detector (task_01).
- `internal/speccheck/coherence.go` holds `stagedDetectors`, which places
  each code at the stage that can decide it, and gains the new code (task_01).
- `.agents/skills/roundfix/references/spec.md` lists every Spec Consistency
  Check code and gains the new one (task_02).
- `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md` carry the
  Roundfix skill's version fields, raised with its content (task_02).
- `.agents/skills/write-tasks/SKILL.md` is the authoring rule's home, and it
  and `skills/write-tasks/SKILL.md`, its mirror, carry the version fields
  raised with it (task_02).
- `.agents/skills/qa-gate/SKILL.md` holds the refusal template, and it and
  `skills/qa-gate/SKILL.md`, its mirror, carry the version fields raised
  with it (task_04).

## What is not governed

- The production and test sources task_01 creates in `internal/speccheck`
  and task_03 changes or creates in `internal/spec` and `internal/daemon`.
- The skill reference mirror `skills/roundfix/references/spec.md`.
- `skills/testdata/owned-skill-versions.json` and
  `docs/references/behavior-surfaces.json`.
- `CONTEXT.md` and the user guides `docs/user-guide/commands/spec.md` and
  `docs/user-guide/context-driven-development.md`.
- Every planning file of this Spec, its adopted references and ADR-0315.

## Sanctioned regeneration

The repository-owned commands resolve the skill mirrors, the skill versions
and the Behavior Surface Record after the edits of task_02 and task_04. This
declaration records the regeneration that follows the approved edits and adds
no source path.

```yaml
command: make skills-sync
```

```yaml
command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
outputs:
  - skills/testdata/owned-skill-versions.json
```

```yaml
command: go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1
outputs:
  - docs/references/behavior-surfaces.json
```

No Baseline module changes, so no module version is recorded. An owned skill
edit leaves the derived Baseline artifacts byte-identical, so `make
baseline-digests` has nothing to rewrite. No test is removed or renamed and no
platform set changes, so the coverage record is not re-recorded. No command
surface is added, so `docs/references/skill-coverage.json` does not change.
No record or mirror is hand-edited.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`, the lint,
  formatter or test-runner configuration, any Baseline module,
  `.roundfixrc.yml`, `docs/agents/setup-context.json`, the Run Database
  schema or the `### QA settlement` section of any skill.
- No change to the Verification prober, the `spec check` command's output
  shape, any other detector, `spec.QAReport` or its readers, the report shape
  detector, or the full QA Report template.
- No test, Verification command or QA row:
  - runs pytest or another External Test Tool;
  - reads or writes the real `~/.roundfix`;
  - reaches the network;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository;
  - records an environment value.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
