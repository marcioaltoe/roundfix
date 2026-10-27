---
task: task_02
spec: 0327-a-spec-commits-the-executables-it-names-and-archives-a-declined-qa-gate
status: pending
type: backend
complexity: medium
---

# Task 02: Archive accepts a Declared QA Decline without an override

## Overview

`roundfix archive` reads a QA Report for every Spec with a Task Graph, so a
Spec whose graph declares `qa: declined` with a `qa_reason` is refused with
`no QA Report found` and can only archive with `--qa-override`, against
ADR-0154 and the Baseline's docs-layout guide. This Task makes archive accept
such a Spec once every Task is completed, refuses the override for it, and
writes the rule into the archive usage, the Roundfix `archive` and
`spec-delivery` references, the `archive-spec` skill, the `archive` guide and
`CONTEXT.md`.

## Requirements

1. MUST answer the Backlog Entry
   [`references/2026-10-26-a-spec-cannot-commit-its-executables-settle-exits-128-and-archive-refuses-a-declined-qa.md`](references/2026-10-26-a-spec-cannot-commit-its-executables-settle-exits-128-and-archive-refuses-a-declined-qa.md)
   of 2026-10-26, defect (b), by implementing API Contracts 5 and 8 of
   `_techspec.md`.
2. MUST change `Archive` in `internal/spec/archive.go` as API Contract 5
   states: without an override, a graph whose `QADeclined` is true archives
   when every Task is completed and reads no QA Report; with an override, a
   declined graph whose Tasks are all completed is refused with the existing
   `QA archive override is not allowed because the Spec already qualifies for
   normal archive`. Every other graph keeps today's path, and the record is
   the one `BuildArchiveRecord` already writes (`no-qa`, no override fields).
3. MUST add `internal/cli/archive_declined_qa_test.go`, through CLI dispatch
   on committed temporary repositories built with the implement workspace
   fixture and `declareImplementQADeclined`, with:
   - `TestArchiveAcceptsADeclaredQADecline`: one completed Task, no QA
     Report; exit `0`, stdout starts `archived <slug> -> `, and the record
     carries `disposition: no-qa` and no `qa_override` line;
   - `TestArchiveRefusesAnOverrideForAQualifyingDeclaredQADecline`: the same
     Spec with `--qa-override --approval maintainer --reason <text>`; exit
     `2` with the already-qualifies refusal and the Spec folder in place;
   - `TestArchiveRefusesADeclaredQADeclineWithAPendingTask`: a second Task
     `pending`; exit `2` naming `Task "task_02" is "pending"`.
4. MUST change `archiveUsage` in `internal/cli/archive.go` as API Contract 8
   states, and the usage pinned in `TestArchiveHelpNamesThePullRequestOption`
   (`internal/cli/archive_pull_request_test.go`) to match; nothing else in
   that test file.
5. MUST write API Contract 8's sentence, verbatim, with the override note,
   into `.agents/skills/roundfix/references/archive.md` (beside the
   eligibility paragraph), `.agents/skills/roundfix/references/spec-delivery.md`
   (its Advance step), a new `### Declared QA decline` section of
   `.agents/skills/archive-spec/SKILL.md` placed outside `### QA settlement`,
   and `docs/user-guide/commands/archive.md`; then run `make skills-sync` and
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
   which raises the version pairs of `archive-spec` and `roundfix`.
6. MUST add **Declared QA Decline** to `CONTEXT.md` through `domain-modeling`
   (a Task Graph's `qa: declined` with a non-empty `qa_reason`: the Spec has
   no QA Task, archive accepts it once every Task is completed with the
   disposition `no-qa`, and an override is refused for it; ADR-0154, with an
   `_Avoid_` line) and revise **Archive Command** with API Contract 8's
   sentence, keeping its `_Avoid_` line.
7. MUST, as the last step after the last guide edit, re-record the Behavior
   Surface fingerprint with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.
   Hand-edit no mirror, version field or record.
8. MUST NOT change the Archive Record schema, `BuildArchiveRecord`, the
   override's fields, any text inside a `### QA settlement` section, settle,
   the Daemon, the Makefile or the CI workflows.

## Subtasks

- [ ] Let a declined graph archive on Task completion and refuse the override for it.
- [ ] Add the three CLI tests and update the usage and its pin.
- [ ] Write the rule into the references, the `archive-spec` skill, the guide and the glossary; sync and record; re-record the Behavior Surface last.

