---
task: task_01
spec: 0323-the-status-the-park-table-and-write-tasks-say-what-to-do-next
status: pending
type: backend
complexity: medium
---

# Task 01: The follow-up status names the manual release classification

## Overview

In the Release Plan's `manual_classification_required` state,
`roundfix follow-up --status` prints an empty version on its `Next release`
line and on its maintainer line, and its JSON leaves `release.reason` null.
This Task fills `release.reason` with the number of commits that need a
manual release classification, gives the maintainer item the command that
classifies them, renders both without a blank field, documents the state in
the follow-up guide and the Roundfix `follow-up` reference, and records the
skill version and the Behavior Surface fingerprint.

## Requirements

1. MUST answer the Backlog Entry
   [`references/2026-10-26-follow-up-status-leaves-the-next-version-blank-under-manual-classification.md`](references/2026-10-26-follow-up-status-leaves-the-next-version-blank-under-manual-classification.md)
   of 2026-10-26 (`_prd.md` → Sources 1) with `_techspec.md` → Interfaces
   and API Contracts 1-2.
2. MUST, in `internal/cli/follow_up_status.go` `runFollowUpStatus`, set
   `release.reason` to `<n> commit(s) since <base> need a manual release
   classification`, `<n>` being `len(plan.Classification.BlockingCommits)`,
   and set the `release_classification` item's detail to `classify them with
   roundfix release plan --from <base> --to <target> --impact
   <none|patch|minor|major> --reason <text>`, `<target>` being
   `plan.Target.Name`.
3. MUST, in `renderFollowUpStatusText`, print
   `Next release: not proposed (manual_classification_required): <reason>`
   for that state and `  release classification: <detail>` for a
   `release_classification` item, keeping every other line's bytes.
4. MUST keep `release.proposed_version` the empty string and every JSON field
   and type unchanged.
5. MUST add `TestFollowUpStatusNamesTheManualClassification` to
   `internal/cli/follow_up_status_test.go` as `_techspec.md` → Testing
   Approach describes (one `feat:` and two `chore:` commits, so the count of
   blocking commits differs from the count of changes), and update
   `TestFollowUpStatusReportsTheNextRelease` only where the Testing Approach
   declares: a non-nil `reason` for `manual` alone, and `manual release
   classification` in its `manual` text check.
6. MUST add the guide paragraph of Interfaces to
   `docs/user-guide/commands/follow-up.md` after the
   `Next release: no unreleased changes` block, and the reference paragraph
   to the canonical `.agents/skills/roundfix/references/follow-up.md` after
   the paragraph ending "it never calls the judge.".
7. MUST run `make skills-sync`, then
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
   (it raises both version fields of the `roundfix` skill and records them),
   then `make skills-sync` again so the mirror matches.
8. MUST, as the Task's last step, after every guide and skill edit,
   re-record the Behavior Surface fingerprint ("re-record after the last
   guide edit") with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.
9. MUST NOT change `internal/releaseplan`, `roundfix release plan` output,
   any other release state's bytes, any other skill, `docs/agents/`, the
   Makefile or the CI workflows.

## Subtasks

- [ ] Fill the reason and the maintainer detail and render both lines.
- [ ] Add the new test and adjust the declared assertions of the existing one.
- [ ] Document the state, sync and record the skill, then re-record the Behavior Surface fingerprint last.

## Acceptance Criteria

- [ ] The manual state prints no blank version in text, and its JSON carries the reason and the classification command.
- [ ] The other release states keep their expectations.
- [ ] The guide, the canonical reference and its mirror describe the state; the skill version and the Behavior Surface record are current.

## Consumers

None. The status JSON keeps its fields and types and only two values of one
state change; no Go package outside `internal/cli` reads the follow-up
report, and the Roundfix skill reads the unchanged `kind` field.

## Context

- instruction: `internal/releaseplan/model.go`
- instruction: `internal/cli/releaseplan_command.go`
- instruction: `skills/skills_test.go`
- instruction: `internal/docscontract/skill_coverage_test.go`
- interface: `internal/cli/follow_up_status.go`
- interface: `internal/cli/follow_up_status_test.go`
- interface: `docs/user-guide/commands/follow-up.md`
- interface: `.agents/skills/roundfix/references/follow-up.md`
- interface: `skills/roundfix/references/follow-up.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`

## Verification

- `grep -q 'func TestFollowUpStatusNamesTheManualClassification' internal/cli/follow_up_status_test.go && go test -timeout 15m -count=1 -run '^TestFollowUpStatus' ./internal/cli && tr -s '[:space:]' ' ' < docs/user-guide/commands/follow-up.md | grep -qF -- 'Next release: not proposed (manual_classification_required): <n> commit(s) since <base> need a manual release classification' && for ref in .agents/skills/roundfix/references/follow-up.md skills/roundfix/references/follow-up.md; do tr -s '[:space:]' ' ' < "$ref" | grep -q -- 'the .release_classification. item of .maintainer. names the' || { printf 'missing manual classification text: %s\n' "$ref" >&2; exit 1; }; done && go test -timeout 10m -count=1 ./skills && go test -timeout 15m -count=1 -tags docscontract -run '^TestTheSkillCoverageMapIsCurrent$' ./internal/docscontract` — expected: exit 0. Before this Task the new test does not exist, so the first check fails; afterwards the new and existing status tests pass, the guide and both reference copies describe the state, the owned-skill version is recorded and the Behavior Surface record is current.

## References

- `_prd.md` → Sources 1; Goals 1-2; Core Feature 1; Success Metrics 1-2
- `_techspec.md` → Current behavior 1-3; Interfaces; API Contracts 1-2; Invariants 1-2; Testing Approach; Build Order 1
