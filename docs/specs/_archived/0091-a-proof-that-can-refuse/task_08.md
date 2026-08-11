---
task: task_08
spec: 0091-a-proof-that-can-refuse
status: completed
type: backend
complexity: low
---

# Task 08: Recognise the absence the installed acpx reports

## Overview

Third corrective Task, closing the QA gate's F-003. Task 04 set out to stop
appending a close error for a Session that never opened, and it works for the
shape it covered: `internal/agent/acpx_runner.go` maps absence from acpx exit
code 4. The installed acpx reports the same absence as exit 1 with `No named
session` on stderr, so the diagnosis still gains cleanup noise about a Session
that never existed. Task 04's two fixtures pass and never exercise the
installed shape.

The finding is `Friction`: the diagnosis stays correct and actionable, it just
ends with a sentence that contradicts itself. It is closed here rather than
deferred because the fix is one recognition rule and its missing fixture, and
because a proof that reports a failure to close something it never opened
teaches the operator to skim the diagnosis.

Minted after the second corrective Task, which the loop's own rule treats as
the point to re-examine decomposition rather than add another. The maintainer
was consulted and chose to continue: the defect is a fixture that did not cover
the installed exit shape, not a graph that was cut wrong.

## Requirements

1. MUST recognise the installed acpx missing-session shape — exit code 1 with
   `No named session` on stderr — as Session absence, alongside the exit code 4
   shape already recognised.
2. MUST NOT append a close error to the diagnosis when the Session is absent by
   either shape, while keeping the cleanup observation recorded.
3. MUST keep every other exit code 1 result classified exactly as it is today:
   only the `No named session` stderr shape changes meaning.
4. MUST cover both shapes in one fixture, so the gap that let this reach QA
   cannot reopen silently and the exit code 4 shape cannot be lost while adding
   the exit code 1 one. Asserting the existing Task 04 test still passes is not
   available as Verification: it passes before this Task runs, which the
   Daemon's pre-work probe refuses as vacuous — it refused exactly that on the
   first attempt at this Task.

## Subtasks

- [ ] Recognise the installed missing-session shape where exit code 4 is
      recognised today.
- [ ] Keep an unrelated exit code 1 classified as it is now.
- [ ] Add the fixture that exercises the installed shape.

## Acceptance Criteria

1. A proof failure that occurs before the Codex disposable Session opens ends
   without a close error, whether acpx reports absence as exit 4 or as exit 1
   with `No named session`.
2. An exit code 1 whose stderr does not report a missing session keeps its
   current classification and message.
3. The cleanup observation itself is still recorded; only the contradictory
   close error is absent.

## Verification

- `GOCACHE="$PWD/.gocache" go test ./internal/agent -run '^TestMissingSessionIsRecognisedFromBothExitShapes$' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestMissingSessionIsRecognisedFromBothExitShapes'` — expected: exits 0. This test does not exist yet; it must assert Session absence for both the exit code 4 shape Task 04 already recognises and the installed exit code 1 with `No named session`, so the new shape provably joins the old one instead of replacing it.
- `GOCACHE="$PWD/.gocache" go test ./internal/agent -run '^TestUnrelatedExitOneKeepsItsClassification$' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestUnrelatedExitOneKeepsItsClassification'` — expected: exits 0. This test does not exist yet; it keeps the recognition narrow, so the fix cannot be "treat every exit 1 as absence".

## References

- QA report: `qa/qa-report-2026-08-10-01.md` — repeated finding F-003.
- Task 04: `task_04.md` — the exit code 4 shape this extends.

## Result

### Implementation

- `CloseSession` now preserves the existing exit 4 missing-Session mapping and
  also maps exit 1 to missing-Session only when stderr contains the installed
  acpx marker `No named session`.
- The combined proof fixture exercises both absence encodings through
  disposable Session cleanup. A separate negative fixture pins an unrelated
  exit 1 to its existing `InfrastructureError` type and exact diagnostic.

### Focused checks

- Red signal before the production change:
  `rtk proxy env GOCACHE=/private/tmp/roundfix-task08-gocache go test ./internal/agent -run '^(TestMissingSessionIsRecognisedFromBothExitShapes|TestUnrelatedExitOneKeepsItsClassification)$' -count=1`
  failed because `installed_missing_session_exit_one` appended the contradictory
  close error. The same run also showed that the first draft of the negative
  fixture had named the wrong current error type; the fixture was corrected to
  the observed `InfrastructureError` contract before the production edit.
- After the production change, that same two-test command passed.
- `rtk proxy env GOCACHE=/private/tmp/roundfix-task08-gocache go test ./internal/agent -run '^(TestDisposableSessionCloseIsNotAppendedWhenTheSessionNeverOpened|TestMissingSessionIsRecognisedFromBothExitShapes|TestUnrelatedExitOneKeepsItsClassification|TestDisposableSessionCloseIsAppendedWhenAnOpenSessionWillNotClose|TestACPXCloseSessionReturnsCloseFailure)$' -count=1`
  passed.
- `rtk proxy env GOCACHE=/private/tmp/roundfix-task08-gocache go test ./internal/agent -count=1`
  passed.

### Acceptance criteria evidence

1. `TestMissingSessionIsRecognisedFromBothExitShapes` drives a proof failure
   through exit 4 and exit 1 with `No named session`, then compares the exact
   leading diagnosis and proves neither case returns an
   `AgentSessionCleanupError`.
2. `TestUnrelatedExitOneKeepsItsClassification` proves exit 1 with
   `close rejected` remains an `InfrastructureError` and compares the complete
   pre-existing close diagnostic byte for byte.
3. Each absence case captures exactly one warning containing the disposable
   Session cleanup observation and the typed missing-Session reason while the
   returned diagnosis remains free of the close error.

The authored `## Verification` commands were not run; the Daemon owns them.
