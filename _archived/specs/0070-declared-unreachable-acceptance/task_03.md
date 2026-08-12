---
task: task_03
spec: 0070-declared-unreachable-acceptance
status: completed
type: backend
complexity: high
---

# Task 03: Archive the declared case and refuse the rest

## Overview

The archive boundary gains exactly one accepting case: a `partial` report whose
only unmet rows are declared unreachable, with no finding-blocked and no
environment-blocked row, and whose declared count is covered by the Spec's own
declarations. It also stamps what remains unproven, so the archive record
carries the debt instead of hiding it.

This Spec widens an archive gate, so this slice is written as refusals with one
exception. The fixtures that must still refuse are worth more than the one that
must now accept.

## Requirements

1. MUST archive a Spec whose newest report is `partial` with
   `rows_blocked_declared` greater than zero, `rows_blocked_finding` zero,
   `rows_blocked_environment` zero, and at least as many declarations as
   declared rows.
2. MUST refuse, with a diagnostic naming the cause, when: any row is
   finding-blocked; any row is environment-blocked; the declared count exceeds
   the Spec's declarations; or the verdict is `fail`.
3. MUST keep environment-blocked rows blocking. Circumstance is not
   unreachability, and this is the refusal most likely to be relaxed by
   accident.
4. MUST stamp the archived Spec with the satisfying action of every declaration
   that covered an unmet row, so a reader of the archive learns what was never
   verified.
5. MUST leave `qa_override` unchanged in meaning, reach, and wording.
6. MUST archive every Spec that archives today under `pass` identically, proven
   over the archived corpus rather than one fixture.
7. SHOULD make each refusal diagnostic name the count that caused it, so the
   operator does not have to open the report to learn which rule fired.

## Subtasks

- [ ] Add the accepting case to the archive precondition.
- [ ] Stamp the unproven actions into the archived Spec.
- [ ] Add the refusal diagnostics naming each cause.
- [ ] Build the fixture matrix over the whole API Contracts table.
- [ ] Assert non-regression over the archived corpus.

## Acceptance Criteria

- [ ] A `partial` report with declared rows only, fully covered by
      declarations, archives and stamps every satisfying action.
- [ ] A `partial` report with any finding-blocked row refuses, naming the
      finding count.
- [ ] A `partial` report with any environment-blocked row refuses, naming the
      environment count.
- [ ] A report declaring more blocked rows than the Spec declares refuses,
      naming the shortfall.
- [ ] A `fail` report refuses exactly as today.
- [ ] A `pass` report archives exactly as today, with no stamp added.
- [ ] `qa_override` archives a genuinely failed Spec, unchanged.
- [ ] Every Spec in the archived corpus still satisfies the archive
      precondition it satisfied before this Task.

## Context

- interface: `internal/spec/archive.go`
- interface: `internal/cli/archive.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/spec -count=1 -run 'Archive' -v | grep -q -- "--- PASS"`
  — expected: exit 0; the archive tests ran and passed.
- `go test ./internal/cli -count=1 -run 'Archive' -v | grep -q -- "--- PASS"`
  — expected: exit 0; the command tests ran and passed.
- `go test ./internal/spec ./internal/cli -count=1` — expected: exit 0.
- `go test -parallel 16 ./... 2>&1 | grep -q "^FAIL" && exit 1 || exit 0`
  — expected: exit 0.
- `if git diff --name-only HEAD | grep -E "^(\.agents|skills)/" | grep -q .; then exit 1; fi`
  — expected: exit 0; this Task touches no skill; the gate contract is task_04's
  bounded scope.

## References

- `_prd.md` → Core Features 3, 4, 5 and 6; Decisions; Success Metrics 2 and 4.
- `_techspec.md` → API Contracts; Build Order 3; Risks & Considerations.
- ADR-0080.

## Result

Implemented the archive exception for a `partial` QA Report whose blocked rows
are exclusively declared unreachable. The archive precondition still accepts
`pass` through its existing path, still refuses `fail`, and now rejects partial
reports with count-bearing diagnostics for finding-blocked,
environment-blocked, or declaration-shortfall states. An accepted declared
case stamps every author-supplied `satisfied-by` action into `_prd.md`
frontmatter as `unproven` before moving the Spec.

