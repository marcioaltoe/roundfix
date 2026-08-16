---
task: task_09
spec: 0096-a-failure-the-agent-can-read
status: pending # pending | in_progress | completed | failed — only implement-task changes this
type: backend
complexity: low
---

# Task 09: Stop the failure cause being shadowed away

## Overview

A Spec written to give an Agent a cause it can read now reports
`verification failed: <nil>`. In `internal/daemon/engine.go` the publish step uses
`metadata, err := req.publishFailedCommand(...)`, which declares a new `err` in the
inner scope and shadows the verifier's error; the failure string then formats the
inner one, which is nil whenever publishing succeeded. Three tests catch it and
the repository gate is red.

## Requirements

1. MUST report the verifier's cause in the Verification failure string, not the
   publish step's error.
2. MUST still return the publish step's error when publishing fails.
3. MUST leave the Repeated metadata and the diagnostics line unchanged.
4. MUST leave the three tests that caught this passing, and `make verify` green.

## Subtasks

- [ ] Stop shadowing the verifier error at the publish call.
- [ ] Keep the publish failure path returning its own error.

## Acceptance Criteria

- [ ] A failed Verification reports the verifier's cause, such as an exit status.
- [ ] A publish failure still returns that error and no outcome.
- [ ] `TestResolveCycleVerificationFailureFailsBatchAndContinues`,
      `TestRunOutcomeDerivedFromUnresolvedIssues` and
      `TestRunResolveVerificationFailureDoesNotCommit` pass.
- [ ] `make verify` exits 0.

## Verification

- `go test -count=1 ./internal/daemon -run '^(TestResolveCycleVerificationFailureFailsBatchAndContinues|TestRunOutcomeDerivedFromUnresolvedIssues)$' -v > /tmp/0096-t09.log 2>&1; s=$?; for t in TestResolveCycleVerificationFailureFailsBatchAndContinues TestRunOutcomeDerivedFromUnresolvedIssues; do grep -q "^--- PASS: $t" /tmp/0096-t09.log || { echo "FAIL: $t"; cat /tmp/0096-t09.log; exit 1; }; done; exit $s` — expected: exits 0 and the log names both passing tests; both fail today with `<nil>`.
- `go test -count=1 ./internal/cli -run '^TestRunResolveVerificationFailureDoesNotCommit$' -v > /tmp/0096-t09-cli.log 2>&1; s=$?; grep -q '^--- PASS: TestRunResolveVerificationFailureDoesNotCommit' /tmp/0096-t09-cli.log || { cat /tmp/0096-t09-cli.log; exit 1; }; exit $s` — expected: exits 0, proving the public journey reports a cause again. Fails today.
- `grep -n 'metadata, err := req.publishFailedCommand' internal/daemon/engine.go && { echo 'the publish step still shadows the verifier error'; exit 1; }; grep -q 'publishFailedCommand' internal/daemon/engine.go || { echo 'the publish call vanished rather than being unshadowed'; exit 1; }; exit 0` — expected: exits 0, proving the shadowing is gone and the call survives. It prints the offending line on failure. Fails today.
- `grep -n 'verification failed: %v' internal/daemon/engine.go > /tmp/0096-t09-fmt.txt; test -s /tmp/0096-t09-fmt.txt || { echo 'the failure string was removed rather than corrected'; exit 1; }; grep -n 'metadata, err :=' internal/daemon/engine.go && { echo 'an inner err still shadows at the publish call'; exit 1; }; exit 0` — expected: exits 0, proving the failure string still exists and no longer formats a shadowed error. Fails today.

## Context

- interface: `internal/daemon/engine.go`

## References

`_techspec.md` → System Architecture, the feedback prompt. `_prd.md` → Core
Feature 1; Goal 1; User Story 1. ADR-0135.
Evidence: this Spec's QA report `qa/qa-report-2026-08-16.md`, finding F-001.
