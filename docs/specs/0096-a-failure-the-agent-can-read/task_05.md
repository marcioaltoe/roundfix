---
task: task_05
spec: 0096-a-failure-the-agent-can-read
status: pending # pending | in_progress | completed | failed — only implement-task changes this
type: backend
complexity: low
---

# Task 05: Say what the run budget bounds where it is set

## Overview

A maintainer setting a maximum Run duration cannot tell from the setting whether
it bounds wall-clock time from the Run's start or is evaluated at Work Item
boundaries. The single measured overrun has no established cause, so the
behaviour stays and its contract is stated where it is configured.

## Requirements

1. MUST state, in the configuration template the tool renders, what the maximum
   Run duration bounds and when it is evaluated.
2. MUST leave the evaluation behaviour unchanged; this Task documents, it does not
   move the check.
3. MUST leave the rendered template valid input to the loader that reads it.

## Subtasks

- [ ] State the contract in the rendered template.
- [ ] Prove the template still loads.

## Acceptance Criteria

- [ ] The rendered configuration states what the budget bounds and when it is
      evaluated.
- [ ] The rendered template round-trips through the loader.
- [ ] No budget evaluation code changed.

## Verification

- `go test -count=1 ./internal/config -run 'TestRenderedConfig' -v > /tmp/0096-t05.log 2>&1; s=$?; grep -q '^--- PASS: TestRenderedConfig' /tmp/0096-t05.log || { cat /tmp/0096-t05.log; exit 1; }; exit $s` — expected: exits 0 and the log names the passing template tests; fails today if no such test exists, in which case this Task adds it.
- `grep -A 3 'max_run_duration' internal/config/config.go | grep -qiE 'bounds|evaluated|wall.clock|work item' || { echo 'the rendered template still does not say what the budget bounds'; grep -B 2 -A 3 'max_run_duration' internal/config/config.go; exit 1; }` — expected: exits 0, proving the contract is stated beside the setting. Fails today.
- `go build -buildvcs=false -o /tmp/0096-t05-roundfix ./cmd/roundfix && /tmp/0096-t05-roundfix profiles show --json > /dev/null 2>&1; git diff --name-only HEAD -- internal/daemon internal/rounds > /tmp/0096-t05-behaviour.txt; test ! -s /tmp/0096-t05-behaviour.txt || { echo 'budget evaluation code changed, which this Task forbids:'; cat /tmp/0096-t05-behaviour.txt; exit 1; }; grep -A 3 'max_run_duration' internal/config/config.go | grep -qiE 'bounds|evaluated' || { echo 'nothing changed at all'; exit 1; }` — expected: exits 0, proving the built command still runs, no evaluation code moved, and the statement landed. Fails today on the last clause.

## Context

- interface: `internal/config/config.go`

## References

`_techspec.md` → Build Order 5; System Architecture, the configuration surface.
`_prd.md` → Core Feature 6; Open Questions. ADR-0137.
