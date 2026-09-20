# Authorization and Task-commit audit

The audit used `git diff-tree --no-commit-id --name-only -r <commit>` for every
authorization and Task commit. Git emitted an fsmonitor IPC warning while
reading paths but returned every requested path and all ancestry checks exited
`0`.

Authorization chronology:

- `2b73ac3f79489a2c772741b978a5dc0a6df68d47` introduced the authorization record.
- `03cfb8ee29c7c0ec9fb822a9defb251d49316709` widened it to the historically governed `internal/spec/archive.go` path.
- Both commits are ancestors of every consuming Task commit. The widening commit is the direct parent of Task 01 and precedes Task 02 in ancestry, even though its recorded wall-clock timestamp is later than the child commit's timestamp.

Task paths:

- Task 01 `7d9e94d8a7cd5802dcea4c2d0b1dfd4cac20b5e4`: its own Task file plus ordinary CLI/spec source and tests; no governed path.
- Task 02 `1335873901dc9f96083e47a37764229960f66177`: its own Task file, ordinary CLI/spec source and tests, and governed `internal/spec/archive.go`; that governed path is expressly bounded.
- Task 03 `60224664f49268383026f9b4a8bee15b89072b10`: its own Task file, ordinary user-guide documentation, and the expressly bounded canonical/mirror skill paths.

ADR-0130 makes the audit judge governed paths rather than treating ordinary Go
source, tests, and documentation as governed by proximity. No Task commit
changed a governed path outside the exact authorization. The authorization and
its widening were separate prior commits; no prerequisite or consequent fix
was folded into a consuming Task commit. No derived version pin changed.
Finally, the canonical and distributed skill files compare byte-identically.

