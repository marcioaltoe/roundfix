---
status: accepted
created_at: 2026-09-24T00:00:00Z
updated_at: 2026-09-24T00:00:00Z
deprecated_at: null
superseded_by: null
---

# Only the frozen authorization opens a red repository gate

The Daemon runs the configured repository Verification command before a Task.
When that command is red, an ordinary Task cannot start. A Task written to
repair the red gate needs an explicit, bounded entry without allowing any
other Task to bypass the precondition.

At Run start, Roundfix resolves the authorization record and freezes its
`precondition_repairs` list. Only Tasks named in that list may enter while the
repository gate is red. A Task file, Agent edit, or Run Worktree cannot add
permission after the Run starts. A named Task settles completed only when all
of its Verification commands pass, including the same configured repository
command that was red on entry.

Tasks not named in the frozen authorization remain blocked by the red
repository gate. This keeps the authorization boundary independent of mutable
worktree contents and makes the repair permission reviewable before the Run
begins.

Deriving permission from a Task file or an Agent edit was rejected: either
would let work inside the Run widen the set of Tasks allowed to enter. Reading
only the authorization resolved at Run start preserves the maintainer's
decision and keeps the gate's entry and settlement rules symmetric.

The accepted cost is that a repair Task omitted from `precondition_repairs`
cannot begin until a new Run resolves an authorization that names it.

## References

- [ADR-0014: Daemon runs Task Verification and settles status](0014-daemon-runs-task-verification-and-settles-status.md)
- [ADR-0057: Daemon exclusively owns Implement Task status](0057-daemon-exclusively-owns-implement-task-status.md)
- [Spec 0158 technical specification](../specs/0158-daemon-verification-and-access-readiness/_techspec.md)
