---
task: task_04
spec: 0203-a-reviewer-that-validates-its-findings-and-remembers-its-rounds
status: completed
type: backend
complexity: high
---

# Task 04: Round 2 continues the round-1 session, and the skill describes the reviewer

## Overview

After task_03, round 2 reads the delta and the recorded round-1 findings, but
it still opens a new reviewer session. This Task keeps a round-1 session open
when findings stand, and continues it in round 2 through acpx on the same
selection. It closes the session at every point the TechSpec names. It
records the ACP session id of each round, so `continued` is measured rather
than assumed.

It also documents everything the Spec changes in the review reference of the
Roundfix skill, its mirror and the review command guide, and raises the skill's
version.

The session id comes from the ACP stream acpx relays: the `sessionId` every
`session/update` notification carries. The session name comes from the
checkout's own record. Nothing here writes the ledger or reads another
checkout.

## Requirements

1. MUST add `ACPSessionID string` to `agent.ExecuteResult`. The acpx runner
   sets it from the `sessionId` of the first `session/update` notification it
   reads in a prompt stream, and leaves it empty when none carries one. Every
   existing `ExecuteResult` field and its value MUST stay as today.
2. MUST keep `reviewSessionRef` random and unchanged. It MUST implement
   `_techspec.md` → The continued session:
   - a round 1 whose outcome is `findings` does not end its session, and
     records `session`, `selection` and `sessionOpen: true`;
   - round 2 prepares the recorded session on the recorded selection, falls
     back to today's loop with fresh names when preparation fails, and always
     ends its session;
   - a lineage change ends the prior open session, and so does a reuse at the
     same head that resolves to `findings-dismissed`;
   - every other review ends each session it prepared.
3. MUST append each round's `ACPSessionID` to `lineage.acpSessionIds`, and
   set `continued` true only when round 1 and round 2 report the same
   non-empty id.
4. MUST describe, in `.agents/skills/roundfix/references/review.md` and in
   `docs/user-guide/commands/review.md`, each of the following:
   - the Delivery Conventions, with the version string
     `roundfix/delivery-conventions/v1`;
   - the finding grammar;
   - validation and the status `dismissed-by-validation`, with its three rules;
   - the fail-closed validator;
   - the Reviewer Lineage and its round-2 delta;
   - the continued session and `continued`;
   - the ceiling and `ceiling-closed`;
   - the `dispose` refusal.

   Both files MUST contain the phrases `roundfix/delivery-conventions/v1`,
   `dismissed-by-validation`, `Reviewer Lineage` and `ceiling-closed`. The
   text MUST go under the existing `roundfix review` material, and MUST NOT
   touch the `### QA settlement` section of any skill.
