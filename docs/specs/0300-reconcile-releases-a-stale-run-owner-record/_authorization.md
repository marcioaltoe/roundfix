---
status: approved
granted: 2026-10-22
action: make the owner proof tell a proven start identity mismatch from an unreadable identity, make the Reconcile Command report a terminal Run whose owner PID now belongs to another process as a Stale Run Owner Record and release that record with --apply without sending any signal, keep every unreadable or unproven identity preserved and Force Stop unchanged, and say so in the reconcile guide, the Roundfix reconcile reference and the glossary
consuming: 0300-reconcile-releases-a-stale-run-owner-record
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/reconcile.md
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0300

On 2026-10-22, at about 00:55 BRT, the maintainer approved cycle AO through
AskUserQuestion. The answers this Spec rests on were:

- scope (multi-select, all four, in this order) "Z: portão e corretiva",
  "Espera de checks", "Retomar item arquivado", "Y: reconcile PID": one Spec
  and one minor each, v0.108 to v0.111; this Spec is Y and delivers fourth;
- authorization: the same rules as cycle AN, that is the Governed Paths each
  Spec declares, and still not the Makefile, `.github/workflows`, `go.mod`,
  `go.sum`, or lint and formatter configuration.

This Spec's design is the adopted Backlog Entry's Expected, which the
maintainer kept: a proven identity mismatch never signals the process now
using the PID, Force Stop's rule stays, reconcile classifies the record as
releasable and `--apply` releases it without sending any signal, and an
unreadable identity stays preserved.

The cycle AN rules that answer carries were approved on 2026-10-21, at about
01:00 BRT ("the same rules as cycle AM"), which carry the cycle AM approval of
2026-10-20, at about 18:15–18:35 BRT ("the same rules as cycle AL"), the
cycle AL approval of 2026-10-20, at about 13:10 BRT ("Mesmas regras do AK
(Recomendado)"), and the cycle AJ grant of 2026-10-19, at about 16:20 BRT,
"Mexer em Governed Paths declarados" (option text: "'Concedo' para os
caminhos que cada Spec declarar (módulos Baseline, registros, skills, testes
governados; fora: Makefile, CI, go.mod, lint, que pedem autorização
específica).") and "Merge squash + release + push". The 2026-09-30 skills
grant ("considere autorizado a ajustar todas as skills se necessário") also
stands.

This record grants `implement`, `commit`, `push`, `pull_request` and `merge`.
It does not grant `release`: the operator releases. No Task, Verification
command or QA row reads or writes the real `~/.roundfix`, runs `roundfix
reconcile --apply` outside a temporary Roundfix Home, sends a signal to a
process it did not start, reaches the network or reads another repository.

The governed set was measured on the authoring branch `docs/plan-0300`, based
on `34b4e246`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` was identical before
  and after, and the probe file never existed in the tree.
- It covered every file the Tasks declare, every planning file of this Spec,
  its adopted reference and the Backlog Entry path it left, ADR-0313,
  `docs/references/skill-coverage.json` and `CONTEXT.md`.
- Three of those files are governed, and all three are changed and listed
  above.

## Why each governed path is unavoidable

- `.agents/skills/roundfix/references/reconcile.md` describes every
  `reconcile` candidate kind and what `--apply` releases, and gains the owner
  record (task_01).
- `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md` carry the
  Roundfix skill's version fields, raised with its content (task_01).

## What is not governed

- The production and test sources task_02 changes or creates in
  `internal/store` and `internal/cli`.
- The skill reference mirror `skills/roundfix/references/reconcile.md`.
- `skills/testdata/owned-skill-versions.json` and
  `docs/references/behavior-surfaces.json`.
- `CONTEXT.md` and the user guide `docs/user-guide/commands/reconcile.md`.
- Every planning file of this Spec, its adopted reference and ADR-0313.

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
  `.roundfixrc.yml`, `docs/agents/setup-context.json`, the Run Database
  schema or the `### QA settlement` section of any skill.
- No change to Force Stop, the Proven Signal operations, the Delivery Queue
  owner proof, `doctor`, or the classification of an Active Run, a live
  matching owner, an unreadable identity or a Run recorded without an owner
  identity.
- No test, Verification command or QA row:
  - reads or writes the real `~/.roundfix`;
  - signals a process other than an unreaped child it started;
  - reaches the network;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository;
  - records an environment value.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
