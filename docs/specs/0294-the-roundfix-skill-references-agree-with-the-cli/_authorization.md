---
status: approved
granted: 2026-10-21
action: make the Roundfix skill's review, events and deliver references, the review and deliver user guides and the Run Event Stream glossary entry state what the CLI does, and add documentation contract tests that read the reviewer's Delivery Conventions and the Run Event Stream categories from the code so the documents cannot drift from them again
consuming: 0294-the-roundfix-skill-references-agree-with-the-cli
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/deliver.md
  - .agents/skills/roundfix/references/events.md
  - .agents/skills/roundfix/references/review.md
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0294

On 2026-10-21, at about 01:00 BRT, the maintainer approved cycle AN through
AskUserQuestion. The answers this Spec rests on were:

- scope "V → W → X (Recomendado)" (option text: "Três Specs: v0.105 a
  v0.107. Y fica no backlog como entrada aberta; são só 3 linhas de ruído no
  reconcile e não custam nada."): V, this Spec, delivers first of the three,
  after the cycle AM Specs;
- authorization: the same rules as cycle AM, that is the Governed Paths each
  Spec declares, and still not the Makefile, `.github/workflows`, `go.mod`,
  `go.sum`, or lint and formatter configuration.

The cycle AM rules that answer carries were approved on 2026-10-20, between
about 18:15 and 18:35 BRT ("the same rules as cycle AL"); they carry the cycle
AJ grant of 2026-10-19, at about 16:20 BRT, "Mexer em Governed Paths
declarados" (option text: "'Concedo' para os caminhos que cada Spec declarar
(módulos Baseline, registros, skills, testes governados; fora: Makefile, CI,
go.mod, lint, que pedem autorização específica).") and "Merge squash +
release + push". The 2026-09-30 skills grant ("considere autorizado a ajustar
todas as skills se necessário") also stands.

This record grants `implement`, `commit`, `push`, `pull_request` and `merge`.
It does not grant `release`: the operator releases. No Task, Verification
command or QA row reads or writes the real `~/.roundfix`, reaches the network,
reads another repository or sends anything to a provider, the judge or a
reviewer.

The governed set was measured on the authoring branch `docs/plan-0294`, based
on `8fb2da1f`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` was identical before
  and after, and the probe file never existed in the tree.
- It covered every file the Tasks declare or read, every planning file of
  this Spec and its adopted reference.
- Five of those files are governed, and all five are changed and listed
  above.

## Why each governed path is unavoidable

- `.agents/skills/roundfix/references/review.md` holds the stale C3 and C5
  conventions and the C5 eligibility sentence (task_01).
- `.agents/skills/roundfix/references/events.md` holds the category list that
  leaves out `usage` (task_01).
- `.agents/skills/roundfix/references/deliver.md` names an archived `_prd.md`
  where the queue accepts an Archive Record (task_01).
- `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md` carry the
  Roundfix skill's version fields, raised with its content (task_01).

## What is not governed

- The documentation contract tests task_02 changes in
  `internal/docscontract`.
- The skill reference mirrors under `skills/roundfix/references/`.
- `skills/testdata/owned-skill-versions.json` and
  `docs/references/behavior-surfaces.json`.
- `CONTEXT.md` and the user guides `docs/user-guide/commands/review.md` and
  `docs/user-guide/commands/deliver.md`.
- Every planning file of this Spec and its adopted reference.

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
surface or covered source changes, so `docs/references/skill-coverage.json`
does not change. No record or mirror is hand-edited.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`, the lint,
  formatter or test-runner configuration, any Baseline module,
  `.roundfixrc.yml`, `docs/agents/setup-context.json` or the `### QA
  settlement` section of any skill.
- No production Go source changes: the reviewer's conventions, the review
  validation, the Run Event Stream, the `events` help text and the Delivery
  Queue keep their behavior.
- No test, Verification command or QA row:
  - reads or writes the real `~/.roundfix`;
  - reaches the network;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository;
  - records an environment value.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
