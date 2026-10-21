---
status: approved
granted: 2026-10-20
action: make the Delivery Queue's Pull Request Record Commit name the Spec by its number so its subject fits 100 characters for every slug, keep adopting a record commit made by hand under any subject, make a queue commit that a commit hook refuses park with the hook, the subject and the hand-commit recovery, and say so in the deliver guide, the Roundfix deliver reference and the glossary
consuming: 0290-the-pull-request-record-commit-fits-commitlint
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/deliver.md
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0290

On 2026-10-20, between about 18:15 and 18:35 BRT, the maintainer approved
cycle AM through AskUserQuestion. The answers this Spec rests on were:

- scope "R → S → T (Recomendado)" (option text: "Três Specs, um minor cada
  (v0.101–v0.103), nesta ordem. R primeiro porque já afeta o oraculum. Mesmas
  regras do AL."): R, this Spec, delivers first;
- this Spec's design "Número + falha nomeada (Recomendado)" (option text: "O
  assunto passa a ser 'docs: record pull request #N in the 0088 Archive
  Record'. Se o hook recusar mesmo assim, o motivo do park cita o hook de
  commit e o contorno, não um delivery-error genérico.");
- authorization: the same rules as cycle AL, that is the Governed Paths each
  Spec declares, and still not the Makefile, `.github/workflows`, `go.mod`,
  `go.sum`, or lint and formatter configuration.

The cycle AL rules that answer carries were approved on 2026-10-20, at about
13:10 BRT ("Mesmas regras do AK (Recomendado)"); they carry the cycle AJ grant
of 2026-10-19, at about 16:20 BRT, "Mexer em Governed Paths declarados"
(option text: "'Concedo' para os caminhos que cada Spec declarar (módulos
Baseline, registros, skills, testes governados; fora: Makefile, CI, go.mod,
lint, que pedem autorização específica).") and "Merge squash + release +
push". The 2026-09-30 skills grant ("considere autorizado a ajustar todas as
skills se necessário") also stands.

This record grants `implement`, `commit`, `push`, `pull_request` and `merge`.
It does not grant `release`: the operator releases. No Task, Verification
command or QA row reads or writes the real `~/.roundfix`, reaches the network,
reads another repository or runs `roundfix deliver` against a real forge.

The governed set was measured on the authoring branch `docs/plan-0290`, based
on `70c78d15`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` was identical before
  and after, and the probe file never existed in the tree.
- It covered every file the Tasks declare, every planning file of this Spec,
  its adopted reference, ADR-0302, the open Backlog Entry this planning change
  writes, `docs/references/skill-coverage.json` and `CONTEXT.md`.
- Three of those files are governed, and all three are changed and listed
  above.

## Why each governed path is unavoidable

- `.agents/skills/roundfix/references/deliver.md` covers every `deliver`
  surface, whose record commit and refusal reason change (task_01).
- `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md` carry the
  Roundfix skill's version fields, raised with its content (task_01).

## What is not governed

- The production and test sources task_02 changes in `internal/cli` and
  `internal/daemon`.
- The skill reference mirror `skills/roundfix/references/deliver.md`.
- `skills/testdata/owned-skill-versions.json` and
  `docs/references/behavior-surfaces.json`.
- `CONTEXT.md` and the user guide `docs/user-guide/commands/deliver.md`.
- Every planning file of this Spec, its adopted reference, ADR-0302 and the
  open Backlog Entry.

## Sanctioned regeneration

The repository-owned commands resolve the skill mirror, the skill version and
the Behavior Surface Record after task_01's edits. This declaration records
the regeneration that follows the approved edits and adds no source path.

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
surface is added and the skill reference changes, so
`docs/references/skill-coverage.json` does not change. No record or mirror is
hand-edited.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`, the lint,
  formatter or test-runner configuration, any Baseline module,
  `.roundfixrc.yml`, `docs/agents/setup-context.json`, `internal/delivery/` or
  the `### QA settlement` section of any skill.
- No change to the archive commit subject, the pull request title, the
  derived merge message, the record proof, the Delivery Engine, the blocker
  constants, the Park Class table, or the Daemon's Task commit refusal.
- The queue never passes `--no-verify` and never rewrites a refused commit.
- No test, Verification command or QA row:
  - reads or writes the real `~/.roundfix`;
  - reaches the network;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository;
  - records an environment value.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
