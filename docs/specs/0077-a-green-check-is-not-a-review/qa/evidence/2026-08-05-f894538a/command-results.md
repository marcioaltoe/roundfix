# Command results — build f894538a

All commands ran from the Run Worktree on 2026-08-05. No command was piped
when its exit status served as gate evidence.

## Repository gate

- `rtk make verify` — exit 0.
  - 3,406 Go tests passed in 26 packages.
  - `TestCheckCorpusBudget` passed in its isolated package run.
  - 4 Skill tests passed.
  - `roundfix skills check` passed for every owned Skill.
  - `bin/roundfix` built successfully.
  - `bin/roundfix spec check` reported no finding for Spec 0077.

## Recorded payload and classification

- `rtk go test ./internal/reviewsource/coderabbit -count=1 -run '^(TestEvidenceHierarchyPrecedence|TestEvidenceRefusalClassTable|TestEvidenceRefusalReasonIsBoundedVerbatim|TestEvidenceRateLimitWithoutAuthoritativeCommentStaysPending|TestEvidenceExpectedHeadRejectsUnboundAndStaleSignals|TestEvidenceRefusalForStaleHeadDoesNotSettleCurrentHead|TestEvidenceRecordedCommitStatusCorpus|TestIssueCommentsMapGitHubRateLimitCommentJSON)$' -v`
  — exit 0; 25 tests passed.
- `rtk go test ./internal/reviewsource/coderabbit -count=1` — exit 0; 77
  tests passed.

The focused source-backed assertions establish:

- Pull Request #107's recorded head
  `c6c14bece33bddf153c81c16029a97537f94d7c9`, green `CodeRabbit` status with
  description `Review rate limited`, and authoritative CodeRabbit comment
  resolve to `skipped`; `Reason` and `Detail` are `Review limit reached`.
- The same green status without the authoritative comment remains `pending`.
- Rate-limit and path-filter refusals, including title-case variants, resolve
  to `skipped` by class.
- Unknown successful check names and output titles stay `pending`.
- A recorded `Review completed` status still resolves `verified`; a completed
  review with unresolved threads remains `reviewed`.
- Earlier-head refusal and completed-review signals settle nothing for the
  current head.

## Assembled watch and event behavior

- `rtk go test ./internal/watch -count=1 -run '^(TestRunUnrecognisedPendingEvidenceStopsBeforeFetchWithDiagnostic|TestRunReviewSkippedStopsBeforeFetch|TestRunReviewSkippedDuringMergeReadyPreservesTerminalEvidence|TestRunReviewEvidenceSharedByPreFetchAndMergeReady)$' -v`
  — exit 0; 4 tests passed.
- `rtk go test ./internal/cli -count=1 -run '^(TestRunWatchReviewSkippedPublishesReasonWithoutArtifactsOrCleanup|TestRunWatchUnrecognisedGreenSignalDiagnosesAndDoesNotPush|TestRunWatchReviewIssuesKnownAfterFetchedZero|TestRunWatchArtifactEvidenceInheritedWithoutCurrentHeadPolling)$' -v`
  — exit 0; 4 tests passed.
- `rtk go test ./internal/runevent -count=1 -run '^(TestProjectStreamEventReviewSkippedOutcome|TestProjectStreamEventPendingUnrecognisedOutcome|TestProjectStreamEventOutcomeContextProjectsReviewIssuesEvidenceAndRecovery)$' -v`
  — exit 0; 3 tests passed.
- `rtk go test ./internal/watch ./internal/runevent ./internal/cli -count=1`
  — exit 0; 1,048 tests passed across 3 packages.

The public CLI runner assertions establish that refused evidence ends
`ReviewSkipped`, prints the refusal and next action, creates no review
artifacts, makes zero review-artifact commits and zero Final Push calls, and
records `evidence_state: skipped`. The unknown-green fixture ends non-Clean,
prints `signal was not recognised`, makes zero Final Push calls, and records
`evidence_state: pending`. Adjacent verified fixtures still reach Clean and
perform their expected single Final Push.

## Tooling and scope

- `rtk make skills-sync-check` — exit 0; 4 tests passed.
- `rtk cmp -s .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md` —
  exit 0.
- `rtk git diff-tree --no-commit-id --name-only -r f894538a -- '*.go'` —
  exit 0 with no output; task_04 changed no Go file.
- `rtk git log --oneline -G 'retry|re-request|retrigger|capacity wait' 0bfa6c4e..HEAD -- internal/reviewsource internal/watch internal/cli internal/runevent`
  — exit 0 with no output; the implementation commits introduced no matching
  retry/re-request/retrigger/capacity-wait policy.
- Exact path inspection with `git diff-tree` showed tasks 01-03 only in their
  Task files and bounded product/test packages. Task 04 changed its Task file,
  the canonical/mirrored Roundfix Skill, and ADR-0081 deterministic files under
  `DERIVED_DIGEST_PATHS`.

## Environment-blocked live journey

The Roundfix QA prompt states: `Pull Request: none open; Pull Request journeys
are environment-blocked.` No attempt was made to infer a Pull Request from the
Run Worktree branch, which is never pushed and has no Pull Request. Equivalent
evidence is the recorded #107 payload replay, assembled public CLI runner,
terminal event projection, positive-review canaries, and full repository gate
above.

## F-001 evidence — built help contradicts the closed default

`rtk proxy ./bin/roundfix watch --help` exited 0 and printed:

```text
With --until-clean, Clean requires accepted Review Source Evidence
on the pushed head — a successful CodeRabbit check or commit status, or a
CodeRabbit APPROVED review — with zero unresolved CodeRabbit threads
```

The synchronized operator Skill at `.agents/skills/roundfix/SKILL.md` states
the shipped rule instead:

```text
verified requires a recognised review-completed current-head CodeRabbit check
or commit status, or a current-head CodeRabbit APPROVED review
```

It also states that an unrecognised successful signal stays `pending` and that
a green check is not evidence that a review ran. The built help therefore
teaches the permissive contract this Spec removes.
