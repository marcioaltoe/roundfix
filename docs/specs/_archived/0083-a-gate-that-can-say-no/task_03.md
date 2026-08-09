---
task: task_03
spec: 0083-a-gate-that-can-say-no
status: completed
type: test
complexity: medium
---

# Task 03: Stop the archived corpus counter from gating authoring

## Overview

The corpus golden pins a per-code finding count over the archived Spec corpus.
Adding an ADR or a Spec moves that count, so ordinary authoring fails the build
and the repair is to edit a recorded number — which teaches nothing and defers
the same tax to the next author. This task keeps the sweep and changes what it
asserts: archived counts are derived and reported, and only a regression in the
active corpus fails.

## Requirements

1. MUST keep sweeping both the active and archived corpora, so the check's
   coverage of the sweep itself is unchanged.
2. MUST stop failing when an archived-corpus count changes because unrelated
   authoring landed; those counts are reported for a human, not asserted.
3. MUST still fail when the active corpus regresses, so a Spec that introduces a
   real consistency defect is still caught.
4. MUST state, in the check's own output or a comment, why archived counts are
   reported rather than asserted, so the next author does not re-pin them.
5. MUST NOT weaken any `spec check` detector itself; only this characterization
   check's assertion changes.
6. MUST change only these repository-relative paths plus this Task file:
   `internal/speccheck/constraints_characterization_test.go` and
   `internal/speccheck/testdata/corpus-golden.json`. Any other changed path
   fails this Task.

## Subtasks

- [x] Keep the existing sweep over both corpora.
- [x] Report archived counts instead of asserting them.
- [x] Keep or add an assertion that fails on an active-corpus regression.
- [x] Record why archived counts are not asserted.
- [x] Confirm the changed-file set matches the declared boundary.

## Acceptance Criteria

- [x] Adding an ADR or a Spec no longer requires editing a recorded count.
- [x] An introduced consistency defect in an active Spec still fails the check,
      proven by observation rather than asserted.
- [x] The archived corpus counts still appear in the check's output.
- [x] No `spec check` detector changed.

## Context

- instruction: `docs/workflow/authorizations/2026-08-07-make-the-gate-honest.md`
- interface: `internal/speccheck/constraints_characterization_test.go`

## Verification

- `go test ./internal/speccheck -run 'Corpus' -count=1 -v > /tmp/task_03-1.log 2>&1 && grep -q '^--- PASS: ' /tmp/task_03-1.log` — expected: exits 0, proving the corpus checks run and pass.
- `go test ./internal/speccheck -count=1` — expected: exits 0.
- `grep -q -iE 'archiv' internal/speccheck/constraints_characterization_test.go` — expected: exits 0, proving the archived-corpus treatment is stated where the next author will read it.
- `(git diff --name-only HEAD; git ls-files --others --exclude-standard) | grep -v -E '^(internal/speccheck/constraints_characterization_test\.go|internal/speccheck/testdata/corpus-golden\.json|docs/specs/0083-a-gate-that-can-say-no/task_03\.md)$' | grep . ; test $? -eq 1` — expected: exits 0, proving no path outside the declared boundary changed.

## References

- `_techspec.md` → Build Order 5; Interfaces: the corpus counter.
- `_prd.md` → Core Feature 3; Goal 3.

## Result

### Implementation

- Kept both repository corpus sweeps in `TestCheckCorpusGolden`, but split
  their outcomes: active counts remain compared with the golden, while
  archived counts are rendered and logged as a derived report.
- Changed the golden to the active-only
  `roundfix-speccheck-corpus/v2` schema. Its update guidance explains that
  historical authoring moves archived counts, and the reader rejects unknown
  fields so archived pins cannot be added back unnoticed.
- Added the same rationale beside the archived report in the test. No
  production `spec check` source or detector changed.

### Focused checks

- Before the implementation, a temporary authorized-fixture mutation changed
  only archived `SC-ADR-RELATED` from 72 to 73. The focused
  `TestCheckCorpusGolden` run exited non-zero while every active count remained
  zero, reproducing the unstable archived-count oracle. The mutation was
  restored before implementation.
- After the implementation,
  `rtk proxy env GOCACHE=<worktree>/.gocache go test ./internal/speccheck -run '^TestCheckCorpusGolden$' -count=1 -v`
  exited zero and printed all 18 archived per-code counts under
  `reported, not asserted`.
- With a temporary contradictory requirement pair added to this active Task,
  the focused corpus test exited non-zero and reported
  `SC-REQUIREMENT-CONTRADICTORY: 1` in the active map. The pair was restored
  immediately; this observed mutation proves the active-corpus assertion still
  detects a real consistency regression.
- The first focused run did not reach the test because the sandbox denied the
  host Go cache. Re-running with the worktree-local `GOCACHE` reached the
  intended test signal; no tracked cache artifact was created.

### Acceptance evidence

1. The v2 golden contains only active counts, and its decoder rejects unknown
   fields; archived authoring has no recorded number to update.
2. The temporary active contradiction made the focused corpus test fail with
   one `SC-REQUIREMENT-CONTRADICTORY` finding.
3. The post-change verbose focused run printed the complete archived count map,
   including the non-zero historical counts.
4. The implementation diff is confined to this characterization test and its
   golden fixture; production detector files under `internal/speccheck` remain
   untouched.

### Daemon verification

Not run in this Agent turn. The Daemon owns every command in `## Verification`
and the Task's terminal status.
