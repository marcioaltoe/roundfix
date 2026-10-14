---
status: approved
granted: 2026-10-13
action: make the History Sanitize Command and the history section of roundfix baseline update write Archive Records from Legacy Archive Folders that keep a readable summary, record legacy override gaps as legacy not recorded, fall back for sources and the creation date, and rebase the links of the outcome and of promoted copies, with the history user guide, the Roundfix skill's archive and baseline references and the glossary describing them
consuming: 0271-history-sanitize-records-keep-their-summary-links-and-dates
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/archive.md
  - .agents/skills/roundfix/references/baseline.md
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0271

On 2026-10-13 the maintainer approved cycle AG through the AskUserQuestion of
about 17:20 BRT (six Specs, order B → A → C → D → E → F) and then "Pode
iniciar" (about 17:25 BRT), to start them under the same operating rules as
cycles AE and AF, including the 2026-10-11 answer "Concedo" for the Governed
Paths each Spec declares. The Makefile and the CI workflows are not covered in
AG. For this Spec, item F, the decision is: records already written stay
byte-identical, and only generation from Legacy Archive Folders changes. On
2026-09-30 the maintainer had granted the skills ("considere autorizado a
ajustar todas as skills se necessário").

The governed set was measured on the authoring branch at `9963ffd3` with
`GovernedPath`, through a `go test -overlay` probe in `internal/speccheck`
that wrote nothing to the repository, against every file the Tasks declare and
every planning file. Four of them are governed.

## Why each governed path is unavoidable

- `.agents/skills/roundfix/references/archive.md` covers the
  `command history sanitize` surface in the Skill Coverage Map, whose source
  `internal/cli/history.go` task_02 changes; task_04 describes the new
  readings there.
- `.agents/skills/roundfix/references/baseline.md` covers the
  `command baseline update` surface, whose source
  `internal/cli/baseline_update_history.go` task_02 changes; task_04 says the
  update's conversions read the same way.
- `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md` carry the
  skill's version fields, raised with its content.

## What is not governed

The Go sources and tests under `internal/spec` and `internal/cli`; the mirrors
`skills/roundfix/references/archive.md` and
`skills/roundfix/references/baseline.md`;
`skills/testdata/owned-skill-versions.json`;
`docs/references/behavior-surfaces.json`; `CONTEXT.md`;
`docs/user-guide/commands/history.md`; the ADR
`docs/adr/0277-a-legacy-archive-folder-conversion-keeps-the-summary-links-and-dates-its-folder-held.md`;
and the new Backlog Entry
`docs/backlog/2026-10-13-an-archive-record-outcome-keeps-links-relative-to-the-removed-spec-folder.md`.

## Sanctioned regeneration

The repository-owned commands resolve the skill mirrors, the skill version and
the Behavior Surface Record. This declaration records the regeneration that
follows the approved edits and adds no source path.

```yaml
command: make skills-sync
```

`go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
records the raised skill version, and
`go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`
re-records the fingerprint of the `history` user guide. No Baseline module
changes, so no module version is recorded, and `make baseline-digests`
rewrote nothing when measured after both skill edits on a scratch copy on
2026-10-13. No record or mirror is hand-edited.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`, the lint,
  formatter or test-runner configuration, any Baseline module, the review
  conventions, or the `### QA settlement` section of any skill.
- No change to an Archive Record already written, to the Archive Command or to
  the Delivery Queue's exact-archive proof.
- No test, Verification command or QA row reads or writes the real
  `~/.roundfix`, reaches the network, or reads another repository.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
