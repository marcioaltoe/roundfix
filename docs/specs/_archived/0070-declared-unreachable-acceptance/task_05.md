---
task: task_05
spec: 0070-declared-unreachable-acceptance
status: completed
type: test
complexity: medium
---

# Task 05: Replay Spec 0058 and hold the corpus unchanged

## Overview

This Spec's first Success Metric is a replay: Spec 0058's QA report and
artifacts, with its release row declared, must archive without `qa_override`
and stamp the release action as unproven. That Spec is why this one exists —
its evidence was complete except for something no gate may ever perform, and
the only exit available spent the mechanism reserved for failed evidence.

The second half is the guard: every Spec that archives today still archives
identically. A widened gate is only safe if the widening is measured.

## Requirements

1. MUST replay Spec 0058's archived report and artifacts with its release row
   declared, and assert it archives without `qa_override`.
2. MUST assert the resulting archive record names the release action that
   remains unproven, rather than dropping it.
3. MUST assert a row declared unreachable that the environment could in fact
   reach is reported as wrongly declared rather than accepted.
4. MUST assert a blocked row with no matching declaration still blocks the
   archive.
5. MUST assert non-regression across the archived Spec corpus: every Spec that
   satisfied the archive precondition before this Spec still satisfies it, and
   no archived artifact is modified.
6. MUST make the replay fixture's provenance explicit — which report it
   reproduces, and that the declaration is added by this Spec rather than
   present in the original.

## Subtasks

- [ ] Build the Spec 0058 replay fixture with its provenance recorded.
- [ ] Assert the archive succeeds and the unproven action is stamped.
- [ ] Assert the wrongly-declared and unmatched-row refusals.
- [ ] Assert corpus non-regression.

## Acceptance Criteria

- [ ] The 0058 replay archives without `qa_override` and stamps the release
      action.
- [ ] The replay fixture records the report path it reproduces and states the
      declaration was added by this Spec.
- [ ] A wrongly declared row is reported, not accepted.
- [ ] An unmatched blocked row still blocks the archive.
- [ ] Every archived Spec still satisfies its archive precondition.
- [ ] No file under `docs/specs/_archived/` is modified.

## Context

- instruction: `docs/agents/spec-routing.md`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/spec -count=1 -run 'Replay|Corpus|Unreachable' -v | grep -q -- "--- PASS"`
  — expected: exit 0; the replay and corpus tests ran and passed.
- `go test ./internal/spec ./internal/cli -count=1` — expected: exit 0.
- `go test -parallel 16 ./... 2>&1 | grep -q "^FAIL" && exit 1 || exit 0`
  — expected: exit 0.
- `git diff --name-only HEAD -- docs/specs/_archived | grep -q . && exit 1 || exit 0`
  — expected: exit 0; no archived Spec file changed.

## References

- `_prd.md` → Success Metrics 1, 2, 3 and 4; Decisions.
- `_techspec.md` → Testing Approach; Build Order 5.
- ADR-0080.

## Result

Added a Spec 0058 replay over the real archived artifact tree. The test copies
that tree into a temporary Spec Root, removes the historical `qa_override`,
adds the Spec-0070 declaration overlay, and supplies a newer typed QA Report
whose only unmet row is the real tagged OIDC release. The replay invokes the
real Archive boundary and parses the resulting PRD frontmatter.

Acceptance evidence:

- Archive without override and stamped release action:
  `TestSpec0058ReplayArchivesDeclaredUnreachableRelease` required a `partial`
  report with one declared-only blocked row, archived without `qa_override`,
  first compared the copied source QA Report byte-for-byte, and then compared
  `unproven` with `a maintainer publishes a tagged release and records the
  run`.
- Explicit provenance:
  `TestSpec0058ReplayRecordsFixtureProvenance` required the fixture to name
  `docs/specs/_archived/0058-npm-trusted-publishing-and-release-preflight/qa/qa-report-2026-08-01-04.md`,
  state that the original PRD had no declaration section, and state that Spec
  0070 added the declaration overlay.
- Wrongly declared row:
  `TestSpec0058ReplayReportsWronglyDeclaredReachableRow` read a replay report
  that records a `wrongly-declared-row finding`, required verdict `fail`, and
  proved Archive refused while leaving the active Spec in place.
- Unmatched blocked row:
  `TestSpec0058ReplayRefusesUnmatchedBlockedRow` omitted the declaration while
  retaining `rows_blocked_declared: 1`; Archive refused and named the count,
  zero declarations, and shortfall one.
- Archived corpus eligibility:
  the existing `TestArchivedPassCorpusRemainsArchiveEligible` reran against
  the full archived pass corpus and continued to accept every prior archive
  precondition with no unproven action.
- Archived corpus immutability:
  the accepted replay hashes every file under `docs/specs/_archived/` before
  and after archiving the temporary copy and requires byte-identical maps.
  `rtk git status --short -- docs/specs/_archived` also reported no changed
  path after the implementation.

Focused checks:

- Before implementation,
  `rtk rg -n "Spec0058Replay|archive-replay-0058" internal/spec` returned no
  matches, establishing that the historical replay and provenance fixture did
  not exist.
- The first focused replay run exercised all four new tests; three passed and
  the provenance assertion exposed a wrapped fixture sentence. After fixing
  that fixture text,
  `rtk sh -c 'GOCACHE=/private/tmp/roundfix-task-05-gocache rtk go test ./internal/spec -run "^TestSpec0058Replay" -count=1'`
  passed all four replay tests.
- After the final replay assertion,
  `rtk sh -c 'GOCACHE=/private/tmp/roundfix-task-05-gocache rtk go test ./internal/spec -run "^Test(Spec0058Replay.*|ArchivedPassCorpusRemainsArchiveEligible|ArchivedQAOverrideCorpusIncludesFailedSpec)$" -count=1'`
  passed all six focused replay and corpus tests.
- `rtk git diff --check` passed.

No follow-up was discovered. The commands under `## Verification` were not
run; the Daemon owns them and Task settlement.
