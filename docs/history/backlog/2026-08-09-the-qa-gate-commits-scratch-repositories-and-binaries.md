---
type: fix # feat | fix | perf | refactor
status: deferred
created: 2026-08-09
spec: null # Spec slug when status: promoted
reason: null # required when status: declined
---

# The QA gate commits scratch repositories and binaries into the Spec

A gate that proves a CLI behaviour builds isolated repositories to run it in, and builds the binary it exercises. Spec 0090's gate committed both into the Spec's evidence directory: nine directories recorded as Gitlinks (mode `160000`) and six executables — `roundfix` twice, `codex-acp` twice, `acpx` twice.

Full text in Git at `a30dc847a037fe584812b4471aa7da1701a9ec59`: `docs/history/backlog/2026-08-09-the-qa-gate-commits-scratch-repositories-and-binaries.md`.