Acceptance evidence:

- Declared-only acceptance and stamp:
  `TestRunArchiveDeclaredUnreachableContract/declared-only_partial_report_archives`
  moved the Spec and compared the parsed `unproven` sequence with both
  satisfying actions. The `surplus_declarations_still_cover_declared_rows`
  case also archived when declarations outnumbered declared rows.
- Finding refusal:
  `TestRunArchiveDeclaredUnreachableContract/finding-blocked_partial_report_refuses`
  kept the Spec active and required stderr to name
  `rows_blocked_finding is 2`.
- Environment refusal:
  `TestRunArchiveDeclaredUnreachableContract/environment-blocked_partial_report_refuses`
  kept the Spec active and required stderr to name
  `rows_blocked_environment is 3`.
- Declaration shortfall:
  `TestRunArchiveDeclaredUnreachableContract/declaration_shortfall_refuses`
  required the diagnostic to name declared count 3, declaration count 1, and
  shortfall 2.
- Fail compatibility:
  `TestRunArchiveDeclaredUnreachableContract/failing_report_refuses_exactly_as_before`
  retained the existing `newest QA Report verdict is "fail"; expected "pass"`
  refusal.
- Pass compatibility:
  `TestRunArchiveDeclaredUnreachableContract/passing_report_remains_unchanged`
  archived a Spec that carried a declaration and observed no `unproven`
  frontmatter. The pre-existing pass case with
  `rows_blocked_environment: 3` also remained green, preserving the command's
  existing pass behavior rather than tightening it as part of this Task.
- QA override compatibility:
  `TestArchivedQAOverrideCorpusIncludesFailedSpec` found the real archived
  `0057-baseline-capability-evidence-and-retention` Spec with
  `qa_override: true` and newest verdict `fail`. The archive-spec contract
  remains byte-unchanged and still states that the override reaches failed or
  missing QA evidence only.
- Archived corpus compatibility:
  `TestArchivedPassCorpusRemainsArchiveEligible` checked all 50 archived Specs
  whose newest report is `pass`; every Task remained completed, the new archive
  precondition accepted every Spec, and none gained an `unproven` action.

Focused checks:

- Before implementation,
  `rtk sh -c 'GOCACHE=/private/tmp/roundfix-task-03-gocache rtk go test ./internal/cli ./internal/spec -run "^Test(RunArchiveDeclaredUnreachableContract|ArchivedPassCorpusRemainsArchiveEligible)$" -count=1'`
  reported the legacy signal: the pass and fail rows passed, while the declared
  acceptance and three count-specific refusals failed; the corpus test also
  lacked the new precondition function.
- After the final implementation and test edits,
  `rtk sh -c 'GOCACHE=/private/tmp/roundfix-task-03-gocache rtk go test ./internal/cli -run "^TestRunArchive" -count=1'`
  passed 17 tests in one package.
- `rtk sh -c 'GOCACHE=/private/tmp/roundfix-task-03-gocache rtk go test ./internal/spec -run "^Test(ArchivedPassCorpusRemainsArchiveEligible|ArchivedQAOverrideCorpusIncludesFailedSpec|UnreachableReadsDeclaredAcceptance|ReadQAReportBlockedCounts)$" -count=1'`
  passed 7 tests in one package.
- The two exact archived-corpus tests passed with verbose evidence: 50 pass
  Specs checked and the failed QA-override Spec named above found.
- `rtk git diff --check` passed.

Follow-up: the repository-owned Roundfix Skill still describes the Archive
Command as pass-only. Updating `.agents/skills/roundfix/SKILL.md` and its mirror
is protected-tooling work outside this Task and outside the Spec's bounded
tooling authorization; route that synchronization through an expressly
authorized follow-up before Pull Request delivery.

The commands under `## Verification` were not run; the Daemon owns them and
Task settlement.
