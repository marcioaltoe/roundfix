---
status: done
created_at: 2026-07-27
updated_at: 2026-09-08
absorbed_by: 0054-tooling-task-and-verification-hygiene
---

# Build hygiene — `go build ./cmd/roundfix` drops a 20 MiB binary at the repository root, and nothing ignores it (2026-07-27)

`.gitignore` ignores `/bin/`, which is where the Makefile writes its build output. It does not ignore `/roundfix`, which is where `go build ./cmd/roundfix` writes when invoked without `-o`. Agents run the bare form routinely as a compile check, so the artifact appears unignored and untracked in a working tree that Daemon commits stage by diff.

Full text in Git at `a30dc847a037fe584812b4471aa7da1701a9ec59`: `docs/history/findings/2026-07-27-bare-go-build-writes-an-untracked-root-binary.md`.
