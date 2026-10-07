---
type: fix # feat | fix | perf | refactor
status: done
created: 2026-08-14
spec: null # Spec slug when status: promoted
reason: null # required when status: declined
---

# Staging a Task commit fails when the Task deleted a file

An `implement` Run ended `Failed` because the Daemon could not commit a Task that had finished its work correctly. The staging command names every path the Task touched, including the ones it **deleted**, and `git add -f -- <path>` does not match a path absent from both the working tree and the index.

Full text in Git at `a30dc847a037fe584812b4471aa7da1701a9ec59`: `docs/history/backlog/2026-08-14-staging-a-task-commit-fails-when-the-task-deleted-a-file.md`.
