---
task: task_04
spec: 0095-a-verification-that-ran-before-anyone-believed-it
status: pending # pending | in_progress | completed | failed — only implement-task changes this
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
