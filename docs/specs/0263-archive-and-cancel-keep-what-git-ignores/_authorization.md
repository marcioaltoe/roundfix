---
status: approved
granted: 2026-10-12
action: make roundfix archive and roundfix cancel refuse while Git ignores a file under the Spec folder, listing each file and writing nothing, and make a failed final folder removal keep every write, exit 1 and print the command that finishes it, with the help, the user guides, the Roundfix and archive-spec skills and the glossary describing them
consuming: 0263-archive-and-cancel-keep-what-git-ignores
paths:
  - .agents/skills/archive-spec/SKILL.md
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/archive.md
  - internal/spec/archive.go
  - skills/archive-spec/SKILL.md
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0263

On 2026-10-12 the maintainer approved cycle AF, answering "Sim, como no AE"
to "As regras de operação do AE valem para o AF?". The cycle AE rules
therefore apply, including the 2026-10-11 answer "Concedo" for the Governed
Paths each Spec declares; the Makefile and the CI workflows are not covered in
AF. For this Spec, item B, the decisions are: `roundfix archive` and
`roundfix cancel` refuse, list each Git-ignored file under the Spec folder
from `git status` with `--ignored`, and write nothing; a failed final folder
removal keeps forward recovery and prints the exact command that finishes it;
symbolic-link destination parents stay accepted. On 2026-09-30 the maintainer
had granted the skills ("considere autorizado a ajustar todas as skills se
necessário").

The governed set was measured on the authoring branch at `94fbf85b` with
`GovernedPath`, through a `go test -overlay` probe in `internal/speccheck`
that wrote nothing to the repository, against every file the Tasks declare
and every planning file. Six of them are governed.

## Why each governed path is unavoidable

- `internal/spec/archive.go` holds `Archive`, whose final folder removal
  task_02 makes return a typed error with its filled result.
- `.agents/skills/roundfix/references/archive.md` describes the Archive and
  Cancel Commands, and `.agents/skills/archive-spec/SKILL.md` the archive
  workflow; task_03 adds the refusal and the recovery to both.
- `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`,
  `skills/archive-spec/SKILL.md` and the canonical archive-spec `SKILL.md`
  carry the skills' version fields, raised with their content.

## What is not governed

The Go sources and tests the Tasks declare under `internal/cli` and
`internal/spec` other than `internal/spec/archive.go`; the mirror
`skills/roundfix/references/archive.md`;
`skills/testdata/owned-skill-versions.json`;
`docs/references/skill-coverage.json` and
`docs/references/behavior-surfaces.json`; `CONTEXT.md`; the user guides under
`docs/user-guide/commands/`; and
`docs/adr/0269-archive-and-cancel-refuse-ignored-spec-files-and-finish-a-failed-removal-forward.md`.

## Sanctioned regeneration

The repository-owned commands resolve the skill mirrors, the skill versions
and the Behavior Surface Record. This declaration records the regeneration
that follows the approved edits and adds no source path.

```yaml
command: make skills-sync
```

`go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
records the raised skill versions, and
`go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`
re-records the fingerprints of the `archive` and `cancel` help and user-guide
surfaces. No Baseline module changes, so no module version is recorded, and
`make baseline-digests` rewrote nothing when measured after a skill edit on
2026-10-12. No record or mirror is hand-edited.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`, the lint,
  formatter or test-runner configuration, any Baseline module, or the
  `### QA settlement` section of any skill.
- No change to the symbolic-link destination parent checks.
- No test, Verification command or QA row reads or writes the real
  `~/.roundfix` or reaches the network.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
