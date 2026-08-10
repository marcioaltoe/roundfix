---
task: task_08
spec: 0091-a-proof-that-can-refuse
status: pending
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
4. MUST cover the installed shape with a fixture, so the gap that let this reach
   QA cannot reopen silently.

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

- `GOCACHE="$PWD/.gocache" go test ./internal/agent -run '^TestMissingSessionIsRecognisedFromTheInstalledExitShape$' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestMissingSessionIsRecognisedFromTheInstalledExitShape'` — expected: exits 0. This test does not exist yet; it must supply exit code 1 with `No named session` and assert the result is Session absence.
- `GOCACHE="$PWD/.gocache" go test ./internal/agent -run '^TestUnrelatedExitOneKeepsItsClassification$' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestUnrelatedExitOneKeepsItsClassification'` — expected: exits 0. This test does not exist yet; it keeps the recognition narrow, so the fix cannot be "treat every exit 1 as absence".
- `GOCACHE="$PWD/.gocache" go test ./internal/agent -run '^TestDisposableSessionCloseIsNotAppendedWhenTheSessionNeverOpened$' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestDisposableSessionCloseIsNotAppendedWhenTheSessionNeverOpened'` — expected: exits 0. This existing test must keep passing, so the new shape joins the old one rather than replacing it.

## References

- QA report: `qa/qa-report-2026-08-10-01.md` — repeated finding F-003.
- Task 04: `task_04.md` — the exit code 4 shape this extends.
