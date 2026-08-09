---
task: task_01
spec: positive-control
status: completed
type: backend
---

# Task 01: Exercise the positive control

## Verification

- `test -f "$(git rev-parse --git-path roundfix-qa-agent-task_01.done)"` — expected: fails before the Agent and passes after its controlled work.
