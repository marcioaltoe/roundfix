---
task: task_01
spec: 0001-capacity-flow
status: pending
type: backend
---

# Task 01: First public capacity slice

## Verification

- `printf 'start task_01\n' >> "$QA_COORD_DIR/verification.log"; sleep 0.25; printf 'end task_01\n' >> "$QA_COORD_DIR/verification.log"` — expected: the complete real shell command runs under Daemon Verification.
