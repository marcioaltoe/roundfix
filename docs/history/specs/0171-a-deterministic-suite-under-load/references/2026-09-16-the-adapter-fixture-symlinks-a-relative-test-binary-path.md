---
status: done
created_at: 2026-09-16
updated_at: 2026-09-25
absorbed_by: 0171-a-deterministic-suite-under-load
---

# ACPX fixtures — The adapter fixture symlinks a relative test binary path (2026-09-16)

Spec 0139 replaced the per-test hard link to the compiled test binary with a
symlink, which removed the macOS kill that the link churn caused. The pre-PR
review of that delivery found a narrower defect the change introduced.

## 1. A relative `argv[0]` produces a link that cannot resolve

- Symptom / evidence:
  - `provisionFakeAdapter` in `internal/agent/acpx_runner_test.go` calls
    `os.Symlink(os.Args[0], path)` and stores whatever `argv[0]` holds.
  - Each adapter is created in its own temporary directory, so a relative target
    resolves against that directory rather than the caller's working directory.
  - Under `go test`, `argv[0]` is an absolute path, so the suite and CI are
    unaffected. Compiling first and running the binary directly, as in
    `go test -c ./internal/agent` followed by `./agent.test`, stores `./agent.test`
    and the fixture cannot execute.
  - The hard link this replaced had no such failure mode: it copied the inode
    rather than a path.
- Root cause: the fixture keeps a path where it previously kept an inode, and
  does not resolve it first.
- Action / suggestion:
  - Resolve `os.Args[0]` to an absolute path before creating the symlink.
  - Spec 0139 is archived and its gate reported, so this repair belongs to the
    next Spec that touches these fixtures rather than to a Task appended beneath
    a settled terminal gate.

## Addendum — 2026-09-25 — Revalidated in triage

Revalidated against main 7a9b6ec6: still holds: acpx_runner_test.go still symlinks os.Args[0]; a compiled test binary run by relative path fails, by absolute path passes. Ranked in the 2026-09-25 triage priority list.

## Addendum — 2026-09-25 — Adopted by Spec 0171

Adopted by Spec 0171-a-deterministic-suite-under-load. Its Task resolves
`os.Args[0]` to an absolute path before `provisionFakeAdapter` links it, and
proves the fixture runs from a compiled test binary started by a relative path.
