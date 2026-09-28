---
task: task_03
spec: 0174-operator-surfaces-that-tell-the-truth
status: pending
type: backend
complexity: high
---

# Task 03: Settlement reads the front matter the derived Verification reads

## Overview

`readQAReport` in `internal/spec/qa.go` bounds a QA Report's front matter with the shared `splitFrontmatter`, which closes at the first `\n---` after an opening `---\n`. A report that opens with two consecutive `---` lines therefore hands YAML a block that starts with a document marker and still carries `verdict: pass`, while the awk reader in `DerivedQAVerification` (`internal/spec/task.go`) closes the front matter at line two, finds no verdict and refuses. Spec 0170's second QA Report had exactly that shape; the Daemon's `settleQAVerdict` in `internal/daemon/task_engine.go` settled the QA Task `completed` and committed it as `(pass)`. The report is written by the QA Agent into the Spec's `qa/` directory; `readQAReport` is read by QA settlement, `roundfix archive`, `roundfix settle` and `roundfix qa-report accept`. Spec 0172, archived at `docs/history/specs/0172-a-qa-gate-that-tells-the-truth/`, made the seed `pending` and refused a hollow pass in `QAReportEligibility`; this Task changes only how the front matter is bounded, below that decision.

## Requirements

1. MUST make `readQAReport` bound the front matter as the derived QA Verification's awk reader does: the first line is exactly `---`; the front matter closes at the first later line that is exactly `---`; inside it exactly one line begins with `verdict:` in column one and its value trimmed of whitespace is non-empty. Lines split on `\n` and compare byte-exactly. The bounded lines are YAML-decoded as today and the body starts after the closing line; the hollow-report check, the verdict switch and the blocked-row rules are unchanged, and `splitFrontmatter` keeps serving every other caller.
2. MUST refuse each broken bound with a `QAReportError` naming it: `QA Report front matter is empty`, `QA Report front matter must open with a "---" first line`, `QA Report front matter has no closing "---" line`, and `QA Report front matter has <n> "verdict:" lines; expected exactly 1`; a front matter with no `verdict:` line keeps `frontmatter has no verdict field`. `TestReadQAReportRefusesAnEmptyFrontMatter`, `TestReadQAReportRefusesAMissingOpeningLine`, `TestReadQAReportRefusesAnUnclosedFrontMatter` and `TestReadQAReportRefusesADuplicatedVerdictLine` pin each message, and `TestReadQAReportReadsAWellFormedFrontMatter` pins the verdict and body of a well-formed report.
3. MUST prove the two readers agree in `internal/spec/qa_frontmatter_test.go`: render `DerivedQAVerification` for a fixture slug, run it with `sh -c` in a temporary repository holding each report as the only report of that slug, and assert that its exit is `0` exactly when `readQAReport` returns no error: `TestQAReportReaderAgreesWithTheDerivedVerification` over a named table (a well-formed `pass`; two consecutive opening `---` lines; two `verdict:` lines; no closing line; an empty verdict value; a blank line before the opening `---`; a `----` closing line; an indented `verdict:` line; CRLF line endings), and `TestQAReportReaderAgreesWithTheDerivedVerificationOnTheArchive` over every QA Report under `docs/history/specs/`, each copied under its own file name. Every fixture whose `verdict:` line has a value carries a supported verdict and consistent blocked-row counts, so readability is the only difference under test.
4. MUST make `settleQAVerdict` return the read error of an unreadable report and settle the QA Task with reason `QA verdict unreadable: <cause>`; `fail`, `pending` and `missing` keep `QA verdict <verdict>`, and a refused `pass` or `partial` keeps `QA verdict <verdict> not accepted: <cause>`. The QA Report commit is still created, with verdict `unreadable`. Update the `unreadable verdict` row of `TestTaskCycleQAVerdictMatrixSettlesRunAndCommitsReport` in `internal/daemon/task_engine_test.go` to the new reason, and put the new daemon tests in `internal/daemon/qa_frontmatter_test.go`: in `TestQASettlementRefusesAnEmptyFrontMatter` a QA Agent that prepends a second `---` line to the seeded report leaves the QA Task `failed` with verdict `unreadable` and a reason containing `QA Report front matter is empty`, and in `TestQASettlementAcceptsAFilledSeededFrontMatter` a QA Agent that fills the seeded report with a `pass` and a Results row leaves it `completed`.
5. MUST add `TestQAReportAcceptRefusesAnEmptyFrontMatter` to `internal/cli/qa_report_test.go` and `TestArchiveRefusesAnEmptyFrontMatter` to `internal/cli/archive_test.go`, each asserting a non-zero exit and a stderr naming `QA Report front matter is empty`, with the files left in place.
6. MUST keep `internal/spec/archive_test.go` byte-identical and `TestArchivedPassCorpusRemainsArchiveEligible` green.
7. MUST state in `.agents/skills/qa-gate/SKILL.md` that the QA Agent keeps the seeded front matter as the report's only one, and that a report whose front matter is empty or duplicated never settles the gate (the phrase `front matter is empty or duplicated`); regenerate `skills/qa-gate/SKILL.md` with `make skills-sync`, then run `make baseline-digests`.
8. MUST state in the `qa-report accept` and `archive` sections of `docs/user-guide/commands.md` and in the QA Report entry of `CONTEXT.md` that a report whose front matter is empty or duplicated is unreadable and refused, using the phrase `front matter is empty or duplicated`.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion, each negative case separate.

