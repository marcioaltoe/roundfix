---
task: task_02
spec: 0181-gates-that-refuse-only-what-someone-can-act-on
status: pending
type: backend
complexity: medium
---

# Task 02: The pre-PR Pull Request row never decides a qualifying partial

## Overview

The authored QA gate runs before any Pull Request exists. Every matrix therefore records its Pull Request row as `blocked (environment: no open Pull Request)`, and `QAReportEligibility` in `internal/spec/qa.go` refuses every declared `partial` that carries it: `rows_blocked_environment is 1; expected 0`. This Task derives that one row from the report's validated Results rows and stops it from deciding a qualifying `partial` (ADR-0167). It also tells the QA Agent so through the Daemon's QA prompt. `roundfix qa-report accept`, QA settlement, `roundfix settle` and `roundfix archive` all read the one eligibility function, so the built binary accepts Spec 0179's archived report that the v0.20.0 binary refuses.

## Requirements

1. MUST add the exported constants `QANoOpenPullRequestStatus` (`blocked (environment: no open Pull Request)`) and `QAPullRequestRowSource` (`Pull Request row`) to `internal/spec/qa.go`, and `RowsBlockedPrePullRequest int` to `QAReport`.
2. MUST derive `RowsBlockedPrePullRequest` in `readQAReport` from the report body alone, with the same fence, heading and table helpers `qaReportHollow` uses:
   - it reads only tables inside `## Results` whose header has both a `Status` and a `Provenance` column, compared trimmed and case-insensitively;
   - it counts a row only when its trimmed status equals `QANoOpenPullRequestStatus` exactly, and one item of its provenance, split on `;` and `,` and trimmed, equals `QAPullRequestRowSource` exactly. A provenance that merely contains the words, such as `not a Pull Request row`, and a status differing in case excuse nothing.

   A table without a `Provenance` column contributes nothing. The frontmatter counts MUST be read and validated exactly as before.
3. MUST change only the environment check of the `partial` branch of `QAReportEligibility`. Let `excused` be the smaller of `RowsBlockedPrePullRequest` and `RowsBlockedEnvironment`. When environment-blocked rows remain beyond `excused`, it refuses:
   - with the unchanged `rows_blocked_environment is %d; expected 0` when `excused` is zero;
   - otherwise with `rows_blocked_environment is %d, %d outside the pre-PR Pull Request row; expected 0 outside it`.

   Every other check, its order and the `pass` branch MUST stay unchanged.
4. MUST replace the resolved no-Pull-Request line of the QA prompt in `internal/agent/spec_prompt.go` with the TechSpec's line. The new line keeps the `Pull Request: none open;` prefix and tells the Agent to name the Pull Request row in the row's provenance, and that the row alone never prevents a qualifying declared partial. `TestBuildQAPromptStatesPullRequestJourneysAreEnvironmentBlockedWhenNoneIsOpen` MUST pin the new line under its existing name.
5. MUST keep `internal/speccheck/mechanical.go`'s count cross-check unchanged: the row still counts in `rows_blocked_environment`.
6. MUST keep `TestQAReportEligibilityKeepsExistingPartialRefusalPrecedence` green without editing it, rename or remove no top-level test, and change no exported function signature.
7. MUST put the new tests in `internal/spec/qa_prepr_row_test.go`.

## Subtasks

- [ ] Derive the pre-PR Pull Request row from the Results rows.
- [ ] Apply it in the `partial` eligibility check with the stated messages.
- [ ] Update the QA prompt line and its pinned test.
- [ ] Add a test for each acceptance criterion, each negative case separate.

## Acceptance Criteria

- [ ] A declared `partial` whose only environment-blocked row is the pre-PR Pull Request row, with the declared rows covered by the PRD's declarations, is eligible.
- [ ] A `partial` with one pre-PR row and one other environment-blocked row refuses with `rows_blocked_environment is 2, 1 outside the pre-PR Pull Request row; expected 0 outside it`.
- [ ] A row with the no-PR status whose provenance does not name the Pull Request row as one of its items (including `not a Pull Request row`), a status differing in case, and a Results table without a provenance column excuse nothing, and the refusal keeps its unchanged message.
- [ ] A `pass` report's eligibility is unchanged, and a `partial` with no declared row still refuses as before.
- [ ] The built binary's `roundfix qa-report accept` exits `0` on Spec 0179's archived `qa-report-2026-09-29.md`.

## Context

- instruction: `docs/adr/0167-the-pre-pr-pull-request-row-never-decides-a-qualifying-partial.md`
- interface: `internal/spec/qa.go`
- creates: `internal/spec/qa_prepr_row_test.go`
- interface: `internal/agent/spec_prompt.go`
- interface: `internal/agent/spec_prompt_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestAPartialQualifiesWhenItsOnlyEnvironmentRowIsThePrePullRequestRow|TestAPartialWithAnotherEnvironmentRowStillRefuses|TestANoPullRequestStatusWithoutThePullRequestSourceStillRefuses|TestAResultsTableWithoutProvenanceGainsNoPrePullRequestRow|TestAPassIsUnchangedByThePrePullRequestRow|TestQAReportEligibilityKeepsExistingPartialRefusalPrecedence|TestBuildQAPromptStatesPullRequestJourneysAreEnvironmentBlockedWhenNoneIsOpen)$" ./internal/spec ./internal/agent 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestAPartialQualifiesWhenItsOnlyEnvironmentRowIsThePrePullRequestRow TestAPartialWithAnotherEnvironmentRowStillRefuses TestANoPullRequestStatusWithoutThePullRequestSourceStillRefuses TestAResultsTableWithoutProvenanceGainsNoPrePullRequestRow TestAPassIsUnchangedByThePrePullRequestRow TestQAReportEligibilityKeepsExistingPartialRefusalPrecedence TestBuildQAPromptStatesPullRequestJourneysAreEnvironmentBlockedWhenNoneIsOpen; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done && tmp="$(mktemp -d)" && trap 'rm -rf "$tmp"' EXIT && go build -buildvcs=false -o "$tmp/roundfix" ./cmd/roundfix && "$tmp/roundfix" qa-report accept docs/history/specs/0179-review-findings-with-evidence-and-no-unselected-providers/qa/qa-report-2026-09-29.md` — expected: exit 0; before this Task the five new named tests do not exist and the built binary refuses the 0179 report with `rows_blocked_environment is 1; expected 0`, so the command fails.

## References

- [_techspec.md](_techspec.md) — The pre-PR Pull Request row
- `_prd.md` → Goal 2; Goal 4; Core Feature 2; Success Metric 2
- `_techspec.md` → API Contract 4; Testing Approach 2
- [references/2026-09-29-the-pull-request-row-blocks-every-qualifying-partial.md](references/2026-09-29-the-pull-request-row-blocks-every-qualifying-partial.md)

## Result
