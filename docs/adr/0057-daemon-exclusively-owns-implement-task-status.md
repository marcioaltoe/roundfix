---
status: accepted
created_at: 2026-07-21T20:44:54Z
updated_at: 2026-08-11T21:35:56Z
deprecated_at: null
superseded_by: null
---

# Daemon exclusively owns Implement Task status

During an Implement Run, the Daemon is the sole writer of a Task's status. It
writes `in_progress` at Task start and writes `completed` or `failed` only
after the applicable Agent, Verification, repair, infrastructure, and
settlement outcome. The Task file remains the durable owner of the status
field; this decision changes the writer, not the storage location.

Initial and Verification Feedback Agents hand back implementation-ready work,
may record focused-check evidence in the Result section, and must not run the
Task's declared `## Verification` commands or claim a terminal Task outcome.
After every successful Agent handoff, the Daemon reloads the Task and
normalizes any Agent-authored status back to its current `in_progress` state.
An Agent-authored `completed` or `failed` value is therefore never a verdict
and cannot bypass Daemon Verification. Genuine Agent execution,
infrastructure, or unreadable-artifact failures may still be settled failed by
the Daemon without a Verification command because no implementation-ready
handoff exists.

This supersedes ADR-0014's permission for the Agent to write Task status while
preserving its Daemon-run Verification and final-settlement requirements. It
also supersedes the pre-Verification Agent-failure shortcut retained by Spec
0024. The existing same-Session Verification Feedback bound, failed Task
Worktree preservation, independent Task continuation, dependency blocking,
and Settle Command recovery remain unchanged.
