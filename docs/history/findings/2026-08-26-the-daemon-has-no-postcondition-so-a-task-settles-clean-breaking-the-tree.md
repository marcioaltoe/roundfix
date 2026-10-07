---
status: done
absorbed_by: 0122-verified-content-and-terminal-settlement
created_at: 2026-08-26
updated_at: 2026-09-08
kind: finding
---

# The Daemon has a green-tree precondition and no postcondition, so a Task settles Clean while breaking the tree

A Task's scoped Verification passed and the Task settled Clean; in the same tree state the Supervisor's `make verify` exited 2 with 14 typecheck errors born from the change (baseline the same day: main exits 0). The Agent classified them as pre-existing and deferred them. The Daemon checks tree health on entry and never on exit, so a Task can hand the next Task a broken tree with a Clean verdict on it.

Full text in Git at `a30dc847a037fe584812b4471aa7da1701a9ec59`: `docs/history/findings/2026-08-26-the-daemon-has-no-postcondition-so-a-task-settles-clean-breaking-the-tree.md`.
