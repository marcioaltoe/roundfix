---
task: task_12
spec: 0103-a-suite-that-leaks-nothing
status: pending # pending | in_progress | completed | failed — only implement-task changes this
type: test
complexity: low
---

# Task 12: Size the stress case to the suite it runs inside

## Overview

The concurrent-exec stress case proves the compiled fixture survives a dense
start. It does so by hard-linking the test binary 512 times and executing every
copy, which passes in isolation and is killed by the host when the authoritative
suite runs it at `-parallel 16`. A stress case that dies of its own weight
measures the machine, not the fix.

## Requirements

1. MUST scale the stress to the parallelism actually available rather than to a
   fixed 16×32.
2. MUST still execute concurrently enough to exercise the link-then-execute
   window the case exists to prove.
3. MUST pass inside the authoritative repository gate, not only in isolation.
4. MUST NOT add a retry, a sleep, or a skip; a reduced size is not a weakened
   assertion.

## Subtasks

- [ ] Scale workers and executions to available parallelism.
- [ ] Keep the concurrency that exercises the window.
- [ ] Prove it under the authoritative gate.

## Acceptance Criteria

- [ ] The stress case passes under `make verify`.
- [ ] It still runs its executions concurrently across more than one worker.
- [ ] No retry, sleep, or skip was introduced.
- [ ] It still fails if a fixture is written and executed rather than linked.

## Verification

- `go test -count=1 -parallel 16 ./internal/agent -run 'TestFixtureBinarySurvivesConcurrentExec' -v > /tmp/0103-t12.log 2>&1; s=$?; grep -q '^--- PASS: TestFixtureBinarySurvivesConcurrentExec' /tmp/0103-t12.log || { cat /tmp/0103-t12.log; exit 1; }; grep -q 'signal: killed' /tmp/0103-t12.log && { echo 'the stress case is still killed under load'; exit 1; }; test $s -eq 0 || { cat /tmp/0103-t12.log; exit 1; }; grep -rq 'GOMAXPROCS\|runtime.NumCPU' internal/agent/acpx_runner_test.go || { echo 'the case passed in isolation but is still sized to a fixed sixteen'; exit 1; }` — expected: exits 0, proving the case survives the parallelism the gate uses. The isolated run passes today, so it is anchored to the sizing that makes it survive the loaded one.
- `grep -rq 'GOMAXPROCS\|runtime.NumCPU' internal/agent/acpx_runner_test.go || { echo 'the stress case still uses a fixed worker count'; exit 1; }; grep -n 'workers *= *16' internal/agent/acpx_runner_test.go && { echo 'the fixed sixteen-worker constant is still there'; exit 1; }; exit 0` — expected: exits 0, proving the size follows the machine. Fails today, where the constant is `workers = 16`.
- `! grep -rn 'time.Sleep\|t.Skip\|retry' internal/agent/acpx_runner_test.go > /tmp/0103-t12-shortcut.txt 2>&1; grep -n 'TestFixtureBinarySurvivesConcurrentExec' -A 40 internal/agent/acpx_runner_test.go | grep -E 'time.Sleep|t.Skip|retry' && { echo 'the stress case was weakened rather than sized'; exit 1; }; grep -rq 'GOMAXPROCS\|runtime.NumCPU' internal/agent/acpx_runner_test.go || { echo 'the stress case was not sized to the machine'; exit 1; }` — expected: exits 0, proving no shortcut was taken and the sizing happened. Fails today on the last clause.
- `make verify > /tmp/0103-t12-verify.log 2>&1; s=$?; grep -q 'FAIL' /tmp/0103-t12-verify.log && { echo 'the authoritative gate is not green:'; grep -B 3 -A 8 'FAIL' /tmp/0103-t12-verify.log | head -40; exit 1; }; test $s -eq 0 || { tail -30 /tmp/0103-t12-verify.log; exit 1; }; grep -rq 'GOMAXPROCS\|runtime.NumCPU' internal/agent/acpx_runner_test.go || { echo 'the gate is green, but the stress case is still sized to a fixed sixteen'; exit 1; }` — expected: exits 0, proving the case passes where it actually has to. The gate's failure is load-dependent and does not reproduce on every machine, so the green result is anchored to the sizing that makes it reliable.

## Context

- interface: `internal/agent/acpx_runner_test.go`

## References

`_techspec.md` → Testing Approach, the compiled fixture. `_prd.md` → Core
Feature 4; Non-Goals, no retry or loosened assertion. ADR-0125.
Evidence: this Spec's QA report `qa/qa-report-2026-08-14-02.md`, finding F-01.
