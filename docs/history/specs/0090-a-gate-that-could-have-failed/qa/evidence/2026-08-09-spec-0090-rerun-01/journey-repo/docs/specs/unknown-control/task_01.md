---
task: task_01
spec: unknown-control
status: pending
type: backend
---

# Task 01: Make the probe unobservable

## Verification

- `test -f "$(git rev-parse --git-path roundfix-qa-agent-task_01.done)"` — expected: the runner cannot find `sh` and records unknown.
