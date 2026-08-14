---
task: task_04
spec: 0095-a-verification-that-ran-before-anyone-believed-it
status: completed # pending | in_progress | completed | failed — only implement-task changes this
type: backend
complexity: medium
---

# Task 04: Refuse a reversed exit condition

## Overview

Some authored commands invert their own intent: they exit zero exactly when the
work is wrong. The measured forms are `grep -c`, which prints a count and exits
zero whenever it prints one; `grep -v | wc -l`, whose status is the pipeline's
last stage; and `test $(…)` with no comparison. Each is nameable without running
anything, so a static refusal catches them in milliseconds.

## Requirements

1. MUST refuse an authored Verification command matching a known reversed form,
   naming the form and the working replacement.
2. MUST cover, each as its own case, the three measured forms: a count-and-exit
   `grep -c`, a filtered count whose status comes from the last pipeline stage,
   and a `test $(…)` with no comparison.
3. MUST NOT refuse a command that uses one of those tools correctly, so
   `grep -q` and a `test` with a comparison pass.
4. MUST skip rather than fail when the artifact it reads is absent, as the
   presence-aware contract requires.
5. MUST register at the tasks stage beside the existing work-independence
   detector.

## Subtasks

- [ ] Add the refusal code and its detector.
- [ ] Cover each measured reversed form and its correct counterpart.
- [ ] Register it at the tasks stage.

## Acceptance Criteria

- [ ] Each of the three measured forms is refused, named, and given a working
      replacement.
- [ ] `grep -q` and a compared `test` are not refused.
- [ ] A Spec without a task graph skips rather than fails.
- [ ] The code appears in the staged registry at the tasks stage.

## Verification

- `grep -q 'SC-VERIFY-INVERTED-EXIT' internal/speccheck/*.go && grep -q 'CodeVerifyInvertedExit' internal/speccheck/coherence.go` — expected: exits 0, proving the code exists and is registered. Fails today on both clauses.
- `go test -count=1 ./internal/speccheck -run 'TestVerifyInvertedExit' -v > /tmp/0095-t04.log 2>&1; s=$?; grep -q '^--- PASS: TestVerifyInvertedExit' /tmp/0095-t04.log || { cat /tmp/0095-t04.log; exit 1; }; exit $s` — expected: exits 0 and the log names the passing test; fails today.
- `grep -c 'PASS' /tmp/0095-t04.log > /tmp/0095-t04-n.txt; test "$(cat /tmp/0095-t04-n.txt)" -ge 4 || { echo "expected at least four passing cases, got $(cat /tmp/0095-t04-n.txt)"; cat /tmp/0095-t04.log; exit 1; }` — expected: exits 0, proving each measured form and at least one correct counterpart are exercised as their own cases rather than one combined assertion.

## References

`_techspec.md` → Build Order 4; Interfaces, the refusal codes; Testing Approach,
the static detectors. `_prd.md` → Core Feature 2; User Story 4. ADR-0093,
ADR-0094.

## Result

Implemented `SC-VERIFY-INVERTED-EXIT` as a tasks-stage Spec Consistency Check
detector. It reports one error per matched Verification command, names the
measured shell form, and supplies an exit-zero replacement. The detector runs
with the existing Task Graph checks, records a skip when `_tasks.md` is absent,
and documents the stable code as Inverted Verification Exit in `CONTEXT.md`.

Focused checks:

- Before implementation,
  `GOCACHE=/tmp/roundfix-0095-task04-go-cache go test ./internal/speccheck -run '^TestVerifyInvertedExit' -count=1`
  exited 1 because `CodeVerifyInvertedExit` and `InvertedExitVerification` did
  not exist.
- After implementation, the same focused test exited 0:
  `ok roundfix/internal/speccheck 0.360s`.
- `GOCACHE=/tmp/roundfix-0095-task04-go-cache go test ./internal/speccheck -count=1`
  exited 0: `ok roundfix/internal/speccheck 0.780s`.
- `rtk make verify-incremental` exited 0. It ran formatting checks, all Go
  package tests, skill sync/checks, and the build; `internal/speccheck` exited
  0 in 1.371s on the final run.
- `git diff --check` exited 0 with no diagnostics.

Acceptance evidence:

- The table-driven `TestVerifyInvertedExit` has distinct refusal cases for
  `grep -c`, `grep -v ... | wc -l`, and bare `test $(...)`; each asserts the
  emitted code, named form, and working replacement.
- The same test has independent permitted cases for `grep -q`, a compared
  `test`, a redirected count followed by a comparison, and a filtered count
  inside a compared `test`.
- `TestVerifyInvertedExitSkipsWithoutTaskGraph` checks the existing
  `no-taskgraph` fixture produces no inverted-exit finding and records the
  detector as skipped for missing `_tasks.md`.
- The staged registry places `CodeVerifyInvertedExit` immediately after
  `CodeVerifyWorkIndependent` at `StageTasks`; the stage-scope tests require it
  to remain absent from earlier-stage findings and named among earlier-stage
  skips.

The Daemon-owned Verification commands were not run in this Agent turn.
