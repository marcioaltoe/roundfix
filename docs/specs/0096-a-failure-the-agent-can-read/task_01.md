---
task: task_01
spec: 0096-a-failure-the-agent-can-read
status: pending # pending | in_progress | completed | failed — only implement-task changes this
type: backend
complexity: medium
---

# Task 01: Say that the command produced nothing

## Overview

A Verification that redirects its output and then fails captures nothing, and the
repair prompt carries a command, an exit status, and no cause. In the measured
Task the Agent spent its one repair turn rewriting its own Task file with a
diagnosis it had deduced. An empty prompt and a prompt about nothing read the
same; only one of them tells the Agent to go look for the output.

## Requirements

1. MUST render, when the captured diagnostic is empty, text stating that the
   command produced no output.
2. MUST name where the output went when the command redirected it, so the Agent
   can read the file rather than infer a cause.
3. MUST render a non-empty diagnostic exactly as it renders today.
4. MUST be supplied by both Daemon call sites that build the repair prompt, so
   the state reaches a real Run rather than only the builder's tests.

## Subtasks

- [ ] Carry the absent-diagnostic state on the feedback.
- [ ] Render it in the prompt, naming the redirect target when there is one.
- [ ] Supply it from both Daemon call sites.

## Acceptance Criteria

- [ ] An empty diagnostic renders text saying the command produced no output.
- [ ] The redirect target is named when the command has one.
- [ ] A non-empty diagnostic renders unchanged.
- [ ] Both `internal/daemon` call sites populate the state.

## Verification

- `go test -count=1 ./internal/agent -run 'TestBuildVerificationRepairPrompt' -v > /tmp/0096-t01.log 2>&1; s=$?; grep -q '^--- PASS: TestBuildVerificationRepairPrompt' /tmp/0096-t01.log || { cat /tmp/0096-t01.log; exit 1; }; grep -q '^--- PASS: TestBuildVerificationRepairPromptStatesAnAbsentDiagnostic' /tmp/0096-t01.log || { echo 'the absent-diagnostic case does not exist'; cat /tmp/0096-t01.log; exit 1; }; exit $s` — expected: exits 0, proving the builder's existing cases still pass and the absent-diagnostic case was added. The case is named exactly: a looser pattern matched the unrelated `...IncludesPathFailureAndNoOutputBody`, which asserts the prompt *omits* the log body and passes on an untouched tree.
- `test -s /tmp/0096-t01.log || { echo 'the suite produced no output'; exit 1; }; grep -qi 'no tests to run' /tmp/0096-t01.log && { echo 'the suite selected no cases'; exit 1; }; grep -rq 'DiagnosticEmpty' internal/agent || { echo 'the feedback carries no absent-diagnostic state'; exit 1; }` — expected: exits 0, refusing a vacuous run and proving the state exists on the feedback. Fails today.
- `n=$(grep -rn 'DiagnosticEmpty' internal/daemon --include='*.go' | grep -v '_test.go' | wc -l | tr -d ' '); test "$n" -ge 2 || { echo "expected both Daemon call sites to supply the state, found $n:"; grep -rn 'DiagnosticEmpty' internal/daemon --include='*.go' | grep -v '_test.go'; exit 1; }` — expected: exits 0, proving the state reaches production rather than only the builder. Fails today, where the count is zero.
- `go test -count=1 ./internal/agent ./internal/daemon > /tmp/0096-t01-regress.log 2>&1; s=$?; grep -q 'FAIL' /tmp/0096-t01-regress.log && { echo 'the prompt or its callers regressed:'; grep -B 3 -A 8 'FAIL' /tmp/0096-t01-regress.log | head -30; exit 1; }; grep -rq 'DiagnosticEmpty' internal/daemon || { echo 'both packages pass, but no call site supplies the state'; exit 1; }; exit $s` — expected: exits 0, proving the prompt and the Daemon that builds it agree, anchored to the call site so a green pair cannot pass untouched.

## Context

- interface: `internal/agent/agent.go`
- interface: `internal/daemon/task_engine.go`

## References

`_techspec.md` → Build Order 1; System Architecture, the feedback prompt.
`_prd.md` → Core Feature 1; Goal 1; User Story 1. ADR-0135, ADR-0111.
