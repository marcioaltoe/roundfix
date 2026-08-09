---
task: task_01
spec: ordinary-failure
status: failed
type: backend
---

# Task 01: Observe a real non-zero verdict

## Verification

- `test -f "$(git rev-parse --git-path roundfix-qa-agent-task_01.done)"` — expected: the no-op Agent leaves a real exit-1 verdict after both turns.