5. MUST raise both version fields of `.agents/skills/roundfix/SKILL.md` by the
   step the version rule in force on its starting commit requires. It MUST
   change no other line of that file. Then it MUST run `make skills-sync` so
   `skills/roundfix/SKILL.md` and `skills/roundfix/references/review.md`
   match, and record the new version with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`.
   It MUST run `make baseline-digests` and name in the Result any file it
   changed.
6. MUST put the new review tests in `internal/cli/review_session_test.go` and
   the runner test in `internal/agent/agent_session_id_test.go`. They use the
   fake runner's `PrepareSession`, `RunPrepared` and `EndSession` records and
   the existing fake ACP stream helper, never a real adapter or `~/.acpx`.
7. MUST keep `TestReviewSessionRefIsUniquePerInvocation`,
   `TestReviewCommandUsesFallbackOnlyWhenSelectionFailsBeforePrompt`,
   `TestReviewSessionReadsWithoutWriting` and
   `TestACPXRunPromptPublishesUpdateLinesAndCapturesStopReason` green. If one
   asserts an `EndSession` call that Requirement 2 removes, it MUST be updated
   without renaming, and named in the Result.

## Subtasks

- [ ] Capture the ACP session id in the acpx runner.
- [ ] Keep, continue and close the reviewer session per the lineage.
- [ ] Record the session ids and `continued`.
- [ ] Document the reviewer in the skill reference and the command guide.
- [ ] Raise and record the skill version, and sync the mirror.
- [ ] Add a test for each acceptance criterion, with each negative case separate.

## Acceptance Criteria

- [ ] `TestACPXRunPromptReportsTheACPSessionID`: a fake stream whose
      notifications carry `sessionId` `s-1` reports `s-1`; a stream without it
      reports an empty id.
- [ ] `TestReviewRoundOneWithFindingsLeavesItsSessionOpen`: no `EndSession`
      for that name, and the record holds `session`, `selection` and
      `sessionOpen: true`. A round 1 ending `reviewed` ends its session.
- [ ] `TestReviewRoundTwoContinuesTheRecordedSession` reproduces Surface
      Transcript 4. The same name is prepared and then ended. With different
      ids in the two rounds, `continued` is false.
- [ ] `TestReviewRoundTwoFallsBackToAFreshSessionWhenPreparationFails`: the
      recorded name is ended, and a fresh `reviewSessionRef` name runs the
      prompt with `continued: false`.
- [ ] `TestReviewLineageChangeEndsTheOpenSession`: a rebase ends the prior
      open session before the new round 1.
- [ ] `TestReviewDismissedReuseEndsTheOpenSession`: after the operator
      dismisses every standing finding, the reuse ends the open session and
      exits `0`.
- [ ] The review reference, its mirror and the command guide carry the four
      phrases. The version fields moved, `make skills-sync-check` passes, and
      `TestEveryOwnedSkillVersionIsRecorded` passes.
- [ ] The existing tests in Requirement 7 pass.

## Context

- interface: `internal/agent/agent.go`
- interface: `internal/agent/acpx_runner.go`
- interface: `internal/agent/acp_stream.go`
- creates: `internal/agent/agent_session_id_test.go`
- interface: `internal/cli/review.go`
- creates: `internal/cli/review_lineage.go`
- interface: `internal/cli/review_test.go`
- creates: `internal/cli/review_session_test.go`
- creates: `.agents/skills/roundfix/references/review.md`
- creates: `skills/roundfix/references/review.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- creates: `docs/user-guide/commands/review.md`
- instruction: `docs/adr/0197-a-pre-pr-reviewer-lineage-spans-at-most-two-rounds.md`
- instruction: `docs/adr/0189-an-owned-skills-version-names-its-content-and-its-minimum-is-the-bundle.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestACPXRunPromptReportsTheACPSessionID|TestACPXRunPromptPublishesUpdateLinesAndCapturesStopReason|TestReviewRoundOneWithFindingsLeavesItsSessionOpen|TestReviewRoundTwoContinuesTheRecordedSession|TestReviewRoundTwoFallsBackToAFreshSessionWhenPreparationFails|TestReviewLineageChangeEndsTheOpenSession|TestReviewDismissedReuseEndsTheOpenSession|TestReviewSessionRefIsUniquePerInvocation|TestReviewCommandUsesFallbackOnlyWhenSelectionFailsBeforePrompt|TestReviewSessionReadsWithoutWriting|TestReviewCeilingBlocksWithoutCallingTheReviewer|TestReviewValidatorDismissesADaemonSettlementAsC2)$" ./internal/agent ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestACPXRunPromptReportsTheACPSessionID TestACPXRunPromptPublishesUpdateLinesAndCapturesStopReason TestReviewRoundOneWithFindingsLeavesItsSessionOpen TestReviewRoundTwoContinuesTheRecordedSession TestReviewRoundTwoFallsBackToAFreshSessionWhenPreparationFails TestReviewLineageChangeEndsTheOpenSession TestReviewDismissedReuseEndsTheOpenSession TestReviewSessionRefIsUniquePerInvocation TestReviewCommandUsesFallbackOnlyWhenSelectionFailsBeforePrompt TestReviewSessionReadsWithoutWriting TestReviewCeilingBlocksWithoutCallingTheReviewer TestReviewValidatorDismissesADaemonSettlementAsC2; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the seven new tests do not exist, so the command fails.
- `for file in .agents/skills/roundfix/references/review.md skills/roundfix/references/review.md docs/user-guide/commands/review.md; do test -f "$file" || { printf 'missing file: %s\n' "$file" >&2; exit 1; }; for phrase in "roundfix/delivery-conventions/v1" "dismissed-by-validation" "Reviewer Lineage" "ceiling-closed"; do tr -s '[:space:]' ' ' < "$file" | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' "$file" "$phrase" >&2; exit 1; }; done; done; base="$(git log -1 --format=%H -- docs/specs/0203-a-reviewer-that-validates-its-findings-and-remembers-its-rounds/_techspec.md)" || exit 1; test -n "$base" || exit 1; tmp="$(mktemp -d)" || exit 1; git show "$base:.agents/skills/roundfix/SKILL.md" > "$tmp/before.md" || exit 1; awk 'substr($0,1,8)=="version:" {print; exit}' "$tmp/before.md" > "$tmp/old"; awk 'substr($0,1,8)=="version:" {print; exit}' .agents/skills/roundfix/SKILL.md > "$tmp/new"; if cmp -s "$tmp/old" "$tmp/new"; then printf 'the skill version did not change\n' >&2; exit 1; fi; make skills-sync-check && go test -count=1 ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$'` — expected: exit 0; before this Task the phrases are absent and the version is unchanged.

## References

- [_prd.md](_prd.md) — Goal 4; User Story 4; Core Features 7 and 9; Success Metric 5; Recorded limits; Declared breaks
- [_techspec.md](_techspec.md) — The continued session; Data Models; Invariants 7 and 8; API Contracts 1–5; Surface Transcript 4; Integration Points; Testing Approach 4; Build Order 4
- ADR-0197; ADR-0017; ADR-0018; ADR-0051; ADR-0189

## Result

Implemented this Task's session lifecycle and documentation slice. Task status
and the authored Verification commands remain Daemon-owned. No commit, push,
Pull Request, Task Graph edit, or other Task-file edit was performed.

### Implementation and acceptance evidence

- `ExecuteResult.ACPSessionID` captures the first `session/update`
  notification's `sessionId`, including on stopped/error prompt returns. All
  prior result assignments, raw protocol output, message collection, usage,
  event publication and stop reason remain unchanged.
  `TestACPXRunPromptReportsTheACPSessionID` covers `s-1`, missing IDs, no
  notifications, and later notifications not replacing the first ID.
- `TestReviewRoundOneWithFindingsLeavesItsSessionOpen` checks both the printed
  and persisted session name, selection and open flag, with no closure.
  Separate reviewed, blocked and validation-dismissed tests check closure.
- `TestReviewRoundTwoContinuesTheRecordedSession` exercises Surface Transcript
  4 through temporary Git commits and an operator fixing disposition. It
  checks the same prepared name and runtime, the delta prompt and recorded
  conversation, two matching ACP IDs, `continued: true`, and closure.
  Separate different-ID and empty-ID tests require `continued: false`.
  `TestReviewRoundTwoContinuesTheRecordedFallbackSelection` checks selection
  index 1 rather than assuming the preferred selection.
- `TestReviewRoundTwoFallsBackToAFreshSessionWhenPreparationFails` checks that
  the recorded name is closed and a fresh random name receives the prompt,
  with `continued: false`. A separate post-prompt failure test checks closure
  without activating fallback; round-2 findings also close their session.
- `TestReviewLineageChangeEndsTheOpenSession` rebases the candidate and checks
  that the prior name closes before the new round-1 preparation, with a fresh
  ACP-ID history. A separate provider-omission test checks closure when the
  policy changes to `none`, without another reviewer prompt.
- `TestReviewDismissedReuseEndsTheOpenSession` records an evidence-backed
  operator dismissal and checks exit 0, the closed flag, no new prompt, and no
  repeated closure on later reuse. A separate standing-reuse test keeps the
  session open.
- The review reference, shipped mirror and command guide describe Delivery
  Conventions, finding grammar, the three validation rules and fail-closed
  validator, Reviewer Lineage/delta, measured continuation, ceiling closure,
  and dispose refusal. A focused Python inspection confirmed all four required
  phrases in each file, identical skill mirrors, and that only the two version
  fields changed in `SKILL.md`; `### QA settlement` is byte-identical to HEAD.
