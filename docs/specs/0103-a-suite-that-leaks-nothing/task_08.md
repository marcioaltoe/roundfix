---
task: task_08
spec: 0103-a-suite-that-leaks-nothing
status: pending # pending | in_progress | completed | failed — only implement-task changes this
type: backend
complexity: low
---

# Task 08: Keep a gate's scratch state out of the evidence

## Overview

A QA gate builds scratch repositories and binaries to prove what it claims. Those
belong to the gate, not to the Spec's evidence: a committed binary is not
something a reader can check, and a scratch repository committed as a gitlink is
a reference to a tree nobody else has.

## Requirements

1. MUST refuse a Spec evidence directory that holds a built binary or a
   submodule-shaped directory reference.
2. MUST name the offending path and what kind it is.
3. MUST skip rather than fail where the Spec has no evidence directory.
4. MUST NOT refuse an ordinary evidence artifact a reader can open.

## Subtasks

- [ ] Add the refusal for binaries and gitlinks under a Spec's evidence.
- [ ] Cover both offending kinds and an ordinary artifact.

## Acceptance Criteria

- [ ] A committed binary under a Spec's evidence is refused and named.
- [ ] A gitlink under a Spec's evidence is refused and named.
- [ ] A text or JSON artifact is not refused.
- [ ] A Spec with no evidence directory skips rather than fails.

## Verification

- `go test -count=1 ./internal/speccheck -run 'TestEvidenceRefusesScratchState' -v > /tmp/0103-t08.log 2>&1; s=$?; grep -q '^--- PASS: TestEvidenceRefusesScratchState' /tmp/0103-t08.log || { cat /tmp/0103-t08.log; exit 1; }; exit $s` — expected: exits 0 and the log names the passing test; fails today, where no such test exists.
- `! grep -qi 'no tests to run' /tmp/0103-t08.log` — expected: exits 0, refusing a vacuous run.
- `grep -c '^--- PASS' /tmp/0103-t08.log > /tmp/0103-t08-n.txt; test "$(cat /tmp/0103-t08-n.txt)" -ge 4 || { echo "expected the binary, gitlink, ordinary-artifact, and absent-directory cases, got $(cat /tmp/0103-t08-n.txt)"; cat /tmp/0103-t08.log; exit 1; }` — expected: exits 0, proving each case runs separately.

## Context

- interface: `internal/speccheck/mechanical.go`

## References

`_techspec.md` → Build Order 8. `_prd.md` → Core Feature 7; Goal 5; User
Story 5; Success Metrics.
