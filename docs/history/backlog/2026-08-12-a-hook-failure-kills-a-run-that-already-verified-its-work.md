---
type: fix # feat | fix | perf | refactor
status: done
created: 2026-08-12
spec: null # Spec slug when status: promoted
reason: null # required when status: declined
---

# A hook failure kills a Run whose work was already verified

The Daemon runs the authoritative Verification and then commits. When the repository's `pre-commit` hook refuses, `lint-staged` reverts and the Run ends `Failed`. The repair loop covers a Verification failure, not a hook failure.

Full text in Git at `a30dc847a037fe584812b4471aa7da1701a9ec59`: `docs/history/backlog/2026-08-12-a-hook-failure-kills-a-run-that-already-verified-its-work.md`.
