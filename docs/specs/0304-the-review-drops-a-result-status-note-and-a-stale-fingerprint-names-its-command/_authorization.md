---
status: approved
granted: 2026-10-22
action: reword Delivery Convention C2 as version 4, make the pre-PR review dismiss a finding that names a Result Status Note against a completed non-QA Task by validation without the sealed validator while every other finding keeps the validator's fail-closed path, make a stale Behavior Surface fingerprint problem name its re-record command, and say so in the review guide, the Roundfix review reference and the glossary
consuming: 0304-the-review-drops-a-result-status-note-and-a-stale-fingerprint-names-its-command
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/review.md
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0304

On 2026-10-22 the maintainer approved cycle AP through AskUserQuestion; the
questions of about 15:25 BRT were answered by about 18:30 BRT. The answers
this Spec rests on were:

- scope (multi-select, all four): "Correções de adotantes (Recomendado)",
  "Ruído C2 e suite docscontract", "Modos follow-up A e B", "Preset de
  baseline Python"; this Spec is the C2 noise item, with the stale
  fingerprint message the operator grouped with it;
- order: "Correções, follow-up, Python (Recomendado)" (option text: "Segue a
  ordem que você definiu para o follow-up: primeiro o que trava outros
  projetos (Light Tier no Fluxus), depois o bom funcionamento, depois as
  features.");
- authorization: the same as cycle AO, that is the Governed Paths each Spec
  declares, and still not this repository's Makefile, `.github/workflows`,
  `go.mod`, `go.sum`, or lint and formatter configuration.

The cycle AO rules that answer carries were approved on 2026-10-22, at about
00:55 BRT ("the same rules as cycle AN"), which carry the cycle AN approval of
2026-10-21, at about 01:00 BRT, the cycle AM, AL and AK approvals, and the
cycle AJ grant of 2026-10-19, at about 16:20 BRT, "Mexer em Governed Paths
declarados" (option text: "'Concedo' para os caminhos que cada Spec declarar
(módulos Baseline, registros, skills, testes governados; fora: Makefile, CI,
go.mod, lint, que pedem autorização específica).") and "Merge squash +
release + push". The 2026-09-30 skills grant ("considere autorizado a ajustar
todas as skills se necessário") also stands.

This record grants `implement`, `commit`, `push`, `pull_request` and `merge`.
It does not grant `release`: the operator releases. No Task, Verification
command or QA row reads or writes the real `~/.roundfix`, runs `roundfix
review` against a real provider, sends a request to a provider, the judge or a
reviewer, reaches the network or reads another repository.

The governed set was measured on the authoring branch `docs/plan-0304`, based
on `42f523c2`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` was identical before
  and after, and the probe file never existed in the tree.
- It covered every file the Tasks declare, every planning file of this Spec,
  its two adopted references, ADR-0317, `docs/references/skill-coverage.json`,
  `CONTEXT.md` and `.roundfixrc.yml`.
- Four of those files are governed. Three are changed and listed above;
  `.roundfixrc.yml` is only read, by task_02's contract test, and is not
  changed.

## Why each governed path is unavoidable

- `.agents/skills/roundfix/references/review.md` carries the Delivery
  Conventions verbatim, which the documentation contract holds equal to the
  code, and gains the dismissal (task_01).
- `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md` carry the
  Roundfix skill's version fields, raised with its content (task_01).

## What is not governed

- The production and test sources task_01 changes or creates in
  `internal/cli`, and the test file task_02 changes in
  `internal/docscontract`.
- The skill reference mirror `skills/roundfix/references/review.md`.
- `skills/testdata/owned-skill-versions.json` and
  `docs/references/behavior-surfaces.json`.
- `CONTEXT.md` and the user guide `docs/user-guide/commands/review.md`.
- Every planning file of this Spec, its adopted references and ADR-0317.

## Sanctioned regeneration

The repository-owned commands resolve the skill mirrors, the skill versions
and the Behavior Surface Record after task_01's edits. This declaration
records the regeneration that follows the approved edits and adds no source
path.

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
  `.roundfixrc.yml`, `docs/agents/setup-context.json` or the `### QA
  settlement` section of any skill.
- No change to the Agent prompt, the `implement-task` skill, the review
  record schema, anchor validation, the validator's prompt, grammar or
  fail-closed handling, or `review dispose`.
- No test, Verification command or QA row:
  - reads or writes the real `~/.roundfix`;
  - reaches the network;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository;
  - records an environment value.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