- The four existing Requirement 7 test names were included in focused checks
  and passed without edits. The existing
  `TestReviewRoundTwoPromptCarriesTheDeltaAndRoundOneFindings` was updated,
  without renaming, to expect one closure after round 2 rather than one per
  round. Its prompt and finding/disposition assertions remain intact.

### Focused checks and regeneration

- Starting evidence: the two required new test files were absent, and the
  existing configured session loop unconditionally ended a prepared session.
- `GOCACHE="$PWD/.gocache" rtk go test ./internal/cli ./internal/agent -run
  'TestReview(RoundOne|RoundTwo|LineageChange|DismissedReuse|StandingReuse|ProviderOmission|Session|CommandUsesFallback)|TestACPXRunPrompt'
  -count=1`: exit 0; 32 tests passed across both packages. This focused selection
  differs from the Task's authored Verification and also covers the existing
  ACP update/stop-reason test and the three existing review-session tests.
- A test invocation using the default macOS Go cache was blocked by cache
  filesystem permissions; the repository-local cache rerun above passed.
- `rtk make skills-sync`: exit 0; regenerated
  `skills/roundfix/SKILL.md` and `skills/roundfix/references/review.md`.
- `rtk go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$'
  -record-skill-versions`: exit 0; recorded the patch step `0.1.6` to `0.1.7`
  in `skills/testdata/owned-skill-versions.json`.
