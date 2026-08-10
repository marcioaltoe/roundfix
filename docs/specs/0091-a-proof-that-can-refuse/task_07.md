---
task: task_07
spec: 0091-a-proof-that-can-refuse
status: pending
type: backend
complexity: medium
---

# Task 07: Stop rejecting a selection the runtime accepted

## Overview

Corrective Task, closing the QA gate's second finding and the fixture debt the
membership verdict exposed.

F-002: `bin/roundfix profiles validate --json` against this repository rejects
the configured `claude/opus/high` with `effective_selection_mismatch` while its
own diagnostic prints requested and observed as the same `opus/high`. The same
tuple passes when validated in isolation. A proof that refuses what it just
observed to match is worse than no proof: it teaches the operator to distrust
the verdict.

Two `internal/cli` detach tests also fail, and they are correct to: their fake
runtime advertises only `gpt-5.6-sol` while the profile they build declares a
`gpt-5.5` fallback, and preflight proves every configured tuple. The membership
verdict from Task 03 now refuses that fallback, naming the advertised set. The
fixtures must offer what the profile declares.

F-001 is out of scope by decision: ADR-0119 accepts whichever refusal fires
first, so an adapter refusing before the membership check is the documented
outcome and not a defect to fix here.

## Requirements

1. MUST make `profiles validate` accept `claude/opus/high` when the observed
   selection matches the requested one, in the repository's own configuration
   and alongside the other four configured tuples.
2. MUST preserve `effective_selection_mismatch` for a genuine mismatch: a
   selection whose observed model or reasoning effort differs from the
   requested one still fails, with the same classification.
3. MUST make the two failing detach fixtures advertise the models their
   profiles declare, without weakening what those tests assert about detach
   behaviour.
4. MUST NOT change the membership verdict introduced by Task 03, nor the
   catalogue read introduced by Task 02.

## Subtasks

- [ ] Find why a matching claude tuple is read as mismatched when validated
      with its siblings but not alone.
- [ ] Keep the mismatch classification for a real divergence.
- [ ] Teach the detach fixtures to advertise the fallback model.

## Acceptance Criteria

1. `bin/roundfix profiles validate --json` reports every configured tuple
   proving, with no `effective_selection_mismatch` among them.
2. A selection whose observed effort differs from the requested one still fails
   with `effective_selection_mismatch`.
3. `TestRunImplementDetachPrintsReportAndCompletesRun` and
   `TestRunImplementDetachSurvivesCallerProcessGroupKill` pass, and each still
   asserts the detach behaviour its name claims.

## Verification

- `GOCACHE="$PWD/.gocache" go test ./internal/cli -run '^TestRunImplementDetachPrintsReportAndCompletesRun$' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestRunImplementDetachPrintsReportAndCompletesRun'` — expected: exits 0. Fails today: the fixture advertises only `gpt-5.6-sol` while its profile declares a `gpt-5.5` fallback.
- `GOCACHE="$PWD/.gocache" go test ./internal/cli -run '^TestRunImplementDetachSurvivesCallerProcessGroupKill$' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestRunImplementDetachSurvivesCallerProcessGroupKill'` — expected: exits 0. Fails today for the same reason.
- `GOCACHE="$PWD/.gocache" go test ./internal/agent -run '^TestProofAcceptsAMatchingSelectionAmongSiblings$' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestProofAcceptsAMatchingSelectionAmongSiblings'` — expected: exits 0. This test does not exist yet; it must prove that a tuple whose observed selection matches the requested one is accepted when proved alongside other configured tuples, which is the shape F-002 reports.
- `GOCACHE="$PWD/.gocache" go test ./internal/agent -run '^TestProofStillRejectsAGenuineEffectiveMismatch$' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestProofStillRejectsAGenuineEffectiveMismatch'` — expected: exits 0. This test does not exist yet; it keeps the mismatch classification honest, so the fix cannot be "stop checking".

## References

- QA report: `qa/qa-report-2026-08-10.md` — findings F-001 and F-002.
- ADR: `../../adr/0119-the-refusal-that-fired-first-is-the-refusal.md` — why
  F-001 is accepted rather than fixed.
