---
task: task_02
spec: 0385-a-queue-item-stops-only-for-its-own-work
status: completed
type: backend
complexity: medium
---

# Task 02: A review finding's anchor may end a sentence

## Overview

`parseReviewFindingAnchor` reads an anchor followed by `.` or `;` when a space
or the end of the finding comes next, so a finding such as
`docs/user-guide/commands/skills.md:34. Failure: …` is validated against the
candidate diff instead of blocking the item as `findings name no file and
line`. Punctuation followed by anything else keeps the anchor unreadable. This
answers the Backlog Entry recorded by the operator on 2026-11-11
(`references/2026-11-11-a-review-finding-anchored-as-path-line-and-a-period-blocks-the-item.md`).

## Requirements

1. MUST characterize first, in the Task's `## Result`: on the starting tree
   `parseReviewFindingAnchor` reads none of `path.go:2. Failure: broken`,
   `path.go:2; Failure: broken` and `path.go:2.`.
2. MUST change `parseReviewFindingAnchor` in
   `internal/cli/review_validation.go` as `_techspec.md` → API Contract 3
   says: after the pattern's match, a next rune `.` or `;` is accepted only
   when it is the text's last rune or the rune after it is a space
   (`unicode.IsSpace`); every other next rune keeps today's rule. MUST NOT
   change `reviewAnchorPattern`, the path capture, the line rules or any
   other function.
3. MUST add to `TestReviewFindingAnchorParsesEachForm` in
   `internal/cli/review_validation_test.go` the nine rows `_techspec.md` →
   Testing Approach lists, with their paths and lines, changing no existing
   row, and add `TestReviewReadsAnAnchorFollowedBySentencePunctuation`,
   calling `t.Parallel()` first, on `validationReviewFixture` with the two
   findings and assertions Testing Approach lists.
4. MUST change `docs/user-guide/commands/review.md` with the first text of API
   Contract 7, word for word, right after "For a comma-listed anchor, only the
   first range counts.".
5. MUST, as its last step and after the last guide edit, re-record the
   Behavior Surface Record with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.

## Subtasks

- [ ] Record the characterization.
- [ ] Accept sentence punctuation in the parser and add the rows and the review test.
- [ ] Add the guide sentence, then re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] An anchor followed by `.` or `;` and then a space or the end is read with its path and lines.
- [ ] `path:2.5`, `path:2.md`, `path:2;x` and `path:2..` stay unreadable.
- [ ] A review answer whose findings end their anchors with punctuation ends with standing findings, not `blocked`.

## Context

- instruction: `docs/adr/0395-a-queue-item-stops-only-for-its-own-work.md`
- instruction: `docs/adr/0196-a-pre-pr-review-finding-parks-only-after-validation.md`
- interface: `internal/cli/review_validation.go`
- interface: `internal/cli/review_validation_test.go`
- interface: `docs/user-guide/commands/review.md`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- `parseReviewFindingAnchor`'s accepted forms — consumers: `./internal/cli`
- The `review` guide sentence and the Behavior Surface Record — consumers: `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 30m -v -run 'ReviewFindingAnchor|ReviewReadsAnAnchor|ReviewBlocksFindings|ReviewFindingsAllDismissed|ReviewDiffIndex|ReviewValidator|ReviewPermission' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in 'TestReviewReadsAnAnchorFollowedBySentencePunctuation' 'TestReviewFindingAnchorParsesEachForm/path.go:2._Failure:_broken' 'TestReviewFindingAnchorParsesEachForm/path.go:2.5_broken'; do printf '%s\n' "$out" | grep -qF -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done` — expected: exit 0. Before this Task the new test and rows do not exist; afterwards they pass with every existing anchor, validation and review permission test of `./internal/cli`.

- `tr -s '[:space:]' ' ' < docs/user-guide/commands/review.md | grep -qF -- 'An anchor followed by a period or a semicolon reads as one followed by a space when a space or the end of the finding comes next' || { echo 'review guide: anchor sentence missing' >&2; exit 1; }; go test -count=1 -timeout 30m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the sentence is absent; afterwards the guide carries it and the documentation contracts, the Behavior Surface Record included, pass.

## References

- `_prd.md` → Goals; User Story 3; Core Features 3 and 5; Success Metric 3; Success Metric 5
- `_techspec.md` → Current behavior 5; API Contracts 3 and 7; Invariant 4; Declared breaks; Testing Approach; Build Order 2
- ADR-0395
- [The adopted Backlog Entry](references/2026-11-11-a-review-finding-anchored-as-path-line-and-a-period-blocks-the-item.md)

## Result

Skills loaded: implement-task, coding-guidelines, golang-cli, golang-testing, testing-boss, systematic-debugging, no-workarounds, tech-writer, evidence-gate.

Starting-tree characterization, before changing the parser:

- Added only the nine specified table rows and the specified review regression
  test, then ran
  `go test -count=1 -timeout 5m -v ./internal/cli -run '^TestReview(FindingAnchorParsesEachForm|ReadsAnAnchorFollowedBySentencePunctuation)$'`.
  Exit 1: `path.go:2. Failure: broken`, `path.go:2; Failure: broken` and
  `path.go:2.` each returned an empty anchor and `readable=false`.
- The backticked range and tab-followed dotted filename also returned
  `readable=false`. All existing table rows and the four new malformed-suffix
  rows passed. The new review test returned exit 2, outcome `blocked`, reason
  `findings name no file and line`.
- The parser's post-match delimiter check rejects both punctuation runes;
  the regular expression already captures the intended path and line range.

Task status and declared Verification remain Daemon-owned.

Implementation and focused evidence:

- Acceptance criterion 1: only `parseReviewFindingAnchor`'s post-match
  delimiter check changed. `.` and `;` require end-of-text or a following
  `unicode.IsSpace` rune. The pattern, path capture, line conversion and
  range validation are unchanged. The same focused command above now exits
  0: all 24 parser rows pass, including the five newly readable forms with
  their expected paths and lines. All 15 original rows remain unchanged.
- Acceptance criterion 2: the four new negative rows (`path.go:2.5 broken`,
  `path.go:2.md broken`, `path.go:2;x broken`, `path.go:2.. broken`) pass by
  remaining unreadable.
- Acceptance criterion 3: `TestReviewReadsAnAnchorFollowedBySentencePunctuation`
  calls `t.Parallel()` first and uses `validationReviewFixture` with
  `review.txt:2.` and backticked `review.txt:2` followed by `;`, each with
  `Failure:`. It passes, asserting `exitRunFailed`, outcome `findings`, the
  original findings text and two standing findings anchored to `review.txt`
  line 2.
- Inserted API Contract 7's first text verbatim (apart from line wrapping)
  immediately after the comma-listed anchor sentence. A focused text check
  confirmed the wording and location, and an original-row comparison confirmed
  exactly nine additions without modifying an existing row.
- No declared Verification, repository Verification, commit, push or PR was
  run. The starting tree's only change was the Daemon's `task_02.md` status.
  No other Task or Task Graph file was edited. No follow-up work identified.
- Final implementation step, after the last guide edit:
  `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`
  exited 0 (`ok roundfix/internal/docscontract`, 0.475s), regenerating
  `docs/references/behavior-surfaces.json` through its sanctioned recorder.
  Only this evidence note was added afterwards.
