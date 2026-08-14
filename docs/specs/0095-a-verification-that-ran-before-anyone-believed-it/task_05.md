---
task: task_05
spec: 0095-a-verification-that-ran-before-anyone-believed-it
status: pending # pending | in_progress | completed | failed — only implement-task changes this
type: backend
complexity: medium
---

# Task 05: Refuse a Verification that reaches outside the repository

## Overview

A Verification that depends on state the repository does not carry produces a
false red: the work is right and the gate fails. Three forms were measured in one
fleet session, and all three cost Tasks that had finished correctly — an
undeclared environment variable, a `test -n "$VAR" &&` guard in front, and a
dependency on a temporary directory or tree snapshot outside the repository.

## Requirements

1. MUST refuse an authored Verification command that references an environment
   variable the repository does not declare.
2. MUST refuse a command that guards itself on an environment variable's
   presence, since an absent variable fails the whole chain rather than the work.
3. MUST refuse a command that depends on a path outside the repository, unless
   the Task itself creates that path.
4. MUST NOT refuse a command that writes to and reads from a temporary path it
   creates within the same command, which is the working redirect form the
   authoring contract teaches.
5. MUST skip rather than fail when the artifact it reads is absent.
6. MUST register at the tasks stage.

## Subtasks

- [ ] Add the refusal code and its detector.
- [ ] Cover each measured non-hermetic form and the working redirect form.
- [ ] Register it at the tasks stage.

## Acceptance Criteria

- [ ] Each of the three measured forms is refused and named.
- [ ] The working redirect form the authoring contract teaches is not refused.
- [ ] A Spec without a task graph skips rather than fails.
- [ ] The code appears in the staged registry at the tasks stage.

## Verification

- `grep -q 'SC-VERIFY-NON-HERMETIC' internal/speccheck/*.go && grep -q 'CodeVerifyNonHermetic' internal/speccheck/coherence.go` — expected: exits 0, proving the code exists and is registered. Fails today on both clauses.
- `go test -count=1 ./internal/speccheck -run 'TestVerifyNonHermetic' -v > /tmp/0095-t05.log 2>&1; s=$?; grep -q '^--- PASS: TestVerifyNonHermetic' /tmp/0095-t05.log || { cat /tmp/0095-t05.log; exit 1; }; exit $s` — expected: exits 0 and the log names the passing test; fails today.
- `grep -c 'PASS' /tmp/0095-t05.log > /tmp/0095-t05-n.txt; test "$(cat /tmp/0095-t05-n.txt)" -ge 4 || { echo "expected at least four passing cases, got $(cat /tmp/0095-t05-n.txt)"; cat /tmp/0095-t05.log; exit 1; }` — expected: exits 0, proving the three refused forms and the permitted redirect form each run as their own case.

## References

`_techspec.md` → Build Order 5; Testing Approach, the static detectors.
`_prd.md` → Core Feature 3; User Story 4. ADR-0093, ADR-0094.