- `rtk make baseline-digests`: exit 0, `changed: false`; no derived file changed.
- `rtk proxy git diff --check`: exit 0; no whitespace diagnostic. Git also
  printed its sandbox fsmonitor IPC warning; it did not prevent inspection.

- Initial `rtk make verify-incremental`: exit 2. Its suite guard detected
  Agent edits made while the command was running; two CLI process-owner tests
  also could not read the sandboxed process table, and the unrelated
  `TestTaskBudgetReasonNamesTheSettlementThatRenewedIt` exceeded its 200 ms
  budget under concurrent checks. No assertions or timing were changed. A
  rerun uses process access and a stationary tree.
- Final `rtk make verify-incremental` with process access and no concurrent
  edits: exit 0. Formatting, `go vet ./...`, `go test -parallel 16 ./...`,
  `skills-sync-check`, skill checks and the binary build passed. The process
  tests and the daemon budget test passed without source/assertion changes.
  This is incremental implementation evidence, not the Daemon's authored
  Verification or terminal Task settlement.

### Changed paths

- `internal/agent/agent.go`
- `internal/agent/acpx_runner.go`
- `internal/agent/agent_session_id_test.go`
- `internal/cli/review.go`
- `internal/cli/review_lineage.go`
- `internal/cli/review_lineage_test.go`
- `internal/cli/review_session_test.go`
- `.agents/skills/roundfix/SKILL.md`
- `.agents/skills/roundfix/references/review.md`
- `skills/roundfix/SKILL.md`
- `skills/roundfix/references/review.md`
- `skills/testdata/owned-skill-versions.json`
- `docs/user-guide/commands/review.md`
- This Task file's `## Result` only; its `in_progress` status was set by the
  Daemon, not this Agent.

Tests use only fake prepared-session runners and the existing fake ACP stream
helper, temporary repositories and temporary Artifact Directories. No real
adapter, reviewer, network, `~/.acpx` or live ledger was accessed. Actual
adapter resumption remains measured by the recorded IDs during real use.

## Recorded paths

The Daemon recorded these paths, which this Task changed without declaring them in `## Context`.

- `internal/cli/review_lineage_test.go`

## Carry-forward provenance

- Source Run: `run_20261001T145917Z_95f0b38476c6ebcc`
- Source commit: `380fb6a28a82b717b9bef22db638fccf1920626c`
