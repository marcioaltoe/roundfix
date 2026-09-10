---
task: task_04
spec: failed-recovery
status: pending
type: backend
complexity: low
---

# Task 04: Repair the finding

## Overview

Represent a corrective Task authored after a failed QA gate.

## Requirements

1. MUST run before the failed gate runs again.
2. MUST leave Task status settlement to the Daemon.

## Acceptance Criteria

- [ ] The Implement Command dispatches this Task.
- [ ] Verification observes the Agent marker and the Daemon settles the Task.

## Verification

- `test -f "$(git rev-parse --git-path 'roundfix-test-agent-task_04.done')"` — the public flow's controlled ACP adapter records that this corrective Task was dispatched.

## References

- `_prd.md` → Goal 1 and User Story 1.
- `_techspec.md` → Implementation Design.