## Acceptance Criteria

- [ ] A declined Spec with every Task completed archives with `no-qa` and no override; a pending Task still refuses.
- [ ] `--qa-override` on such a Spec is refused as already qualifying.
- [ ] The usage, references, skill, guide and glossary carry the rule, and the owned-skill and Behavior Surface records are current.

## Consumers

- Archive eligibility of a Spec Task Graph — consumers: `./internal/cli`, `./internal/spec`

## Context

- instruction: `docs/adr/0154-a-qa-archive-override-records-user-authority-not-a-pass.md`
- instruction: `.agents/skills/domain-modeling/SKILL.md`
- instruction: `internal/cli/archive_output_test.go`
- interface: `internal/spec/archive.go`
- interface: `internal/cli/archive.go`
- interface: `internal/cli/archive_pull_request_test.go`
- creates: `internal/cli/archive_declined_qa_test.go`
- interface: `.agents/skills/archive-spec/SKILL.md`
- interface: `skills/archive-spec/SKILL.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `.agents/skills/roundfix/references/archive.md`
- interface: `.agents/skills/roundfix/references/spec-delivery.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/roundfix/references/archive.md`
- interface: `skills/roundfix/references/spec-delivery.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/user-guide/commands/archive.md`
- interface: `CONTEXT.md`
- interface: `docs/references/behavior-surfaces.json`

## Verification

- `out="$(go test -count=1 -timeout 15m -v -run '^(TestArchiveAcceptsADeclaredQADecline|TestArchiveRefusesAnOverrideForAQualifyingDeclaredQADecline|TestArchiveRefusesADeclaredQADeclineWithAPendingTask|TestArchiveHelpNamesThePullRequestOption)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestArchiveAcceptsADeclaredQADecline TestArchiveRefusesAnOverrideForAQualifyingDeclaredQADecline TestArchiveRefusesADeclaredQADeclineWithAPendingTask TestArchiveHelpNamesThePullRequestOption; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < internal/cli/archive.go | grep -qF -- 'the Task Graph declares qa: declined with a qa_reason' || { echo 'archive usage does not name the decline' >&2; exit 1; }; go test -count=1 -timeout 15m ./internal/spec` — expected: exit 0. Before this Task the three decline tests do not exist and the usage does not name the decline, so the first pass check fails. Afterwards the decline archives, the override and the pending Task are refused, the usage pin matches, and the `internal/spec` archive suites still pass.
- `for file in .agents/skills/roundfix/references/archive.md skills/roundfix/references/archive.md .agents/skills/roundfix/references/spec-delivery.md skills/roundfix/references/spec-delivery.md .agents/skills/archive-spec/SKILL.md skills/archive-spec/SKILL.md docs/user-guide/commands/archive.md; do tr -s '[:space:]' ' ' < "$file" | grep -qF -- 'archives once every Task is completed, without a QA Report and without an override' || { printf '%s lacks the decline rule\n' "$file" >&2; exit 1; }; done; grep -q '^### Declared QA decline' .agents/skills/archive-spec/SKILL.md || { echo 'archive-spec lacks its Declared QA decline section' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Declared QA Decline**:' || { echo 'CONTEXT.md lacks **Declared QA Decline**' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'A Declared QA Decline archives on Task completion alone, without an override.' || { echo 'Archive Command entry not revised' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Archive Command**:' || { echo 'CONTEXT.md lacks **Archive Command**' >&2; exit 1; }; diff -r .agents/skills/archive-spec skills/archive-spec && diff -r .agents/skills/roundfix skills/roundfix && go test -timeout 15m -count=1 -run '^(TestEveryOwnedSkillVersionIsRecorded|TestSettlementGuidanceIsOneTable)$' ./skills && go test -timeout 15m -count=1 -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task no reference, skill or guide states the decline rule, so the first phrase check fails. Afterwards the skills and their mirrors agree, the owned-skill versions are recorded, the `### QA settlement` table is unchanged, and the documentation contracts pass with the Behavior Surface record re-recorded last.

## References

- `_prd.md` → Goals 4 and 5; Core Features 4 and 5; Success Metrics 5 and 7; Glossary
- `_techspec.md` → Current behavior 7 and 9; API Contracts 5 and 8; Declared breaks 2 and 5; Invariant 3; Testing Approach; Build Order 2
- ADR-0154
