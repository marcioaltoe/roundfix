---
status: approved
granted: 2026-10-11
action: make the Archive Record name its delivery pull request through the Archive Command and a Delivery Queue Pull Request Record Commit, recover a squash-dropped source_revision from the pull request head ref in every reader, warn when an archive's source_revision is outside the default branch, read only real records from the archive root, and describe it in the skills, user guides and glossary
consuming: 0256-an-archive-record-that-survives-a-squash-merge
paths:
  - .agents/skills/archive-spec/SKILL.md
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/archive.md
  - .agents/skills/roundfix/references/deliver.md
  - .agents/skills/roundfix/references/review.md
  - .agents/skills/roundfix/references/runs.md
  - .agents/skills/roundfix/references/spec-delivery.md
  - skills/archive-spec/SKILL.md
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0256

On 2026-09-30 the maintainer asked for unattended work through every release
of the program. Of the skills the maintainer said "considere autorizado a
ajustar todas as skills se necessário". Of Baseline sources, guides and
`.roundfixrc.yml` the maintainer said "Autorizar os dois".

On 2026-10-11 the maintainer decided cycle AE. For the squash item the
decision was "PR + refs/pull, com aviso": the record always stores
`pull_request`, readers fetch `refs/pull/N/head` when the commit is missing,
`roundfix archive` warns when `source_revision` is not an ancestor of the
default branch, no new ref is pushed and archived records stay
byte-identical. For the review item the decision was that `review` treats as
an Archive Record only `NNNN-<slug>.md` with
`schema: roundfix/archive-record/v1`, and that every archive-root reader is
swept. To the request for grants, "Grants dos caminhos governados", the
maintainer answered "Concedo" for the Governed Paths each Spec declares.
The Makefile and `.github/workflows` are not covered for this Spec, and none
is changed.

No live provider or forge call is authorized by this record. Tests,
Verification and QA use temporary repositories, a local bare repository as
the forge and temporary homes. No adopter repository is read or written.

The governed set was measured with `GovernedPath` on the authoring branch at
`7a8b7f70`. The probe was a `go test -overlay` test in `internal/speccheck`
that wrote nothing to the repository. It ran against every file the Tasks
declare, and the nine paths above are governed.

## Why each governed path is unavoidable

- `.agents/skills/roundfix/references/archive.md`, `review.md`, `runs.md`,
  `deliver.md` and `spec-delivery.md` describe the commands whose behavior
  changes: the archive flag and warning, the review's record identity and
  recovery, the Run causes recovery, the Pull Request Record Commit and the
  hand delivery step. The CLI guide requires the skill to change in the same
  pull request, and the Skill Coverage Map names these references as the
  covering skill files.
- `.agents/skills/archive-spec/SKILL.md` gains the squash-merge section.
- `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md` and
  `skills/archive-spec/SKILL.md`: an owned skill's content changes only with
  its version, and the record command raises both version fields.

## What is not governed

These paths are ordinary: every Go file and test the Tasks declare under
`internal/spec`, `internal/cli`, `internal/delivery`, `internal/runcause`,
`internal/worktree` and `internal/speccheck`; `CONTEXT.md`; the user guides
under `docs/user-guide/commands/`; `docs/references/behavior-surfaces.json`;
`skills/testdata/owned-skill-versions.json`; and the reference mirrors under
`skills/roundfix/references/`.

## Sanctioned regeneration

```yaml
command: make skills-sync
```

`go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
raises and records the owned skills' versions, and
`go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`
re-records the Behavior Surface Record. Neither output is hand-edited.

## Limits

- No new dependency in `go.mod`, and no change to the Makefile, the lint,
  formatter or test-runner configuration, `.roundfixrc.yml`, the CI workflows
  or any Baseline module or derived Baseline file.
- No change to the `### QA settlement` section of any skill, to the governed
  archive cut in `internal/spec/archive.go`, to archived Specs, existing
  Archive Records or existing history entries.
- No test, Verification command or QA row opens a network connection, reads
  or writes the real `~/.roundfix`, or pushes a ref.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
