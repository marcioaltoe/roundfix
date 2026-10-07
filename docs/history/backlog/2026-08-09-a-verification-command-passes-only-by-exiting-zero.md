---
type: fix # feat | fix | perf | refactor
status: done
created: 2026-08-09
spec: null # Spec slug when status: promoted
reason: null # required when status: declined
---

# A Verification command passes only by exiting zero, and task authoring does not say so

The Daemon runs each Task's Verification through `sh -c` and treats a non-zero exit as failure — `internal/daemon/daemon.go:101`. Nothing in the task-authoring contract states that, and its absence is easy to author straight past.

Full text in Git at `a30dc847a037fe584812b4471aa7da1701a9ec59`: `docs/history/backlog/2026-08-09-a-verification-command-passes-only-by-exiting-zero.md`.
