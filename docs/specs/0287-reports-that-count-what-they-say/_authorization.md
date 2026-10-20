---
status: approved
granted: 2026-10-20
action: print the number of packages the Linux Lane container selected on the lane's pass line, measure each Run's artifacts in the storage report under the Artifact Root the Run wrote them to, carry the list or paragraph a colon introduces into an Archive Record summary, and describe each change in the user guides, the Roundfix skill references and the glossary
consuming: 0287-reports-that-count-what-they-say
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/archive.md
  - .agents/skills/roundfix/references/deliver.md
  - .agents/skills/roundfix/references/storage.md
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0287

On 2026-10-20, at about 06:30 BRT, the maintainer approved cycle AK through
AskUserQuestion. The answers this Spec rests on were:

- next cycle "Duas Specs + um lote (Recomendado)": M (0285) → N (0286) → O
  (0287), one minor each;
- authorization "Sim, mesmas regras do AJ (Recomendado)": the Governed Paths
  each Spec declares, and still not the Makefile, `.github/workflows`,
  `go.mod`, `go.sum`, or lint and formatter configuration.

The cycle AJ rules that answer carries were approved on 2026-10-19, at about
16:20 BRT, among them "Mexer em Governed Paths declarados" (option text:
"'Concedo' para os caminhos que cada Spec declarar (módulos Baseline,
registros, skills, testes governados; fora: Makefile, CI, go.mod, lint, que
pedem autorização específica).") and "Merge squash + release + push". The
2026-09-30 skills grant ("considere autorizado a ajustar todas as skills se
necessário") also stands.

This record grants `implement`, `commit`, `push`, `pull_request` and `merge`.
It does not grant `release`: the operator releases. No Task, Verification
command or QA row pulls or builds an image, changes `.roundfixrc.yml`, or
reads the real `~/.roundfix`.

The governed set was measured on the authoring branch `docs/plan-0287`, based
on `9448aa64`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` was identical before
  and after, and the probe file never existed in the tree.
- It covered every file the Tasks declare, every planning file of this Spec,
  its three adopted Backlog Entries, the open Backlog Entry this planning
  change writes, `docs/references/skill-coverage.json` and `.roundfixrc.yml`.
- Six of those files are governed. Five are changed and listed above.
  `.roundfixrc.yml` is governed too, but no Task changes it, so it is not
  bounded.

## Why each governed path is unavoidable

- `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md` carry the
  Roundfix skill's version fields, raised with its content (task_01, task_02,
  task_03).
- `.agents/skills/roundfix/references/deliver.md` covers `deliver`, whose
  Linux Lane pass line now counts the packages that ran (task_01).
- `.agents/skills/roundfix/references/storage.md` covers `storage report`,
  which now measures the resolved Artifact Root (task_02).
- `.agents/skills/roundfix/references/archive.md` covers `archive` and
  `history sanitize`, whose record summary carries the list its colon
  introduces (task_03).

## What is not governed

- The production and test sources the Tasks change in `internal/cli`,
  `internal/store` and `internal/spec`.
- The skill reference mirrors under `skills/roundfix/references/`.
- `skills/testdata/owned-skill-versions.json` and
  `docs/references/behavior-surfaces.json`.
- `CONTEXT.md` and the user guides `docs/user-guide/commands/deliver.md`,
  `storage-report.md`, `archive.md` and `history.md`.
- Every planning file of this Spec, its adopted references, and the open
  Backlog Entry.

## Sanctioned regeneration

The repository-owned commands resolve the skill mirror, the skill version and
the Behavior Surface Record. This declaration records the regeneration that
follows the approved edits and adds no source path.

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
surface is added, so `docs/references/skill-coverage.json` does not change. No
record or mirror is hand-edited.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`, the lint,
  formatter or test-runner configuration, any Baseline module,
  `.roundfixrc.yml`, `docs/agents/setup-context.json` or the `### QA
  settlement` section of any skill.
- No change to the Linux Lane's image, init, mounts, package selection,
  network and pull policy, cap, outcome classification, blockers or failure
  lines; to the storage report's sections, columns or reconciliation; to what
  `roundfix gc sanitize` or Run Retention may remove; or to any Archive Record
  already written.
- The storage report stays read-only and creates no directory.
- No test, Verification command or QA row:
  - reads or writes the real `~/.roundfix`;
  - reaches the network or pulls or builds an image;
  - needs Docker, except the one QA row that runs the real lane only when the
    Run sandbox reaches a Docker daemon and is declared unreachable otherwise;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository;
  - records an environment value.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