## Acceptance Criteria

- [ ] A QA Report the derived QA Verification's reader refuses is unreadable to `readQAReport`, and one it accepts stays readable, over the fixture table and every archived QA Report.
- [ ] A report with an empty front matter settles the QA Task `failed` with verdict `unreadable` and a reason naming the cause, and `qa-report accept` and `archive` refuse it.
- [ ] A seeded report filled with a `pass` and a Results row still settles `completed`, and every archived `pass` Spec stays archive-eligible.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/spec/qa.go`
- creates: `internal/spec/qa_frontmatter_test.go`
- interface: `internal/daemon/task_engine.go`
- interface: `internal/daemon/task_engine_test.go`
- creates: `internal/daemon/qa_frontmatter_test.go`
- interface: `internal/cli/qa_report_test.go`
- interface: `internal/cli/archive_test.go`
- interface: `.agents/skills/qa-gate/SKILL.md`
- interface: `skills/qa-gate/SKILL.md`
- interface: `docs/user-guide/commands.md`
- interface: `CONTEXT.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestQAReportReaderAgreesWithTheDerivedVerification|TestQAReportReaderAgreesWithTheDerivedVerificationOnTheArchive|TestReadQAReportRefusesAnEmptyFrontMatter|TestReadQAReportRefusesAMissingOpeningLine|TestReadQAReportRefusesAnUnclosedFrontMatter|TestReadQAReportRefusesADuplicatedVerdictLine|TestReadQAReportReadsAWellFormedFrontMatter|TestArchivedPassCorpusRemainsArchiveEligible|TestQASettlementRefusesAnEmptyFrontMatter|TestQASettlementAcceptsAFilledSeededFrontMatter|TestTaskCycleQAVerdictMatrixSettlesRunAndCommitsReport|TestQAReportAcceptRefusesAnEmptyFrontMatter|TestArchiveRefusesAnEmptyFrontMatter)$" ./internal/spec ./internal/daemon ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestQAReportReaderAgreesWithTheDerivedVerification TestQAReportReaderAgreesWithTheDerivedVerificationOnTheArchive TestReadQAReportRefusesAnEmptyFrontMatter TestReadQAReportRefusesAMissingOpeningLine TestReadQAReportRefusesAnUnclosedFrontMatter TestReadQAReportRefusesADuplicatedVerdictLine TestReadQAReportReadsAWellFormedFrontMatter TestArchivedPassCorpusRemainsArchiveEligible TestQASettlementRefusesAnEmptyFrontMatter TestQASettlementAcceptsAFilledSeededFrontMatter TestTaskCycleQAVerdictMatrixSettlesRunAndCommitsReport TestQAReportAcceptRefusesAnEmptyFrontMatter TestArchiveRefusesAnEmptyFrontMatter; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done && grep -q "front matter is empty or duplicated" .agents/skills/qa-gate/SKILL.md && grep -q "front matter is empty or duplicated" docs/user-guide/commands.md && grep -q "front matter is empty or duplicated" CONTEXT.md && diff -r .agents/skills/qa-gate skills/qa-gate >/dev/null` — expected: exit 0; before this Task none of the new named tests exists and no guide states the front-matter rule, so the command fails.

## References

- [_techspec.md](_techspec.md) — The QA Report front matter
- `_prd.md` → Goal 4; Core Feature 3; Success Metric 3
- `_techspec.md` → API Contract 6; Testing Approach 3
